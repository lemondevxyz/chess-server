package game

import "github.com/toms1441/chess/internal/board"

type (
	StructUpdateMessage struct {
		Message string `json:"message"`
	}

	StructUpdatePromotion struct {
		Player uint8       `json:"player"`
		Dst    board.Point `json:"dst"`
	}
)
