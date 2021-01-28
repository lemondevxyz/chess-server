package board

type Board [8][8]*Piece

// NewBoard creates a new board with the default placement.
func NewBoard() *Board {
	b := Board{}
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
			b[x][y] = &Piece{
				T:      v,
				Player: 1,
				X:      x,
				Y:      y,
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
			b[x][y] = &Piece{
				T:      v,
				Player: 2,
				X:      x,
				Y:      y,
			}
		}
	}

	return &b
}

// String method returns a string. makes it easier to debug
func (b *Board) String() (str string) {
	for k, s := range b {
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

func (b *Board) Move(p *Piece, dst Point) bool {

	if p.CanGo(dst.X, dst.Y) {
		if p != nil {
			b[p.X][p.Y] = nil

			p.X = dst.X
			p.Y = dst.Y

			b[dst.X][dst.Y] = p
			return true
		}
	}

	return false
}
