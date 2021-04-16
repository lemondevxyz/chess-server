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
	ID string `json:"id" validate:"required"`
}

// [U]
type GameModel struct {
	// which pieces are yours
	P1    bool         `json:"p1"`
	Board *board.Board `json:"board"`
}

// [O]
type MoveModel struct {
	ID  *int8        `json:"id" validate:"required"`
	Dst *board.Point `json:"dst" validate:"required"`
}

// [O] sent as response to http
type PossibleModel struct {
	ID     *int8         `json:"id,omitempty" validate:"required"` // [C]
	Points *board.Points `json:"points,omitempty"`                 // [U]
}

// [U]
type TurnModel struct {
	P1 bool `json:"p1"`
}

// [O]
type PromoteModel struct {
	ID   int   `json:"id" validate:"required"`
	Kind uint8 `json:"kind" validiate:"required"`
}

// [U]
type PromotionModel PromoteModel /*struct {
	Type uint8       `json:"type"`
	Dst  board.Point `json:"dst"`
}*/

// [O]
type CastlingModel struct {
	Src *int `json:"src" validate:"required"`
	Dst *int `json:"dst" validate:"required"`
}

// [U]
type CheckmateModel TurnModel

// [O]
type DoneModel struct {
	P1 bool `json:"p1" validate:"required"`
}
