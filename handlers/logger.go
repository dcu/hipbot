package handlers

import (
	"fmt"
	"github.com/dcu/hipbot/xmpp"
	"time"
)

type LoggerHandler struct {
}

func (logger *LoggerHandler) Matches(message *xmpp.Chat) bool {
	return true
}

func (logger *LoggerHandler) Process(client *xmpp.Client, roomId string, message *xmpp.Chat) {
	time := time.Now().Format("2006-01-02T15:04:05.999999-07:00")
	fmt.Printf("[%s] %s: %s \n", time, message.Remote, message.Text)
}
