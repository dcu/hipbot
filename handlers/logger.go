package handlers

import (
	"fmt"
	"github.com/daneharrigan/hipchat"
	"time"
)

type LoggerHandler struct {
}

func (logger *LoggerHandler) Matches(message *hipchat.Message) bool {
	return true
}

func (logger *LoggerHandler) Process(client *hipchat.Client, roomId string, message *hipchat.Message) {
	time := time.Now().Format("2006-01-02T15:04:05.999999-07:00")
	fmt.Printf("[%s] %s: %s \n", time, message.From, message.Body)
}
