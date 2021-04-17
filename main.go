package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/gobwas/ws"
	"github.com/gorilla/mux"
	"github.com/toms1441/chess-server/internal/rest"
)

const apiver = "v1"

func debug_game() {
	fmt.Println("endless loop mode")
	if debug != "yes" {
		fmt.Printf("debug state: %s\n", debug)
	}
	for {
		x := rest.ClientChannel()
		cl1 := <-x
		time.Sleep(time.Second)
		go func() {
			cn, _, _, err := ws.Dial(context.Background(), "ws://localhost:8080/api/v1/ws")
			if err != nil {
				fmt.Printf("ws.Dial: %s\n", err)
			}

			for {
				b := make([]byte, 2048)
				_, err := cn.Read(b)
				if err != nil {
					panic(err)
				}
			}
		}()
		cl2 := <-x
		if !p1 {
			cl1, cl2 = cl2, cl1
		}

		id, _ := cl2.Invite(cl1.PublicID, rest.InviteLifespan)
		time.Sleep(time.Millisecond * 10)
		cl1.AcceptInvite(id)

		var err error
		switch debug {
		case "castling":
			err = debugCastling(cl1.Client(), cl2.Client())
		case "checkmate":
			err = debugCheckmate(cl1.Client(), cl2.Client())
		case "promotion":
			err = debugPromotion(cl1.Client(), cl2.Client())
		}

		if err != nil {
			panic(err)
		}
	}
}

func main() {
	if debug != "no" {
		// go debug_game()
	}

	rout := mux.NewRouter()
	{ // api routes
		api := rout.PathPrefix("/api/" + apiver).Subrouter()

		api.HandleFunc("/cmd", rest.CmdHandler).Methods("POST", "OPTIONS")
		api.HandleFunc("/invite", rest.InviteHandler).Methods("POST", "OPTIONS")
		api.HandleFunc("/accept", rest.AcceptInviteHandler).Methods("POST", "OPTIONS")
		api.HandleFunc("/ws", rest.WebsocketHandler).Methods("GET", "OPTIONS")
		api.HandleFunc("/avali", rest.GetAvaliableUsersHandler).Methods("GET", "OPTIONS")
		api.HandleFunc("/possib", rest.PossibHandler).Methods("POST", "OPTIONS")

		api.HandleFunc("/protect", func(w http.ResponseWriter, r *http.Request) {
			_, err := rest.GetUser(r)
			if err != nil {
				w.WriteHeader(http.StatusForbidden)
				w.Write(nil)
			}

			w.WriteHeader(http.StatusOK)
			w.Write(nil)
		}).Methods("GET")
	}

	var proto string
	var port string
	if debug != "no" {
		proto = "tcp"
		port = ":8080"
	} else {
		proto = "unix"
		port = "http.sock"

		os.Remove(port)
		os.Remove("ws.sock")
	}

	color.New(color.FgBlue).Println("Listening on", port)

	listen, err := net.Listen(proto, port)
	if err != nil {
		panic(err)
	}

	if proto == "unix" {
		ws, err := net.Listen("unix", "ws.sock")
		if err != nil {
			panic(err)
		}

		go rest.WebsocketServe(ws)

		os.Chmod("ws.sock", 0777)
		os.Chmod("http.sock", 0777)
	}

	defer listen.Close()

	http.Serve(listen, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := &rest.Context{
			ResponseWriter: w,
		}

		method := color.New(color.BgMagenta, color.Bold).Sprint(" " + r.Method + " ")
		path := color.New(color.BgBlue).Sprint(" " + r.URL.Path + " ")

		if debug != "no" {
			ctx.Header().Add("Access-Control-Allow-Origin", "*")
			ctx.Header().Add("Access-Control-Allow-Headers", "Content-Type, Accept, Authorization")
			ctx.Header().Add("Access-Control-Allow-Methods", "GET, POST")
		}
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
