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
	b      board.Board
}

func NewGame(cl1, cl2 *Client) *Game {

	cl1.num = 1
	cl2.num = 2

	g := &Game{
		cs:   [2]*Client{cl1, cl2},
		turn: make(chan int),
	}

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

	return g
}

func (g *Game) Do(c Command) {
}

func (g *Game) Update(c *Client, u Update) error {
	if u.ID == UpdateBoard {
		d, err := json.Marshal(g.b)
		if err != nil {
			return err
		}

		u.Data = d
	}

	body, err := json.Marshal(u)
	if err != nil {
		return err
	}

	_, err = c.w.Write(body)
	if err != nil {
		return err
	}

	return nil
}

func (g *Game) UpdateAll(u Update) error {

	err := g.Update(g.cs[0], u)
	if err != nil {
		return err
	}

	err = g.Update(g.cs[1], u)
	if err != nil {
		return err
	}

	return nil
}
