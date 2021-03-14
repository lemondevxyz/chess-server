package main

import (
	"github.com/toms1441/chess-server/internal/board"
	"github.com/toms1441/chess-server/internal/order"
	"github.com/toms1441/chess-server/internal/rest"
)

func promotionDebug(us1, us2 *rest.User) {
	getOrder := func(src, dst board.Point) order.Order {
		return getOrder(getMoveData(src, dst))
	}

	for i := 1; i < 6; i++ {
		oldx := 7 - i
		if oldx < 0 {
			oldx = oldx * -1
		}

		newx := oldx - 1
		newy := 0
		if i == 5 {
			newy = 1
		}
		us2.Client().Do(getOrder(board.Point{oldx, 0}, board.Point{newx, newy}))

		oldx = i
		newx = oldx + 1
		newy = 7
		if i == 5 {
			newy = 6
		}
		us1.Client().Do(getOrder(board.Point{oldx, 7}, board.Point{newx, newy}))
	}
}
