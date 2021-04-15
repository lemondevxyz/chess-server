package game

import (
	"io"
	"sync"

	"github.com/toms1441/chess-server/internal/board"
	"github.com/toms1441/chess-server/internal/order"
)

// Client is a struct used for the server to communicate to the client.
type Client struct {
	// W where to write updates
	W   io.Writer
	p1  bool // player 1 or 2??
	id  string
	g   *Game
	mtx sync.Mutex
}

// Do executes a command. It automatically checks if the player is in a game, or if the command's ID is invalid.
// Use of cbs[cmd.ID] is discouraged.
func (c *Client) Do(cmd order.Order) error {
	if c.g == nil {
		return ErrGameNil
	}

	x, ok := cbs[cmd.ID]
	if !ok {
		return ErrCommandNil
	}

	c.mtx.Lock()
	err := x(c, cmd)
	c.mtx.Unlock()

	if c.g != nil {
		if c.g.done { // we cannot do this in switch turn
			// cause it would freeze the program
			c.g.close()
		}
	}

	return err
}

// Game returns the pointer to client's game
func (c *Client) Game() *Game {
	return c.g
}

// LeaveGame leaves the game for client. It's generally used for testing, and doesn't send a order.Done message after it finishes.
// Use of this function in production is generally discouraged, as it could freeze the game
func (c *Client) LeaveGame() {
	g := c.g
	if g == nil {
		return
	}

	x := c.g.cs[board.GetInversePlayer(c.p1)]
	g.Update(x, order.Order{
		ID:        order.Done,
		Parameter: x.p1,
	})

	c.g.close()
}

// P1 returns if the client is player one or two.
func (c *Client) P1() bool {
	return c.p1
}
