package board

import (
	"sort"
	"testing"
)

type pointSlice []Point

func (ps pointSlice) Len() int {
	return len(ps)
}

func (ps pointSlice) Less(i, j int) bool {
	p, o := ps[i], ps[j]
	if p.X == o.X {
		return p.Y < o.Y
	}

	return p.X < o.X
}

func (ps pointSlice) Swap(i, j int) {
	ps[i], ps[j] = ps[j], ps[i]
}

func (ps pointSlice) Merge(sp ...pointSlice) pointSlice {
	ret := append(pointSlice{}, ps...)
	for _, v := range sp {
		ret = append(ret, v...)
	}

	return ret
}

func (ps pointSlice) In(pnt Point) bool {
	for _, v := range ps {
		if v.Equal(pnt) {
			return true
		}
	}

	return false
}

func (ps pointSlice) Remove(i int) {
	if i >= len(ps) || i == -1 {
		return
	}
}

func (ps pointSlice) Clean() pointSlice {
	for k := len(ps) - 1; k >= 0; k-- {
		v := ps[k]
		if !v.Valid() {
			ps[len(ps)-1], ps[k] = ps[k], ps[len(ps)-1]
			ps = ps[:len(ps)-1]
		}
	}

	return ps
}

type displayBoard [8][8]Piece

func (b displayBoard) String() (str string) {
	for x := 0; x < 8; x++ {
		for y := 0; y < 8; y++ {
			pec := b[x][y]

			if !pec.Valid() {
				str += " "
			} else {
				str += pec.ShortString()
			}
			str += " "
		}

		str += "\n"
	}

	return str
}

const t_empty = Bishop
const t_points = King
const t_point = Knight

func board_set(p Point, ps pointSlice) *displayBoard {
	b := &displayBoard{}
	for x := 0; x < 8; x++ {
		for y := 0; y < 8; y++ {
			b[x][y] = Piece{T: t_empty}
		}
	}

	for _, v := range ps {
		b[v.X][v.Y] = Piece{T: t_points}
	}
	b[p.X][p.Y] = Piece{T: t_point}

	return b
}

func in_board(b *displayBoard, want pointSlice) bool {
	for _, v := range want {
		p := b[v.X][v.Y]
		if !p.Valid() || p.T != t_points {
			return false
		}
	}

	return true
}

func out_board(b *displayBoard, want pointSlice) bool {
	for x, v := range b {
		for y, b := range v {
			if b.Valid() && b.T != t_empty && b.T != t_point {
				found := false

				for _, p := range want {
					if p.X == int8(x) && p.Y == int8(y) {
						if b.T == t_points {
							found = true
							break
						} else {
							found = false
							break
						}
					}
				}

				if !found {
					//fmt.Printf("x: %d | y: %d | %v", x, y, b.String())
					return true
				}
			}
		}
	}

	return false
}

func generate_test(t *testing.T, name string, p Point, ps Points, want []Point) {
	sp := pointSlice{}
	for _, v := range ps {
		sp = append(sp, v)
	}

	b := board_set(p, sp)

	sort.Sort(sp)
	sort.Sort(pointSlice(want))

	t.Logf("src: %v", p)
	t.Logf("have: %v", ps)
	t.Logf("want: %v", want)

	t.Logf("\n%s", b.String())

	if !in_board(b, want) {
		t.Fatalf("%s is not predictable. points not inside selection", name)
	} else if out_board(b, want) {
		t.Fatalf("%s is not predictable. points outside selection", name)
	}
}

func TestPointsMerge(t *testing.T) {
	ps := pointSlice{
		{4, 3},
	}
	ps.Merge(pointSlice{{3, 4}})

	if !ps.In(Point{4, 3}) || ps.In(Point{3, 4}) {
		t.Fatalf("Merge does not merge...")
	}
}

func TestPointsClean(t *testing.T) {
	ps := pointSlice{
		{1, -4},
		{-6, 0},
		{1, 1},
	}
	ps = ps.Clean()

	t.Log(ps)
	if ps.In(Point{1, -4}) || ps.In(Point{-6, 0}) || !ps.In(Point{1, 1}) {
		t.Fatalf("Clean does not remove invalid points")
	}
}

