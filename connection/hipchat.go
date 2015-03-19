package connection

import (
	"flag"
	"fmt"
	"github.com/daneharrigan/hipchat"
	"github.com/dcu/hipbot/handlers"
	"github.com/dcu/hipbot/shared"
	"os"
	"time"
)

var (
	HandlersList = []handlers.Handler{
		&handlers.TimeHandler{},
		&handlers.LoggerHandler{},
		&handlers.RubyHandler{},
		&handlers.GiphyHandler{},
	}
)

func startListeningMessages(client *hipchat.Client) {
	roomId := *shared.Config.Room + "@conf.hipchat.com"
	client.Status("chat")
	client.Join(roomId, *shared.Config.FullName)

	go func() {
		for message := range client.Messages() {
			for _, handler := range HandlersList {
				if handler.Matches(message) {
					handler.Process(client, roomId, message)
				}
			}
		}
	}()
}

// from https://help.hipchat.com/knowledgebase/articles/64377-xmpp-jabber-support-details
// Connections are dropped after 150s of inactivity. We suggest sending a single space (" ") as keepalive data every 60 seconds
func startPinger(client *hipchat.Client) {
	ticker := time.NewTicker(60 * time.Second)

	go func() {
		for {
			<-ticker.C
			client.Say("", "Pinger", " ")
		}
	}()
}

func Start() {
	resource := "bot"
	client, err := hipchat.NewClient(*shared.Config.Username, *shared.Config.Password, resource)
	if err != nil {
		fmt.Printf("Client Error: %s\n", err)
		flag.Usage()
		os.Exit(1)
	}

	startListeningMessages(client)
	startPinger(client)
}
