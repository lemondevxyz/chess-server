package board

import "testing"

func TestGetKing(t *testing.T) {
	const ourking = 28
	const ournumber = 1

	const theirking = 4
	const theirnumber = 2

	have := GetKing(ournumber)
	if ourking != have {
		t.Fatalf("GetQueen: want: %d - have: %d", ourking, have)
	}
	have = GetKing(theirnumber)
	if theirking != have {
		t.Fatalf("GetQueen: want: %d - have: %d", theirking, have)
	}
}

func TestGetQueen(t *testing.T) {
	const ourqueen = 27
	const ournumber = 1

	const theirqueen = 3
	const theirnumber = 2

	have := GetQueen(ournumber)
	if ourqueen != have {
		t.Fatalf("GetQueen: want: %d - have: %d", ourqueen, have)
	}
	have = GetQueen(theirnumber)
	if theirqueen != have {
		t.Fatalf("GetQueen: want: %d - have: %d", theirqueen, have)
	}
}

func verify_by_2int(t *testing.T, num, w1, w2 int, fn func(uint8) [2]int) {
	vals := fn(uint8(num))
	h1, h2 := vals[0], vals[1]

	if w1 != h1 || w2 != h2 {
		t.Fatalf("want: %d / have: %d", w1, h1)
		t.Fatalf("want: %d / have: %d", w2, h2)
	}
}

func TestGetBishops(t *testing.T) {
	num, r1, r2 := 1, 26, 29
	verify_by_2int(t, num, r1, r2, GetBishops)

	num, r1, r2 = 2, 2, 5
	verify_by_2int(t, num, r1, r2, GetBishops)
}

func TestGetKnights(t *testing.T) {
	num, r1, r2 := 1, 25, 30
	verify_by_2int(t, num, r1, r2, GetKnights)

	num, r1, r2 = 2, 1, 6
	verify_by_2int(t, num, r1, r2, GetKnights)
}

func TestGetRooks(t *testing.T) {
	num, r1, r2 := 1, 24, 31
	verify_by_2int(t, num, r1, r2, GetRooks)

	num, r1, r2 = 2, 0, 7
	verify_by_2int(t, num, r1, r2, GetRooks)
}

func TestGetRange(t *testing.T) {
	const ournumber = 1
	const theirnumber = 2

	ourarr := []int{}
	for i := 16; i < 32; i++ {
		ourarr = append(ourarr, i)
	}

	theirarr := []int{}
	for i := 0; i < 16; i++ {
		theirarr = append(theirarr, i)
	}

	ourwant := GetRange(ournumber)
	theirwant := GetRange(theirnumber)

	for i := 0; i < 16; i++ {
		ourhave := ourarr[i]
		theirhave := theirarr[i]

		if ourhave != ourwant[i] {
			t.Fatalf("getRange does not match want. want: %d | have: %d", ourwant[i], ourhave)
		}
		if theirhave != theirwant[i] {
			t.Fatalf("getRange does not match want. want: %d | have: %d", theirwant[i], theirhave)
		}
	}
}

func TestGetStartRow(t *testing.T) {
	const player, enemy uint8 = 1, 2

	brd := NewBoard()
	_, rook, err := brd.Get(Point{0, GetStartRow(player)})
	if err != nil {
		t.Fatalf("error: %s", err)
	}

	if rook.T != Rook || rook.Player != player {
		t.Log(rook)
		t.Fatalf("Piece is not rook, or the player does not match")
	}

	_, rook, err = brd.Get(Point{0, GetStartRow(enemy)})
	if err != nil {
		t.Fatalf("error: %s", err)
	}

	if rook.T != Rook || rook.Player != enemy {
		t.Log(rook)
		t.Fatalf("Piece is not rook, or the player does not match")
	}
}

func TestGetPawnRow(t *testing.T) {
	const player, enemy uint8 = 1, 2

	brd := NewBoard()
	_, rook, err := brd.Get(Point{0, GetPawnRow(player)})
	if err != nil {
		t.Fatalf("error: %s", err)
	}

	if rook.T != PawnF || rook.Player != player {
		t.Log(rook)
		t.Fatalf("Piece is not rook, or the player does not match")
	}

	_, rook, err = brd.Get(Point{0, GetPawnRow(enemy)})
	if err != nil {
		t.Fatalf("error: %s", err)
	}

	if rook.T != PawnB || rook.Player != enemy {
		t.Log(rook)
		t.Fatalf("Piece is not rook, or the player does not match")
	}
}

func TestBelongsTo(t *testing.T) {
	const ournumber = 2
	ourvalues := [16]int8{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}

	for _, v := range ourvalues {
		if !BelongsTo(v, ournumber) {
			t.Fatalf("%d should belong to id: '%d'", v, ournumber)
		}
	}

	const theirnumber = 1
	theirvalues := [16]int8{16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31}
	for _, v := range theirvalues {
		if !BelongsTo(v, theirnumber) {
			t.Fatalf("%d should belong to id: '%d'", v, theirnumber)
		}
	}
}
