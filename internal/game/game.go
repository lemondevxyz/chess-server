package game

import (
	"encoding/json"

	"github.com/toms1441/chess-server/internal/board"
	"github.com/toms1441/chess-server/internal/order"
)

type Game struct {
	cs        map[bool]*Client // refers to p1
	turn      bool             // refers to p1
	done      bool
	b         *board.Board
	canCastle map[bool]bool // refers to p1
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

	cl1.p1 = true
	cl2.p1 = false

	g := &Game{
		cs: map[bool]*Client{
			true:  cl1,
			false: cl2,
		},
		turn: false,
		canCastle: map[bool]bool{
			true:  true,
			false: true,
		},
	}

	cl1.g, cl2.g = g, g

	g.b = board.NewBoard()

	g.b.Listen(func(id int, p board.Piece, src board.Point, dst board.Point) {
		if p.Kind == board.Pawn {
			if dst.Y == 7 || dst.Y == 0 {
				c := g.cs[p.P1]
				if c != nil {

					x := order.PromoteModel{
						ID: id,
					}

					g.Update(c, order.Order{
						ID:        order.Promote,
						Parameter: x,
					})
				}
			}
		} else if p.Kind == board.King || p.Kind == board.Rook {
			g.canCastle[p.P1] = false
		}
	})

	return g, nil
}

// SwitchTurn called after a player ends their turn, to notify the other player.
func (g *Game) SwitchTurn() {

	bef := g.turn
	aft := !g.turn

	if g.b.FinalCheckmate(aft) {
		g.UpdateAll(order.Order{
			ID:        order.Done,
			Parameter: bef,
		})

		g.done = true

		return
	}

	g.turn = aft

	x, _ := json.Marshal(order.TurnModel{
		P1: aft,
	})

	if g.b.Checkmate(aft) {
		g.Update(g.cs[aft], order.Order{
			ID:        order.Checkmate,
			Parameter: aft,
		})
		g.Update(g.cs[bef], order.Order{
			ID:        order.Checkmate,
			Parameter: aft,
		})
	}
	g.UpdateAll(order.Order{ID: order.Turn, Data: x})
}

// IsTurn returns if it's the client's turn this time
func (g *Game) IsTurn(c *Client) bool {
	return c.p1 == g.turn
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

// UpdateAll sends the update to all of the players. Difference between this and calling update individually is the data does not get re-marshalized.
// Use this whenever the data is the same between the two players
func (g *Game) UpdateAll(u order.Order) error {
	if g.cs[false] == nil || g.cs[true] == nil {
		return ErrClientNil
	}

	if u.Data == nil {
		x, ok := ubs[u.ID]
		if !ok {
			return ErrUpdateNil
		}

		err := x(g.cs[true], &u)
		if err != nil {
			return err
		}
	}

	body, err := json.Marshal(u)
	if err != nil {
		return err
	}

	g.cs[true].W.Write(body)
	g.cs[false].W.Write(body)

	return nil
}

// Board returns the actual board.
func (g *Game) Board() *board.Board {
	return g.b
}

// Close closes the game, and cleans up the clients
func (g *Game) close() {
	do := func(cl *Client) {
		if cl.g != nil {
			cl.mtx.Lock()
			cl.g = nil
			cl.mtx.Unlock()
		}
	}

	c1 := g.cs[true]
	c2 := g.cs[false]
	if c1 != nil {
		do(c1)
	}
	if c2 != nil {
		do(c2)
	}

	delete(g.cs, false)
	delete(g.cs, true)
}
