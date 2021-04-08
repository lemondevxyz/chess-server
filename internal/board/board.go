// Package board provides game-logic for chess, without the need of interaction from the user.
package board

import (
	"encoding/json"
)

// MoveEvent is a function called post-movement of a piece, ret is a boolean representing the validity of the move.
type MoveEvent func(p Piece, src Point, dst Point, ret bool)

type Board struct {
	// data [8][8]*Piece
	data [32]Piece
	// move event listener
	ml []MoveEvent
}

// NewBoard creates a new board with the default placement.
func NewBoard() *Board {
	b := Board{
		ml:   []MoveEvent{},
		data: [32]Piece{},
	}

	row := [32]uint8{
		0:  Rook,
		1:  Knight,
		2:  Bishop,
		3:  Queen,
		4:  King,
		5:  Bishop,
		6:  Knight,
		7:  Rook,
		8:  PawnB,
		9:  PawnB,
		10: PawnB,
		11: PawnB,
		12: PawnB,
		13: PawnB,
		14: PawnB,
		15: PawnB,
		16: PawnF,
		17: PawnF,
		18: PawnF,
		19: PawnF,
		20: PawnF,
		21: PawnF,
		22: PawnF,
		23: PawnF,
		24: Rook,
		25: Knight,
		26: Bishop,
		27: Queen,
		28: King,
		29: Bishop,
		30: Knight,
		31: Rook,
	}

	for k, s := range row {
		player := 2
		if k >= 16 {
			player = 1
		}

		x := int8(k / 8)
		y := int8(k % 8)

		if k >= 16 {
			x += 4
		}

		b.data[k] = Piece{
			T:      s,
			Player: uint8(player),
			Pos:    Point{x, y},
		}
	}

	return &b
}

func (b Board) Copy() *Board {
	o := Board{
		ml:   []MoveEvent{},
		data: [32]Piece{},
	}

	for k, s := range b.data {
		o.data[k] = s
	}

	return &o
}

func (b Board) string(def string) (str string) {
	for i := 0; i < (8 * 8); i++ {
		x := int8(i % 8)
		y := int8(i / 8)

		char := def

		pos := Point{x, y}
		for _, v := range b.data {
			if v.Pos.Equal(pos) {
				char = v.ShortString()
			}
		}

		str += char + " "
		if i != 0 && x == 0 {
			str += "\n"
		}
	}

	return str

}

// String method returns a string. makes it easier to debug
func (b Board) String() (str string) {
	return b.string(" ")
}

// Listen returns adds a callback that gets called pre and post movement.
func (b *Board) Listen(callback MoveEvent) {
	b.ml = append(b.ml, callback)
}

// Set sets a piece in the board without game-logic interfering.
func (b *Board) Set(i int, pos Point) error {
	if i >= len(b.data) {
		return ErrInvalidPoint
	}

	if !pos.Valid() {
		b.data[i].T = Empty
		b.data[i].Pos = Point{-1, -1}
	} else {
		b.data[i].Pos = pos
	}

	return nil
}

// Get returns a piece and it's index. Or otherwise -1, an empty piece and an error.
func (b Board) Get(src Point) (int, Piece, error) {
	for k, v := range b.data {
		if v.Pos.Equal(src) {
			if !v.Valid() {
				return k, v, ErrInvalidPoint
			}

			return k, v, nil
		}
	}

	return -1, Piece{}, ErrEmptyPiece
}

func (b Board) GetByIndex(i int) (Piece, error) {
	if i >= len(b.data) {
		return Piece{}, ErrInvalidPoint
	}

	return b.data[i], nil
}

