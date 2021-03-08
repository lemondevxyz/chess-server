// Package board provides game-logic for chess, without the need of interaction from the user.
package board

import (
	"encoding/json"
)

// MoveEvent is a function called post-movement of a piece, ret is a boolean representing the validity of the move.
type MoveEvent func(p *Piece, src Point, dst Point, ret bool)

type Board struct {
	data [8][8]*Piece
	// move event listener
	ml []MoveEvent
}

// NewBoard creates a new board with the default placement.
func NewBoard() *Board {
	b := Board{
		ml: []MoveEvent{},
	}

	row := [2][8]uint8{
		{
			Rook,
			Knight,
			Bishop,
			King,
			Queen,
			Bishop,
			Knight,
			Rook,
		},
		{
			PawnB,
			PawnB,
			PawnB,
			PawnB,
			PawnB,
			PawnB,
			PawnB,
			PawnB,
		},
	}

	for x, s := range row {
		for y, v := range s {
			// x := k + 6
			b.data[x][y] = &Piece{
				T:      v,
				Player: 1,
				Pos:    Point{x, y},
			}
		}
	}

	row[0], row[1] = row[1], row[0]
	for k, s := range row {
		for y, v := range s {
			if v == PawnB {
				v = PawnF
			}

			x := k + 6
			b.data[x][y] = &Piece{
				T:      v,
				Player: 2,
				Pos:    Point{x, y},
			}
		}
	}

	return &b
}

// String method returns a string. makes it easier to debug
func (b *Board) String() (str string) {
	for k, s := range b.data {
		if k != 0 {
			str += "\n"
		}

		for _, v := range s {
			if v == nil {
				str += "  "
			} else {
				str += v.ShortString() + " "
			}
		}
	}

	return str
}

// Listen returns adds a callback that gets called pre and post movement.
func (b *Board) Listen(callback MoveEvent) {
	b.ml = append(b.ml, callback)
}

// Set sets a piece in the board without game-logic interfering.
func (b *Board) Set(p *Piece) {
	if p != nil {
		if p.T == Empty {
			b.data[p.Pos.X][p.Pos.Y] = nil
		} else {
			b.data[p.Pos.X][p.Pos.Y] = p
		}
	}
}

// Get returns a piece
func (b *Board) Get(src Point) *Piece {
	return b.data[src.X][src.Y]
}

// Possib is the same as Piece.Possib, but with removal of illegal moves.
func (b *Board) Possib(p *Piece) Points {
	ps := p.Possib()
	if p.T != Knight && p.T != PawnB && p.T != PawnF {
		orix, oriy := p.Pos.X, p.Pos.Y

		// starting from x, y this function loops through possible points
		// afterwards it changes the value via op function which receives x, y and modifies them
		// in-case it encountered a piece in the way it wait to finish and removes all the following points
		loop := func(x, y int, op func(int, int) (int, int)) {
			rm := false
			for i := 0; i < 8; i++ {
				pnt := Point{x, y}
				if !ps.In(pnt) || !pnt.Valid() {
					break
				}
				if !rm {
					// encountered piece in the way
					o := b.Get(pnt)
					if o != nil {
						if p.Player == o.Player {
							index := ps.Index(o.Pos)
							if index >= 0 {
								ps[index] = ps[len(ps)-1]
								ps = ps[:len(ps)-1]
							}
						}
						rm = true
					} else {
						// this direction cannot possibly have following points
						break
					}
				} else {
					// start deleting following points, cause we reached a piece in the way
					index := ps.Index(pnt)
					if index >= 0 {
						ps[index] = ps[len(ps)-1]
						ps = ps[:len(ps)-1]
					}
				}
				x, y = op(x, y)
			}
		}

		x, y := orix, oriy
		// normal direction
		{
			x, y = Up(orix, oriy)
			loop(x, y, Up)
			x, y = Down(orix, oriy)
			loop(x, y, Down)
			x, y = Left(orix, oriy)
			loop(x, y, Left)
			x, y = Right(orix, oriy)
			loop(x, y, Right)
		}

		// combination direction
		{
			x, y = UpLeft(orix, oriy)
			loop(x, y, UpLeft)
			x, y = UpRight(orix, oriy)
			loop(x, y, UpRight)
			x, y = DownLeft(orix, oriy)
			loop(x, y, DownLeft)
			x, y = DownRight(orix, oriy)
			loop(x, y, DownRight)
		}
	}

	return ps
}

