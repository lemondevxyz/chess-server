package game

import "encoding/json"

// Update is a communication structure from the server to the client, while Command is from the client to the server.
type Update struct {
	ID   uint8           `json:"id"`
	Data json.RawMessage `json:"data"`
	// used inside this package only
	parameter interface{}
}

// UpdateCallback sets the data for the update, since some updates are quite repetitive.
type UpdateCallback func(c *Client, u *Update) error

// Tests are bundled with command tests.
const (
	// UpdateBoard is an update for the board, this happens whenever a player moves a piece.
	UpdateBoard uint8 = iota + 1
	// UpdatePromotion happens whenever a pawn reaches the end of their board.
	UpdatePromotion
	// UpdatePause is sent whenever one of the players wants to pause the game for the other player to confirm, and sent another time to confirm game pause or opposite.
	UpdatePause
	// UpdateMessage whenever a player sends a message
	UpdateMessage
)

// redundant updates go here
// as well as verification
var ubs = map[uint8]UpdateCallback{
	UpdateBoard: func(c *Client, u *Update) error {
		body, err := json.Marshal(c.g.b)
		if err != nil {
			return err
		}

		u.Data = body

		return nil
	},

	UpdatePromotion: func(c *Client, u *Update) error {
		x, ok := u.parameter.(StructUpdatePromotion)
		if !ok {
			return ErrUpdateParameter
		}

		dst := x.Dst
		p := c.g.b.Get(dst)

		if p.Player != x.Player {
			return ErrIllegalPromotion
		}

		return nil
	},
}
