// Copyright 2016 The Mellium Contributors.
// Use of this source code is governed by the BSD 2-clause
// license that can be found in the LICENSE file.

package xmpp

import (
	"context"
	"crypto/tls"
	"encoding/xml"
	"fmt"
	"io"

	"mellium.im/xmlstream"
	"mellium.im/xmpp/internal/ns"
)

// StartTLS returns a new stream feature that can be used for negotiating TLS.
// If cfg is nil, a default configuration is used that uses the domainpart of
// the sessions local address as the ServerName.
func StartTLS(cfg *tls.Config) StreamFeature {
	return StreamFeature{
		Name:       xml.Name{Local: "starttls", Space: ns.StartTLS},
		Prohibited: Secure,
		List: func(ctx context.Context, e xmlstream.TokenWriter, start xml.StartElement) (req bool, err error) {
			if err = e.EncodeToken(start); err != nil {
				return true, err
			}
			startRequired := xml.StartElement{Name: xml.Name{Space: "", Local: "required"}}
			if err = e.EncodeToken(startRequired); err != nil {
				return true, err
			}
			if err = e.EncodeToken(startRequired.End()); err != nil {
				return true, err
			}
			return true, e.EncodeToken(start.End())
		},
		Parse: func(ctx context.Context, d *xml.Decoder, start *xml.StartElement) (bool, interface{}, error) {
			parsed := struct {
				XMLName  xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-tls starttls"`
				Required struct {
					XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-tls required"`
				}
			}{}
			err := d.DecodeElement(&parsed, start)
			return parsed.Required.XMLName.Local == "required" && parsed.Required.XMLName.Space == ns.StartTLS, nil, err
		},
		Negotiate: func(ctx context.Context, session *Session, data interface{}) (SessionState, io.ReadWriter, error) {
			conn := session.Conn()
			state := session.State()
			r := session.TokenReader()
			defer r.Close()
			d := xml.NewTokenDecoder(r)

			// If no TLSConfig was specified, use a default config.
			if cfg == nil {
				cfg = &tls.Config{
					ServerName: session.LocalAddr().Domain().String(),
					MinVersion: tls.VersionTLS12,
				}
			}

			var rw io.ReadWriter
			if (state & Received) == Received {
				fmt.Fprint(conn, `<proceed xmlns='urn:ietf:params:xml:ns:xmpp-tls'/>`)
				rw = tls.Server(conn, cfg)
			} else {
				// Select starttls for negotiation.
				fmt.Fprint(conn, `<starttls xmlns='urn:ietf:params:xml:ns:xmpp-tls'/>`)

				// Receive a <proceed/> or <failure/> response from the server.
				t, err := d.Token()
				if err != nil {
					return 0, nil, err
				}
				switch tok := t.(type) {
				case xml.StartElement:
					switch {
					case tok.Name.Space != ns.StartTLS:
						return 0, nil, fmt.Errorf("xmpp: unknown namespace during TLS negotiation")
					case tok.Name.Local == "proceed":
						// Skip the </proceed> token.
						if err = d.Skip(); err != nil {
							return 0, nil, err
						}
						rw = tls.Client(conn, cfg)
					case tok.Name.Local == "failure":
						// Skip the </failure> token.
						if err = d.Skip(); err != nil {
							return 0, nil, err
						}
						return 0, nil, fmt.Errorf("xmpp: receiver indicated that TLS negotiation failed")
					default:
						return 0, nil, fmt.Errorf("xmpp: unknown element during TLS negotiation")
					}
				default:
					return 0, nil, fmt.Errorf("xmpp: disallowed XML sent during TLS negotiation")
				}
			}
			return Secure, rw, nil
		},
	}
}
