package board

import (
	"testing"
	"time"
)

// placement test
func TestNewBoard(t *testing.T) {
	u := [32]uint8{
		Rook, Knight, Bishop, Queen, King, Bishop, Knight, Rook,
		PawnB, PawnB, PawnB, PawnB, PawnB, PawnB, PawnB, PawnB,
		PawnF, PawnF, PawnF, PawnF, PawnF, PawnF, PawnF, PawnF,
		Rook, Knight, Bishop, Queen, King, Bishop, Knight, Rook,
	}

	b := NewBoard()

	t.Log(b.data)

	ps := []Point{}
	for y := int8(0); y < 2; y++ {
		for x := int8(0); x < 8; x++ {
			//t.Log(y, x)
			ps = append(ps, Point{x, y})
		}
	}
	for y := int8(6); y < 8; y++ {
		for x := int8(0); x < 8; x++ {
			ps = append(ps, Point{x, y})
			// ps.Insert(Point{x, y})
		}
	}

	// t.Logf("want array: %v", ps)

	for x := 0; x < 32; x++ {
		if b.data[x].T != u[x] {
			t.Fatalf("rows(types) are not setup properly: %d | want: %d, have: %d", x, u[x], b.data[x].T)
		}
		if !b.data[x].Pos.Equal(ps[x]) {
			t.Fatalf("rows(position) are not setup properly: %d | want: %s, have: %s", x, ps[x].String(), b.data[x].Pos.String())
		}
	}

	t.Logf("\n%s", b.String())
}

func TestBoardCopy(t *testing.T) {
	brd := NewBoard()
	drb := brd.Copy()

	if brd.String() != drb.String() {
		t.Fatalf("board.Copy doesnt copy well")
	}

	brd.Move(19, Point{3, 4})

	t.Logf("\n%s", brd.String())
	t.Logf("\n%s", drb.String())
	if brd.String() == drb.String() {
		t.Fatalf("board.Copy original move applies to copy of board")
	}
}

func TestBoardListen(t *testing.T) {
	b := NewBoard()

	valid := make(chan bool)
	invalid := make(chan bool)

	ok := make(chan bool)
	go func() {
		a := <-valid
		b := <-invalid

		if a && b {
			ok <- true
		}
	}()

	x, y := false, false
	// func(p Piece, src Point, dst Point, ret bool)
	b.Listen(func(_ Piece, _, _ Point, ret bool) {
		t.Log(ret)
		if ret {
			valid <- true
		} else {
			invalid <- true
		}
	})

	b.Move(9, Point{1, 3})
	b.Move(9, Point{1, 2})

	select {
	case <-time.After(time.Millisecond * 20):
		t.Fatalf("Listen does not listen. pre: %t - post: %t", x, y)
	case <-ok:
		break
	}

}

func TestBoardSet(t *testing.T) {
	b := NewBoard()

	pos := Point{1, 1}
	const place = 2

	b.Set(place, pos)
	if b.data[place].Pos != pos {
		t.Fatalf("Set does not work")
	}
}

func TestBoardMove(t *testing.T) {
	b := NewBoard()

	if !b.Move(11, Point{3, 3}) {
		t.Fatalf("CanGo failed")
	}

	if b.data[10].T != PawnB {
		t.Fatalf("Pawn didn't move")
	}
}

func TestBoardMovePawn(t *testing.T) {

	b := NewBoard()

	// dark 1
	// light 1
	d1 := 11
	l1 := 20
	l2 := 21

	t.Log(b.data[d1], b.data[l1], b.data[l2])

	b.Move(d1, Point{3, 3})
	b.Move(l1, Point{4, 4})
	b.Move(l2, Point{5, 5})

	t.Log(b.data[d1], b.data[l1], b.data[l2])

	if !b.Move(d1, Point{4, 4}) {
		t.Fatalf("backward pawn should kill other pawn")
	} else if !b.Move(l2, Point{4, 4}) {
		t.Fatalf("forward pawn should kill other pawn")
	}

	d2 := 12
	b.Move(d2, Point{3, 4})
	if b.Move(d2, Point{4, 4}) {
		t.Fatalf("pawn cannot takeout pawn in front of it")
	}

}

