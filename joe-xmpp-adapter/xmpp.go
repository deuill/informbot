package xmpp

import (
	// Standard library
	"context"
	"crypto/rand"
	"crypto/tls"
	"encoding/xml"
	"fmt"
	"io"
	"strings"

	// Third-party packages
	"github.com/go-joe/joe"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"mellium.im/sasl"
	"mellium.im/xmlstream"
	"mellium.im/xmpp"
	"mellium.im/xmpp/dial"
	"mellium.im/xmpp/jid"
	"mellium.im/xmpp/stanza"
)

// DefaultAuthMechanisms represents the list of SASL authentication mechanisms this client is allowed
// to use in server authentication.
var defaultAuthMechanisms = []sasl.Mechanism{
	sasl.Plain,
	sasl.ScramSha1,
	sasl.ScramSha1Plus,
}

// Config represents required and optional configuration values used in setting up the XMPP bot client.
type Config struct {
	// Required configuration.
	JID      string
	Password string

	// Optional configuration.
	NoTLS       bool // Whether to disable TLS connection to the XMPP server.
	UseStartTLS bool // Whether or not connection will be allowed to be made over StartTLS.

	// Other fields.
	Logger *zap.Logger // The instance to use for emitting log messages.
}

// Client represents an active XMPP session against a server, and configuration for handling messages
// against a Joe instance.
type Client struct {
	brain   *joe.Brain    // The mediator between this adapter and other handlers.
	session *xmpp.Session // The active XMPP session.
	logger  *zap.Logger   // The logger instance to use, defaults to a global logger used by Joe.
}

// Send wraps the given text in a message stanza and sets the recipient to the given channel, which
// is expected to be a JID (bare for direct messages). A error is returned if the channel JID does
// not parse, or if the message fails to send for any reason.
func (c *Client) Send(msg, channel string) error {
	jid, err := jid.Parse(channel)
	if err != nil {
		return errors.Wrap(err, "parsing JID failed")
	}

	// Determine whether this is a direct or group-chat message from the resource part of the JID,
	// which is only set if the message was originally sent as part of a group-chat.
	var kind = stanza.ChatMessage
	if jid.Resourcepart() != "" {
		msg = jid.Resourcepart() + ", " + msg
		jid, kind = jid.Bare(), stanza.GroupChatMessage
	}

	c.logger.Debug("Sending message",
		zap.String("jid", jid.String()),
		zap.String("type", string(kind)))

	return c.session.Send(context.Background(),
		xmlstream.Wrap(
			xmlstream.Wrap(
				xmlstream.Token(xml.CharData(msg)),
				xml.StartElement{Name: xml.Name{Local: "body"}},
			),
			xml.StartElement{
				Name: xml.Name{Local: "message"},
				Attr: []xml.Attr{
					{Name: xml.Name{Local: "id"}, Value: randomID()},
					{Name: xml.Name{Local: "to"}, Value: jid.String()},
					{Name: xml.Name{Local: "type"}, Value: string(kind)},
				},
			},
		),
	)
}

// GroupInfo represents information needed for joining a MUC, either automatically or as part of an
// invite (direct or mediated).
type GroupInfo struct {
	Channel  jid.JID `xml:"-"`
	Password string  `xml:"password`
	Invite   struct {
		From jid.JID `xml:"from,attr"`
	} `xml:"invite"`
}

// MessageStanza represents an XMPP message stanza, commonly used for transferring chat messages
// among users or group-chats.
type MessageStanza struct {
	// Base, common fields.
	stanza.Message
	Body string `xml:"body"`

	// Additional, optional fields.
	Group GroupInfo `xml:"x"`
}

// HandleInvite responds to the given invite (direct or mediated) with an 'available' presence,
// which allows the client to participate in MUCs.
func (c *Client) HandleInvite(w xmlstream.TokenWriter, info *GroupInfo) error {
	jid, err := info.Channel.WithResource(c.session.LocalAddr().Localpart())
	if err != nil {
		return errors.Wrap(err, "setting JID for MUC failed")
	}

	_, err = xmlstream.Copy(w, xmlstream.Wrap(
		xmlstream.Wrap(
			xmlstream.MultiReader(
				xmlstream.Wrap(
					xmlstream.Token(xml.CharData(info.Password)),
					xml.StartElement{Name: xml.Name{Local: "password"}},
				),
				xmlstream.Wrap(nil, xml.StartElement{
					Name: xml.Name{Local: "history"},
					Attr: []xml.Attr{
						{Name: xml.Name{Local: "maxchars"}, Value: "0"},
					},
				}),
			),
			xml.StartElement{
				Name: xml.Name{Local: "x"},
				Attr: []xml.Attr{
					{Name: xml.Name{Local: "xmlns"}, Value: "http://jabber.org/protocol/muc"},
				},
			},
		),
		stanza.Presence{
			ID:   randomID(),
			Type: stanza.AvailablePresence,
			To:   jid,
		}.StartElement(),
	))

	if err != nil {
		return errors.Wrap(err, "setting presence for MUC failed")
	}

	return nil
}

