package board

import (
	"fmt"
	"sort"
)

type Point struct {
	X int8 `json:"x"`
	Y int8 `json:"y"`
}

func (p Point) String() string {
	return fmt.Sprintf("%d:%d", p.X, p.Y)
}

type Points map[string]Point

func abs(i int8) int8 {
	if i < 0 {
		return i * -1
	}

	return i
}

func (ps Points) String() (str string) {
	sli := []string{}
	for k := range ps {
		sli = append(sli, k)
	}

	sort.Strings(sli)
	for _, v := range sli {
		str += v + ", "
	}

	if len(str) >= 2 {
		return "[" + str[:len(str)-2] + "]"
	}

	return "[]"
}

// Clean removes all out of bounds points, and duplicate poitns
func (ps Points) Clean() {
	// remove invalid points
	for k, pnt := range ps {
		if !pnt.Valid() {
			delete(ps, k)
		}
	}
}

// Merge merges ps with all
func (ps Points) Merge(all ...Points) (ret Points) {
	ret = Points{}
	for k, s := range ps {
		ret[k] = s
	}

	for _, v := range all {
		for k, s := range v {
			ret[k] = s
		}
	}

	return ret
}

// In checks if dst is in ps
func (ps Points) In(dst Point) bool {
	_, ok := ps[dst.String()]
	return ok
}

// Delete deletes the pnt
func (ps Points) Delete(dst Point) {
	delete(ps, dst.String())
}

// Outside generates points that are outside of ps
func (ps Points) Outside() Points {
	ret := Points{}
	for x := 0; x < 8; x++ {
		for y := 0; y < 8; y++ {
			p := Point{int8(x), int8(y)}
			if !ps.In(p) {
				ret[p.String()] = p
			}
		}
	}

	return ret
}

// Insert adds a point to Points
func (ps Points) Insert(sp ...Point) {
	for _, v := range sp {
		ps[v.String()] = v
	}
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
func (p Point) Diagonal() Points {
	var x, y int8 = 0, 0
	res := p.X - p.Y
	if res > 0 {
		x = res
	} else {
		y = abs(res)
	}

	//orix, oriy := x, y

	ret := Points{}
	ret.Insert(Point{int8(x), int8(y)})
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

		ret.Insert(o)
	}

	// this part took me a bit to figure it out
	x, y = 0, 7
	res = p.X + p.Y
	if res < 7 {
		y = res
	} else if res > 7 {
		x = x + (res - 7)
	}

	ret.Insert(Point{x, y})

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

		ret.Insert(o)
	}

	ret.Clean()
	return ret
}

// Horizontal generates horizontal points
func (p Point) Horizontal() Points {
	ret := Points{}
	for i := int8(0); i < 8; i++ {
		if p.X == i {
			continue
		}

		ret.Insert(Point{i, p.Y})

	}

	return ret
}

// Vertical generates vertical points
func (p Point) Vertical() Points {
	ret := Points{}
	for i := int8(0); i < 8; i++ {
		if p.Y == i {
			continue
		}

		ret.Insert(Point{p.X, i})
	}

	return ret
}

// Square generates square points
func (p Point) Square() Points {
	ps := Points{}
	ps.Insert(
		Point{p.X + 1, p.Y + 1},
		Point{p.X + 1, p.Y},
		Point{p.X + 1, p.Y - 1},
		Point{p.X, p.Y + 1},
		Point{p.X, p.Y - 1},
		Point{p.X - 1, p.Y + 1},
		Point{p.X - 1, p.Y},
		Point{p.X - 1, p.Y - 1},
	)
	ps.Clean()

	return ps
}

// Knight generates [2, 1] and [1, 2] points
func (p Point) Knight() Points {
	ps := Points{}
	ps.Insert(
		Point{p.X + 2, p.Y + 1},
		Point{p.X + 2, p.Y - 1},
		Point{p.X - 2, p.Y + 1},
		Point{p.X - 2, p.Y - 1},

		Point{p.X + 1, p.Y + 2},
		Point{p.X - 1, p.Y + 2},
		Point{p.X + 1, p.Y - 2},
		Point{p.X - 1, p.Y - 2},
	)
	ps.Clean()

	return ps
}

// Forward generates a point forward. Forward being up -1
func (p Point) Forward() Points {
	ps := Points{}
	ps.Insert(Point{p.X, p.Y - 1})
	ps.Clean()

	return ps
}

// Backward generates a point backward. Backward being down +1
func (p Point) Backward() Points {
	ps := Points{}
	ps.Insert(Point{p.X, p.Y + 1})
	ps.Clean()

	return ps
}

// Corner generates [+1, +1], [+1, -1], [-1, +1] and [-1, -1].
func (p Point) Corner() Points {
	ps := Points{}

	ps.Insert(
		Point{p.X + 1, p.Y + 1},
		Point{p.X + 1, p.Y - 1},
		Point{p.X - 1, p.Y + 1},
		Point{p.X - 1, p.Y - 1},
	)

	ps.Clean()

	return ps
}

// The following is a collection of modifier functions, that take in x,y and modify them to go in certain directions...
// Also the use of x, y values(instead of Point) makes these more comprehensible
func UpLeft(x, y int8) (int8, int8)    { return x - 1, y - 1 }
func UpRight(x, y int8) (int8, int8)   { return x + 1, y - 1 }
func DownLeft(x, y int8) (int8, int8)  { return x - 1, y + 1 }
func DownRight(x, y int8) (int8, int8) { return x + 1, y + 1 }

func Up(x, y int8) (int8, int8)    { return x, y - 1 }
func Down(x, y int8) (int8, int8)  { return x, y + 1 }
func Left(x, y int8) (int8, int8)  { return x - 1, y }
func Right(x, y int8) (int8, int8) { return x + 1, y }
