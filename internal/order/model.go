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
	ID  int8        `json:"id"`
	Dst board.Point `json:"dst"`
}

// [O] sent as response to http
type PossibleModel struct {
	ID     *int8         `json:"src,omitempty"`    // [C]
	Points *board.Points `json:"points,omitempty"` // [U]
}

// [U]
type TurnModel struct {
	Player uint8 `json:"player"`
}

// [O]
type PromoteModel struct {
	ID   int   `json:"id"`
	Type uint8 `json:"type"`
}

// [U]
type PromotionModel PromoteModel /*struct {
	Type uint8       `json:"type"`
	Dst  board.Point `json:"dst"`
}*/

// [O]
type CastlingModel struct {
	Src int `json:"src"`
	Dst int `json:"dst"`
}

// [U]
type CheckmateModel TurnModel /* struct {
	Player uint8 `json:"player"`
}*/

// [O]
type DoneModel struct {
	// Result represents the result of the match
	// it equal to the winning player's number
	Result uint8 `json:"result"`
}
