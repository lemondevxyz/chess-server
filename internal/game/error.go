package game

import "errors"

var (
	ErrGameNil          = errors.New("game is nil")
	ErrPieceNil         = errors.New("piece is nil")
	ErrIllegalMove      = errors.New("illegal move")
	ErrIllegalPromotion = errors.New("illegal promotion")
)
