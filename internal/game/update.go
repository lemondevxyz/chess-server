package game

import (
	"encoding/json"

	"github.com/toms1441/chess-server/internal/model"
)

// Update is a communication structure from the server to the client, while Command is from the client to the server.

// UpdateCallback sets the data for the update, since some updates are quite repetitive.
type UpdateCallback func(c *Client, u *model.Order) error

// Tests are bundled with command tests.

// redundant updates go here
// as well as verification for certain updates.
var ubs = map[uint8]UpdateCallback{
	model.OrMove: func(c *Client, u *model.Order) error {
		x, ok := u.Parameter.([]byte)
		if !ok {
			return ErrUpdateParameter
		}
		u.Data = x

		return nil
	},
	model.OrPromote: func(c *Client, u *model.Order) error {
		x, ok := u.Parameter.(model.PromoteOrder)
		if !ok {
			return ErrUpdateParameter
		}

		var err error
		u.Data, err = json.Marshal(x)
		if err != nil {
			u.Data = nil
			return err
		}

		return nil
	},
	model.OrPromotion: func(c *Client, u *model.Order) error {
		x, ok := u.Parameter.(model.PromotionOrder)
		if !ok {
			return ErrUpdateParameter
		}

		var err error
		u.Data, err = json.Marshal(x)
		if err != nil {
			return err
		}

		return nil
	},
	model.OrCheckmate: func(c *Client, u *model.Order) error {
		x, ok := u.Parameter.(bool)
		if !ok {
			return ErrUpdateParameter
		}

		body, err := json.Marshal(model.CheckmateOrder{
			P1: x,
		})
		if err != nil {
			return err
		}

		u.Data = body
		return nil
	},
	// lamo laziness
	model.OrDone: func(c *Client, u *model.Order) error {
		x, ok := u.Parameter.(bool)
		if !ok {
			return ErrUpdateParameter
		}

		var err error
		u.Data, err = json.Marshal(model.DoneOrder{
			P1: x,
		})
		if err != nil {
			return err
		}

		return nil
	},
}
