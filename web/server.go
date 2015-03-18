package web

import (
	"fmt"
	"net/http"
	"os"
)

func GetIndex(response http.ResponseWriter, request *http.Request) {
	fmt.Fprintln(response, "Nothing to see here.")
}

func Start() {
	http.HandleFunc("/", GetIndex)

	port := os.Getenv("PORT")
	if port == "" {
		port = "4000"
	}

	http.ListenAndServe(":"+port, nil)
}
