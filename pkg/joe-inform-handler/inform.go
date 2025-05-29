package inform

import (
	// Standard library
	"bytes"
	"context"
	"os/exec"
	"strings"
	"text/template"

	// Third-party packages
	"github.com/go-joe/joe"
	"github.com/pkg/errors"
)

// The prefix used for all keys stored.
const keyPrefix = "org.deuill.informbot"

type Inform struct {
	sessions map[string]*Session // A list of open sessions, against their authors.

	bot    *joe.Bot // The initialized bot to read commands from and send responses to.
	config *Config  // The configuration for the Inform bot.
}

func (n *Inform) SayTemplate(channel string, template *template.Template, data interface{}) error {
	var buf bytes.Buffer
	if err := template.Execute(&buf, data); err != nil || buf.Len() == 0 {
		n.bot.Say(channel, messageUnknownError)
		return err
	}

	n.bot.Say(channel, buf.String())
	return nil
}

func (n *Inform) Handle(ctx context.Context, ev joe.ReceiveMessageEvent) error {
	// Validate event data.
	if ev.AuthorID == "" {
		n.bot.Say(ev.Channel, messageUnknownError)
		return nil
	}

	// Check for stored rule-set against author ID, and send welcome message if none was found.
	var author, authorKey = &Author{}, keyPrefix + ".author." + ev.AuthorID
	if ok, err := n.bot.Store.Get(authorKey, author); err != nil {
		n.bot.Say(ev.Channel, messageUnknownError)
		return err
	} else if !ok {
		if err := n.SayTemplate(ev.Channel, templateWelcome, nil); err != nil {
			return errors.Wrap(err, "failed storing author information")
		}

		// Create and store new Author representation.
		author = NewAuthor(ev.AuthorID)
		if err := n.bot.Store.Set(authorKey, author); err != nil {
			n.bot.Say(ev.Channel, messageUnknownError)
			return err
		}
	}

	// Check for open session, and handle command directly if not prefixed.
	var cmd = ev.Text
	if n.sessions[author.ID] != nil {
		if !strings.HasPrefix(ev.Text, author.Options.Prefix) {
			if err := n.sessions[author.ID].Run(cmd); err != nil {
				n.bot.Say(ev.Channel, messageRunError, err)
				return err
			}
			n.bot.Say(ev.Channel, n.sessions[author.ID].Output())
			return nil
		} else {
			cmd = ev.Text[len(author.Options.Prefix):]
		}
	}

	// Handle meta-commands.
	var fields = strings.Fields(cmd)
	if len(fields) == 0 {
		return nil
	} else if len(fields) == 1 {
		cmd = fields[0]
	} else if len(fields) >= 2 {
		cmd = strings.Join(fields[:2], " ")
	}

	switch strings.ToLower(cmd) {
	case "help", "h":
		return n.SayTemplate(ev.Channel, templateHelp, nil)
	case "story", "stories", "story list", "list stories", "s":
		return n.SayTemplate(ev.Channel, templateStoryList, author)
	case "story add", "stories add", "add stories":
		if len(fields) < 4 {
			n.bot.Say(ev.Channel, messageUnknownStory)
		} else if story, err := author.AddStory(fields[2], fields[3]); err != nil {
			n.bot.Say(ev.Channel, messageInvalidStory, err)
		} else if err := story.Compile(ctx, n.config); err != nil {
			n.bot.Say(ev.Channel, "TODO: Compilation error: "+err.Error())
			return err
		} else if err = n.bot.Store.Set(authorKey, author); err != nil {
			n.bot.Say(ev.Channel, messageUnknownError)
			return err
		} else {
			n.bot.Say(ev.Channel, messageAddedStory, fields[2])
		}
		return nil
	case "story remove", "stories remove", "story rem", "stories rem":
		// TODO: Check for active session.
		if len(fields) < 3 {
			n.bot.Say(ev.Channel, messageUnknownStory)
		} else if err := author.RemoveStory(fields[2]); err != nil {
			n.bot.Say(ev.Channel, messageInvalidStory, err)
		} else if err = n.bot.Store.Set(authorKey, author); err != nil {
			n.bot.Say(ev.Channel, messageUnknownError)
			return err
		} else {
			n.bot.Say(ev.Channel, messageRemovedStory, fields[2])
		}
		return nil
	case "story start", "stories start":
		if len(fields) < 3 {
			n.bot.Say(ev.Channel, messageUnknownStory)
		} else if story, err := author.GetStory(fields[2]); err != nil {
			n.bot.Say(ev.Channel, messageInvalidStory, err)
		} else if _, ok := n.sessions[author.ID]; ok {
			n.bot.Say(ev.Channel, "TODO: Stop session before starting")
		} else if sess, err := NewSession(story); err != nil {
			n.bot.Say(ev.Channel, messageInvalidSession, err)
			return err
		} else if err = sess.Start(ctx, n.config); err != nil {
			n.bot.Say(ev.Channel, "TODO: Cannot start session: %s", err)
			return err
		} else {
			n.bot.Say(ev.Channel, messageStartedSession, fields[2], author.Options.Prefix)
			n.bot.Say(ev.Channel, sess.Output())
			n.sessions[author.ID] = sess
		}
		return nil
	case "story end", "stories end":
		if n.sessions[author.ID] == nil {
			n.bot.Say(ev.Channel, "TODO: No active session")
		} else {
			n.bot.Say(ev.Channel, "TODO: Stopped session")
			delete(n.sessions, author.ID)
		}
		return nil
	case "option", "options", "option list", "list options", "o":
		return n.SayTemplate(ev.Channel, templateOptionList, author)
	case "option set", "options set", "set option", "set options":
		if len(fields) < 4 {
			n.bot.Say(ev.Channel, messageUnknownOption)
		} else if err := author.SetOption(fields[2], fields[3]); err != nil {
			n.bot.Say(ev.Channel, messageInvalidOption, err)
		} else if err = n.bot.Store.Set(authorKey, author); err != nil {
			n.bot.Say(ev.Channel, messageUnknownError)
			return err
		} else {
			n.bot.Say(ev.Channel, messageSetOption, fields[2], fields[3])
		}
		return nil
	}

	return n.SayTemplate(ev.Channel, templateUnknownCommand, cmd)
}

