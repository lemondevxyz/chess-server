package main

import (
	"github.com/gin-gonic/gin"
	"github.com/toms1441/chess/serv/internal/board"
	"github.com/toms1441/chess/serv/internal/game"
)

func main() {
	g := game.Game{}
	b := board.NewBoard()

	r := gin.Default()

	r.Run(":8080")
}
