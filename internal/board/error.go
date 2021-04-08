package board

import "fmt"

var (
	ErrInvalidPoint = fmt.Errorf("point out of bounds")
	ErrEmptyPiece   = fmt.Errorf("empty piece")
)
