package board

const (
	DirUp uint8 = 1 << iota
	DirDown
	DirLeft
	DirRight
)

func Set(b, flag uint8) uint8    { return b | flag }
func Clear(b, flag uint8) uint8  { return b &^ flag }
func Toggle(b, flag uint8) uint8 { return b ^ flag }
func Has(b, flag uint8) bool     { return b&flag != 0 }

type Point struct {
	X int
	Y int
}

type Points []Point

func abs(i int) int {
	if i < 0 {
		return i * -1
	}

	return i
}

// Clean removes all out of bounds points
func (ps Points) Clean() (ret Points) {
	ret = ps
	for i := len(ps) - 1; i >= 0; i-- {
		p := ret[i]
		if !p.Valid() {
			ret = append(ret[:i], ret[i+1:]...)
		}
	}

	return
}

// Merge merges ps with all
func (ps Points) Merge(all ...Points) (ret Points) {
	ret = append(ret, ps...)
	for _, v := range all {
		ret = append(ret, v...)
	}

	return ret
}

// In checks if dst is in ps
func (ps Points) In(dst Point) bool {
	for _, v := range ps {
		if v.Equal(dst) {
			return true
		}
	}

	return false
}

// Outside generates points that are outside of ps
func (ps Points) Outside() (ret Points) {
	for x := 0; x < 8; x++ {
		for y := 0; y < 8; y++ {
			p := Point{x, y}
			if !ps.In(p) {
				ret = append(ret, p)
			}
		}
	}

	return ret
}

// Equal asserts if p is equal to o
func (p Point) Equal(o Point) bool {
	return p.X == o.X && p.Y == o.Y
}

// Valid return false when out of bounds
func (p Point) Valid() bool {
	return !(p.X > 7 || p.Y > 7 || p.X < 0 || p.Y < 0)
}

/*
// IsDiagonalUpRight returns true if dst is UpRight of p
func (p Point) IsDiagonalUpRight(dst Point) bool {
	x := dst.X - p.X
	y := dst.Y - p.Y
	if x > 0 && y > 0 && x == y {
		return true
	}

	return false
}

// IsDiagonalUpLeft returns true if dst is UpLeft of p
func (p Point) IsDiagonalUpLeft(dst Point) bool {
	x := dst.X - p.X
	y := dst.Y - p.Y

	if abs(y) == x && y < 0 && x > 0 {
		return true
	}

	return false
}

// IsDiagonalUpRight returns true if dst is DownRight of p
func (p Point) IsDiagonalDownRight(dst Point) bool {
	x := dst.X - p.X
	y := dst.Y - p.Y
	if abs(x) == y && x < 0 && y > 0 {
		return true
	}

	return false
}

// IsDiagonalDownLeft returns true if dst is DownLeft of p
func (p Point) IsDiagonalDownLeft(dst Point) bool {
	x := dst.X - p.X
	y := dst.Y - p.Y

	if y < 0 && x < 0 && x == y {
		return true
	}

	return false
}

// IsLeft returns if the dst is left of the point
func (p Point) IsLeft(dst Point) bool {
	if dst.X != p.X || dst.Y == p.Y {
		return false
	}

	return dst.Y-p.Y < 0
}

// IsRight returns if the dst is right of the point
func (p Point) IsRight(dst Point) bool {
	if dst.X != p.X || dst.Y == p.Y {
		return false
	}

	return dst.Y-p.Y > 0
}

// IsUp returns if the dst is up of the point
func (p Point) IsUp(dst Point) bool {
	if dst.X == p.X || dst.Y != p.Y {
		return false
	}

	return dst.X-p.X < 0
}

// IsDown returns if the dst is downwards of the point
func (p Point) IsDown(dst Point) bool {
	if dst.X == p.X || dst.Y != p.Y {
		return false
	}

	return dst.X-p.X > 0
}
*/

func (p Point) Direction(dst Point) (d uint8) {
	x := dst.X - p.X
	if x < 0 {
		d = Set(d, DirUp)
	} else if x > 0 {
		d = Set(d, DirDown)
	}

	y := dst.Y - p.Y
	if y < 0 {
		d = Set(d, DirLeft)
	} else if y > 0 {
		d = Set(d, DirRight)
	}

	return
}

// Diagonal generates diagonal points
func (p Point) Diagonal() (ret Points) {
	x := 7 - p.X
	y := 7 - p.Y

	diff := 0
	if x > y {
		diff = x
	} else {
		diff = y
	}

	x = p.X + diff
	y = p.Y + diff
	ret = append(ret, Point{x, y})

	for i := 0; i < 8; i++ {
		x--
		y--

		if p.X == x && p.Y == y {
			continue
		}

		p := Point{x, y}
		if !p.Valid() {
			break
		}

		ret = append(ret, p)
	}

	y = p.Y + diff
	x = p.X - diff

	ret = append(ret, Point{x, y})
	for i := 0; i < 8; i++ {
		x++
		y--

		if p.X == x && p.Y == y {
			continue
		}

		p := Point{x, y}
		if !p.Valid() {
			break
		}
		ret = append(ret, p)
	}

	return ret.Clean()
}

// Horizontal generates horizontal points
func (p Point) Horizontal() (ret Points) {
	for i := 0; i < 8; i++ {
		if p.Y == i {
			continue
		}

		ret = append(ret, Point{p.X, i})
	}

	return ret
}

// Vertical generates vertical points
func (p Point) Vertical() (ret Points) {
	for i := 0; i < 8; i++ {
		if p.X == i {
			continue
		}

		ret = append(ret, Point{i, p.Y})
	}

	return ret
}

// Square generates square points
func (p Point) Square() Points {
	ps := Points{
		{p.X + 1, p.Y + 1},
		{p.X + 1, p.Y},
		{p.X + 1, p.Y - 1},
		{p.X, p.Y + 1},
		{p.X, p.Y - 1},
		{p.X - 1, p.Y + 1},
		{p.X - 1, p.Y},
		{p.X - 1, p.Y - 1},
	}.Clean()

	return ps
}

// Knight generates [2, 1] and [1, 2] points
func (p Point) Knight() Points {
	ps := Points{
		{p.X + 2, p.Y + 1},
		{p.X + 2, p.Y - 1},
		{p.X - 2, p.Y + 1},
		{p.X - 2, p.Y - 1},

		{p.X + 1, p.Y + 2},
		{p.X - 1, p.Y + 2},
		{p.X + 1, p.Y - 2},
		{p.X - 1, p.Y - 2},
	}.Clean()

	return ps
}

// Forward generates a point forward. Forward being up -1
func (p Point) Forward() Points {
	ps := Points{
		{p.X - 1, p.Y},
	}.Clean()

	return ps
}

// Backward generates a point backward. Backward being down +1
func (p Point) Backward() Points {
	ps := Points{
		{p.X + 1, p.Y},
	}.Clean()

	return ps
}

// Left generates a point to the left.
func (p Point) Left() Points {
	ps := Points{
		{p.X, p.Y - 1},
	}.Clean()

	return ps
}

// Right generates a point to the right.
func (p Point) Right() Points {
	ps := Points{
		{p.X, p.Y + 1},
	}.Clean()

	return ps
}

// Corner generates [+1, +1], [+1, -1], [-1, +1] and [-1, -1].
func (p Point) Corner() Points {
	ps := Points{
		{p.X + 1, p.Y + 1},
		{p.X + 1, p.Y - 1},
		{p.X - 1, p.Y + 1},
		{p.X - 1, p.Y - 1},
	}.Clean()

	return ps
}
