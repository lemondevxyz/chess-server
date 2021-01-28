package board

import "testing"

type TestCase struct {
	src Point
	dst []Point
	ok  bool
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
		p.src = v.src
		for _, p := range v {
			p.CanGo(p.dst)
		}
	}
}
