package game

import (
	"encoding/json"
	"fmt"

	"github.com/toms1441/chess-server/internal/order"
)

// Command is a communication structure sent from the client to the server.
// Data should be encoded in JSON, and each command has it's own parameters.
/*
type Command struct {
	ID   uint8           `validate:"required" json:"id"`
	Data json.RawMessage `validate:"required" json:"data"`
}
*/

type CommandCallback func(c *Client, o order.Order) error

/*
const (
	// CmdPiece is used whenever a player wants to move one of their pieces.
	// Data parameters are `{src: {x: 3, y: 3}, dst: {x: 5, y: 3}}`
	CmdPiece uint8 = iota + 1
	// CmdPromotion is what happens when a pawn reaches the end of the board, allowing the player to replace the pawn with a dead piece.
	// Data parameters are `{id: 1}` - for example pawn
	// Piece IDs are stored in the board package
	CmdPromotion
	// CmdPauseGame is used to pause the game, the other player has to also send it to confirm.
	// Data parameters are `` - well, empty.
	CmdPauseGame
	// CmdSendMessage is used to send chat messages to the other player.
	// Data parameters are `{message: "hello world"}`
	CmdMessage
)
*/

var cbs = map[uint8]CommandCallback{
	order.Move: func(c *Client, o order.Order) error {
		g := c.g
		if !g.IsTurn(c) {
			return ErrIllegalTurn
		}

		s := &order.MoveModel{}

		err := json.Unmarshal(o.Data, s)
		if err != nil {
			return err
		}

		p := g.b.Get(s.Src)
		if p == nil {
			return ErrPieceNil
		}

		ret := g.b.Move(p, s.Dst)
		if ret == false {
			return ErrIllegalMove
		}

		g.SwitchTurn()

		return g.UpdateAll(order.Order{
			ID:   order.Move,
			Data: o.Data,
		})
	},
	order.Promote: func(c *Client, o order.Order) error {
		g := c.g
		if !g.IsTurn(c) {
			return ErrIllegalTurn
		}

		s := &order.PromoteModel{}

		err := json.Unmarshal(o.Data, s)
		if err != nil {
			return err
		}

		dps := g.b.DeadPieces(c.num)
		_, ok := dps[s.Type]

		if !ok {
			return ErrIllegalPromotion
		}

		p := g.b.Get(s.Src)
		if p == nil {
			return ErrPieceNil
		}

		p.T = s.Type
		g.SwitchTurn()

		return g.UpdateAll(order.Order{
			ID:   order.Promotion,
			Data: o.Data,
		})
	},
	/* TODO: implement later, specifically after chess is working
	CmdPauseGame: func(c *Client, m *Command) error {
		g := c.g
		if g == nil {
			return ErrGameNil
		}

		s := []struct {
			Player1 bool
			Player2 bool
		}{}

		return nil
	},
	*/
	order.Message: func(c *Client, o order.Order) error {
		g := c.g
		s := &order.MessageModel{}

		err := json.Unmarshal(o.Data, s)
		if err != nil {
			return err
		}

		s.Message = fmt.Sprintf("[Player %d]: %s", c.num, s.Message)
		data, err := json.Marshal(s)
		if err != nil {
			return err
		}

		return g.UpdateAll(order.Order{
			ID:   order.Message,
			Data: data,
		})
	},
}
