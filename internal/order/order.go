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
	// Possibility is received whenever a player wants to know their possible moves. [C]
	Possibility
	// Possible is sent whenever a player requests it via Possibility. [U]
	Possible
	// Turn is sent whenever a player moves a piece / special cases such as a promotion. [U]
	Turn
	// Promote is received from a player, to change a pawn that reached the end of the board to a dead piece. [O]
	// When sent to a player, it's an indication that he needs to promote his piece. And if the player sends it, we notify both players using Promotion.
	Promote
	// Promotion is sent whenever a player promotes it's pawn [U]
	Promotion
	// Pause is sent/received whenever a player wants to/pauses the game. [O]
	Pause
	// Message is sent/received whenver a player sends/receives a message. [O]
	Message
	// Done is sent whenever a game ends. [U]
	Done
)
