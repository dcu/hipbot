package handlers

import (
	"github.com/dcu/hipbot/xmpp"
	"time"
)

type TimeHandler struct {
}

func (timeHandler *TimeHandler) Matches(message *xmpp.Chat) bool {
	return message.Text == "time"
}

func (timeHandler *TimeHandler) Process(client *xmpp.Client, roomId string, message *xmpp.Chat) {
	now := time.Now().String()
	client.Send(xmpp.Chat{
		Remote: roomId,
		Type:   "groupchat",
		Text:   "The time is: " + now,
	})
}
