package game

import (
	"encoding/json"
	"fmt"

	"github.com/go-playground/validator"
	"github.com/toms1441/chess-server/internal/board"
	"github.com/toms1441/chess-server/internal/order"
)

// Command is a communication structure sent from the client to the server.
// Data needs to be encoded in JSON, and each command has it's own parameters. Defined in order/model.go

type CommandCallback func(c *Client, o order.Order) error

var cbs map[uint8]CommandCallback
var validate = validator.New()

func init() {
	cbs = map[uint8]CommandCallback{
		order.Move: func(c *Client, o order.Order) error {
			g := c.g
			if !g.IsTurn(c) {
				return ErrIllegalTurn
			}

			if c.inPromotion() {
				return ErrInPromotion
			}

			s := &order.MoveModel{}

			err := json.Unmarshal(o.Data, s)
			// unmarshal the order
			if err != nil {
				return fmt.Errorf("json.Unmarshal: %w", err)
			}

			err = validate.Struct(s)
			if err != nil {
				return fmt.Errorf("validation: %w", err)
			}

			if !board.BelongsTo(s.ID, c.p1) {
				return ErrIllegalMove
			}

			pec, err := g.b.GetByIndex(int(s.ID))
			// check that piece is valid
			if err != nil || !pec.Valid() {
				return ErrPieceNil
			}

			// disallow enemy moving ally pieces
			if pec.P1 != c.p1 {
				return ErrIllegalMove
			}

			// do the order
			ret := g.b.Move(int(s.ID), s.Dst)
			if ret == false {
				return ErrIllegalMove
			}

			// first off update about the move...
			err = g.UpdateAll(order.Order{
				ID:   order.Move,
				Data: o.Data,
			})
			if err != nil {
				return err
			}

			// then, if it's not a promotion switch turns...
			if !(s.Dst.Y == 7 || s.Dst.Y == 0) {
				// promotion
				g.SwitchTurn()
			} else {
				if pec.Kind != board.Pawn {
					g.SwitchTurn()
				}
			}

			return nil
		},
		order.Promote: func(c *Client, o order.Order) error {
			g := c.g
			s := &order.PromoteModel{}

			err := json.Unmarshal(o.Data, s)
			// unmarshal the order
			if err != nil {
				return fmt.Errorf("json.Unmarshal: %w", err)
			}

			err = validate.Struct(s)
			if err != nil {
				return fmt.Errorf("validation: %w", err)
			}

			pec, err := g.b.GetByIndex(s.ID)
			if err != nil {
				return board.ErrEmptyPiece
			}
			if pec.Kind != board.Pawn {
				return ErrIllegalPromotion
			}
			if s.Kind != board.Bishop && s.Kind != board.Knight && s.Kind != board.Rook && s.Kind != board.Queen {
				// only allow [bishop, knight, rook, queen]
				return ErrIllegalPromotion
			}

			pec.Kind = s.Kind

			err = g.UpdateAll(order.Order{
				ID: order.Promotion,
				Parameter: order.PromotionModel{
					ID:   s.ID,
					Kind: s.Kind,
				},
			})
			if err != nil {
				return err
			}

			g.SwitchTurn()
			return nil
		},
		order.Castling: func(c *Client, o order.Order) error {
			if !c.g.IsTurn(c) {
				return ErrIllegalTurn
			}
			if !c.g.canCastle[c.p1] {
				return ErrIllegalCastling
			}

			if c.inPromotion() {
				return ErrInPromotion
			}

			cast := order.CastlingModel{}
			err := json.Unmarshal(o.Data, &cast)
			// unmarshal the order
			if err != nil {
				return fmt.Errorf("json.Unmarshal: %w", err)
			}

			err = validate.Struct(cast)
			if err != nil {
				return fmt.Errorf("validation: %w", err)
			}

			kingid := board.GetKing(c.p1)
			rookid := 0

			rid := board.GetRooks(c.p1)
			r1, r2 := rid[0], rid[1]
			if (kingid != cast.Src && kingid != cast.Dst) || cast.Src != r1 && cast.Dst != r1 && cast.Src != r2 && cast.Dst != r2 {
				fmt.Println("debug 3")
				return ErrIllegalCastling
			}

			if cast.Src == r1 || cast.Dst == r1 {
				rookid = r1
			} else if cast.Src == r2 || cast.Dst == r2 {
				rookid = r2
			}

			brd := c.g.b
			pecrook, err := brd.GetByIndex(rookid)
			if err != nil {
				return board.ErrEmptyPiece
			}
			pecking, err := brd.GetByIndex(kingid)
			if err != nil {
				return board.ErrEmptyPiece
			}

			minx, maxx := pecrook.Pos.X, pecking.Pos.X
			if minx > maxx {
				minx, maxx = maxx, minx
			}

			y := board.GetStartRow(c.p1)
			for x := minx; x < maxx; x++ {
				if x == 0 || x == 4 || x == 7 { // skip king and rook
					continue
				}

				_, _, err := brd.Get(board.Point{x, y})
				if err == nil {
					return ErrIllegalCastling
				}
			}

			if minx == 4 {
				brd.Set(rookid, board.Point{5, y})
				brd.Set(kingid, board.Point{6, y})
			} else if minx == 0 {
				brd.Set(rookid, board.Point{3, y})
				brd.Set(kingid, board.Point{2, y})
			}

			body, err := json.Marshal(order.CastlingModel{
				Src: kingid,
				Dst: rookid,
			})
			if err != nil {
				return err
			}

			err = c.g.UpdateAll(order.Order{
				ID:   order.Castling,
				Data: body,
			})
			if err != nil {
				return err
			}

			c.g.SwitchTurn()

			return nil
		},
		order.Done: func(c *Client, o order.Order) error {
			oth := board.GetInversePlayer(c.p1)

			c.g.done = true

			return c.g.UpdateAll(order.Order{
				ID:        order.Done,
				Parameter: oth,
			})
		},
	}

}
