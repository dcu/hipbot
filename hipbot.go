package main

import (
	"flag"
	"fmt"
	"github.com/daneharrigan/hipchat"
	"github.com/dcu/hipbot/handlers"
	"github.com/dcu/hipbot/shared"
	"os"
)

var (
	configFileFlag *string
)

var (
	HandlersList = []handlers.Handler{
		&handlers.TimeHandler{},
		&handlers.LoggerHandler{},
	}
)

func init() {
	configFileFlag = flag.String("config", "config.yml", "Config file")
	shared.Config.Username = flag.String("username", "", "Username")
	shared.Config.Password = flag.String("password", "", "Password")
	shared.Config.FullName = flag.String("full_name", "", "Full Name")
	shared.Config.Room = flag.String("room", "", "Hipchat room")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Println("Get the config params from here: https://<team>.hipchat.com/account/xmpp")
	}
}

func loadConfig() {
	flag.Parse()

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

func main() {
	loadConfig()

	roomId := *shared.Config.Room + "@conf.hipchat.com"
	resource := "bot"
	client, err := hipchat.NewClient(*shared.Config.Username, *shared.Config.Password, resource)
	if err != nil {
		fmt.Printf("Client Error: %s\n", err)
		flag.Usage()
		os.Exit(1)
	}

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
