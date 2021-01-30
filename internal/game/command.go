package game

import (
	"encoding/json"

	"github.com/toms1441/chess/internal/board"
)

// Command is a communication structure sent from the client to the server.
// Data should be encoded in JSON, and each command have their own parameters.
type Command struct {
	ID   uint8
	Data []byte
}

type CommandBallback func(c *Client, m *Command) error

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
	CmdSendMessage
)

var cbs = map[uint8]CommandBallback{
	CmdPiece: func(c *Client, m *Command) error {
		g := c.g
		if g == nil {
			return ErrGameNil
		}

		data := struct {
			Src board.Point `json:"src"`
			Dst board.Point `json:"dst"`
		}{}

		err := json.Unmarshal(m.Data, &data)
		if err != nil {
			return err
		}

		p := g.b.Get(data.Src)
		if p == nil {
			return ErrPieceNil
		}

		ret := g.b.Move(p, data.Dst)
		if ret == false {
			return ErrIllegalMove
		}

		err = g.UpdateAll(Update{
			ID: UpdateBoard,
		})
		if err != nil {
			return err
		}

		return nil
	},
	CmdPromotion: func(c *Client, m *Command) error {
		g := c.g
		if g == nil {
			return ErrGameNil
		}

		data := struct {
			Src board.Point `json:"src"`
			ID  uint8       `json:"id"`
		}{}

		err := json.Unmarshal(m.Data, &data)
		if err != nil {
			return err
		}

		dps := g.b.DeadPieces()
		_, ok := dps[data.ID]
		if !ok {
			return ErrIllegalPromotion
		}

		p := g.b.Get(data.Src)
		if p == nil {
			return ErrPieceNil
		}
		p.T = data.ID

		err = g.UpdateAll(Update{
			ID: UpdateBoard,
		})
		if err != nil {
			return err
		}

		return nil
	},
	CmdPauseGame: func(c *Client, m *Command) error {
		g := c.g
		if g == nil {
			return ErrGameNil
		}

		return nil
	},
}
