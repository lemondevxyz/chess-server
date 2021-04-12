package board

import "testing"

func piece_test(t *testing.T, name string, p Piece, ps Points) {
	for _, v := range ps {
		if !p.CanGo(v) {
			t.Logf("%s", p.Pos)
			t.Fatalf("%s | want: true, have: false - Pos: %v", name, v)
		}
	}

	for _, v := range ps.Outside() {
		if p.CanGo(v) {
			t.Logf("%s", p.Pos)
			t.Fatalf("%s | want: false, have: true - Pos: %v", name, v)
		}
	}
}

func TestCanGoOutOfBound(t *testing.T) {
	p := Piece{
		Kind: PawnB,
		Pos:  Point{0, 0},
	}

	if p.CanGo(Point{-1, 0}) || p.CanGo(Point{0, -1}) || p.CanGo(Point{-1, -1}) {
		t.Fatalf("CanGo: allows out of bounds")
	}
}

func TestCanGoEqual(t *testing.T) {
	p := Piece{
		Kind: PawnB,
		Pos:  Point{0, 0},
	}

	if p.CanGo(p.Pos) {
		t.Fatalf("CanGo: allow same position movement")
	}
}

func TestPawnB(t *testing.T) {
	p := Piece{
		Kind: PawnB,
		Pos:  Point{1, 1},
	}

	if !p.CanGo(Point{1, 3}) {
		t.Fatalf("Pawn cannot go two steps at the beginning")
	} else if !p.CanGo(Point{1, 2}) {
		t.Fatalf("Pawn cannot go a normal step at the beginning")
	} else if p.CanGo(Point{1, 0}) {
		t.Fatalf("Backward Pawn can go Forwards")
	} else if p.CanGo(Point{2, 1}) {
		t.Fatalf("Backward Pawn can move horizontally")
	}

	p.Pos = Point{2, 1}
	if p.CanGo(Point{4, 1}) {
		t.Fatalf("Pawn can go two steps after the spawn position")
	}
}

func TestPawnF(t *testing.T) {
	p := Piece{
		Kind: PawnF,
		Pos:  Point{1, 6},
	}

	if !p.CanGo(Point{1, 4}) {
		t.Fatalf("Pawn cannot go two steps at the beginning")
	} else if !p.CanGo(Point{1, 5}) {
		t.Fatalf("Pawn cannot go a normal step at the beginning")
	} else if p.CanGo(Point{1, 6}) {
		t.Fatalf("Forward Pawn can go Backwards")
	} else if p.CanGo(Point{2, 6}) {
		t.Fatalf("Forward Pawn can move horizontally")
	}

	p.Pos = Point{5, 1}
	if p.CanGo(Point{3, 1}) {
		t.Fatalf("Pawn can go two steps after the spawn position")
	}
}

func TestRook(t *testing.T) {
	pos := Point{4, 3}

	ps := Points{}
	ps.Insert(
		Point{7, 3},
		Point{6, 3},
		Point{5, 3},
		Point{3, 3},
		Point{2, 3},
		Point{1, 3},
		Point{0, 3},
		Point{4, 7},
		Point{4, 6},
		Point{4, 5},
		Point{4, 4},
		Point{4, 2},
		Point{4, 1},
		Point{4, 0},
	)

	piece_test(t, "Rook", Piece{
		Kind: Rook,
		Pos:  pos,
	}, ps)
}

func TestQueen(t *testing.T) {
	pos := Point{4, 3}
	ps := Points{}
	ps.Insert(
		Point{7, 3},
		Point{6, 3},
		Point{5, 3},
		Point{3, 3},
		Point{2, 3},
		Point{1, 3},
		Point{0, 3},
		Point{4, 7},
		Point{4, 6},
		Point{4, 5},
		Point{4, 4},
		Point{4, 2},
		Point{4, 1},
		Point{4, 0},

		Point{7, 6},
		Point{6, 5},
		Point{5, 4},
		Point{3, 2},
		Point{2, 1},
		Point{1, 0},

		Point{7, 0},
		Point{6, 1},
		Point{5, 2},
		Point{3, 4},
		Point{2, 5},
		Point{1, 6},
		Point{0, 7},
	)

	piece_test(t, "Queen", Piece{
		Kind: Queen,
		Pos:  pos,
	}, ps)
}

func TestKing(t *testing.T) {
	pos := Point{4, 3}
	ps := Points{}
	ps.Insert(
		Point{5, 2},
		Point{5, 3},
		Point{5, 4},

		Point{4, 4},
		Point{4, 2},

		Point{3, 2},
		Point{3, 3},
		Point{3, 4},
	)

	piece_test(t, "King", Piece{
		Kind: King,
		Pos:  pos,
	}, ps)
}
