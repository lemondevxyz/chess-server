package rest

import "errors"

var (
	ErrInvalidInvite = errors.New("invalid invite")
	ErrInviteRate    = errors.New("already invited player. please wait")
)
