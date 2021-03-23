package main

import (
	"fmt"

	"github.com/toms1441/chess-server/internal/board"
	"github.com/toms1441/chess-server/internal/order"
	"github.com/toms1441/chess-server/internal/rest"
)

func checkmateDebug(us1, us2 *rest.User) {
	getOrder := func(src, dst board.Point) order.Order {
		return getOrder(getMoveData(src, dst))
	}

	fmt.Println(us2.Client().Do(getOrder(board.Point{6, 5}, board.Point{5, 5})))
	fmt.Println(us1.Client().Do(getOrder(board.Point{1, 4}, board.Point{3, 4})))
	fmt.Println(us2.Client().Do(getOrder(board.Point{6, 6}, board.Point{4, 6})))
	fmt.Println(us1.Client().Do(getOrder(board.Point{0, 3}, board.Point{4, 7})))
}
