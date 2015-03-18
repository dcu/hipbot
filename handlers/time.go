package handlers

import (
	"github.com/daneharrigan/hipchat"
	"time"
)

type TimeHandler struct {
}

func (timeHandler *TimeHandler) Matches(message *hipchat.Message) bool {
	return message.Body == "@hipbot time"
}

func (timeHandler *TimeHandler) Process(client *hipchat.Client, roomId string, message *hipchat.Message) {
	now := time.Now().String()
	client.Say(roomId, "Bot", "The time is: "+now)
}
