package main

import "fmt"

type Piece uint8

func (p Piece) String() string {
	v := pieceString[p]
	if len(v) > 0 {
		return v
	} else {
		return fmt.Sprintf("%d", p)
	}
}

func (p Piece) SmollString() string {
	v := pieceChar[p]
	if len(v) > 0 {
		return v
	} else {
		return fmt.Sprintf("%d", p)
	}
}

const (
	Empty Piece = iota
	Pawn
	Bishop
	Knight
	Rook
	Queen
	King
)

var (
	pieceString = map[Piece]string{
		Empty:  "Empty",
		Pawn:   "Pawn",
		Bishop: "Bishop",
		Knight: "Knight",
		Rook:   "Rook",
		Queen:  "Queen",
		King:   "King",
	}
	pieceChar = map[Piece]string{
		Empty:  " ",
		Pawn:   "P",
		Bishop: "B",
		Knight: "T",
		Rook:   "R",
		Queen:  "Q",
		King:   "K",
	}
)

type PlayerPiece struct {
	Piece  Piece
	Player int
	X      int
	Y      int
}

func (pp PlayerPiece) String() string {
	return pp.Piece.String()
}

type Board [8][8]PlayerPiece

func (b Board) String() (str string) {
	for k, p := range b {
		if k != 0 {
			str += "\n"
		}

		for _, v := range p {
			str += v.Piece.SmollString() + " "
		}
	}

	return str
}

func newBoard() Board {

	b := Board{}
	row := [2][8]Piece{
		{
			Rook,
			Knight,
			Bishop,
			King,
			Queen,
			Bishop,
			Knight,
			Rook,
		},
		{
			Pawn,
			Pawn,
			Pawn,
			Pawn,
			Pawn,
			Pawn,
			Pawn,
			Pawn,
		},
	}

	for k, s := range row {
		for l, v := range s {
			b[k][l] = PlayerPiece{
				Piece:  v,
				Player: 1,
			}
		}
	}

	row[0], row[1] = row[1], row[0]
	for k, s := range row {
		for l, v := range s {
			b[k+6][l] = PlayerPiece{
				Piece:  v,
				Player: 1,
			}
		}
	}

	return b

}

func main() {
	fmt.Println(newBoard())
}
