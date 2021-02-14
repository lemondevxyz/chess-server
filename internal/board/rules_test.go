package board

import "testing"

func TestEqual(t *testing.T) {
	src := Point{1, 1}
	dst := Point{1, 1}

	if !Equal(src, dst) {
		t.Fatalf("Equal: bad code")
	}

	dst.X++
	if Equal(src, dst) {
		t.Fatalf("Equal: bad code")
	}
}

// We'll test against every possible valid case
// and some invalid cases
func TestForward(t *testing.T) {
	ok := func(res bool, src, dst Point) {
		if Forward(src, dst) != res {
			t.Fatalf("src: %v, dst: %v - Failed. want: %v", src, dst, res)
		}
	}

	src := Point{2, 2}
	dst := Point{1, 0}
	ok(true, src, dst)

	dst = Point{0, 1}
	ok(true, src, dst)

	dst = Point{1, 1}
	ok(true, src, dst)

	dst = Point{3, 0}
	ok(false, src, dst)

	dst = Point{2, 2}
	ok(false, src, dst)

}

func TestBackward(t *testing.T) {
	ok := func(res bool, src, dst Point) {
		if Backward(src, dst) != res {
			t.Fatalf("src: %v, dst: %v - Failed. want: %v", src, dst, res)
		}
	}

	src := Point{0, 0}
	dst := Point{1, 0}
	ok(true, src, dst)

	dst = Point{0, 1}
	ok(true, src, dst)

	dst = Point{1, 1}
	ok(true, src, dst)

	dst = Point{-1, 0}
	ok(false, src, dst)

	dst = Point{0, 0}
	ok(false, src, dst)

}

func TestVertical(t *testing.T) {
	ok := func(res bool, src, dst Point) {
		if Vertical(src, dst) != res {
			t.Fatalf("src: %v, dst: %v - Failed. want: %v", src, dst, res)
		}
	}

	src := Point{5, 5}
	dst := Point{5, 9}
	ok(true, src, dst)

	dst = Point{5, 0}
	ok(true, src, dst)

	dst = Point{6, 5}
	ok(false, src, dst)

}

func TestHorizontal(t *testing.T) {
	ok := func(res bool, src, dst Point) {
		if Horizontal(src, dst) != res {
			t.Fatalf("src: %v, dst: %v - Failed. want: %v", src, dst, res)
		}
	}

	src := Point{5, 5}
	dst := Point{9, 5}
	ok(true, src, dst)

	dst = Point{0, 5}
	ok(true, src, dst)

	dst = Point{5, 6}
	ok(false, src, dst)

}

func TestDiagonal(t *testing.T) {
	ok := func(res bool, src, dst Point) {
		if Diagonal(src, dst) != res {
			t.Fatalf("src: %v, dst: %v - Failed. want: %v", src, dst, res)
		}
	}

	src := Point{3, 3}
	dst := Point{2, 2}
	ok(true, src, dst)

	dst = Point{4, 4}
	ok(true, src, dst)

	dst = Point{9, 9}
	ok(true, src, dst)

	dst = Point{4, 3}
	ok(false, src, dst)

}

func TestWithin(t *testing.T) {

	area := Point{1, 1}
	ok := func(res bool, src, dst Point) {
		if Within(area, src, dst) != res {
			t.Fatalf("src: %v, dst: %v - Failed. want: %v", src, dst, res)
		}
	}

	src := Point{0, 0}
	dst := Point{1, 1}
	ok(true, src, dst)

	dst = Point{1, -1}
	ok(true, src, dst)

	dst = Point{-1, 1}
	ok(true, src, dst)

	dst = Point{-1, -1}
	ok(true, src, dst)

	dst = Point{0, 1}
	ok(false, src, dst)

	area = Point{2, 1}
	dst = Point{2, 1}
	ok(true, src, dst)

	dst = Point{-2, 1}
	ok(true, src, dst)

	dst = Point{1, 2}
	ok(false, src, dst)

}

func TestSquare(t *testing.T) {
	ok := func(res bool, src, dst Point) {
		if Square(src, dst) != res {
			t.Fatalf("src: %v, dst: %v - Failed. want: %v", src, dst, res)
		}
	}

	src := Point{0, 0}
	dsts := []Point{
		{1, -1},
		{1, 0},
		{1, 1},
		{0, 1},
		{0, -1},
		{-1, 0},
		{0, -1},
		{-1, -1},
	}
	for _, dst := range dsts {
		ok(true, src, dst)
	}

	dst := Point{0, 0}
	ok(false, src, dst)

	dst = Point{-2, -2}
	ok(false, src, dst)
}
