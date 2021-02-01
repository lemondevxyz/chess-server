package game

import "github.com/toms1441/chess/internal/board"

type (
	ModelCmdPiece struct {
		Src board.Point `json:"src"`
		Dst board.Point `json:"dst"`
	}

	ModelCmdPromotion struct {
		Src board.Point `json:"src"`
		ID  uint8       `json:"id"`
	}

	ModelCmdMessage struct {
		Message string `json:"message"`
	}
)
