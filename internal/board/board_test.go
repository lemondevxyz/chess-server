package board

import (
	"testing"
	"time"
)

// TODO: 	refactor this whole test
// 			possibly implement generic functions, like empty and try.
// 			a test should not be 500+ SLOC

// placement test
func TestNewBoard(t *testing.T) {
	u := [32]uint8{
		Rook, Knight, Bishop, Queen, King, Bishop, Knight, Rook,
		Pawn, Pawn, Pawn, Pawn, Pawn, Pawn, Pawn, Pawn,
		Pawn, Pawn, Pawn, Pawn, Pawn, Pawn, Pawn, Pawn,
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
		if b.data[x].Kind != u[x] {
			t.Fatalf("rows(types) are not setup properly: %d | want: %d, have: %d", x, u[x], b.data[x].Kind)
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

	ok := make(chan bool)

	x, y := false, false
	// func(p Piece, src Point, dst Point, ret bool)
	b.Listen(func(_ int8, _ Piece, _, _ Point) {
		t.Log("ay")
		ok <- true
		t.Log("ay again")
	})

	go b.Move(9, Point{1, 2})

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

	// lamo
	if !b.data[11].Pos.Equal(Point{3, 3}) {
		t.Fatalf("Pawn didn't move")
	}
}

func TestBoardMovePawn(t *testing.T) {

	b := NewBoard()

	// dark 1
	// light 1
	d1 := int8(11)
	l1 := int8(20)
	l2 := int8(21)

	t.Log(b.data[d1], b.data[l1], b.data[l2])

	b.Move(d1, Point{3, 3})
	b.Move(l1, Point{4, 4})
	b.Move(l2, Point{5, 5})

	t.Log(b.data[d1], b.data[l1], b.data[l2])

	if !b.Move(d1, Point{4, 4}) {
		t.Fatalf("backward pawn should kill other pawn")
	} else if !b.Move(l2, Point{4, 4}) {
		t.Log(b.data[d1], b.data[l1], b.data[l2])
		ps, _ := b.Possib(l2)
		t.Log(ps)
		t.Fatalf("forward pawn should kill other pawn")
	}

	d2 := int8(12)
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
		nm1 int8
		nm2 int8
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
	generic := func(id int8) {
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

	for i := int8(0); i < 16; i++ {
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

	id := int8(25)
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
			if cep.P1 == pec.P1 {
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

func TestBoardCheckmate(t *testing.T) {
	// testing top 10 fast checkmates: https://www.chess.com/article/view/fastest-chess-checkmates
	// black always wins
	empty := func(brd *Board, id int8) {
		if brd == nil {
			t.Fatalf("board is nil")
		}
		brd.Set(id, Point{-1, -1})
	}
	try := func(brd *Board, id int8, pnt Point) {
		pec, err := brd.GetByIndex(id)
		if err != nil {
			t.Fatalf("no piece at %d", id)
		}
		if !brd.Move(id, pnt) {
			t.Fatalf("perfectly legal move is failing. piece: %s | src: %s - dst: %s", pec.String(), pec.Pos.String(), pnt.String())
		}
	}
	{ // BUG: checkmate works through other pieces
		// const ids
		const ourpawn1 = 21
		const ourpawn2 = 22
		const theirpawn = 12
		const theirqueen = 3
		brd := NewBoard()

		try(brd, ourpawn1, Point{5, 5})
		try(brd, theirpawn, Point{4, 3})
		try(brd, ourpawn2, Point{6, 5})
		try(brd, theirqueen, Point{7, 4})
		// R N B   K B N R
		// P P P P   P P P
		//
		//         P
		//               Q
		//           P P
		// P P P P P     P
		// R N B Q K B N R
		// t.Logf("\n%s", brd)
		if brd.Checkmate(true) {
			t.Fatalf("bad checkmate, over other piece")
		}
	}
	{ // BUG: knight can checkmate escapable king
		// this basically tests out if the king can escape by itself, without the help of ally pieces...
		const ourking = 4
		const theirknight = 30
		brd := NewBoard()

		empty(brd, 1)  // Point{1, 0} Knight
		empty(brd, 2)  // Point{2, 0} Bishop
		empty(brd, 3)  // Point{3, 0} Queen
		empty(brd, 11) // Point{3, 1} PawnB
		empty(brd, 12) // Point{4, 1} PawnB

		try(brd, ourking, Point{4, 1})
		try(brd, ourking, Point{4, 2})

		empty(brd, 27) // Point{3, 7} Queen
		empty(brd, 29) // Point{5, 7} Bishop
		empty(brd, 20) // Point{4, 6} PawnF

		brd.Set(16, Point{2, 2}) // Point{0, 6} PawnF

		try(brd, theirknight, Point{5, 5})
		try(brd, theirknight, Point{6, 3})

		if brd.FinalCheckmate(false) {
			t.Logf("\n%s", brd)
			t.Fatalf("knight can final checkmate escapable king")
		}
	}
	{ // checkmate but not a final checkmate
		const ourpawn1 = 16
		const ourpawn2 = 21
		const ourpawn3 = 22
		const ourking = 28
		const theirqueen = 3
		const theirpawn = theirqueen + 8 + 1

		brd := NewBoard()

		try(brd, ourpawn1, Point{0, 4})
		try(brd, theirpawn, Point{4, 3})
		try(brd, ourpawn2, Point{5, 5})
		try(brd, theirqueen, Point{7, 4})

		// t.Logf("\n%s", brd)

		if !brd.Checkmate(true) {

			t.Logf("\n%s", brd.data[theirqueen])
			ps, _ := brd.Possib(theirqueen)
			t.Logf("possib: %s", ps)

			t.Logf("\n%s", brd.data[ourking])

			t.Fatalf("checkmate - want: true | have: false")
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
		try(brd, ourpawn3, Point{6, 5})

		if brd.FinalCheckmate(true) {
			t.Fatalf("final checkmate - want: false | have: true")
		}
	}
	{ // fool's pawn
		brd := NewBoard()

		try(brd, 21, Point{5, 5}) // Point{5, 6} PawnF
		try(brd, 12, Point{4, 3}) // Point{4, 1} PawnB
		try(brd, 22, Point{6, 4}) // Point{6, 6} PawnF
		try(brd, 3, Point{7, 4})  // Point{3, 0}

		if !brd.Checkmate(true) {
			t.Fatalf("no checkmate")
		}
		if !brd.FinalCheckmate(true) {
			t.Fatalf("no final checkmate")
		}
	}
	{ // scholar's mate
		brd := NewBoard()

		try(brd, 20, Point{4, 4}) // Point{4, 6} PawnF
		try(brd, 12, Point{4, 3}) // Point{4, 1} PawnB
		try(brd, 18, Point{2, 4}) // Point{2, 6} PawnF
		try(brd, 5, Point{2, 3})  // Point{5, 0} Bishop
		try(brd, 25, Point{2, 5}) // {1, 7} Knight
		try(brd, 3, Point{7, 4})
		try(brd, 30, Point{5, 5})
		try(brd, 3, Point{5, 6})
		/*
			try(brd, Point{0, 3}, Point{4, 7})
			try(brd, Point{7, 6}, Point{5, 5})
			try(brd, Point{4, 7}, Point{6, 5})
		*/

		//  R N B   K   N R
		//  P P P P   P P P
		//
		//      B   P
		//      P   P
		//      N     N
		//  P P   P   Q P P
		//  R   B Q K B   R
		//
		// t.Logf("\n%s", brd)
		if !brd.Checkmate(true) {
			t.Fatalf("no checkmate")
		}

		if !brd.FinalCheckmate(true) {
			t.Fatalf("no final checkmate")
		}
	}
}
