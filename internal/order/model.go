package order

import (
	"github.com/toms1441/chess-server/internal/board"
)

// [U]
type CredentialsModel struct {
	Token    string `json:"token"`
	PublicID string `json:"public_id"`
}

// [O]
type InviteModel struct {
	ID string `json:"id"`
}

// [U]
type GameModel struct {
	// which pieces are yours
	Player uint8        `json:"player"`
	Board  *board.Board `json:"board"`
}

// [O]
type MoveModel struct {
	Src board.Point `json:"src"`
	Dst board.Point `json:"dst"`
}

// [O] sent as response to http
type PossibleModel struct {
	Src    *board.Point  `json:"src,omitempty"`    // [C]
	Points *board.Points `json:"points,omitempty"` // [U]
}

// [U]
type TurnModel struct {
	Player uint8 `json:"player"`
}

// [O]
type PromoteModel struct {
	Src  board.Point `json:"src"`
	Type uint8       `json:"type"`
}

// [U]
type PromotionModel struct {
	Type uint8       `json:"type"`
	Dst  board.Point `json:"dst"`
}

type CastlingModel struct {
	Src board.Point `json:"src"`
	Dst board.Point `json:"dst"`
}

// [O]
type MessageModel struct {
	Message string `json:"message"`
}

// [U]
type DoneModel struct {
	Result int8 `json:"result"`
	// -1 == you lost
	// 0 == draw/stalemate
	// 1 == you won
}