// Possib is the same as Piece.Possib, but with removal of illegal moves.
func (b Board) Possib(id int) (Points, error) {
	pec, err := b.GetByIndex(id)
	if err != nil {
		return nil, ErrEmptyPiece
	}

	ps := pec.Possib()
	switch pec.T {
	case Knight: // disallow movement to allies
		for k, v := range ps {
			_, cep, err := b.Get(v)
			if err == nil {
				if cep.Player == pec.Player {
					delete(ps, k)
				}
			}
		}
	case PawnF, PawnB: // disallow movement in front of piece
		for k, pnt := range ps {
			_, _, err = b.Get(pnt)
			if err == nil { // piece in the way ...
				delete(ps, k)
			}
		}

		// if our move has a piece in the way then cancel
		// also if we're at 6 or 1, then allow movement to 4 or 3
		x := pec.Pos.X - 1
		if pec.T == PawnB {
			x = pec.Pos.X + 1
		}

		sp := Points{}
		sp.Insert(Point{x, pec.Pos.Y - 1})
		sp.Insert(Point{x, pec.Pos.Y + 1})

		sp.Clean()

		for k, v := range sp {
			_, cep, err := b.Get(v)
			// is there a piece
			if err == nil {
				// is it the enemy's
				if cep.Player != pec.Player {
					// then don't remove this move from the possible moves
					continue
				}
			}

			// empty piece or piece is ours
			// so remove it
			// pawn cannot kill it's friend
			delete(sp, k)
		}

		ps = ps.Merge(ps, sp)
	case King:
		// King's possibilities are it's square possibilities but if enemy threatens a point, it goes away.
		//
		// i.e if going to point kills the king, then he cannot go there.
		//
		// Now when calculating King's possibilities we need to check every enemy's piece possibilities and if they threaten the king.
		// But when checking for another king, the other king will do the same and we have an infinite loop.
		// Therefore we skip calling board.Possib for king, and just call piece.Board(for the king only).

		// no negative side effects come from this, cause if first king can kill second king, and second king goes in first king's range.
		// then that's an easy win!

		// Also make sure to ignore tasty but lethal bait

		// collect enemy's possible points
		// then check if it crosses paths with king's
		sp := Points{}
		for k, s := range b.data {
			if s.Valid() { // always use protection
				if s.Player != pec.Player {
					if s.T != Empty {
						if s.T == King {
							// no infinite loop
							for _, pnt := range s.Possib() {
								sp.Insert(pnt)
							}
						} else {
							pnts, err := b.Possib(k)
							if err != nil {
								for _, pnt := range pnts {
									sp.Insert(pnt)
								}
							}
						}
					}
				}
			}
		}

		for k, v := range ps {
			_, cep, err := b.Get(v)
			// disallow replacing ally piece
			if sp.In(v) || (err == nil && cep.Player == pec.Player) {
				delete(ps, k)
			}
		}

		// a nice preauction, create a copy of board and try king moves.
		// if they land us in a nasty checkmate then discard them
		drb := b.Copy()
		drb.Set(id, Point{-1, -1})
		for k, v := range ps {
			// move could disallow back movement to king's original position
			// so we use set
			drb.Set(id, v)
			if drb.Checkmate(pec.Player) {
				// if that move is checkmattable then discard it
				delete(ps, k)
			}
		}

	default:
		orix, oriy := pec.Pos.X, pec.Pos.Y

		// starting from x, y this function loops through possible points
		// afterwards it changes the value via op function which receives x, y and modifies them
		// in-case it encountered a piece in the way it wait to finish and removes all the following points
		loop := func(x, y int8, op func(int8, int8) (int8, int8)) {
			rm := false
			for i := 0; i < 8; i++ {
				pnt := Point{int8(x), int8(y)}
				if !ps.In(pnt) || !pnt.Valid() {
					break
				}

				if !rm {
					// encountered piece in the way
					_, cep, err := b.Get(pnt)
					if err == nil {
						if pec.Player == cep.Player {
							ps.Delete(cep.Pos)
						}
						rm = true
					} /*else {
						// this direction cannot possibly have following points
						break
					}*/
				} else {
					// start deleting following points, cause we reached a piece in the way
					ps.Delete(pnt)
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

	ps.Clean()
	return ps, nil
}

// FinalCheckmate returns true if the player cannot save themselves. The game ends right after.
// This primarily checks for the Possib moves and if an ally can jump in to save the king.
func (b Board) FinalCheckmate(player uint8) bool {
	// Piece maximum number of moves:
	// - Pawn: 3 killable or 2 at start or just 1
	// - King: 8
	// - Bishop: 13
	// - Rook: 14
	// - Knight: 8
	// - Queen: (Bishop) + Rook = 27
	// - Alltogether: 73
	// This function works by doing every single move(from the checkmatted player) on another board, and checking if was still a checkmate.
	// In-case not - this returns false, otherwise it returns true.
	if !b.Checkmate(player) {
		return false
	}

	exist := false
	ret := true

	for k, s := range b.data {
		if s.Valid() {
			if s.Player == player {
				if s.T == King {
					exist = true
					// continue
					// this prevents the king from defending itself
				}

				if !ret {
					continue
				}

				oldpos := s.Pos
				b.Set(k, Point{-1, -1}) // erase the old piece

				ps, err := b.Possib(k)
				if err != nil {
					panic(err)
				}

				for _, v := range ps {
					s.Pos = v
					b.Set(k, v)
					if !b.Checkmate(player) {
						ret = false
					} else {
						b.Set(k, Point{-1, -1})
					}
				}

				s.Pos = oldpos
				b.Set(k, oldpos)
			}
		}
	}

	if !exist {
		return true
	}

	return ret
}

// Checkmate returns true if the player has been checkmatted
func (b Board) Checkmate(player uint8) bool {
	var king Piece
	if player == 1 {
		king = b.data[4]
	} else if player == 2 {
		king = b.data[28]
	}

	if !king.Valid() {
		// no king, automatically win
		return true
	}

	for k, s := range b.data {
		if s.Valid() {
			if s.Player != player {
				possib := s.Possib()
				if s.T != King { // avoid infinite loop
					var err error
					possib, err = b.Possib(k)

					if err != nil {
						continue
					}
				}

				if possib.In(king.Pos) {
					return true
				}
			}
		}
	}

	return false
}

// Move moves a piece from it's original position to the destination. Returns true if it did, or false if it didn't.
func (b *Board) Move(id int, dst Point) (ret bool) {
	pec, err := b.GetByIndex(id)
	if err != nil {
		ret = false
		return
	}

	defer func() {
		src := pec.Pos

		if ret {
			b.data[id].Pos = dst

			// if there's a piece there
			// then delete it
			_, _, err := b.Get(src)
			if err != nil {
				b.Set(id, Point{-1, -1})
			}
		}

		for _, v := range b.ml {
			v(pec, src, dst, ret)
		}
	}()

	if err == nil && pec.Valid() {
		// can we legally go there, i.e is it in the possible combinations??
		// so for example bishop cannot go horizontally
		if pec.CanGo(dst) {
			di, cep, err := b.Get(dst)
			// is there a piece in the destination??
			if cep.Valid() && err == nil {
				// is the piece's an enemy
				if pec.Player != cep.Player {
					// is it not a pawn(cause pawns cannot enemy forward or backward of them)
					if pec.T != PawnF && pec.T != PawnB {
						ps, err := b.Possib(di)
						if err != nil {
							ret = false
						} else {
							ret = ps.In(dst)
						}
					}
				}
			} else {
				// no piece in the destination
				ps, err := b.Possib(id)
				if err != nil {
					ret = false
				} else {
					ret = ps.In(dst)
				}
			}
		} else {
			if pec.T == PawnF || pec.T == PawnB {
				ps, err := b.Possib(id)
				if err != nil {
					ret = false
				} else {
					ret = ps.In(dst)
				}
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

	return nil
}
