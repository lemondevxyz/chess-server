package main

import (
	"github.com/gin-gonic/gin"
	"github.com/toms1441/chess/internal/board"
	"github.com/toms1441/chess/internal/game"
	"github.com/toms1441/chess/internal/rest"
)

func main() {
	g := game.Game{}
	b := board.NewBoard()

	r := gin.Default()

	rest.Init(r.Group("/api/v0/"))

	r.Run(":8080")
}