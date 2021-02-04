package game

import "errors"

var (
	ErrClientNil        = errors.New("client is nil")
	ErrGameNil          = errors.New("game is nil")
	ErrCommandNil       = errors.New("command is nil")
	ErrPieceNil         = errors.New("piece is nil")
	ErrIllegalTurn      = errors.New("illegal turn")
	ErrIllegalMove      = errors.New("illegal move")
	ErrIllegalPromotion = errors.New("illegal promotion")
	ErrUpdateNil        = errors.New("update is nil")
	ErrUpdateTimeout    = errors.New("update write timeout")
	ErrUpdateParameter  = errors.New("update parameter is invalid")
)
