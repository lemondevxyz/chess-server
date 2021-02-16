package board

import (
	"testing"
)

type TestCase struct {
	src Point
	dst []Point
	ok  bool
}

func generate(ignore []Point) []Point {
	max := 8
	ret := []Point{}

	for x := 0; x < max; x++ {
		for y := 0; y < max; y++ {
			ret = append(ret, Point{x, y})
		}
	}

	for _, v := range ignore {
		for i := len(ret) - 1; i >= 0; i-- {
			if Equal(v, ret[i]) {
				ret = append(ret[:i], ret[i+1:]...)
			}
		}
	}

	return ret
}

func TestCanGoOutOfBound(t *testing.T) {
	p := Piece{
		T: PawnB,
		X: 0,
		Y: 0,
	}

	if p.CanGo(-1, 0) || p.CanGo(0, -1) || p.CanGo(-1, -1) {
		t.Fatalf("CanGo: allows out of bounds")
	}
}

func TestCanGoPawn(t *testing.T) {

	p := Piece{
		T: PawnB,
	}

	tc := []TestCase{
		{
			src: Point{1, 1},
			dst: []Point{
				{3, 1},
				{2, 1},
			},
			ok: true,
		},
		{
			src: Point{2, 1},
			dst: []Point{
				{3, 1},
			},
			ok: true,
		},
	}

	for _, v := range tc {
		p.X, p.Y = v.src.X, v.src.Y
		for _, d := range v.dst {
			if p.CanGo(d.X, d.Y) != v.ok {
				t.Fatalf("src: %v - dst: %v - want: %v", v.src, d, v.ok)
			}
		}

		for _, d := range generate(v.dst) {
			if p.CanGo(d.X, d.Y) != !v.ok {
				t.Fatalf("src: %v - dst: %v - want: %v", v.src, d, false)
			}
		}
	}

	p.T = PawnF
	tc = []TestCase{
		{
			src: Point{6, 2},
			dst: []Point{
				{5, 2},
				{4, 2},
			},
			ok: true,
		},
		{
			src: Point{2, 1},
			dst: []Point{
				{1, 1},
			},
			ok: true,
		},
	}

	for _, v := range tc {
		p.X, p.Y = v.src.X, v.src.Y
		for _, d := range v.dst {
			if v.ok != p.CanGo(d.X, d.Y) {
				t.Fatalf("src: %v - dst: %v - want: %v", v.src, d, v.ok)
			}
		}

		for _, d := range generate(v.dst) {
			if p.CanGo(d.X, d.Y) != !v.ok {
				t.Fatalf("src: %v - dst: %v - want: %v", v.src, d, false)
			}
		}
	}

}

func TestCanGoBishop(t *testing.T) {

	p := Piece{
		T: Bishop,
	}

	tc := []TestCase{
		{
			src: Point{4, 4},
			dst: []Point{
				{7, 7},
				{6, 6},
				{5, 5},
				{3, 3},
				{2, 2},
				{1, 1},
				{0, 0},

				{1, 7},
				{2, 6},
				{3, 5},
				{5, 3},
				{6, 2},
				{7, 1},
			},
			ok: true,
		},
	}

	for _, v := range tc {
		p.X, p.Y = v.src.X, v.src.Y
		for _, d := range v.dst {
			if v.ok != p.CanGo(d.X, d.Y) {
				t.Fatalf("src: %v - dst: %v - want: %v", v.src, d, v.ok)
			}
		}

		for _, d := range generate(v.dst) {
			if p.CanGo(d.X, d.Y) != !v.ok {
				t.Fatalf("src: %v - dst: %v - want: %v", v.src, d, false)
			}
		}
	}
}

