package handlers

import (
	"github.com/daneharrigan/hipchat"
)

type Handler interface {
	Matches(message *hipchat.Message) bool
	Process(client *hipchat.Client, roomJid string, message *hipchat.Message)
}
