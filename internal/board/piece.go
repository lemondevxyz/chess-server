package board

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
	// Moves horizontally, vertically and diagonally
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
	Player uint8
	// T the piece type
	T uint8
	// X place in array
	X int
	// Y place in array
	Y int
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
		Empty:  "",
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
func (p *Piece) CanGo(x, y int) bool {

	limit := 8
	if x >= limit && y >= limit {
		if x < 0 && y < 0 {
			// out of bounds
			return false
		}
	}

	if p.X == x && p.Y == y {
		return false
	}

	src, dst := Point{
		X: p.X,
		Y: p.Y,
	}, Point{
		X: x,
		Y: y,
	}
	// i.e starting point
	switch p.T {
	// Only horizontally, can't move back
	// 2 points at start, 1 point after that
	case PawnF, PawnB:
		area := Point{
			X: 1,
			Y: 0,
		}

		ok := false
		if p.T == PawnF {
			ok = Forward(src, dst)
		} else if p.T == PawnB {
			ok = Backward(src, dst)
		}

		if p.X == 1 || p.X == 6 {
			// can move two units
			return ok && (Within(Point{X: 2, Y: 0}, src, dst) || Within(area, src, dst))
		} else {
			return ok && Within(area, src, dst)
		}
	// Only diagonal
	case Bishop:
		return Diagonal(src, dst)
	// Move within [2, 1] or [1, 2]
	case Knight:
		area := Point{X: 2, Y: 1}
		return Within(area, src, dst) || Within(Swap(area), src, dst)
	// horizontal or vertical
	case Rook:
		return Horizontal(src, dst) || Vertical(src, dst)
	// move within square or diagonal or horizontal or vertical
	case Queen:
		return Vertical(src, dst) || Horizontal(src, dst) || Diagonal(src, dst) || Square(src, dst)
	// move within square
	case King:
		return Square(src, dst)
	}

	return false
}
