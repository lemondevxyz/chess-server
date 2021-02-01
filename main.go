package main

import (
	"fmt"

	"github.com/toms1441/chess/internal/board"
)

func main() {
	b := board.NewBoard()

	body, err := b.MarshalJSON()

	b = &board.Board{}
	err = b.UnmarshalJSON(body)
	fmt.Println(b, err)
}
