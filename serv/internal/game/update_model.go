package game

import "github.com/toms1441/chess/internal/board"

type (
	ModelUpdateMessage struct {
		Message string `json:"message"`
	}

	ModelUpdatePromotion struct {
		Player uint8       `json:"player"`
		Dst    board.Point `json:"dst"`
	}

	ModelUpdateTurn struct {
		Player uint8 `json:"player"`
	}
)