// TODO: make diagonal tests more reliable...
func TestDiagonal(t *testing.T) {

	p := Point{4, 3}
	generate_test(t, "diagonal top", p, p.Diagonal(), []Point{
		// +, +
		{7, 6},
		{6, 5},
		{5, 4},
		// -, -
		{3, 2},
		{2, 1},
		{1, 0},
		// -, +
		{3, 4},
		{2, 5},
		{1, 6},
		{0, 7},
		// +, -
		{5, 2},
		{6, 1},
		{7, 0},
	})

	p = Point{5, 2}
	generate_test(t, "diagonal left", p, p.Diagonal(), []Point{
		// +, +
		{7, 4},
		{6, 3},
		// -, -
		{4, 1},
		{3, 0},
		// +, -
		{7, 0},
		{6, 1},
		// -, +
		{4, 3},
		{3, 4},
		{2, 5},
		{1, 6},
		{0, 7},
	})

	p = Point{5, 3}
	generate_test(t, "diagonal center", p, p.Diagonal(), []Point{
		// +, +
		{7, 5},
		{6, 4},
		// -, -
		{4, 2},
		{3, 1},
		{2, 0},
		// +, -
		{7, 1},
		{6, 2},
		// -, +
		{4, 4},
		{3, 5},
		{2, 6},
		{1, 7},
	})

	p = Point{5, 4}
	generate_test(t, "diagonal right", p, p.Diagonal(), []Point{
		// +, +
		{7, 6},
		{6, 5},
		// -, -
		{4, 3},
		{3, 2},
		{2, 1},
		{1, 0},
		// +, -
		{7, 2},
		{6, 3},
		// -, +
		{4, 5},
		{3, 6},
		{2, 7},
	})

	p = Point{6, 3}
	generate_test(t, "diagonal bottom", p, p.Diagonal(), []Point{
		// +, +
		{7, 4},
		// -, -
		{5, 2},
		{4, 1},
		{3, 0},
		// +, -
		{7, 2},
		// -, +
		{5, 4},
		{4, 5},
		{3, 6},
		{2, 7},
	})

	p = Point{3, 3}
}

func TestHorizontal(t *testing.T) {
	p := Point{4, 3}
	generate_test(t, "horizontal", p, p.Horizontal(), []Point{
		{4, 0},
		{4, 1},
		{4, 2},
		{4, 4},
		{4, 5},
		{4, 6},
		{4, 7},
	})
}

func TestVertical(t *testing.T) {
	p := Point{4, 3}
	generate_test(t, "vertical", p, p.Vertical(), []Point{
		{0, 3},
		{1, 3},
		{2, 3},
		{3, 3},
		{5, 3},
		{6, 3},
		{7, 3},
	})
}

func TestSquare(t *testing.T) {
	p := Point{4, 3}
	generate_test(t, "square", p, p.Square(), []Point{
		{5, 4},
		{5, 3},
		{5, 2},
		{4, 4},
		{4, 2},
		{3, 4},
		{3, 3},
		{3, 2},
	})

	// out of bounds?
	p = Point{7, 7}
	generate_test(t, "square", p, p.Square(), []Point{
		{6, 7},
		{7, 6},
		{6, 6},
	})
}

func TestKnight(t *testing.T) {
	p := Point{4, 3}
	generate_test(t, "knight", p, p.Knight(), []Point{
		{6, 4},
		{6, 2},
		{2, 4},
		{2, 2},
		{5, 5},
		{5, 1},
		{3, 5},
		{3, 1},
	})

	p = Point{0, 1}
	t.Log(p.Knight())
	generate_test(t, "knight", p, p.Knight(), []Point{
		{1, 3},
		{2, 0},
		{2, 2},
	})
}

func TestCorner(t *testing.T) {
	p := Point{4, 3}
	generate_test(t, "corner", p, p.Corner(), []Point{
		{5, 4},
		{5, 2},
		{3, 4},
		{3, 2},
	})
}

func TestDirection(t *testing.T) {
	p := Point{4, 3}
	ps := pointSlice{
		{3, 2},
		{3, 3},
		{3, 4},
		{4, 4},
		{5, 4},
		{5, 3},
		{5, 2},
		{4, 2},
	}
	dirs := []uint8{
		Set(DirUp, DirLeft),
		DirUp,
		Set(DirUp, DirRight),
		DirRight,
		Set(DirDown, DirRight),
		DirDown,
		Set(DirDown, DirLeft),
		DirLeft,
	}

	for i, v := range ps {
		d := p.Direction(v)
		if d != dirs[i] {
			t.Fatalf("%d failed. want: %d, have: %d", i, dirs[i], d)
		}
	}
}

/*
func TestSmaller(t *testing.T) {
	p := Point{4, 3}
	dst := Point{3, 3}
	dst2 := Point{4, 2}

	if p.Smaller(dst) != true || p.Smaller(dst2) != true {
		t.Fatalf("smaller bad logic")
	}
}

func TestBigger(t *testing.T) {
	p := Point{4, 3}
	dst := Point{5, 3}
	dst2 := Point{4, 4}

	if p.Bigger(dst) != true || p.Bigger(dst2) != true {
		t.Fatalf("smaller bad logic")
	}
}
*/
