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
				if !b.data[v.nm1].Pos.Equal(v.dst) || b.data[v.nm2].T != Empty {
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

/*
// check if some pieces can move over pieces that are in the way...
func TestBoardMoveInTheWay(t *testing.T) {
	b := NewBoard()

	// try moving the rook through the pawn
	p := b.Get(Point{7, 0})
	// rook should not be able to move one bit if it's in the start
	for _, v := range p.Possib() {
		if b.Move(p, v) {
			t.Logf("%s", v.String())
			t.Logf("\n%s", b.String())
			t.Fatalf("rook can move over other pieces")
		}
	}

	// try moving knight through pawn
	p = b.Get(Point{7, 1})
	for _, v := range p.Possib() {

		b = NewBoard()
		p = b.Get(Point{7, 1})
		if !v.Valid() {
			continue
		}

		o := b.Get(v)
		want := true
		if o != nil {
			if o.Player == p.Player {
				want = false
			}
		}

		have := b.Move(p, v)
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

	// try moving bishop through pawn
	p = b.Get(Point{7, 2})
	for _, v := range p.Possib() {
		if b.Move(p, v) {
			t.Logf("%s", v.String())
			t.Logf("\n%s", b.String())
			t.Fatalf("bishop can move over other pieces")
		}
	}

	//p = b.Get(Point{7, })
	// try moving king through other pieces
	p = b.Get(Point{7, 3})
	for _, v := range p.Possib() {
		if b.Move(p, v) {
			t.Logf("%s", v.String())
			t.Logf("\n%s", b.String())
			t.Fatalf("king can move over other pieces")
		}
	}

	p = b.Get(Point{7, 4})
	for _, v := range p.Possib() {
		if b.Move(p, v) {
			t.Logf("%s", v.String())
			t.Logf("\n%s", b.String())
			t.Fatalf("queen can move over other pieces")
		}
	}

	// try moving king to a threatened
	{
		brd := NewBoard()
		for i := 0; i < 8; i++ {
			brd.Set(&Piece{Pos: Point{1, i}, T: Empty})
			brd.Set(&Piece{Pos: Point{6, i}, T: Empty})
		}

		pec := brd.Get(Point{7, 4})
		brd.Set(&Piece{Pos: Point{7, 4}, T: Empty})
		brd.Set(&Piece{Pos: Point{6, 3}, T: King})
		t.Log(brd.Possib(pec))

		t.Logf("\n%s", brd)
	}
}

// while above makes sure that no "special" piece can move over from it's starting position, this one tests past bugs.
func TestBoardPieceBugInTheWay(t *testing.T) {
	{ // bishop shouldn't skip over enemy pawn and killing knight
		brd := NewBoard()
		pec := brd.Get(Point{6, 4})

		pec.T = Empty
		pos := Point{4, 2}

		brd.Set(pec)
		brd.Set(&Piece{
			T:   Bishop,
			Pos: pos,
		})

		t.Logf("bishop possible moves: %s", pos.String())
		t.Log(brd.Possib(brd.Get(Point{4, 2})))
		if brd.Move(brd.Get(Point{4, 2}), Point{0, 6}) {
			t.Fatalf("bishop can skip enemy pawn and kill knight")
		}
	}
	{ // knight cannot override nearby pawn, but it's in the possible moves
		brd := NewBoard()

		pos := Point{7, 6}
		pec := brd.Get(pos)

		t.Logf("knight possible moves: %s", pos.String())
		possib := brd.Possib(pec)
		t.Log(possib)

		if possib.In(Point{6, 4}) {
			t.Fatalf("knight possible moves is killing nearby pawn")
		}
	}
	{ // pawn possiblity needs to include killable pieces
		brd := NewBoard()

		pos := Point{6, 4}
		pec := brd.Get(pos)

		pos.X -= 2
		if !brd.Move(pec, pos) {
			t.Fatalf("major fault within move")
		}

		pos = Point{1, 3}
		pec = brd.Get(pos)

		pos.X += 2
		if !brd.Move(pec, pos) {
			t.Fatalf("major fault within move")
		}

		pos = Point{4, 4}
		pec = brd.Get(pos)
		if !brd.Possib(pec).In(Point{3, 3}) {
			t.Fatalf("pawn does not include killable pieces")
		}
	}
	{ // queen possible moves should not include it's fellow allies
		brd := NewBoard()

		pos := Point{7, 4}
		pec := brd.Get(pos)

		sp := Points{
			{7, 3},
			{7, 5},
			{6, 4},
			{6, 3},
			{6, 5},
		}

		ps := brd.Possib(pec)
		for _, v := range sp {
			if ps.In(v) {
				t.Fatalf("queen possible moves includes it's fellow allies. point: %s", v.String())
			}
		}
	}
	{ // pawn backward shouldn't have X-2 at X = 1, pawn forward shouldn't have X+2 at X = 6
		brd := NewBoard()

		pawnf := &Piece{
			Pos:    Point{1, 1},
			T:      PawnF,
			Player: 1,
		}
		pawnb := &Piece{
			Pos:    Point{6, 1},
			T:      PawnB,
			Player: 2,
		}

		brd.Set(pawnb)
		brd.Set(pawnf)

		if brd.Possib(pawnb).In(Point{4, 1}) {
			t.Fatalf("backward pawn can go to 4, 1")
		}
		if brd.Possib(pawnf).In(Point{3, 1}) {
			t.Fatalf("backward pawn can go to 3, 1")
		}
	}
}

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
