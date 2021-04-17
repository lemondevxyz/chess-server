package board

import (
	"errors"
)

var (
	ErrInvalidID    = errors.New("id is invalid")
	ErrInvalidPoint = errors.New("invalid point")
	ErrEmptyPiece   = errors.New("piece is empty")
)