// Default executable names for required runtime dependencies.
const (
	defaultInform7   = "/usr/libexec/ni"
	defaultInform6   = "/usr/libexec/inform6"
	defaultDumbFrotz = "/usr/bin/dfrotz"
)

type Config struct {
	// Required attributes.
	Bot *joe.Bot // The bot handler.

	// Optional attributes.
	Inform7   string // The path to the `ni` Inform 7 compiler.
	Inform6   string // The path to the `inform6` Inform 6 compiler.
	DumbFrotz string // The path to the `dumb-frotz` interpreter.
}

func New(conf Config) (*Inform, error) {
	if conf.Bot == nil {
		return nil, errors.New("bot given is nil")
	} else if conf.Bot.Store == nil {
		return nil, errors.New("storage module required")
	}

	// Set default paths for runtime dependencies, if needed.
	if conf.Inform7 == "" {
		conf.Inform7 = defaultInform7
	}
	if conf.Inform6 == "" {
		conf.Inform6 = defaultInform6
	}
	if conf.DumbFrotz == "" {
		conf.DumbFrotz = defaultDumbFrotz
	}

	// Verify and expand paths for runtime dependencies.
	if i7, err := exec.LookPath(conf.Inform7); err != nil {
		return nil, errors.Wrap(err, "Inform 7 compiler not found")
	} else if i6, err := exec.LookPath(conf.Inform6); err != nil {
		return nil, errors.Wrap(err, "Inform 6 compiler not found")
	} else if frotz, err := exec.LookPath(conf.DumbFrotz); err != nil {
		return nil, errors.Wrap(err, "Frotz interpreter not found")
	} else {
		conf.Inform7, conf.Inform6, conf.DumbFrotz = i7, i6, frotz
	}

	return &Inform{
		bot:      conf.Bot,
		config:   &conf,
		sessions: make(map[string]*Session),
	}, nil
}
