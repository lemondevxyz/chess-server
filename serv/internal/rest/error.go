package rest

import "errors"

var (
	ErrClient = errors.New("invalid token")
	ErrGame   = errors.New("invalid game")
)
