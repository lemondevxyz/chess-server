package board

import (
	"testing"
)

type TestCase struct {
	src Point
	dst []Point
	ok  bool
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
		{
			src: Point{2, 1},
			dst: []Point{
				{4, 1},
				{2, 2},
			},
			ok: false,
		},
	}

	for _, v := range tc {
		p.X, p.Y = v.src.X, v.src.Y
		for _, d := range v.dst {
			if p.CanGo(d.X, d.Y) != v.ok {
				t.Fatalf("src: %v - dst: %v - want: %v", v.src, d, v.ok)
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
		{
			src: Point{2, 1},
			dst: []Point{
				{3, 1},
			},
			ok: false,
		},
	}

	for _, v := range tc {
		p.X, p.Y = v.src.X, v.src.Y
		for _, d := range v.dst {
			if v.ok != p.CanGo(d.X, d.Y) {
				t.Fatalf("src: %v - dst: %v - want: %v", v.src, d, v.ok)
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
			src: Point{2, 1},
			dst: []Point{
				{4, 3},
				{3, 2},
				{1, 2},
				{1, 0},
			},
			ok: true,
		},
		{
			src: Point{2, 1},
			dst: []Point{
				{2, 1},
				{2, 2},
				{3, 1},
			},
			ok: false,
		},
	}

	for _, v := range tc {
		p.X, p.Y = v.src.X, v.src.Y
		for _, d := range v.dst {
			if v.ok != p.CanGo(d.X, d.Y) {
				t.Fatalf("src: %v - dst: %v - want: %v", v.src, d, v.ok)
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
		{
			src: Point{4, 3},
			dst: []Point{
				{6, 5},
				{6, 1},
				{2, 3},
			},
			ok: false,
		},
	}

	for _, v := range tc {
		p.X, p.Y = v.src.X, v.src.Y
		for _, d := range v.dst {
			if v.ok != p.CanGo(d.X, d.Y) {
				t.Fatalf("src: %v - dst: %v - want: %v", v.src, d, v.ok)
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
				{4, 7},
				{4, 3},
				{2, 4},
				{4, 5},
				{1, 4},
			},
			ok: true,
		},
		{
			src: Point{4, 4},
			dst: []Point{
				{5, 1},
				{6, 7},
				{6, 5},
				{6, 1},
				{2, 3},
			},
			ok: false,
		},
	}

	for _, v := range tc {
		p.X, p.Y = v.src.X, v.src.Y
		for _, d := range v.dst {
			if v.ok != p.CanGo(d.X, d.Y) {
				t.Fatalf("src: %v - dst: %v - want: %v", v.src, d, v.ok)
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
				{6, 4},
				{3, 4},
				{2, 4},
				// vertical
				{4, 6},
				{4, 3},
				{4, 2},
				// diagonal
				{5, 5},
				{7, 7},
				{1, 1},
			},
			ok: true,
		},
		{
			src: Point{4, 4},
			dst: []Point{
				{5, 1},
				{6, 7},
				{6, 5},
				{6, 1},
				{2, 3},
				{7, 0},
			},
			ok: false,
		},
	}

	for _, v := range tc {
		p.X, p.Y = v.src.X, v.src.Y
		for _, d := range v.dst {
			if v.ok != p.CanGo(d.X, d.Y) {
				t.Fatalf("src: %v - dst: %v - want: %v", v.src, d, v.ok)
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
		{
			src: Point{4, 4},
			dst: []Point{
				{5, 1},
				{6, 7},
				{6, 5},
				{6, 1},
				{2, 3},
				{7, 0},
			},
			ok: false,
		},
	}

	for _, v := range tc {
		p.X, p.Y = v.src.X, v.src.Y
		for _, d := range v.dst {
			if v.ok != p.CanGo(d.X, d.Y) {
				t.Fatalf("src: %v - dst: %v - want: %v", v.src, d, v.ok)
			}
		}
	}

}
