package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/toms1441/chess-server/internal/rest"
)

func main() {

	r := mux.NewRouter()
	r.HandleFunc("/cmd", rest.CmdHandler)
	r.HandleFunc("/ws", rest.WebsocketHandler)

	http.ListenAndServe(":8080", r)
}
