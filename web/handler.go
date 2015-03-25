package web

import (
	"github.com/dcu/hipbot/shared"
	"net/http"
)

type Handler struct {
	HandlerFunc http.HandlerFunc
}

func NewHandler(handlerFunc http.HandlerFunc) *Handler {
	return &Handler{HandlerFunc: handlerFunc}
}

func (handler *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	if r.FormValue("api_key") != *shared.Config.ApiKey {
		w.WriteHeader(400)
		w.Write([]byte("Bad Request"))
		return
	}

	handler.HandlerFunc(w, r)
}
