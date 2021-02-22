package board

type Point struct {
	X int
	Y int
}

type Points []Point

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
