// Package board provides game-logic for chess, without the need of interaction from the user.

// brd == An instance of board
// drb == A second instance of board
// OR(if there are two points)
// src == An instance of a point
// dst == A second instance of point
// id == An ID of a piece
// di == A second id of a piece
// pec == An instance of piece
// cep == A second of instance piece
package board

import (
	"encoding/json"
)

// MoveEvent is a function called post-movement of a piece.
type MoveEvent func(id int, pec Piece, src Point, dst Point)

type Board struct {
	// data [8][8]*Piece
	data [32]Piece
	// move event listener
	ml []MoveEvent
}

// NewBoard creates a new board with the default placement.
func NewBoard() *Board {
	brd := Board{
		ml:   []MoveEvent{},
		data: [32]Piece{},
	}

	row := [32]uint8{
		// 0 -> 7
		// 0   | 1    | 2    | 3    | 4    | 5    | 6    | 7
		// 0,0 | 1, 0 | 2, 0 | 3, 0 | 4, 0 | 5, 0 | 6, 0 | 7, 0
		Rook, Knight, Bishop, Queen, King, Bishop, Knight, Rook,
		// 8 -> 15
		// 8   | 9    | 10   | 11   | 12   | 13   | 14   | 15
		// 0,1 | 1, 1 | 2, 1 | 3, 1 | 4, 1 | 5, 1 | 6, 1 | 7, 1
		Pawn, Pawn, Pawn, Pawn, Pawn, Pawn, Pawn, Pawn,
		// 16 -> 23
		// 16  | 17   | 18   | 19   | 20   | 21   | 22   | 23
		// 0,6 | 1, 6 | 2, 6 | 3, 6 | 4, 6 | 5, 6 | 6, 6 | 7, 6
		Pawn, Pawn, Pawn, Pawn, Pawn, Pawn, Pawn, Pawn,
		// 24 -> 31
		// 24  | 25   | 26   | 27   | 28   | 29   | 30   | 31
		// 0,7 | 1, 7 | 2, 7 | 3, 7 | 4, 7 | 5, 7 | 6, 7 | 7, 7
		Rook, Knight, Bishop, Queen, King, Bishop, Knight, Rook,
	}

	for k, s := range row {
		p1 := false
		if k >= 16 {
			p1 = true
		}

		x := int8(k % 8)
		y := int8(k / 8)

		if k >= 16 {
			y += 4
		}

		brd.data[k] = Piece{
			Kind: s,
			P1:   p1,
			Pos:  Point{x, y},
		}
	}

	return &brd
}

func (brd Board) Copy() *Board {
	drb := Board{
		ml:   []MoveEvent{},
		data: [32]Piece{},
	}

	for k, s := range brd.data {
		drb.data[k] = s
	}

	return &drb
}

// String method returns a string. makes it easier to debug
func (brd Board) String() (str string) {
	for i := 0; i < (8 * 8); i++ {
		x := int8(i % 8)
		y := int8(i / 8)
		if i != 0 && x == 0 {
			str += "\n"
		}

		char := " "

		pos := Point{x, y}
		for _, pec := range brd.data {
			if pec.Pos.Equal(pos) {
				char = pec.ShortString()
			}
		}

		str += char + " "
	}

	return
}

// Listen returns adds a callback that gets called pre and post movement.
func (brd *Board) Listen(callback MoveEvent) {
	brd.ml = append(brd.ml, callback)
}

// Set sets a piece's position in the board without game-logic interfering.
func (brd *Board) Set(id int, pos Point) error {
	if !IsIDValid(int8(id)) {
		return ErrInvalidID
	}

	if !pos.Valid() {
		brd.data[id].Pos = Point{-1, -1}
	} else {
		// also if there's a piece in that position, erase it!
		di, _, err := brd.Get(pos)
		if err == nil {
			brd.data[di].Pos = Point{-1, -1}
		}

		brd.data[id].Pos = pos
	}

	return nil
}

// SetKind sets a piece's kind
func (brd *Board) SetKind(id int, kind uint8) error {
	if !IsIDValid(int8(id)) {
		return ErrInvalidID
	}

	brd.data[id].Kind = kind

	return nil
}

