package board

import "fmt"

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
	X int `json:"x"`
	Y int `json:"y"`
}

func (p Point) String() string {
	return fmt.Sprintf("%d:%d", p.X, p.Y)
}

type Points []Point

func (ps Points) Len() int {
	return len(ps)
}

func (ps Points) Less(i, j int) bool {
	p, o := ps[i], ps[j]
	if p.X == o.X {
		return p.Y < o.Y
	}

	return p.X < o.X
}

func (ps Points) Swap(i, j int) {
	ps[i], ps[j] = ps[j], ps[i]
}

func abs(i int) int {
	if i < 0 {
		return i * -1
	}

	return i
}

// Clean removes all out of bounds points
func (ps Points) Clean() (ret Points) {
	ret = ps
	// remove invalid poitns
	for i := len(ps) - 1; i >= 0; i-- {
		p := ret[i]
		if !p.Valid() {
			ret = append(ret[:i], ret[i+1:]...)
		}
	}

	exist := map[string]struct{}{}
	for i := len(ret) - 1; i >= 0; i-- {
		p := ret[i]

		_, ok := exist[p.String()]
		if ok {
			ret = append(ret[:i], ret[i+1:]...)
		} else {
			exist[p.String()] = struct{}{}
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
	return ps.Index(dst) != -1
}

func (ps Points) Index(dst Point) int {
	for k, v := range ps {
		if v.Equal(dst) {
			return k
		}
	}

	return -1
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
	x, y := 0, 0
	res := p.X - p.Y
	if res > 0 {
		x = res
	} else {
		y = abs(res)
	}

	//orix, oriy := x, y

	ret = append(ret, Point{x, y})
	for i := 0; i < 8; i++ {
		x++
		y++

		o := Point{x, y}

		if p.Equal(o) {
			continue
		}
		if !o.Valid() {
			break
		}

		ret = append(ret, o)
	}

	// this part took me a bit to figure it out
	x, y = 0, 7
	res = p.X + p.Y
	if res < 7 {
		y = res
	} else if res > 7 {
		x = x + (res - 7)
	}

	ret = append(ret, Point{x, y})

	for i := 0; i < 8; i++ {
		x++
		y--

		o := Point{x, y}

		if p.Equal(o) {
			continue
		}
		if !o.Valid() {
			break
		}

		ret = append(ret, o)
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

func (p Point) Increase(dir uint8) Point {
	x, y := p.X, p.Y
	if Has(dir, DirUp) {
		x--
	} else if Has(dir, DirDown) {
		x++
	}
	if Has(dir, DirLeft) {
		y--
	} else if Has(dir, DirRight) {
		y++
	}

	return Point{x, y}
}

// The following is a collection of generic functions, that start from x,y and return a new point from that perspective.
// Also the use of x, y values(instead of Point) makes these more comprehensible
func UpLeft(x, y int) (int, int)    { return x - 1, y - 1 }
func UpRight(x, y int) (int, int)   { return x - 1, y + 1 }
func DownLeft(x, y int) (int, int)  { return x + 1, y - 1 }
func DownRight(x, y int) (int, int) { return x + 1, y + 1 }

func Up(x, y int) (int, int)    { return x - 1, y }
func Down(x, y int) (int, int)  { return x + 1, y }
func Left(x, y int) (int, int)  { return x, y - 1 }
func Right(x, y int) (int, int) { return x, y + 1 }

// Smaller returns true if dst is smaller than src. Smaller compares x to x, and then y to y.
/*
func (p Point) Smaller(dst Point) bool {
	if p.Equal(dst) {
		return false
	}

	x, y := p.X, p.Y
	if x > dst.X {
		return true
	}

	return y > dst.Y
}

// Bigger returns true if dst is bigger than src
func (p Point) Bigger(dst Point) bool {
	if p.Equal(dst) {
		return false
	}

	return !p.Smaller(dst)
}
*/
