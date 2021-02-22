package main

import (
	"fmt"

	"github.com/toms1441/chess-server/internal/board"
)

/*
func bishop(p board.Point) []board.Point {
	x, y := p.X, p.Y
	orix, oriy := x, y

	diff := 0
	if x > y {
		diff = 8 - x
	} else if x < y {
		diff = 8 - y
	} else if x == y {
		diff = 8 - x
	}

	x += diff
	y += diff

	sl := []board.Point{}

	sl = append(sl, board.Point{orix, oriy})
	for i := 0; i < 8; i++ {
		x = x - 1
		y = y - 1

		if x == -1 || y == -1 || x == 8 || y == 8 {
			break
		}

		if x == orix && y == oriy {
			continue
		}

		sl = append(sl, board.Point{x, y})
	}

	x = 0
	y = oriy + diff
	sl = append(sl, board.Point{x, y})

	for i := 0; i < 8; i++ {
		y--
		x++

		if x == -1 || y == -1 || x == 8 || y == 8 {
			break
		}

		if x == orix && y == oriy {
			continue
		}

		sl = append(sl, board.Point{x, y})
	}

	return sl
}
*/

func knight(p board.Point) []board.Point {
	x, y := p.X, p.Y
	ps := []board.Point{
		{x + 2, y + 1},
		{x + 2, y - 1},

		{x - 2, y + 1},
		{x - 2, y - 1},

		{x + 1, y + 2},
		{x + 1, y - 2},

		{x - 1, y + 2},
		{x - 1, y - 2},
	}

	for i := len(ps) - 1; i >= 0; i-- {
		v := ps[i]

		if v.X < 0 || v.X > 7 || v.Y < 0 || v.Y > 7 {
			ps = append(ps[:i], ps[i+1:]...)
		}
	}

	return ps
}

func main() {
	b := board.NewBoard()
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			b.Set(&board.Piece{
				T: board.Bishop,
				X: i,
				Y: j,
			})
		}
	}

	p := board.Point{1, 2}
	b.Set(&board.Piece{X: p.X, Y: p.Y, T: board.Knight})

	sl := knight(p)
	for _, v := range sl {
		fmt.Println(v)
		b.Set(&board.Piece{
			T: board.King,
			X: v.X,
			Y: v.Y,
		})
	}

	fmt.Println(b)
}
