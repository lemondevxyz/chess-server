package board

// absolute value of i
func abs(i int) int {
	if i < 0 {
		return i * -1
	}

	return i
}

const (
	Empty uint8 = iota
	PawnF       // move forward, start from 1
	PawnB       // move backward, start from 6
	Bishop
	Knight
	Rook
	Queen
	King
)

type Piece struct {
	Player int
	T      uint8
	X      int
	Y      int
}

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
