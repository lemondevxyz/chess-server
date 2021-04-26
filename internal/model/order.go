package model

import (
	"encoding/json"

	"github.com/toms1441/chess-server/internal/board"
)

// Order package provides models and ids for Commands and Updates.
// Each order is described as either an Update[U], a Command[C] or Both[O].

type Order struct {
	ID   uint8           `json:"id" validate:"required"`
	Data json.RawMessage `json:"data" validate:"required"`
	// Parameter primarily used in game.
	Parameter interface{} `json:"-"`
}

const (
	// Credentials is sent whenever a user connects to the websocket server. [U]
	OrCredentials uint8 = iota + 1
	// Invite is sent whenever a user receives an invite to a game. [U]
	OrInvite
	// Game is sent whenever a game starts. Sent by the invite handler. [U]
	OrGame
	// Move is sent/received whenever a player's piece moves. If the src is from the king/rook, if the dst is from the rook/king.
	OrMove
	// Turn is sent whenever a player moves a piece / special cases such as a promotion. [U]
	OrTurn
	// Promote is received from a player, to change a pawn that reached the end of the board to a dead piece. [O]
	// When sent to a player, it's an indication that he needs to promote his piece. And if the player sends it, we notify both players using Promotion.
	OrPromote
	// Promotion is sent whenever a player promotes it's pawn [U]
	OrPromotion
	// Castling is the act of switching the king and the rook's positions. This is only legal when the king and the rook haven't moved, and nothing is in between them. [O]
	OrCastling
	// Checkmate is sent whenever the king is in danger and needs to move.
	OrCheckmate
	// Done is sent whenever a game ends, or when the player wants to leave the game. [O]
	OrDone
)

// [U]
type CredentialsOrder struct {
	Profile Profile `json:"profile"`
	Token   string  `json:"token"`
}

// [O]
type InviteOrder struct {
	// Profile is only used for updates...
	Profile Profile `json:"profile" validate:"required"`
}

// [U]
type GameOrder struct {
	// which pieces are yours
	P1 bool `json:"p1,omitempty"`
	// Profile is the other player's profile
	Profile Profile      `json:"profile,omitempty"`
	Brd     *board.Board `json:"brd"`
}

// [O]
type MoveOrder struct {
	ID  int8        `json:"id"`
	Dst board.Point `json:"dst"`
}

// [U]
type TurnOrder struct {
	P1 bool `json:"p1"`
}

// [O]
type PromoteOrder struct {
	ID   int8  `json:"id"`
	Kind uint8 `json:"kind"`
}

// [U]
type PromotionOrder PromoteOrder

// [O]
type CastlingOrder struct {
	Src int8 `json:"src"`
	Dst int8 `json:"dst"`
}

// [U]
type CheckmateOrder TurnOrder

// [O]
type DoneOrder struct {
	P1 bool `json:"p1"`
}
