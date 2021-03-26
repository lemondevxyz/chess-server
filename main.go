package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/fatih/color"
	"github.com/gorilla/mux"
	"github.com/toms1441/chess-server/internal/rest"
	"github.com/toms1441/chess-server/internal/rest/headless"
)

const port = ":8080"

// just to build without debug
// via build.sh
var debug = "yes"

func debug_game(debugValue Debug, solo bool) {
	ch := rest.ClientChannel()
	go func() {
		var us2 *rest.User

		us1 := <-ch
		if solo {
			go headless.NewClient("ws://localhost" + port + "/ws")
		}

		us2 = <-ch

		id, err := us2.Invite(us1.PublicID, rest.InviteLifespan)
		if err != nil {
			fmt.Printf("error: %s\n", err.Error())
		} else {
			us1.AcceptInvite(id)
		}

		// a little delay isn't bad
		time.Sleep(time.Second * 1)

		switch debugValue {
		case debugCastling:
			castlingDebug(us2, us1)
		case debugPromote:
			promotionDebug(us2, us1)
		case debugCheckmate:
			checkmateDebug(us2, us1)
		}
	}()
}

func main() {
	rout := mux.NewRouter()

	if debug == "yes" {
		debug_game(debugPromote, true)
	}

	rout.HandleFunc("/cmd", rest.CmdHandler).Methods("POST", "OPTIONS")
	rout.HandleFunc("/invite", rest.InviteHandler).Methods("POST", "OPTIONS")
	rout.HandleFunc("/accept", rest.AcceptInviteHandler).Methods("POST", "OPTIONS")
	rout.HandleFunc("/ws", rest.WebsocketHandler).Methods("GET", "OPTIONS")
	rout.HandleFunc("/avali", rest.GetAvaliableUsersHandler).Methods("GET", "OPTIONS")
	rout.HandleFunc("/possib", rest.PossibHandler).Methods("POST", "OPTIONS")

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
