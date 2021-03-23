package main

import (
	"encoding/json"

	"github.com/toms1441/chess-server/internal/board"
	"github.com/toms1441/chess-server/internal/order"
)

type Debug uint8

const (
	// automatically create game whenever a user connects
	debugInvite Debug = iota + 1
	debugCastling
	debugPromote
	debugCheckmate
)

func getMoveData(src board.Point, dst board.Point) []byte {
	body, _ := json.Marshal(order.MoveModel{
		Src: src,
		Dst: dst,
	})

	return body
}

func getOrder(data []byte) order.Order {
	return order.Order{
		ID:   order.Move,
		Data: data,
	}
}
