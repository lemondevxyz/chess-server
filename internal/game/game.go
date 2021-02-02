package game

import (
	"encoding/json"

	"github.com/toms1441/chess/internal/board"
)

type Game struct {
	id     string
	cs     [2]*Client
	turn   chan int
	listen []chan int
	done   bool
	b      *board.Board
	cmd    map[string]interface{} // command-specific data
}

func NewGame(cl1, cl2 *Client) (*Game, error) {

	if cl1 == nil || cl2 == nil || cl1.W == nil || cl2.W == nil {
		return nil, ErrClientNil
	}

	cl1.num = 1
	cl2.num = 2

	g := &Game{
		cs:   [2]*Client{cl1, cl2},
		turn: make(chan int),
		cmd:  map[string]interface{}{},
	}

	cl1.g = g
	cl2.g = g

	go func(g *Game) {
		for !g.done {
			select {
			case a := <-g.turn:
				var c *Client
				if a == 1 {
					c = g.cs[0]
				} else if a == 2 {
					c = g.cs[1]
				}

				if c == nil {
					break
				}
			}
		}
	}(g)

	g.b = board.NewBoard()

	g.b.Listen(func(p *board.Piece, src board.Point, dst board.Point, ret bool) {
		if ret {
			if p.T == board.PawnB || p.T == board.PawnF {
				if dst.X == 7 || dst.X == 1 {
					c := g.cs[p.Player-1]
					if c != nil {
						g.Update(c, Update{
							ID: UpdatePromotion,
							parameter: ModelUpdatePromotion{
								Player: p.Player,
								Dst:    dst,
							},
						})
					}
				}
			}
		}
	})

	return g, nil
}

func (g *Game) Update(c *Client, u Update) error {
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

	/*
		cherr := make(chan error)
		go func() {
			_, err = c.W.Write(body)
			cherr <- err
		}()

		select {
		case <-time.After(time.Second * 10):
			return ErrUpdateTimeout
		case err := <-cherr:
			return err
		}

		return nil
	*/

	go func() {
		c.W.Write(body)
	}()

	return nil
}

func (g *Game) UpdateAll(u Update) error {
	err := g.Update(g.cs[0], u)
	if err != nil {
		return err
	}

	return g.Update(g.cs[1], u)
}
