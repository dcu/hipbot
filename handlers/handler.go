package handlers

import (
	"github.com/dcu/hipbot/xmpp"
)

type Handler interface {
	Matches(message *xmpp.Chat) bool
	Process(client *xmpp.Client, roomJid string, message *xmpp.Chat)
}