func TestCanGoKnight(t *testing.T) {

	p := Piece{
		T: Knight,
	}

	tc := []TestCase{
		{
			src: Point{4, 3},
			dst: []Point{
				{6, 4},
				{6, 2},
				{5, 5},
				{5, 1},
				{2, 4},
				{2, 2},
				{3, 5},
				{3, 1},
			},
			ok: true,
		},
	}

	for _, v := range tc {
		p.X, p.Y = v.src.X, v.src.Y
		for _, d := range v.dst {
			if v.ok != p.CanGo(d.X, d.Y) {
				t.Fatalf("src: %v - dst: %v - want: %v", v.src, d, v.ok)
			}
		}

		for _, d := range generate(v.dst) {
			if p.CanGo(d.X, d.Y) != !v.ok {
				t.Fatalf("src: %v - dst: %v - want: %v", v.src, d, false)
			}
		}
	}

}

func TestCanGoRook(t *testing.T) {

	p := Piece{
		T: Rook,
	}

	tc := []TestCase{
		{
			src: Point{4, 4},
			dst: []Point{
				{7, 4},
				{6, 4},
				{5, 4},
				{3, 4},
				{2, 4},
				{1, 4},
				{0, 4},

				{4, 7},
				{4, 6},
				{4, 5},
				{4, 3},
				{4, 2},
				{4, 1},
				{4, 0},
			},
			ok: true,
		},
	}

	for _, v := range tc {
		p.X, p.Y = v.src.X, v.src.Y
		for _, d := range v.dst {
			if v.ok != p.CanGo(d.X, d.Y) {
				t.Fatalf("src: %v - dst: %v - want: %v", v.src, d, v.ok)
			}
		}

		for _, d := range generate(v.dst) {
			if p.CanGo(d.X, d.Y) != !v.ok {
				t.Fatalf("src: %v - dst: %v - want: %v", v.src, d, false)
			}
		}
	}

}

func TestCanGoQueen(t *testing.T) {

	p := Piece{
		T: Queen,
	}

	tc := []TestCase{
		{
			src: Point{4, 4},
			dst: []Point{
				// square
				{3, 3},
				{3, 4},
				{3, 5},
				{4, 3},
				{4, 5},
				{5, 3},
				{5, 4},
				{5, 5},
				// horizontal
				{7, 4},
				{6, 4},
				{5, 4},
				{3, 4},
				{2, 4},
				{1, 4},
				{0, 4},
				// vertical
				{4, 7},
				{4, 6},
				{4, 5},
				{4, 3},
				{4, 2},
				{4, 1},
				{4, 0},
				// diagonal
				{7, 7},
				{6, 6},
				{5, 5},
				{3, 3},
				{2, 2},
				{1, 1},
				{0, 0},
				{1, 7},
				{2, 6},
				{3, 5},
				{5, 3},
				{6, 2},
				{7, 1},
			},
			ok: true,
		},
	}

	for _, v := range tc {
		p.X, p.Y = v.src.X, v.src.Y
		for _, d := range v.dst {
			if v.ok != p.CanGo(d.X, d.Y) {
				t.Fatalf("src: %v - dst: %v - want: %v", v.src, d, v.ok)
			}
		}

		for _, d := range generate(v.dst) {
			if p.CanGo(d.X, d.Y) != !v.ok {
				t.Fatalf("src: %v - dst: %v - want: %v", v.src, d, false)
			}
		}
	}

}

func TestCanGoKing(t *testing.T) {

	p := Piece{
		T: King,
	}

	tc := []TestCase{
		{
			src: Point{4, 4},
			dst: []Point{
				// square
				{3, 3},
				{3, 4},
				{3, 5},
				{4, 3},
				{4, 5},
				{5, 3},
				{5, 4},
				{5, 5},
			},
			ok: true,
		},
	}

	for _, v := range tc {
		p.X, p.Y = v.src.X, v.src.Y
		for _, d := range v.dst {
			if v.ok != p.CanGo(d.X, d.Y) {
				t.Fatalf("src: %v - dst: %v - want: %v", v.src, d, v.ok)
			}
		}

		for _, d := range generate(v.dst) {
			if p.CanGo(d.X, d.Y) != !v.ok {
				t.Fatalf("src: %v - dst: %v - want: %v", v.src, d, false)
			}
		}
	}

}