// Get returns a piece and it's index. Or otherwise -1, an empty piece and an error.
func (brd Board) Get(src Point) (int, Piece, error) {
	if !src.Valid() {
		return -1, Piece{}, ErrInvalidPoint
	}

	for k, v := range brd.data {
		if v.Pos.Equal(src) {
			if !v.Valid() {
				return k, v, ErrInvalidPoint
			}

			return k, v, nil
		}
	}

	return -1, Piece{}, ErrEmptyPiece
}

func (brd Board) GetByIndex(i int) (Piece, error) {
	if i >= len(brd.data) {
		return Piece{}, ErrInvalidPoint
	}

	return brd.data[i], nil
}

// Possib is the same as Piece.Possib, but with removal of illegal moves.
func (brd Board) Possib(id int) (Points, error) {
	pec, err := brd.GetByIndex(id)
	if err != nil {
		return nil, ErrEmptyPiece
	}

	ps := pec.Possib()
	switch pec.Kind {
	case Knight: // disallow movement to allies
		for k, v := range ps {
			_, cep, err := brd.Get(v)
			if err == nil {
				if cep.P1 == pec.P1 {
					delete(ps, k)
				}
			}
		}
	case Pawn: // disallow movement in front of piece
		for k, pnt := range ps {
			_, _, err = brd.Get(pnt)
			if err == nil { // piece in the way ...
				delete(ps, k)
			}
		}

		// if our move has a piece in the way then cancel
		// also if we're at 6 or 1, then allow movement to 4 or 3
		y := pec.Pos.Y + 1
		if pec.P1 {
			y = pec.Pos.Y - 1
		}

		sp := Points{}
		sp.Insert(Point{pec.Pos.X + 1, y})
		sp.Insert(Point{pec.Pos.X - 1, y})

		sp.Clean()

		for k, v := range sp {
			_, cep, err := brd.Get(v)
			// is there a piece
			if err == nil {
				// is it the enemy's
				if cep.P1 != pec.P1 {
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
		for _, id := range GetRange(GetInversePlayer(pec.P1)) {
			cep := brd.data[id]
			if cep.Valid() { // always use protection
				if cep.Kind == King {
					// no infinite loop
					for _, pnt := range cep.Possib() {
						sp.Insert(pnt)
					}
				} else {
					pnts, err := brd.Possib(id)
					if err != nil {
						for _, pnt := range pnts {
							sp.Insert(pnt)
						}
					}
				}
			}
		}

		for k, v := range ps {
			_, cep, err := brd.Get(v)
			// disallow replacing ally piece
			if sp.In(v) || (err == nil && cep.P1 == pec.P1) {
				delete(ps, k)
			}
		}

		// a nice preauction, create a copy of board and try king moves.
		// if they land us in a nasty checkmate then discard them
		drb := brd.Copy()
		drb.Set(id, Point{-1, -1})
		for k, v := range ps {
			// move could disallow back movement to king's original position
			// so we use set
			drb.Set(id, v)
			if drb.Checkmate(pec.P1) {
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

					_, cep, err := brd.Get(pnt)
					if err == nil {
						if pec.P1 == cep.P1 {
							// it's our piece, erase it as a point we can through
							ps.Delete(cep.Pos)
						}
						// set to remove all following points from now on
						rm = true
					}
				} else {
					// start deleting following points, cause we reached a piece in the way
					ps.Delete(pnt)
				}
				x, y = op(x, y)
			}
		}

		x, y := orix, oriy
		// normal direction
		// these are basically modifier functions, called in loop. So for example, say we start at 0, 0 and we call up.
		// the loop will call up, so the points will increase incremently:
		// [{1,0},{2,0},{3,0}....] and so on
		// until it reaches and out of bounds point, then it stops....
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
func (brd Board) FinalCheckmate(p1 bool) bool {
	// Piece maximum number of moves:
	// - Pawn: 3 killable or 2 at start or just 1
	// - King: 8
	// - Bishop: 13
	// - Rook: 14
	// - Knight: 8
	// - Queen: (Bishop) + Rook = 27
	// - Alltogether: 73
	// This function works by doing every single move(from the checkmatted player) on another board, and checking if was still a checkmate.
	if !brd.Checkmate(p1) {
		return false
	}

	// validate player number, and get the king's id
	kingid := GetKing(p1)
	if kingid == -1 {
		return true
	}

	// set the king from the id
	king := brd.data[kingid]
	if !king.Pos.Valid() {
		// if he's dead, then it's a final checkmate
		return true
	}

	final := true
	// GetRange returns the id range for the player
	// so we loop over all piece allies
	for _, id := range GetRange(p1) {
		pec := brd.data[id]
		// if it's dead then skip it
		if pec.Valid() {
			// well if we already established that it's not final
			// then stop
			if !final {
				break
			}

			// get all possible moves for that specific piece
			ps, _ := brd.Possib(id)
			oldpos := pec.Pos
			for _, pnt := range ps {
				// s.Pos = v
				// using move instead of Set is because move has extra game-logic, that Set could break.
				// like for example pawns
				if brd.canMove(id, pnt) {
					brd.Set(id, pnt)
					if !brd.Checkmate(p1) {
						return false
					}
				}
			}

			brd.Set(id, oldpos)
		}
	}

	return final
}

// Checkmate returns true if the player has been checkmatted.
// While FinalCheckmate checks if the checkmatted player can do anything about it, this basically check if there is a checkmate.
// It does not care whether the checkmatted player can escape
func (brd Board) Checkmate(p1 bool) bool {
	var king Piece
	id := GetKing(p1)
	if id != -1 {
		king = brd.data[id]
	} else {
		// invalid player number
		return true
	}

	if !king.Valid() {
		// no king, automatically win
		return true
	}

	enemyids := GetInversePlayer(p1)
	for _, id := range GetRange(enemyids) {
		pec := brd.data[id]
		if pec.Valid() {
			ps := pec.Possib()
			if pec.Kind != King { // avoid infinite loop
				var err error
				ps, err = brd.Possib(id)

				if err != nil {
					continue
				}
			}

			if ps.In(king.Pos) {
				return true
			}
		}
	}

	return false
}

// i plan to replace this function by the Possib function...
func (brd *Board) canMove(id int, dst Point) (valid bool) {
	pec, err := brd.GetByIndex(id)
	if err != nil || !pec.Valid() || dst.Equal(pec.Pos) {
		valid = false
		return
	}

	// can we legally go there, i.e is it in the possible combinations??
	// so for example bishop cannot go horizontally
	/*
		_, cep, err := brd.Get(dst)
		// is there a piece in the destination??
		if cep.Valid() && err == nil {
			// is the piece's an enemy
			if pec.P1 != cep.P1 {
				// is it not a pawn(cause pawns cannot enemy forward or backward of them)
				if pec.Kind != Pawn {
					ps, err := brd.Possib(id)
					if err != nil {
						valid = false
					} else {
						valid = ps.In(dst)
					}
				}
			}
		} else {
			// no piece in the destination
			ps, err := brd.Possib(id)
			if err != nil {
				valid = false
			} else {
				valid = ps.In(dst)
			}
		}
	*/

	ps, _ := brd.Possib(id)
	return ps.In(dst)
}

// Move moves a piece from it's original position to the destination. Returns true if it did, or false if it didn't.
func (brd *Board) Move(id int, dst Point) bool {
	pec, _ := brd.GetByIndex(id)
	valid := brd.canMove(id, dst)

	if valid {
		src := pec.Pos

		// if there's a piece there
		// then delete it
		di, _, err := brd.Get(dst)
		if err == nil {
			brd.Set(di, Point{-1, -1})
		}

		brd.data[id].Pos = dst

		for _, v := range brd.ml {
			v(id, pec, src, dst)
		}
	}

	return valid
}

func (brd Board) MarshalJSON() ([]byte, error) {
	body, err := json.Marshal(brd.data)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (brd *Board) UnmarshalJSON(body []byte) error {
	brd.ml = []MoveEvent{}
	if err := json.Unmarshal(body, &brd.data); err != nil {
		return err
	}

	return nil
}