// probs overkill but why not
// basically test if other pieces can kill enemy pieces
func TestBoardMoveKill(t *testing.T) {

	b := NewBoard()
	cs := []struct {
		nm1 int
		nm2 int
		src Point
		dst Point
		t   uint8
	}{
		// Bishop
		// Knight
		// Rook
		// Queen
		// King
		{
			nm1: 2,
			nm2: 29,
			src: Point{3, 3},
			dst: Point{4, 4},
			t:   Bishop,
		},
		{
			nm1: 1,
			nm2: 25,
			src: Point{3, 3},
			dst: Point{4, 5},
			t:   Knight,
		},
		{
			nm1: 0,
			nm2: 24,
			src: Point{3, 3},
			dst: Point{4, 3},
			t:   Rook,
		},
		{
			nm1: 7,
			nm2: 31,
			src: Point{3, 3},
			dst: Point{3, 5},
			t:   Rook,
		},
		{
			nm1: 3,
			nm2: 27,
			src: Point{3, 3},
			dst: Point{3, 4},
			t:   Queen,
		},
		{
			nm1: 4,
			nm2: 28,
			src: Point{3, 3},
			dst: Point{3, 4},
			t:   King,
		},
	}

	for k, v := range cs {
		b.Set(v.nm1, v.src)
		b.Set(v.nm2, v.dst)

		if !b.Move(v.nm1, v.dst) {
			t.Logf("index: %d", k)
			t.Logf("\n%s", b)
			t.Fatalf("test cordinates are invalid. src: %d - dst: %d - type: %d", v.src, v.dst, v.t)
			return
		} else {
			if b.data[v.nm1].Pos.Equal(v.src) {
				t.Fatalf("move doesn't actually move")
			} else {
				if !b.data[v.nm1].Pos.Equal(v.dst) || b.data[v.nm2].Pos.Valid() {
					t.Fatalf("move doesn't replace enemy")
				}
			}
		}

		b.Set(v.nm1, Point{-1, -1})
		b.Set(v.nm2, Point{-1, -1})
	}

	b.Move(0, Point{2, 1})
	b.Move(1, Point{3, 2})

	if b.Move(0, Point{3, 2}) {
		t.Fatalf("ally killed")
	}

}

// check if some pieces can move over pieces that are in the way...
// it's a tiny bit overkill to test every piece type, but reliablity is worth every price
func TestBoardMoveInTheWay(t *testing.T) {
	generic := func(id int) {
		brd := NewBoard()
		pec, _ := brd.GetByIndex(id)

		for _, v := range pec.Possib() {
			z, _ := brd.Possib(id)
			if brd.Move(id, v) {
				t.Logf("pos: %s", v.String())
				t.Logf("board.Possib: %s", z)
				t.Logf("board:\n%s", brd.String())
				t.Fatalf("%s can move over other pieces", pec.Name())
			}
		}
	}

	for i := 0; i < 16; i++ {
		x := i
		if i == 1 || i == 6 || i == 9 || i == 14 { // skip over knights
			continue
		}
		if i >= 8 {
			x += 16
		}

		generic(x)
	}

	b := NewBoard()

	id := 25
	pec, _ := b.GetByIndex(id)

	for _, v := range pec.Possib() {

		b = NewBoard()
		pec, _ = b.GetByIndex(id)
		if !v.Valid() {
			continue
		}

		_, cep, err := b.Get(v)
		want := true
		if err == nil {
			if cep.Player == pec.Player {
				want = false
			}
		}

		have := b.Move(id, v)
		if have != want {
			t.Logf("%s", v.String())
			t.Logf("\n%s", b.String())
			t.Logf("want: %t, have: %t", want, have)
			if want == false {
				t.Fatalf("knight can replace nearby pieces")
			} else {
				t.Fatalf("knight cannot skip over pieces")
			}
		}
	}
}

