package main

import (
	"flag"
	"fmt"
	"github.com/daneharrigan/hipchat"
	"github.com/dcu/hipbot/handlers"
	"github.com/dcu/hipbot/shared"
	"github.com/dcu/hipbot/web"
	"os"
)

var (
	configFileFlag *string
)

var (
	HandlersList = []handlers.Handler{
		&handlers.TimeHandler{},
		&handlers.LoggerHandler{},
		&handlers.RubyHandler{},
	}
)

func init() {
	configFileFlag = flag.String("config", "config.yml", "Config file")
	shared.Config.Username = flag.String("username", os.Getenv("HIPBOT_USERNAME"), "Username")
	shared.Config.Password = flag.String("password", os.Getenv("HIPBOT_PASSWORD"), "Password")
	shared.Config.FullName = flag.String("full_name", os.Getenv("HIPBOT_FULL_NAME"), "Full Name")
	shared.Config.Room = flag.String("room", os.Getenv("HIPBOT_ROOM"), "Hipchat room")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Println("Get the config params from here: https://<team>.hipchat.com/account/xmpp")
	}
}

func loadConfig() {
	flag.Parse()

	if shared.Config.IsConfigured() {
		// Already configured.
		return
	}

	if _, err := os.Stat(*configFileFlag); os.IsNotExist(err) {
		shared.WriteSampleFile(*configFileFlag)
	}

	if !shared.Config.IsConfigured() {
		shared.ParseConfig(*configFileFlag, shared.Config)
	}

	if !shared.Config.IsConfigured() {
		flag.Usage()
	}
}

func startListeningMessages(client *hipchat.Client) {
	roomId := *shared.Config.Room + "@conf.hipchat.com"
	client.Status("chat")
	client.Join(roomId, *shared.Config.FullName)
	for message := range client.Messages() {
		for _, handler := range HandlersList {
			if handler.Matches(message) {
				handler.Process(client, roomId, message)
			}
		}
	}
}

func main() {
	loadConfig()

	resource := "bot"
	client, err := hipchat.NewClient(*shared.Config.Username, *shared.Config.Password, resource)
	if err != nil {
		fmt.Printf("Client Error: %s\n", err)
		flag.Usage()
		os.Exit(1)
	}

	go startListeningMessages(client)
	web.Start()
}
