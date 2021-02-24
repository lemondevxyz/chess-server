package board

import "testing"

func piece_test(t *testing.T, name string, p Piece, ps Points) {
	for _, v := range ps {
		if !p.CanGo(v) {
			t.Fatalf("%s | want: true, have: false - pos: %v", name, v)
		}
	}

	for _, v := range ps.Outside() {
		if p.CanGo(v) {
			t.Fatalf("%s | want: false, have: true - pos: %v", name, v)
		}
	}
}

func TestCanGoOutOfBound(t *testing.T) {
	p := Piece{
		T:   PawnB,
		pos: Point{0, 0},
	}

	if p.CanGo(Point{-1, 0}) || p.CanGo(Point{0, -1}) || p.CanGo(Point{-1, -1}) {
		t.Fatalf("CanGo: allows out of bounds")
	}
}

func TestCanGoEqual(t *testing.T) {
	p := Piece{
		T:   PawnB,
		pos: Point{0, 0},
	}

	if p.CanGo(p.pos) {
		t.Fatalf("CanGo: allow same position movement")
	}
}

func TestPawnB(t *testing.T) {
	p := Piece{
		T:   PawnB,
		pos: Point{1, 1},
	}

	if !p.CanGo(Point{3, 1}) {
		t.Fatalf("Pawn cannot go two steps at the beginning")
	} else if !p.CanGo(Point{2, 1}) {
		t.Fatalf("Pawn cannot go a normal step at the beginning")
	} else if p.CanGo(Point{0, 1}) {
		t.Fatalf("Backward Pawn can go Forwards")
	}

	p.pos = Point{2, 1}
	if p.CanGo(Point{4, 1}) {
		t.Fatalf("Pawn can go two steps after the spawn position")
	}
}

func TestPawnF(t *testing.T) {
	p := Piece{
		T:   PawnF,
		pos: Point{6, 1},
	}

	if !p.CanGo(Point{4, 1}) {
		t.Fatalf("Pawn cannot go two steps at the beginning")
	} else if !p.CanGo(Point{5, 1}) {
		t.Fatalf("Pawn cannot go a normal step at the beginning")
	} else if p.CanGo(Point{7, 1}) {
		t.Fatalf("Forward Pawn can go Backwards")
	}

	p.pos = Point{5, 1}
	if p.CanGo(Point{3, 1}) {
		t.Fatalf("Pawn can go two steps after the spawn position")
	}
}

func TestRook(t *testing.T) {
	pos := Point{4, 3}

	piece_test(t, "Rook", Piece{
		T:   Rook,
		pos: pos,
	}, Points{
		{7, 3},
		{6, 3},
		{5, 3},
		{3, 3},
		{2, 3},
		{1, 3},
		{0, 3},
		{4, 7},
		{4, 6},
		{4, 5},
		{4, 4},
		{4, 2},
		{4, 1},
		{4, 0},
	})
}

// no need - generators tests for this
//func TestKnight(t *testing.T) {
//}

// no need - generators tests for this
/*
func TestBishop(t *testing.T) {
	pos := Point{4, 3}

	piece_test(t, "Bishop", Piece{
		T:   Bishop,
		pos: pos,
	}, Points{
		{7, 6},
		{6, 5},
		{5, 4},
		{3, 2},
		{2, 1},
		{1, 0},

		{7, 0},
		{6, 1},
		{5, 2},
		{3, 4},
		{2, 5},
		{1, 6},
		{0, 7},
	})
}
*/

func TestQueen(t *testing.T) {
	pos := Point{4, 3}
	piece_test(t, "Queen", Piece{
		T:   Queen,
		pos: pos,
	}, Points{
		{7, 3},
		{6, 3},
		{5, 3},
		{3, 3},
		{2, 3},
		{1, 3},
		{0, 3},
		{4, 7},
		{4, 6},
		{4, 5},
		{4, 4},
		{4, 2},
		{4, 1},
		{4, 0},

		{7, 6},
		{6, 5},
		{5, 4},
		{3, 2},
		{2, 1},
		{1, 0},

		{7, 0},
		{6, 1},
		{5, 2},
		{3, 4},
		{2, 5},
		{1, 6},
		{0, 7},
	})
}

func TestKing(t *testing.T) {
	pos := Point{4, 3}
	piece_test(t, "King", Piece{
		T:   King,
		pos: pos,
	}, Points{
		{5, 2},
		{5, 3},
		{5, 4},

		{4, 4},
		{4, 2},

		{3, 2},
		{3, 3},
		{3, 4},
	})
}
