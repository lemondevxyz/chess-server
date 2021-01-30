package game

import "io"

// Client is a struct used for the server to communicate to the client.
type Client struct {
	// W where to write updates
	w   io.WriteCloser
	num uint8 // player 1 or 2??
	// ID the ID used to authenticate commands
	id string
	// G the underlying game
	g *Game
}

func (c Client) Send() {
}
