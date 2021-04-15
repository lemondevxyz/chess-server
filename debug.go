package main

import (
	"encoding/json"
	"fmt"

	"github.com/toms1441/chess-server/internal/board"
	"github.com/toms1441/chess-server/internal/game"
	"github.com/toms1441/chess-server/internal/order"
)

// where c1 is p1
func debugCastling(cl1, cl2 *game.Client) (err error) {
	list := []order.MoveModel{
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

	// game := cl1.Game()
	p1 := true

	for k, v := range list {
		body, err := json.Marshal(v)
		if err != nil {
			return fmt.Errorf("body: %s\nerror: %s\nindex: %d", string(body), err.Error(), k)
		}

		if p1 {
			err = cl1.Do(order.Order{
				ID:   order.Move,
				Data: body,
			})
		} else {
			err = cl2.Do(order.Order{
				ID:   order.Move,
				Data: body,
			})
		}
		if err != nil {
			return fmt.Errorf("body: %s\nerror: %s\nindex: %d", string(body), err, k)
		}

		p1 = !p1
	}

	return nil
}