// HandleMessage parses the given MessageStanza, validating its contents and responding either as a
// direct message, or as a group-chat mention, depending on the intent. HandleMessage will also handle
// invites to group-chats, joining these automatically and with no confirmation needed.
//
// By default, only messages prepended with the local part of the client JID will be responded to in
// group-chats; this is to avoid handling messages where this is not wanted. Such mentions will be,
// in turn, responded to with a mention for the sending user.
//
// Currently, only mediated invites (XEP-0045) are handled, and rooms are not re-joined if the client
// closes its connection to the server.
func (c *Client) HandleMessage(w xmlstream.TokenWriter, msg *MessageStanza) error {
	var authorID = msg.From.Bare().String()
	var channel = msg.From.Bare().String()

	switch msg.Type {
	case stanza.GroupChatMessage:
		// Don't handle messages that aren't intended for us.
		n := strings.ToLower(c.session.LocalAddr().Localpart())
		if len(msg.Body) <= len(n) || strings.ToLower(msg.Body[:len(n)]) != n {
			return nil
		}

		channel = msg.From.String()
		msg.Body = strings.Trim(msg.Body[len(n):], " ,:")
		fallthrough
	case stanza.ChatMessage:
		// Do not attempt to handle empty or invalid messages.
		if msg.Body == "" {
			return nil
		}

		c.brain.Emit(joe.ReceiveMessageEvent{
			ID:       msg.ID,
			Text:     msg.Body,
			AuthorID: authorID,
			Channel:  channel,
			Data:     msg,
		})
	default:
		// Check if message is a mediated MUC invite, and join MUC if so.
		if !msg.Group.Invite.From.Equal(jid.JID{}) {
			msg.Group.Channel = msg.From.Bare()
			return c.HandleInvite(w, &msg.Group)
		}
	}

	return nil
}

// PresenceStanza represents an XMPP presence stanza, commonly used for communicating
// availability.
type PresenceStanza struct {
	// Base, common fields.
	stanza.Presence
}

// HandlePresence parses the given PresenceStanza and responds (usually to the affirmative),
// depending on the presence type, e.g. for subscription requests, HandlePresence will automatically
// subscribe and respond. Any errors returned in parsing on responding will be returned.
func (c *Client) HandlePresence(w xmlstream.TokenWriter, p *PresenceStanza) error {
	var err error

	// Handle presence stanza based on type.
	switch p.Type {
	case stanza.SubscribePresence:
		// Respond to subscription requests automatically.
		_, err = xmlstream.Copy(w, stanza.Presence{
			ID:   randomID(),
			Type: stanza.SubscribedPresence,
			To:   p.From,
		}.Wrap(nil))
	}

	if err != nil {
		return err
	}

	return nil
}

// HandleXMPP parses incoming XML tokens and calls a corresponding handler, e.g. HandleMessage, for
// the stanza type represented. Unhandled stanza types will be ignored with no error returned.
func (c *Client) HandleXMPP(t xmlstream.TokenReadEncoder, start *xml.StartElement) error {
	var stanza interface{}
	var err error

	switch start.Name.Local {
	case "message":
		stanza = &MessageStanza{}
	case "presence":
		stanza = &PresenceStanza{}
	default:
		c.logger.Debug("Ignoring unknown stanza type", zap.String("type", start.Name.Local))
		return nil // Unknown stanza type, do not handle.
	}

	err = xml.NewTokenDecoder(t).DecodeElement(&stanza, start)
	if err != nil && err != io.EOF {
		c.logger.Error("Decoding element failed", zap.Error(err))
		return nil
	}

	switch start.Name.Local {
	case "message":
		if err := c.HandleMessage(t, stanza.(*MessageStanza)); err != nil {
			c.logger.Error("Handling message failed", zap.Error(err))
		}
	case "presence":
		if err := c.HandlePresence(t, stanza.(*PresenceStanza)); err != nil {
			c.logger.Error("Handling presence failed", zap.Error(err))
		}
	}

	return nil
}

// RegisterAt sets the Joe Brain instance for the XMPP client.
func (c *Client) RegisterAt(brain *joe.Brain) {
	c.brain = brain
}

// Close shuts down the active XMPP session and server connection, returning an error if the process
// fails at any point.
func (c *Client) Close() error {
	if err := c.session.Close(); err != nil {
		return err
	}

	if err := c.session.Conn().Close(); err != nil {
		return err
	}

	return nil
}

// Adapter initializes an XMPP client connection according to configuration given, and returns a Joe
// module, usable in calls to joe.New(), or an error if any occurs.
func Adapter(conf Config) joe.Module {
	return joe.ModuleFunc(func(joeConf *joe.Config) error {
		// Parse and set up JID.
		id, err := jid.Parse(conf.JID)
		if err != nil {
			return errors.Wrap(err, "parsing JID failed")
		}

		var ctx = context.Background()
		var dialer = &dial.Dialer{NoTLS: conf.NoTLS}

		// Initialze connection according to configuration.
		conn, err := dialer.Dial(ctx, "tcp", id)
		if err != nil {
			return errors.Wrap(err, "establishing connection failed")
		}

		// Enable optional features and initialize client session, according to configuration.
		features := []xmpp.StreamFeature{xmpp.BindResource()}
		if conf.UseStartTLS {
			features = append(features, xmpp.StartTLS(true, &tls.Config{ServerName: id.Domain().String()}))
		}
		if conf.Password != "" {
			features = append(features, xmpp.SASL("", conf.Password, defaultAuthMechanisms...))
		}

		sess, err := xmpp.NewClientSession(ctx, id, conn, false, features...)
		if err != nil {
			return errors.Wrap(err, "establishing session failed")
		}

		var c = &Client{session: sess, logger: conf.Logger}
		if c.logger == nil {
			c.logger = joeConf.Logger(id.Network())
		}

		// Send initial presence to let the server know we want to receive messages.
		err = c.session.Send(ctx, stanza.Presence{Type: stanza.AvailablePresence}.Wrap(nil))
		if err != nil {
			return errors.Wrap(err, "setting initial presence failed")
		}

		go c.session.Serve(c)
		joeConf.SetAdapter(c)

		return nil
	})
}

// RandomID returns a cryptographically secure, 16-byte random string, useful for adding to stanzas
// for uniquely identifying them.
func randomID() string {
	var buf = make([]byte, 16)
	rand.Reader.Read(buf)
	return fmt.Sprintf("%x", buf)[:16]
}
