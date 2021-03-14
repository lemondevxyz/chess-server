package main

import (
	"github.com/toms1441/chess-server/internal/board"
	"github.com/toms1441/chess-server/internal/order"
	"github.com/toms1441/chess-server/internal/rest"
)

func castlingDebug(us1, us2 *rest.User) {
	getOrder := func(src, dst board.Point) order.Order {
		return getOrder(getMoveData(src, dst))
	}

	// move all pawns forward
	for i := 0; i <= 7; i++ {
		us2.Client().Do(getOrder(board.Point{6, i}, board.Point{4, i}))
		us1.Client().Do(getOrder(board.Point{1, i}, board.Point{3, i}))
	}

	// move knight
	us2.Client().Do(getOrder(board.Point{7, 6}, board.Point{5, 7}))
	us1.Client().Do(getOrder(board.Point{0, 6}, board.Point{2, 7}))
	us2.Client().Do(getOrder(board.Point{7, 1}, board.Point{5, 0}))
	us1.Client().Do(getOrder(board.Point{0, 1}, board.Point{2, 0}))

	// move bishops
	us2.Client().Do(getOrder(board.Point{7, 5}, board.Point{6, 4}))
	us1.Client().Do(getOrder(board.Point{0, 5}, board.Point{1, 4}))
	us2.Client().Do(getOrder(board.Point{7, 2}, board.Point{6, 3}))
	us1.Client().Do(getOrder(board.Point{0, 2}, board.Point{1, 3}))

	// move queen
	us2.Client().Do(getOrder(board.Point{7, 3}, board.Point{6, 2}))
	us1.Client().Do(getOrder(board.Point{0, 3}, board.Point{1, 2}))

}
