package order

import "encoding/json"

// Order package provides models and ids for Commands and Updates.
// Each order is described as either an Update[U], a Command[C] or Both[O].

type Order struct {
	ID   uint8           `json:"id"`
	Data json.RawMessage `json:"data"`
	// Parameter primarily used in game.
	Parameter interface{} `json:"-"`
}

const (
	// Credentials is sent whenever a user connects to the websocket server. [U]
	Credentials uint8 = iota + 1
	// Invite is sent whenever a user receives an invite to a game. [U]
	Invite
	// Game is sent whenever a game starts. Sent by the invite handler. [U]
	Game
	// Move is sent/received whenever a player's piece moves. [O]
	Move
	// Turn is sent whenever a player moves a piece. [U]
	Turn
	// Promotion is sent whenever a player's pawn reaches the end of the board. Where x equals 7 or equals 0. [U]
	Promotion
	// Promote is received from a player, to change a pawn that reached the end of the board to a dead piece. [C]
	Promote
	// Pause is sent/received whenever a player wants to/pauses the game. [O]
	Pause
	// Message is sent/received whenver a player sends/receives a message. [O]
	Message
	// Done is sent whenever a game ends. [U]
	Done
)
