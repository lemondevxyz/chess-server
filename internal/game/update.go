package game

import (
	"encoding/json"

	"github.com/toms1441/chess-server/internal/order"
)

// Update is a communication structure from the server to the client, while Command is from the client to the server.
/*
type Update struct {
	ID   uint8           `json:"id"`
	Data json.RawMessage `json:"data"`
	// used inside this package only
	parameter interface{}
}
*/

// UpdateCallback sets the data for the update, since some updates are quite repetitive.
type UpdateCallback func(c *Client, u *order.Order) error

// Tests are bundled with command tests.
/*
const (
	// UpdateBoard is an update for the board, this happens whenever a player moves a piece.
	UpdateBoard uint8 = iota + 1
	// UpdatePromotion happens whenever a pawn reaches the end of their board.
	UpdatePromotion
	// UpdatePause is sent whenever one of the players wants to pause the game for the other player to confirm, and sent another time to confirm game pause or opposite.
	UpdatePause
	// UpdateMessage whenever a player sends a message
	UpdateMessage
	// UpdateTurn it's your turn pal
	UpdateTurn
	// UpdateInvite sent whenever a player recieved an invite.
	UpdateInvite
	// UpdateCredentials sent whenever a player connects to websocket
	UpdateCredentials
)
*/

// redundant updates go here
// as well as verification
var ubs = map[uint8]UpdateCallback{
	order.Move: func(c *Client, u *order.Order) error {
		if c.g == nil {
			return ErrGameNil
		}

		body, err := json.Marshal(c.g.b)
		if err != nil {
			return err
		}

		u.Data = body

		return nil
	},
	// verification
	order.Promotion: func(c *Client, u *order.Order) error {
		x, ok := u.Parameter.(order.PromotionModel)
		if !ok {
			return ErrUpdateParameter
		}

		dst := x.Dst
		p := c.g.b.Get(dst)
		if p.Player != x.Player {
			return ErrIllegalPromotion
		}

		var err error
		u.Data, err = json.Marshal(x)
		if err != nil {
			return err
		}

		return nil
	},
	// lamo laziness
	order.Done: func(c *Client, u *order.Order) error {
		x, ok := u.Parameter.(int8)
		if !ok {
			return ErrUpdateParameter
		}

		var err error
		u.Data, err = json.Marshal(order.DoneModel{
			Result: x,
		})
		if err != nil {
			return err
		}

		return nil
	},
}
