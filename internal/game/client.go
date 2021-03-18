package game

import (
	"io"
	"sync"

	"github.com/toms1441/chess-server/internal/order"
)

// Client is a struct used for the server to communicate to the client.
type Client struct {
	// W where to write updates
	W   io.Writer
	num uint8 // player 1 or 2??
	id  string
	// G the underlying game
	g   *Game
	mtx sync.Mutex
}

func (c *Client) Do(cmd order.Order) error {
	//fmt.Println("lock")
	c.mtx.Lock()

	if c.g == nil {
		return ErrGameNil
	}

	x, ok := cbs[cmd.ID]
	if !ok {
		return ErrCommandNil
	}

	err := x(c, cmd)

	c.mtx.Unlock()
	if c.g != nil {
		if c.g.done { // we cannot do this in switch turn
			// cause it would freeze the program if testing
			c.g.Close()
		}
	}

	return err
}

func (c *Client) Game() *Game {
	return c.g
}

func (c *Client) LeaveGame() {
	g := c.g
	if g == nil {
		return
	}

	x := g.cs[0]
	if x == c {
		x = g.cs[1]
	}

	upd := order.Order{ID: order.Done, Parameter: int8(1)}
	c.g.Update(x, upd)

	upd.Parameter = int8(-1)
	c.g.Update(c, upd)

	c.g.Close()

}

func (c *Client) Number() uint8 {
	return c.num
}
