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
			Queen,
			King,
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
				Player: 2,
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
				Player: 1,
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
func (b Board) Get(src Point) *Piece {
	return b.data[src.X][src.Y]
}

// Possib is the same as Piece.Possib, but with removal of illegal moves.
func (b Board) Possib(pec *Piece) Points {
	ps := pec.Possib()
	switch pec.T {
	case Knight: // disallow movement to allies
		for i := len(ps) - 1; i >= 0; i-- {

			v := ps[i]
			cep := b.Get(v)
			if cep != nil {
				if cep.Player == pec.Player {
					ps[i] = ps[len(ps)-1]
					ps = ps[:len(ps)-1]
				}
			}
		}
	case PawnF, PawnB: // disallow movement in front of piece
		for i := len(ps) - 1; i >= 0; i-- {
			pnt := ps[i]
			if b.Get(pnt) != nil { // piece in the way ...
				ps[i] = ps[len(ps)-1]
				ps = ps[:len(ps)-1]
			}
		}

		// if our move has a piece in the way then cancel
		// also if we're at 6 or 1, then allow movement to 4 or 3
		x := pec.Pos.X - 1
		if pec.T == PawnB {
			x = pec.Pos.X + 1
		}

		sp := Points{
			{x, pec.Pos.Y - 1},
			{x, pec.Pos.Y + 1},
		}
		sp = sp.Clean()

		for i := len(sp) - 1; i >= 0; i-- {
			v := sp[i]

			cep := b.Get(v)
			// is there a piece
			if cep != nil {
				// is it the enemy's
				if cep.Player != pec.Player {
					// then don't remove this move from the possible moves
					continue
				}
			}

			// empty piece or piece is ours
			// so remove it
			// pawn cannot kill it's friend
			sp[i] = sp[len(sp)-1]
			sp = sp[:len(sp)-1]
		}

		ps = ps.Merge(ps, sp)
	default:
		orix, oriy := pec.Pos.X, pec.Pos.Y

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
					cep := b.Get(pnt)
					if cep != nil {
						if pec.Player == cep.Player {
							index := ps.Index(cep.Pos)
							if index >= 0 {
								ps[index] = ps[len(ps)-1]
								ps = ps[:len(ps)-1]
							}
						}
						rm = true
					} /*else {
						// this direction cannot possibly have following points
						break
					}*/
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

	return ps.Clean()
}

func (b Board) Checkmate(player uint8) bool {
	var king *Piece
	for _, s := range b.data {
		for _, v := range s {
			if v != nil {
				if v.T == King && v.Player == player {
					king = v
					break
				}
			}
		}
	}

	for _, s := range b.data {
		for _, v := range s {
			if v != nil {
				if v.Player != player {
					possib := v.Possib()
					if possib.In(king.Pos) {
						return true
					}
				}
			}
		}
	}

	return false
}

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
		// can we legally go there, i.e is it in the possible combinations??
		// so for example bishop cannot go horizontally
		if p.CanGo(dst) {
			// is there a piece in the destination??
			if o != nil && o.T != Empty {
				// is the piece's an enemy
				if p.Player != o.Player {
					// is it not a pawn(cause pawns cannot enemy forward or backward of them)
					if p.T != PawnF && p.T != PawnB {
						ret = b.Possib(p).In(dst)
					}
				}
			} else {
				// no piece in the destination
				ret = b.Possib(p).In(dst)
			}
		} else {
			if p.T == PawnF || p.T == PawnB {
				ret = b.Possib(p).In(dst)
			}
		}
	}

	return
}

// MarshalJSON json.Marshaler
func (b Board) MarshalJSON() ([]byte, error) {
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
