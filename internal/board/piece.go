package board

import (
	"encoding/json"
)

const (
	Empty uint8 = iota
	// Pawn Forward -> 1, 0 - 2, 0
	PawnF
	// Pawn Backward -> 6, 0 -> 5, 0
	PawnB
	// Bishop
	// Moves diagonally
	Bishop
	// Knight
	// Moves with [2, 1] or [1, 2]
	Knight
	// Rook
	// Moves horizontally, vertically.
	Rook
	// Queen
	// Moves horizontally, vertically, diagonally, and within square.
	Queen
	// King
	// Moves within square.
	King
)

type Piece struct {
	// Player could be any number, but mostly [1, 2]
	Player uint8 `json:"player"`
	// T the piece type
	T uint8 `json:"type"`
	// pos
	pos Point
}

func (p *Piece) Valid() bool {
	if p.T >= PawnF && p.T <= King {
		return true
	}

	return false
}

// ShortString produces a one-character string to represent the piece. Used for debugging.
func (p *Piece) ShortString() string {
	strings := map[uint8]string{
		Empty:  " ",
		PawnF:  "P",
		PawnB:  "P",
		Bishop: "B",
		Knight: "N",
		Rook:   "R",
		Queen:  "Q",
		King:   "K",
	}

	return strings[p.T]
}

// String representation of the type
func (p *Piece) String() string {
	strings := map[uint8]string{
		Empty:  "Empty",
		PawnF:  "Pawn",
		PawnB:  "Pawn",
		Bishop: "Bishop",
		Knight: "Knight",
		Rook:   "Rook",
		Queen:  "Queen",
		King:   "King",
	}

	return strings[p.T]
}

// CanGo does validation for the piece. Each piece has it's own rules.
func (p *Piece) CanGo(dst Point) bool {
	if dst.Equal(p.pos) || !dst.Valid() {
		return false
	}

	src := p.pos
	// i.e starting point
	switch p.T {
	// Only horizontally, can't move back
	// 2 points at start, 1 point after that
	case PawnF, PawnB:
		ps := Points{}
		if p.T == PawnF {
			ps = src.Forward()
		} else {
			ps = src.Backward()
		}
		if src.X == 1 || src.X == 6 {
			ps = append(ps, Point{X: src.X - 2, Y: src.Y})
			ps = append(ps, Point{X: src.X + 2, Y: src.Y})

			ps.Clean()
		}

		return ps.In(dst)
	// Only diagonal
	case Bishop:
		return src.Diagonal().In(dst)
	// Move within [2, 1] or [1, 2]
	case Knight:
		return src.Knight().In(dst)
	// horizontal or vertical
	case Rook:
		return src.Horizontal().
			Merge(src.Vertical()).
			In(dst)
	// move within square or diagonal or horizontal or vertical
	case Queen:
		return src.Horizontal().
			Merge(src.Vertical()).
			Merge(src.Square()).
			Merge(src.Diagonal()).
			In(dst)
	// move within square
	case King:
		return src.Square().In(dst)
	}

	return false
}

// MarshalJSON json.Marshaler
func (p *Piece) MarshalJSON() ([]byte, error) {
	x := struct {
		P uint8 `json:"player"`
		T uint8 `json:"type"`
	}{p.Player, p.T}

	body, err := json.Marshal(x)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// UnmarshalJSON json.Unmarshaler
func (p *Piece) UnmarshalJSON(b []byte) error {
	x := struct {
		P uint8 `json:"player"`
		T uint8 `json:"type"`
	}{p.Player, p.T}

	err := json.Unmarshal(b, x)
	if err != nil {
		return err
	}

	return nil
}