// while above makes sure that no "special" piece can move over from it's starting position, this one tests past bugs.
func TestBoardPieceBugInTheWay(t *testing.T) {
	{ // bishop shouldn't skip over enemy pawn and killing knight
		// constant ids
		const ourpawn = 20   // {4, 6}
		const ourbishop = 29 // {5, 7}

		brd := NewBoard()
		// first remove our pawn at {6, 4}
		brd.Set(ourpawn, Point{-1, -1})

		pos := Point{2, 4}
		// move bishop to 2, 4
		brd.Set(ourbishop, pos)

		// t.Logf("\n%s", brd.String())
		if brd.Move(ourbishop, Point{5, 0}) {
			t.Fatalf("bishop can skip enemy pawn and kill knight")
		}
	}
	{ // knight cannot override nearby pawn, but it's in the possible moves
		const ourknight = 30 // at {6, 7}

		brd := NewBoard()

		// pec, _ := brd.GetByIndex(ourknight)
		// t.Logf("knight position: %s", pec.Pos)
		possib, _ := brd.Possib(ourknight)
		// t.Logf("knight possible moves: %s", possib)

		if possib.In(Point{4, 6}) {
			t.Fatalf("knight possible moves is killing nearby pawn")
		}
	}
	{ // pawn possiblity needs to include killable pieces
		const ourpawn = 20   // at {4, 6}
		const enemypawn = 11 // at {3, 1}

		brd := NewBoard()
		pec, _ := brd.GetByIndex(ourpawn)

		pec.Pos.Y -= 2
		if !brd.Move(ourpawn, pec.Pos) {
			t.Fatalf("major fault within move")
		}

		pec, _ = brd.GetByIndex(enemypawn)

		pec.Pos.Y += 2
		if !brd.Move(enemypawn, pec.Pos) {
			t.Fatalf("major fault within move")
		}

		possib, _ := brd.Possib(ourpawn)
		if !possib.In(Point{3, 3}) {
			t.Fatalf("pawn does not include killable pieces")
		}
	}
	{ // queen possible moves should not include it's fellow allies
		const ourqueen = 27
		brd := NewBoard()

		// pec, _ := brd.GetByIndex(ourqueen)

		ps := Points{}
		ps.Insert(
			Point{7, 3},
			Point{7, 5},
			Point{6, 4},
			Point{6, 3},
			Point{6, 5},
		)

		sp, _ := brd.Possib(ourqueen)
		for _, v := range sp {
			if ps.In(v) {
				t.Fatalf("queen possible moves includes it's fellow allies. point: %s", v.String())
			}
		}
	}
	{ // pawn backward shouldn't have X-2 at X = 1, pawn forward shouldn't have X+2 at X = 6
		const pawnb = 9
		const pawnf = 17

		brd := NewBoard()

		brd.Set(pawnb, Point{1, 6})
		brd.Set(pawnf, Point{1, 1})

		// t.Logf("\n%s", brd)

		psb, _ := brd.Possib(pawnb)
		if psb.In(Point{1, 4}) {
			t.Fatalf("backward pawn can go to 1, 4")
		}

		psf, _ := brd.Possib(pawnf)
		if psf.In(Point{1, 3}) {
			t.Fatalf("forward pawn can go to 3, 1")
		}
	}
}

