package game

import (
	"encoding/json"
	"io"
	"sync"

	"github.com/kjk/betterguid"
	"github.com/toms1441/chess-server/internal/board"
	"github.com/toms1441/chess-server/internal/model"
)

type Game struct {
	// cs is a map of p1 and !p1 linking them to a client pointer
	cs map[bool]*Client
	// turn is a flip flop of p1. SwitchTurn
	turn bool
	// done is set whenever the game ends
	done bool
	// listenDone is a channel that gets closed whenever the game ends
	listenDone chan struct{}
	// b is the board used for the game
	brd *board.Board
	// canCastle is a map containing if each player could castle or not. canCastle gets set to false whenever the king or either rook moves...
	canCastle map[bool]bool
	// spectators is a map of ids assigned to io writers.
	// Spectators cannot send commands, and only have access to the following updates:
	// OrMove, OrTurn, OrPromotion, OrCastling, OrCheckmate, OrDone
	// All spectator operations should be non-blocking, and should be ignored if they fail
	spectators map[string]io.Writer
	mtx        sync.Mutex
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
		listenDone: make(chan struct{}),
	}

	cl1.g, cl2.g = g, g

	g.brd = board.NewBoard()

	g.brd.Listen(func(id int8, p board.Piece, src board.Point, dst board.Point) {
		if p.Kind == board.Pawn {
			if dst.Y == 7 || dst.Y == 0 {
				c := g.cs[p.P1]
				if c != nil {

					x := model.PromoteOrder{
						ID: id,
					}

					g.Update(c, model.Order{
						ID:        model.OrPromote,
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
	// old turn
	bef := g.turn
	// new turn
	aft := !g.turn

	// well, do we have a final checkmate on the other player?
	if g.brd.FinalCheckmate(aft) {
		// if so, gg
		g.UpdateAll(model.Order{
			ID:        model.OrDone,
			Parameter: bef,
		})

		g.done = true

		return
	}

	// change the turn
	g.mtx.Lock()
	g.turn = aft
	g.mtx.Unlock()

	x, _ := json.Marshal(model.TurnOrder{
		P1: aft,
	})

	if g.brd.Checkmate(aft) {
		g.UpdateAll(model.Order{
			ID:        model.OrCheckmate,
			Parameter: aft,
		})
	}
	g.UpdateAll(model.Order{ID: model.OrTurn, Data: x})
}

// IsTurn returns if it's the client's turn this time
func (g *Game) IsTurn(c *Client) bool {
	if c == nil {
		return false
	}

	return c.p1 == g.turn
}

// Update is used to send updates to the client, such as a movement of a piece.
func (g *Game) Update(c *Client, u model.Order) error {
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
func (g *Game) UpdateAll(u model.Order) error {
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

	go func() {
		for _, writer := range g.spectators {
			writer.Write(body)
		}
	}()

	return nil
}

// Board returns the actual board.
func (g *Game) Board() *board.Board { return g.brd }

// ListenForDone returns a channel that gets closed when the game ends.
func (g *Game) ListenForDone() chan struct{} { return g.listenDone }

// AddSpectator adds a spectator to the list of spectators in the game, and returns it's id.
func (g *Game) AddSpectator(spectator io.Writer) string {
	id := betterguid.New()

	g.mtx.Lock()
	defer g.mtx.Unlock()

	g.spectators[id] = spectator
	go func() {
		body, _ := json.Marshal(model.GameOrder{Brd: g.brd})
		body, _ = json.Marshal(model.Order{ID: model.OrGame, Data: body})

		spectator.Write(body)
	}()

	return id
}

// RmSpectator removes a spectator from the list of spectators, and is safe if the id is invalid
func (g *Game) RmSpectator(id string) {
	g.mtx.Lock()
	defer g.mtx.Unlock()

	delete(g.spectators, id)
}

// Close closes the game, and cleans up any data assigned to the clients or the game struct. It does not send a message to clients indicating that the game is closed
func (g *Game) close() {
	defer g.mtx.Unlock()
	g.mtx.Lock()
	go func() {
		close(g.listenDone)
	}()

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
