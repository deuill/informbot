package main

import (
	// Standard library
	"os"

	// Internal packages
	"github.com/deuill/informbot/joe-inform-handler"
	"github.com/deuill/informbot/joe-xmpp-adapter"

	// Third-party packages
	"github.com/go-joe/file-memory"
	"github.com/go-joe/joe"
)

func main() {
	bot := joe.New(
		"inform",
		xmpp.Adapter(xmpp.Config{
			JID:         os.Getenv("INFORMBOT_JID"),
			Password:    os.Getenv("INFORMBOT_PASSWORD"),
			NoTLS:       os.Getenv("INFORMBOT_NO_TLS") == "true",
			UseStartTLS: os.Getenv("INFORMBOT_USE_STARTTLS") == "true",
		}),
		file.Memory("store.json"),
	)

	in, err := inform.New(inform.Config{Bot: bot})
	if err != nil {
		bot.Logger.Fatal(err.Error())
	}

	bot.Brain.RegisterHandler(in.Handle)
	if err := bot.Run(); err != nil {
		bot.Logger.Fatal(err.Error())
	}
}
