package board

import (
	"testing"
)

// placement test
func TestNewBoard(t *testing.T) {
	u := [2][8]uint8{
		{Rook, Knight, Bishop, King, Queen, Bishop, Knight, Rook},
		{PawnB, PawnB, PawnB, PawnB, PawnB, PawnB, PawnB, PawnB},
	}

	b := NewBoard()
	for x := 0; x < 2; x++ {
		for y := 0; y < len(b); y++ {
			if b[x][y].T != u[x][y] {
				t.Fatalf("Top rows are not setup properly: [%d, %d]", x, y)
			}
		}
	}

	u[0], u[1] = u[1], u[0]
	for x := 0; x < 2; x++ {
		for y := 0; y < len(b); y++ {
			v := u[x][y]
			if v == PawnB {
				v = PawnF
			}

			if b[x+6][y].T != v {
				t.Fatalf("Bottom rows are not setup properly: [%d, %d]", x, y)
			}
		}
	}
}

// somethings do not need tests
//func TestBoardString(t *testing.T) {
//	str := `R N B K Q B N R
//P P P P P P P P
//
//
//
//
//P P P P P P P P
//R N B K Q B N R`
//
//	x, y := NewBoard().String(), str
//
//	fmt.Println(x)
//	fmt.Println(y)
//	fmt.Println(strings.Compare(x, y))
//
//	if x != y {
//		t.Fatalf("strings are layoutted wrong")
//	}
//}

func TestBoardMove(t *testing.T) {
	b := NewBoard()

	x := b[1][3]

	if !b.Move(x, Point{
		X: 3,
		Y: 3,
	}) {
		t.Fatalf("CanGo failed")
	}

	if b[3][3].T != PawnB {
		t.Fatalf("Pawn didn't move")
	}
	t.Log("\n" + b.String())

}
