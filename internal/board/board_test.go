package board

import (
	"testing"
	"time"
)

// placement test
func TestNewBoard(t *testing.T) {
	u := [2][8]uint8{
		{Rook, Knight, Bishop, King, Queen, Bishop, Knight, Rook},
		{PawnB, PawnB, PawnB, PawnB, PawnB, PawnB, PawnB, PawnB},
	}

	b := NewBoard()
	for x := 0; x < 2; x++ {
		for y := 0; y < len(b.data); y++ {
			if b.data[x][y].T != u[x][y] {
				t.Fatalf("Top rows are not setup properly: [%d, %d]", x, y)
			}
		}
	}

	u[0], u[1] = u[1], u[0]
	for x := 0; x < 2; x++ {
		for y := 0; y < len(b.data); y++ {
			v := u[x][y]
			if v == PawnB {
				v = PawnF
			}

			if b.data[x+6][y].T != v {
				t.Fatalf("Bottom rows are not setup properly: [%d, %d]", x, y)
			}
		}
	}

	t.Logf("\n%s", b.String())
}

/* somethings do not need tests
func TestBoardString(t *testing.T) {
	str := `R N B K Q B N R
P P P P P P P P




P P P P P P P P
R N B K Q B N R`

	x, y := NewBoard().String(), str

	fmt.Println(x)
	fmt.Println(y)
	fmt.Println(strings.Compare(x, y))

	if x != y {
		t.Fatalf("strings are layoutted wrong")
	}
}
*/

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
	b.Listen(func(_ *Piece, _, _ Point, ret bool) {
		if ret {
			valid <- true
		} else {
			invalid <- true
		}
	})

	b.Move(b.data[1][1], Point{3, 1})
	b.Move(b.data[3][1], Point{7, 1})

	select {
	case <-time.After(time.Millisecond * 20):
		t.Fatalf("Listen does not listen. pre: %t - post: %t", x, y)
	case <-ok:
		break
	}

}

func TestBoardSet(t *testing.T) {
	b := NewBoard()

	p := &Piece{
		X: 1,
		Y: 1,
		T: Bishop,
	}

	b.Set(p)
	if b.data[p.X][p.Y] != p {
		t.Fatalf("Set does not work")
	}
}

func TestBoardMove(t *testing.T) {
	b := NewBoard()

	x := b.data[1][3]

	if !b.Move(x, Point{
		X: 3,
		Y: 3,
	}) {
		t.Fatalf("CanGo failed")
	}

	if b.data[3][3].T != PawnB {
		t.Fatalf("Pawn didn't move")
	}
}

func TestBoardMovePawn(t *testing.T) {

	b := NewBoard()
	x := b.data[1][3]
	y := b.data[6][4]

	b.Move(x, Point{3, 3})
	b.Move(y, Point{4, 4})

	z := b.data[6][5]
	b.Move(z, Point{5, 5})

	if !b.Move(x, Point{4, 4}) {
		t.Fatalf("backward pawn should kill other pawn")
	} else if !b.Move(z, Point{4, 4}) {
		t.Fatalf("forward pawn should kill other pawn")
	}

	x = b.data[1][4]
	b.Move(x, Point{3, 4})
	if b.Move(x, Point{4, 4}) {
		t.Fatalf("pawn cannot takeout pawn in front of it")
	}

}

// probs overkill but why not
// basically test if other pieces can kill enemy pieces
func TestBoardMoveOthers(t *testing.T) {

	b := NewBoard()
	cs := []struct {
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
			src: Point{3, 3},
			dst: Point{4, 4},
			t:   Bishop,
		},
		{
			src: Point{3, 3},
			dst: Point{4, 5},
			t:   Knight,
		},
		{
			src: Point{3, 3},
			dst: Point{4, 3},
			t:   Rook,
		},
		{
			src: Point{3, 3},
			dst: Point{3, 4},
			t:   Rook,
		},
		{
			src: Point{3, 3},
			dst: Point{3, 4},
			t:   Queen,
		},
		{
			src: Point{3, 3},
			dst: Point{3, 4},
			t:   King,
		},
	}

	for _, v := range cs {
		b.data[v.src.X][v.src.Y] = &Piece{
			Player: 1,
			T:      v.t,
			X:      v.src.X,
			Y:      v.src.Y,
		}
		b.data[v.dst.X][v.dst.Y] = &Piece{
			Player: 2,
			T:      v.t,
			X:      v.dst.X,
			Y:      v.dst.Y,
		}

		x := b.data[v.src.X][v.src.Y]
		if !b.Move(b.data[v.src.X][v.src.Y], v.dst) {
			t.Fatalf("test cordinates are invalid. src: %d - dst: %d - type: %d", v.src, v.dst, v.t)
			return
		} else {
			if b.data[v.src.X][v.src.Y] != nil {
				t.Fatalf("move doesn't actually move")
			} else {
				if b.data[v.dst.X][v.dst.Y] != x {
					t.Fatalf("move doesn't replace enemy")
				}
			}
		}
	}

	b.data[0][0] = &Piece{
		Player: 1,
		T:      PawnF,
		X:      0,
		Y:      0,
	}
	b.data[1][1] = &Piece{
		Player: 1,
		T:      PawnB,
		X:      1,
		Y:      1,
	}

	if b.Move(b.data[0][0], Point{1, 1}) {
		t.Fatalf("ally killed")
	}

}
