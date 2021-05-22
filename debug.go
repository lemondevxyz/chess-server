package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gobwas/ws"
	"github.com/toms1441/chess-server/internal/board"
	"github.com/toms1441/chess-server/internal/game"
	"github.com/toms1441/chess-server/internal/model"
	"github.com/toms1441/chess-server/internal/rest"
)

var debug = "normal"

const p1 = true

var solo = true

func debug_game() {
	fmt.Println("endless loop mode")
	if debug != "yes" {
		fmt.Printf("debug mode: %s\n", debug)
	}

	connect := func() {
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
	}

	if debug == "watchable" || debug == "invite" {
		if debug == "invite" {
			solo = false
		}

		go connect()
	}

	for {
		x := rest.ClientChannel()
		cl1 := <-x
		if solo {
			if debug == "watchable" {
				solo = false
			}

			go connect()

			time.Sleep(time.Second)
		}
		cl2 := <-x
		if !p1 {
			cl1, cl2 = cl2, cl1
		}

		if debug != "invite" {
			err := cl2.Invite(model.InviteOrder{
				Profile: cl1.Profile,
			}, rest.InviteLifespan)
			if err != nil {
				panic(err)
			}
			time.Sleep(time.Millisecond * 10)
			cl1.AcceptInvite(cl2.Profile.ID + "_" + cl2.Profile.Platform)
		}

		var err error
		switch debug {
		case "castling":
			err = debugCastling(cl1.Client(), cl2.Client())
		case "checkmate":
			err = debugCheckmate(cl1.Client(), cl2.Client())
		case "promotion":
			err = debugPromotion(cl1.Client(), cl2.Client())
		case "watchable":
			err = debugWatchable(cl1.Client(), cl2.Client())
		case "invite":
			err = debugInvite(cl2, cl1)
		}

		if err != nil {
			panic(err)
		}
	}
}

func doMove(cl1, cl2 *game.Client, list []model.MoveOrder) error {
	p1 := true

	for k, v := range list {
		body, err := json.Marshal(v)
		if err != nil {
			return fmt.Errorf("body: %s\nerror: %s\nindex: %d", string(body), err.Error(), k)
		}

		o := model.Order{
			ID:   model.OrMove,
			Data: body,
		}
		if p1 {
			err = cl1.Do(o)
		} else {
			err = cl2.Do(o)
		}
		if err != nil {
			return fmt.Errorf("body: %s\nerror: %s\nindex: %d", string(body), err, k)
		}

		p1 = !p1
	}

	return nil
}

// where c1 is p1
func debugCastling(cl1, cl2 *game.Client) (err error) {
	list := []model.MoveOrder{
		// pawns
		{16, board.Point{0, 4}},
		{8, board.Point{0, 3}},
		{17, board.Point{1, 4}},
		{9, board.Point{1, 3}},
		{18, board.Point{2, 4}},
		{10, board.Point{2, 3}},
		{19, board.Point{3, 4}},
		{11, board.Point{3, 3}},
		{20, board.Point{4, 4}},
		{12, board.Point{4, 3}},
		{21, board.Point{5, 4}},
		{13, board.Point{5, 3}},
		{22, board.Point{6, 4}},
		{14, board.Point{6, 3}},
		{23, board.Point{7, 4}},
		{15, board.Point{7, 3}},
		// knight
		{25, board.Point{2, 5}},
		{1, board.Point{2, 2}},
		{30, board.Point{7, 5}},
		{6, board.Point{7, 2}},
		// bishop
		{26, board.Point{3, 6}},
		{2, board.Point{3, 1}},
		{29, board.Point{6, 6}},
		{5, board.Point{6, 1}},
		// queen
		{27, board.Point{4, 6}},
		{3, board.Point{2, 1}},
	}
	// const
	if !p1 {
		list = append(list, model.MoveOrder{27, board.Point{5, 6}})
	}

	return doMove(cl1, cl2, list)
}

func debugCheckmate(cl1, cl2 *game.Client) error {
	var list []model.MoveOrder
	if !p1 {
		list = []model.MoveOrder{
			{21, board.Point{5, 5}},
			{12, board.Point{4, 3}},
			{22, board.Point{6, 4}},
		}
	} else {
		list = []model.MoveOrder{
			{20, board.Point{4, 4}},
			{13, board.Point{5, 3}},
			{30, board.Point{7, 5}},
			{14, board.Point{6, 3}},
		}
	}

	return doMove(cl1, cl2, list)
}

func debugPromotion(cl1, cl2 *game.Client) error {
	list := []model.MoveOrder{
		{17, board.Point{1, 4}},
		{15, board.Point{7, 3}},
		{17, board.Point{1, 3}},
		{15, board.Point{7, 4}},
		{17, board.Point{1, 2}},
		{15, board.Point{7, 5}},
		{17, board.Point{0, 1}},
		{15, board.Point{6, 6}},
	}
	if !p1 {
		list = append(list, model.MoveOrder{25, board.Point{0, 5}})
	}

	return doMove(cl1, cl2, list)
}

func debugWatchable(cl1, cl2 *game.Client) error {
	rev := false
	for {
		var list []model.MoveOrder
		if !rev {
			list = []model.MoveOrder{
				{25, board.Point{0, 5}},
				{1, board.Point{0, 2}},
			}
		} else {
			list = []model.MoveOrder{
				{25, board.Point{1, 7}},
				{1, board.Point{1, 0}},
			}
		}

		doMove(cl1, cl2, list)
		time.Sleep(time.Second * 1)
		rev = !rev
	}
}

func debugInvite(cl1 *rest.User, cl2 *rest.User) error {
	time.Sleep(time.Millisecond * 100)
	return cl2.Invite(model.InviteOrder{
		Profile: cl1.Profile,
	}, rest.InviteLifespan)
}
