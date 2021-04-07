package game

import (
	"encoding/json"
	"fmt"

	"github.com/toms1441/chess-server/internal/board"
	"github.com/toms1441/chess-server/internal/order"
)

type Game struct {
	id        string
	cs        [2]*Client
	turn      uint8
	done      bool
	b         *board.Board
	canCastle map[uint8]bool
}

// NewGame creates a game for client 1 and client 2(cl1, cl2). It fails whenever the clients are already in a game, or one of them is nil.
// Note: you need to call game.SwitchTurn() to get an actual game going, as NewGame does not update the clients with the player turn.
func NewGame(cl1, cl2 *Client) (*Game, error) {

	if cl1 == nil || cl2 == nil || cl1.W == nil || cl2.W == nil {
		return nil, ErrClientNil
	}

	if cl1.g != nil || cl2.g != nil {
		return nil, ErrGameIsNotNil
	}

	cl1.num = 1
	cl1.id = fmt.Sprintf("Player %d", cl1.num)

	cl2.num = 2
	cl2.id = fmt.Sprintf("Player %d", cl2.num)

	g := &Game{
		cs:   [2]*Client{cl1, cl2},
		turn: 0,
		canCastle: map[uint8]bool{
			1: true,
			2: true,
		},
	}

	cl1.g, cl2.g = g, g

	g.b = board.NewBoard()

	g.b.Listen(func(p *board.Piece, src board.Point, dst board.Point, ret bool) {
		if ret {
			if p.T == board.PawnB || p.T == board.PawnF {
				if dst.X == 7 || dst.X == 0 {
					c := g.cs[p.Player-1]
					if c != nil {

						x := order.PromoteModel{
							Src: dst,
						}

						g.Update(c, order.Order{
							ID:        order.Promote,
							Parameter: x,
						})
					}
				}
			} else if p.T == board.King || p.T == board.Rook {
				g.canCastle[p.Player] = false
			}
		}
	})

	//g.SwitchTurn()

	return g, nil
}

// SwitchTurn called after a player ends their turn, to notify the other player.
func (g *Game) SwitchTurn() {

	bef := g.turn
	aft := g.turn
	if g.turn == 1 {
		aft = 2
	} else {
		aft = 1
	}
	if g.b.Checkmate(aft) {
		g.Update(g.cs[aft-1], order.Order{
			ID:        order.Checkmate,
			Parameter: aft,
		})
		g.Update(g.cs[bef-1], order.Order{
			ID:        order.Checkmate,
			Parameter: aft,
		})
	}
	if g.b.FinalCheckmate(aft) {
		upd := order.Order{ID: order.Done, Parameter: int8(1)}
		g.Update(g.cs[bef-1], upd)

		upd.Parameter = int8(-1)
		g.Update(g.cs[aft-1], upd)

		g.done = true

		return
	}

	g.turn = aft

	x, _ := json.Marshal(order.TurnModel{
		Player: aft,
	})

	g.UpdateAll(order.Order{ID: order.Turn, Data: x})
}

// IsTurn returns if it's the client's turn this time
func (g *Game) IsTurn(c *Client) bool {
	return c.num == g.turn
}

// Update is used to send updates to the client, such as a movement of a piece.
func (g *Game) Update(c *Client, u order.Order) error {
	if c == nil {
		return ErrClientNil
	}

	if u.Data == nil {
		x, ok := ubs[u.ID]
		if !ok {
			return ErrUpdateNil
		}

		err := x(c, &u)
		if err != nil {
			return err
		}
	}

	body, err := json.Marshal(u)
	if err != nil {
		return err
	}

	c.W.Write(body)

	return nil
}

// UpdateAll sends the update to all of the players.
func (g *Game) UpdateAll(u order.Order) error {
	if g.cs[0] == nil || g.cs[1] == nil {
		return ErrClientNil
	}

	if u.Data == nil {
		x, ok := ubs[u.ID]
		if !ok {
			return ErrUpdateNil
		}

		err := x(g.cs[0], &u)
		if err != nil {
			return err
		}
	}

	body, err := json.Marshal(u)
	if err != nil {
		return err
	}

	g.cs[0].W.Write(body)
	g.cs[1].W.Write(body)

	return nil
}

// Board returns the actual board.
func (g *Game) Board() *board.Board {
	return g.b
}

// Close closes the game, and cleans up the clients
func (g *Game) close() {
	do := func(cl *Client) {
		cl.mtx.Lock()
		cl.g = nil
		cl.mtx.Unlock()
	}

	if g.cs[0] != nil {
		do(g.cs[0])
		g.cs[0] = nil
	}
	if g.cs[1] != nil {
		do(g.cs[1])
		g.cs[1] = nil
	}
}
