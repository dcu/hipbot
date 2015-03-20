package connection

import (
	"crypto/tls"
	"fmt"
	"github.com/dcu/hipbot/handlers"
	"github.com/dcu/hipbot/shared"
	"github.com/dcu/hipbot/xmpp"
	"time"
)

var (
	Client *xmpp.Client
)

var (
	HandlersList = []handlers.Handler{
		&handlers.TimeHandler{},
		&handlers.LoggerHandler{},
		&handlers.RubyHandler{},
		&handlers.GiphyHandler{},
	}
)

func processMessage(roomId string, message *xmpp.Chat) {
	for _, handler := range HandlersList {
		if handler.Matches(message) {
			handler.Process(Client, roomId, message)
		}
	}
}

func startListeningMessages() {
	roomId := *shared.Config.Room + "@conf.hipchat.com"

	go func() {
		for {
			chat, err := Client.Recv()
			if err != nil {
				println("Error:", err.Error())
				connect()
				continue
			}

			switch v := chat.(type) {
			case xmpp.Chat:
				processMessage(roomId, &v)
			case xmpp.Presence:
				fmt.Println("PRESENCE:", v.From, v.Show)
			}
		}
	}()

}

// from https://help.hipchat.com/knowledgebase/articles/64377-xmpp-jabber-support-details
// Connections are dropped after 150s of inactivity. We suggest sending a single space (" ") as keepalive data every 60 seconds
func startPinger() {
	ticker := time.NewTicker(60 * time.Second)

	go func() {
		for {
			<-ticker.C
			Client.PingR()
		}
	}()
}

//func Start() {
//resource := "bot"
//client, err := hipchat.NewClient(*shared.Config.Username, *shared.Config.Password, resource)
//if err != nil {
//fmt.Printf("Client Error: %s\n", err)
//flag.Usage()
//os.Exit(1)
//}
//Client = client

//startListeningMessages()
//startPinger()
//}

func connect() {
	server := "chat.hipchat.com"
	xmpp.DefaultConfig = tls.Config{
		ServerName:         server,
		InsecureSkipVerify: false,
	}

	var talk *xmpp.Client
	var err error
	options := xmpp.Options{
		Host:          server + ":5222",
		User:          *shared.Config.Username + "@chat.hipchat.com",
		Password:      *shared.Config.Password,
		NoTLS:         true,
		Debug:         false,
		Session:       false,
		Resource:      "bot",
		Status:        "chat",
		StatusMessage: "Available",
	}

	println("Connecting to HipChat...")
	talk, err = options.NewClient()
	if err != nil {
		println("Error:", err.Error())
		connect()
		return
	}

	roomId := *shared.Config.Room + "@conf.hipchat.com"
	talk.JoinMUC(roomId, *shared.Config.FullName)

	Client = talk
}

func Start() {
	connect()
	startListeningMessages()
	startPinger()
}
