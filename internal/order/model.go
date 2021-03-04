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

// [U]
type TurnModel struct {
	Player uint8 `json:"player"`
}

// [U]
type PromotionModel struct {
	Player uint8       `json:"player"`
	Dst    board.Point `json:"dst"`
}

// [C]
type PromoteModel struct {
	Src  board.Point `json:"src"`
	Type uint8       `json:"type"`
}

// [C]
type PauseModel bool

// PausedModel equals to 0 when a single player wants to pause.
// Or equals to 1 when the other player declines. If the other player accepts it equals 2.
type PausedModel uint8

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
