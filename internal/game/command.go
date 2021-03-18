package game

import (
	"encoding/json"
	"fmt"

	"github.com/toms1441/chess-server/internal/board"
	"github.com/toms1441/chess-server/internal/order"
)

// Command is a communication structure sent from the client to the server.
// Data needs to be encoded in JSON, and each command has it's own parameters. Defined in order/model.go

type CommandCallback func(c *Client, o order.Order) error

var cbs map[uint8]CommandCallback

func init() {
	cbs = map[uint8]CommandCallback{
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

			pec := g.b.Get(s.Src)
			if pec == nil {
				return ErrPieceNil
			}

			if pec.Player != c.num {
				return ErrIllegalMove
			}

			cep := g.b.Get(s.Dst)
			// if the destination contains a piece
			if cep != nil {
				// if it's in the same team as a the player
				if cep.Player == pec.Player {
					// if one is a king, and the other is rook
					if (cep.T == board.King && pec.T == board.Rook) || (cep.T == board.Rook || pec.T == board.King) {
						//fn := cbs[order.Castling]
						/* mtx.Lock freezes this
						return c.Do(order.Order{
							ID:   order.Castling,
							Data: o.Data,
						})
						*/

						return cbs[order.Castling](c, o)
					}
				}
			}

			ret := g.b.Move(pec, s.Dst)
			if ret == false {
				return ErrIllegalMove
			}

			//fmt.Println("switch turn")
			if !(s.Dst.X == 7 || s.Dst.X == 0) {
				// promotion
				g.SwitchTurn()
			} else {
				if pec.T != board.PawnF && pec.T != board.PawnB {
					g.SwitchTurn()
				}
			}

			return g.UpdateAll(order.Order{
				ID:   order.Move,
				Data: o.Data,
			})
		},
		order.Promote: func(c *Client, o order.Order) error {
			g := c.g
			s := &order.PromoteModel{}

			err := json.Unmarshal(o.Data, s)
			if err != nil {
				return err
			}

			if s.Src.X != 0 && s.Src.X != 7 {
				return ErrIllegalPromotion
			}

			/* i'm not yet sure about this part
			dps := g.b.DeadPieces(c.num)
			_, ok := dps[s.Type]
			if !ok {
				return ErrIllegalPromotion
			}
			*/

			p := g.b.Get(s.Src)
			if p == nil {
				return ErrPieceNil
			}
			if p.T != board.PawnF && p.T != board.PawnB {
				return ErrIllegalPromotion
			}

			p.T = s.Type
			g.SwitchTurn()

			return g.UpdateAll(order.Order{
				ID: order.Promotion,
				Parameter: order.PromotionModel{
					Dst:  s.Src,
					Type: s.Type,
				},
			})
		},
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
		order.Castling: func(c *Client, o order.Order) error {
			g := c.g
			if !g.IsTurn(c) {
				return ErrIllegalTurn
			}
			if !g.canCastle[c.num] {
				//fmt.Println("here 1")
				return ErrIllegalCastling
			}

			x := 7
			if c.num == 2 {
				x = 0
			}

			kingy := 4

			cast := &order.CastlingModel{}
			err := json.Unmarshal(o.Data, cast)
			if err != nil {
				//fmt.Println("here 2")
				return err
			}

			src := cast.Src
			dst := cast.Dst

			ok := func() bool {
				if src.X != x || src.Y != kingy {
					return false
				}
				if dst.X != x || (dst.Y != 0 && dst.Y != 7) {
					//fmt.Println("here 3")
					return false
				}

				return true
			}

			if !ok() {
				src, dst = dst, src
				if !ok() {
					// neither is a king or a rook
					return ErrIllegalCastling
				}
			}

			miny := dst.Y
			if miny == kingy {
				miny = src.Y
			}

			maxy := kingy
			if miny > maxy {
				miny, maxy = maxy, miny
			}

			for y := miny; y < miny; y++ {
				// king's position
				if y == kingy {
					continue
				}

				pec := g.b.Get(board.Point{x, y})
				// pieces that are in the way
				if pec != nil && pec.T != board.Empty {
					//fmt.Println(pec, pec.Pos, "here 4")
					return ErrIllegalCastling
				}
			}

			rooky := dst.Y

			rook, king := g.b.Get(board.Point{x, rooky}), g.b.Get(board.Point{x, kingy})
			if rook == nil || king == nil || rook.T != board.Rook || king.T != board.King { // somehow ??
				//fmt.Println("here 5")
				return ErrIllegalCastling
			}

			g.b.Set(&board.Piece{
				Pos: board.Point{x, rooky},
				T:   board.Empty,
			})
			g.b.Set(&board.Piece{
				Pos: board.Point{x, kingy},
				T:   board.Empty,
			})

			if rooky == 0 {
				kingy = 2
				rooky = 3
			} else if rooky == 7 {
				kingy = 6
				rooky = 5
			}

			king.Pos.Y = kingy
			rook.Pos.Y = rooky

			g.b.Set(king)
			g.b.Set(rook)

			g.SwitchTurn()

			return g.UpdateAll(order.Order{
				ID:   order.Castling,
				Data: o.Data,
			})
		},
	}

}