/*
func (b *Board) Possib(p *Piece) Points {
	ps := p.Possib()
	if p.T != Knight && p.T != PawnB && p.T != PawnF {
		dir := p.Pos.Direction(dst)

		x, y := p.Pos.X, p.Pos.Y
		stopat := Point{-1, -1}

		// here we see if there's a piece in the way..
		// if so then we remove that point from the possible points we could go through
		for i := 0; i < 8; i++ {
			// direction movement

			pos := Point{x, y}
			if !pos.Valid() {
				break
			}
			if pos.Equal(dst) {
				break
			} else {
				o := b.Get(pos)
				if o != nil && o.T != Empty {
					// there's a piece in the way
					stopat = pos
					break
				}
			}
		}

		// this means there was a piece in the way
		if stopat.Valid() {
			// start from the piece in the way, and cancel all the "next" moves
			x, y = stopat.X, stopat.Y
			// so for example, if we want to go to (4, 3) and there is a piece in (2, 1) - (3, 2) would not be possible
			// as well as (4,3) and so on.

			// therefore we don't really need
			var fn func(Point) bool
			if Has(dir, DirDown) {
				fn = stopat.Bigger
			} else if Has(dir, DirUp) {
				fn = stopat.Smaller
			} else if Has(dir, DirLeft) {
				fn = stopat.Smaller
			} else if Has(dir, DirRight) {
				fn = stopat.Bigger
			}

			for i := len(ps) - 1; i >= 0; i-- {
				v := ps[i]
				if fn(v) {
					ps = append(ps[:i], ps[i+1:]...)
				}
			}
		}
	}

	return ps
}
*/

// Move moves a piece from it's original position to the destination. Returns true if it did, or false if it didn't.
func (b *Board) Move(p *Piece, dst Point) (ret bool) {
	defer func() {
		src := p.Pos
		if p != nil && ret {
			b.data[p.Pos.X][p.Pos.Y] = nil

			p.Pos.X = dst.X
			p.Pos.Y = dst.Y

			b.data[dst.X][dst.Y] = p
		}

		for _, v := range b.ml {
			v(p, src, dst, ret)
		}
	}()

	if p != nil {
		if b.Get(p.Pos) != p {
			return
		}

		o := b.Get(dst)
		if p.CanGo(dst) {
			if o != nil && o.T != Empty {
				// if it's not pawn and there's a piece in dst.
				// then kill it
				// if it's pawn then don't move
				if p.Player != o.Player {
					if p.T != PawnF && p.T != PawnB {
						ret = b.Possib(p).In(dst)
					}
				}
			} else {
				ret = b.Possib(p).In(dst)
			}
		} else {
			if p.T == PawnB || p.T == PawnF {
				x := p.Pos.X
				if p.T == PawnB {
					x++
				} else if p.T == PawnF {
					x--
				}

				ps := Points{
					{x, p.Pos.Y + 1},
					{x, p.Pos.Y - 1},
				}
				if ps.In(dst) {
					if o != nil && o.T != Empty && o.Player != p.Player {
						ret = true
					}
				}
			}
		}
	}

	return
}

// MarshalJSON json.Marshaler
func (b *Board) MarshalJSON() ([]byte, error) {
	body, err := json.Marshal(b.data)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// UnmarshalJSON json.Unmarshaler
func (b *Board) UnmarshalJSON(body []byte) error {
	b.ml = []MoveEvent{}

	err := json.Unmarshal(body, &b.data)
	if err != nil {
		return err
	}

	size := len(b.data)
	for x := 0; x < size; x++ {
		for y := 0; y < size; y++ {
			p := b.data[x][y]
			if p != nil {
				p.Pos.X = x
				p.Pos.Y = y
			}
		}
	}

	return nil
}

// DeadPieces returns all the dead pieces
func (b Board) DeadPieces(player uint8) map[uint8]uint8 {
	x := map[uint8]uint8{
		PawnF:  8,
		PawnB:  8,
		Bishop: 2,
		Knight: 2,
		Rook:   2,
		King:   1,
		Queen:  1,
	}

	for _, s := range b.data {
		for _, v := range s {
			if v != nil && v.Player == player {
				_, ok := x[v.T]
				if ok {
					x[v.T]--
					if x[v.T] == 0 {
						delete(x, v.T)
					}
				}
			}
		}
	}

	return x
}
