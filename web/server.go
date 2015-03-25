package web

import (
	"fmt"
	"github.com/dcu/hipbot/connection"
	"github.com/dcu/hipbot/shared"
	"github.com/dcu/hipbot/xmpp"
	"net/http"
	"os"
)

func GetIndex(response http.ResponseWriter, request *http.Request) {
	fmt.Fprintln(response, "Nothing to see here.")
}

func PostNotification(response http.ResponseWriter, request *http.Request) {
	request.ParseForm()

	roomId := *shared.Config.Room + "@conf.hipchat.com"
	connection.Client.Send(xmpp.Chat{
		Remote: roomId,
		Type:   "groupchat",
		Text:   request.Form["body"][0],
	})
}

func Start() {
	http.Handle("/", NewHandler(GetIndex))
	http.Handle("/notifications", NewHandler(PostNotification))

	port := os.Getenv("PORT")
	if port == "" {
		port = "4000"
	}

	http.ListenAndServe(":"+port, nil)
}
