package main

import (
	"fmt"
	"net/http"

	"github.com/fatih/color"
	"github.com/gorilla/mux"
	"github.com/toms1441/chess-server/internal/rest"
)

const port = ":8080"

func main() {

	rout := mux.NewRouter()

	rout.HandleFunc("/cmd", rest.CmdHandler).Methods("POST", "OPTIOS")
	rout.HandleFunc("/invite", rest.InviteHandler).Methods("POST", "OPTIONS")
	rout.HandleFunc("/accept", rest.AcceptInviteHandler).Methods("POST", "OPTIONS")
	rout.HandleFunc("/ws", rest.WebsocketHandler).Methods("GET", "OPTIONS")
	rout.HandleFunc("/avali", rest.GetAvaliableUsersHandler).Methods("GET", "OPTIONS")

	rout.HandleFunc("/protect", func(w http.ResponseWriter, r *http.Request) {
		_, err := rest.GetUser(r)
		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			w.Write(nil)
		}

		w.WriteHeader(http.StatusOK)
		w.Write(nil)
	}).Methods("GET")

	rout.PathPrefix("/pub").Handler(http.StripPrefix("/pub", http.FileServer(http.Dir("./static/"))))

	color.New(color.FgBlue).Println("Listening on port", port)

	http.ListenAndServe(port, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := &rest.Context{
			ResponseWriter: w,
		}

		method := color.New(color.BgMagenta, color.Bold).Sprint(" " + r.Method + " ")
		path := color.New(color.BgBlue).Sprint(" " + r.URL.Path + " ")

		ctx.Header().Add("Access-Control-Allow-Origin", "*")
		ctx.Header().Add("Access-Control-Allow-Headers", "Content-Type, Accept, Authorization")
		ctx.Header().Add("Access-Control-Allow-Methods", "GET, POST")
		if r.Method == "OPTIONS" {
			ctx.WriteHeader(http.StatusOK)
		} else {
			rout.ServeHTTP(ctx, r)
		}

		code := ""
		sta := ctx.GetStatus()
		if sta <= 299 && sta >= 200 {
			code = color.New(color.BgGreen, color.Bold).Sprintf(" %d ", sta)
		} else if sta >= 400 && sta <= 499 {
			code = color.New(color.BgYellow, color.Bold).Sprintf(" %d ", sta)
		} else if sta >= 500 && sta <= 511 {
			code = color.New(color.BgRed, color.Bold).Sprintf(" %d ", sta)
		} else {
			code = color.New(color.Reset).Sprintf(" %d ", sta)
		}

		fmt.Printf("%s%s%s\n", method, path, code)
	}))
}
