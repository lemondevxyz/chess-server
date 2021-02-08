package game

import (
	"encoding/json"
	"fmt"
)

// Command is a communication structure sent from the client to the server.
// Data should be encoded in JSON, and each command has it's own parameters.
type Command struct {
	ID   uint8           `validate:"required" json:"id"`
	Data json.RawMessage `validate:"required" json:"data"`
}

type CommandBallback func(c *Client, m Command) error

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

var cbs = map[uint8]CommandBallback{
	CmdPiece: func(c *Client, m Command) error {
		g := c.g
		if g == nil {
			return ErrGameNil
		}

		s := &ModelCmdPiece{}

		err := json.Unmarshal(m.Data, s)
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

		return g.UpdateAll(Update{
			ID: UpdateBoard,
		})
	},
	CmdPromotion: func(c *Client, m Command) error {
		g := c.g
		if g == nil {
			return ErrGameNil
		}

		s := &ModelCmdPromotion{}

		err := json.Unmarshal(m.Data, s)
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

		return g.UpdateAll(Update{
			ID: UpdateBoard,
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
	CmdMessage: func(c *Client, m Command) error {
		g := c.g
		if g == nil {
			return ErrGameNil
		}

		s := &ModelCmdMessage{}

		err := json.Unmarshal(m.Data, s)
		if err != nil {
			return err
		}

		s.Message = fmt.Sprintf("[Player] %d: %s", c.num, s.Message)

		return g.UpdateAll(Update{
			ID:   UpdateMessage,
			Data: m.Data,
		})
	},
}
