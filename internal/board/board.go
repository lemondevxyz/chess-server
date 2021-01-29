// Package board provides game-logic for chess, without the need of interaction from the user.
package board

type MoveEvent func(p *Piece, dst Point, pre bool)

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
			b.data[x][y] = &Piece{
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
		b.data[p.X][p.Y] = p
	}
}

// Move moves a piece from it's original position to the destination. Returns true if it did, or false if it didn't.
func (b *Board) Move(p *Piece, dst Point) (ret bool) {
	defer func() {
		for _, v := range b.ml {
			v(p, dst, false)
		}

		if ret {
			b.data[p.X][p.Y] = nil

			p.X = dst.X
			p.Y = dst.Y

			b.data[dst.X][dst.Y] = p
		}
	}()

	for _, v := range b.ml {
		v(p, dst, true)
	}

	if p != nil {
		x := b.data[dst.X][dst.Y]
		if p.CanGo(dst.X, dst.Y) {
			if x != nil {
				if x.T != PawnB && x.T != PawnF {
					if p.T != PawnB && p.T != PawnF {
						if p.Player != x.Player {
							ret = true
						}
					}
				}
			} else {
				ret = true
			}
		} else {
			if p.T == PawnB || p.T == PawnF {
				x := p.X
				y := p.Y
				if p.T == PawnF {
					x--
				} else if p.T == PawnB {
					x++
				}

				if dst.X == x {
					oldy := y
					// other piece
					o := b.data[x][y+1]
					i := b.data[x][y-1]
					if o != nil && o.T != Empty && o.Player != p.Player {
						y = y + 1
					} else if i != nil && i.T != Empty && i.Player != p.Player {
						y = y - 1
					}

					if oldy != y {
						ret = true
					}
				}
			}
		}
	}

	return
}