/*
func TestBoardCheckmate(t *testing.T) {
	// testing top 10 fast checkmates: https://www.chess.com/article/view/fastest-chess-checkmates
	// black always wins
	empty := func(brd *Board, src Point) {
		brd.Set(&Piece{Pos: src, T: Empty})
	}
	try := func(brd *Board, src Point, pnt Point) {
		pec := brd.Get(src)
		if pec == nil {
			t.Fatalf("invalid piece at point: %s", src.String())
		}

		if !brd.Move(pec, pnt) {
			t.Fatalf("perfectly legal move is failing. piece: %s | src: %s - dst: %s", pec.String(), pec.Pos.String(), pnt.String())
		}
	}
	{ // BUG: checkmate works through other pieces
		brd := NewBoard()

		try(brd, Point{6, 5}, Point{5, 5})
		try(brd, Point{1, 4}, Point{3, 4})
		try(brd, Point{6, 6}, Point{5, 6})
		try(brd, Point{0, 3}, Point{4, 7})
		// R N B   K B N R
		// P P P P   P P P
		//
		//         P
		//               Q
		//           P P
		// P P P P P     P
		// R N B Q K B N R
		// t.Logf("\n%s", brd)

		if brd.Checkmate(1) {
			t.Fatalf("bad checkmate, over other piece")
		}
	}
	{ // checkmate but not a final checkmate
		brd := NewBoard()

		try(brd, Point{6, 5}, Point{5, 5})
		try(brd, Point{1, 4}, Point{3, 4})
		try(brd, Point{6, 0}, Point{4, 0})
		try(brd, Point{0, 3}, Point{4, 7})

		if !brd.Checkmate(1) {
			t.Fatalf("no checkmate")
		}
		// R N B   K B N R
		// P P P P   P P P
		//
		//         P
		// P             Q
		//           P
		//   P P P P   P P
		// R N B Q K B N R
		// t.Logf("\n%s", brd)
		try(brd, Point{6, 6}, Point{5, 6})

		if brd.FinalCheckmate(1) {
			t.Fatalf("somehow final checkmate")
		}
	}
	{ // fool's pawn
		brd := NewBoard()

		try(brd, Point{6, 5}, Point{5, 5})
		try(brd, Point{1, 4}, Point{3, 4})
		try(brd, Point{6, 6}, Point{4, 6})
		try(brd, Point{0, 3}, Point{4, 7})

		// R N B   K B N R
		// P P P P   P P P
		//
		//         P
		//             P Q
		//           P
		// P P P P P     P
		// R N B Q K B N R
		// t.Logf("\n%s", brd)

		if !brd.Checkmate(1) {
			t.Fatalf("no checkmate")
		}
		if !brd.FinalCheckmate(1) {
			t.Fatalf("no final checkmate")
		}
	}
	{ // scholar's mate
		brd := NewBoard()

		try(brd, Point{6, 4}, Point{4, 4})
		try(brd, Point{1, 4}, Point{3, 4})
		try(brd, Point{6, 2}, Point{4, 2})
		try(brd, Point{0, 5}, Point{3, 2})
		try(brd, Point{7, 1}, Point{5, 2})
		try(brd, Point{0, 3}, Point{4, 7})
		try(brd, Point{7, 6}, Point{5, 5})
		try(brd, Point{4, 7}, Point{6, 5})

		//  R N B   K   N R
		//  P P P P   P P P
		//
		//      B   P
		//      P   P
		//      N     N
		//  P P   P   Q P P
		//  R   B Q K B   R
		//
		//  t.Logf("\n%s", brd)
		if !brd.Checkmate(1) {
			t.Fatalf("no checkmate")
		}

		if !brd.FinalCheckmate(1) {
			t.Fatalf("no final checkmate")
		}
	}
	{ // bird's opening
		brd := NewBoard()

		try(brd, Point{6, 5}, Point{4, 5})
		try(brd, Point{1, 4}, Point{3, 4})
		try(brd, Point{4, 5}, Point{3, 4})
		try(brd, Point{1, 3}, Point{2, 3})
		try(brd, Point{3, 4}, Point{2, 3})
		try(brd, Point{0, 5}, Point{2, 3})
		try(brd, Point{7, 1}, Point{5, 2})
		try(brd, Point{0, 3}, Point{4, 7})

		// R N B   K   N R
		// P P P     P P P
		//       B
		//
		//               Q
		//     N
		// P P P P P   P P
		// R   B Q K B N R
		// t.Logf("\n%s", brd)

		// this one requires seeing 4 moves ahead
		// technically it's not a final checkmate, but if the player knows this strat then it is.
		// this could be a future enhancement, for the server to memorize all these strategies and stuff, but for now it should fail.

		if !brd.Checkmate(1) {
			t.Fatalf("no checkmate")
		}

		if brd.FinalCheckmate(1) {
			t.Fatalf("no final checkmate")
		}
	}
	{ // italian game smothered mate
		brd := NewBoard()

		try(brd, Point{6, 4}, Point{4, 4})
		try(brd, Point{1, 4}, Point{3, 4})
		try(brd, Point{7, 6}, Point{5, 5})
		try(brd, Point{0, 1}, Point{2, 2})
		try(brd, Point{7, 5}, Point{4, 2})
		try(brd, Point{2, 2}, Point{4, 3})
		try(brd, Point{5, 5}, Point{3, 4})
		try(brd, Point{0, 3}, Point{3, 6})
		try(brd, Point{3, 4}, Point{1, 5})
		try(brd, Point{3, 6}, Point{6, 6})
		try(brd, Point{7, 7}, Point{7, 5})
		try(brd, Point{6, 6}, Point{4, 4})
		try(brd, Point{4, 2}, Point{6, 4})
		try(brd, Point{4, 3}, Point{5, 5})

		// R   B   K B N R
		// P P P P   N P P
		//
		//
		//         Q
		//           N
		// P P P P B P   P
		// R N B Q K R
		// t.Logf("\n%s", brd)

		// this once also requires seeing more than 1 move ahead.

		if !brd.Checkmate(1) {
			t.Fatalf("no checkmate")
		}

		if brd.FinalCheckmate(1) {
			t.Fatalf("no final checkmate")
		}
	}
	// TODO: implement the rest

	// this section tests previous bugs
	// and if they still exist
	{
		// queen at 6,5
		// and king at 6,4
		// but somehow it's final

		brd := NewBoard()
		empty(brd, Point{1, 4})
		empty(brd, Point{6, 4})
		empty(brd, Point{6, 5})

		queen := brd.Get(Point{0, 3})
		brd.Move(queen, Point{4, 7})
		king := brd.Get(Point{7, 4})
		brd.Move(king, Point{6, 4})
		brd.Move(queen, Point{6, 5})

		// R N B   K B N R
		// P P P P   P P P
		//
		//
		//
		//
		// P P P P K Q P P
		// R N B Q   B N R
		// t.Logf("\n%s", brd.String())

		if brd.FinalCheckmate(king.Player) {
			t.Fatalf("queen at 6:5, king at 6:4. but somehow final checkmate")
		}
	}
	{ // somehow king can end the game by checkmating king
		// even if the king can escape
		brd := NewBoard()

		empty(brd, Point{0, 1})
		empty(brd, Point{0, 2})
		empty(brd, Point{0, 3})
		empty(brd, Point{1, 3})
		empty(brd, Point{1, 4})

		try(brd, Point{0, 4}, Point{1, 4})
		try(brd, Point{1, 4}, Point{2, 4})

		empty(brd, Point{7, 3})
		empty(brd, Point{7, 5})
		//empty(brd, Point{7, 6})
		empty(brd, Point{6, 4})

		brd.Set(&Piece{
			Pos:    Point{2, 2},
			T:      PawnF,
			Player: 1,
		})

		try(brd, Point{7, 6}, Point{5, 5})
		try(brd, Point{5, 5}, Point{3, 6})

		if brd.FinalCheckmate(2) {
			t.Fatalf("knight can final checkmate escapable king")
		}

		t.Logf("\n%s", brd.String())
	}
}
*/
