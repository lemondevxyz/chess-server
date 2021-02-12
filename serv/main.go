package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/toms1441/chess/serv/internal/rest"
)

const dir = "./static/"

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", rest.WebsocketHandler)

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(dir))))

	http.ListenAndServe(":6969", r)
}
