package game

// Update is a communication structure from the server to the client, while Command is from the client to the server.
type Update struct {
	ID   uint8
	Data []byte
}

const (
	// UpdateBoard is an update for the board, this happens whenever a player moves a piece.
	// Data parameters are `[board_array]`
	UpdateBoard uint8 = iota + 1
	// UpdatePick happens whenever a pawn reaches the end of their board.
	// Data parameters are `` - empty.
	UpdatePromtion
	// UpdatePause is sent whenever one of the players wants to pause the game for the other player to confirm, and sent another time to confirm game pause or opposite.
	// Data parameters are `{ type: 0 }` - for the player to confirm it
	// and `{ type: 1 }` - game is paused
	// and lastly `{ type: 2 }` - other player declined
	UpdatePause
	// UpdateMessage whenever a player sends a message
	// Data parameters are `{player: 1, message: "hello world"}`
	// Player 0 would mean a message from game
	UpdateMessage
)
