package main

import (
	// Standard library
	"context"
	"os"
	"os/signal"

	// Internal packages
	"go.deuill.org/informbot/pkg/joe-inform-handler"
	"go.deuill.org/informbot/pkg/joe-xmpp-adapter"

	// Third-party packages
	"github.com/go-joe/file-memory"
	"github.com/go-joe/joe"
)

func main() {
	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	bot := joe.New(
		"inform",
		xmpp.Adapter(ctx, xmpp.Config{
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
