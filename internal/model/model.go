package model

import "github.com/toms1441/chess-server/internal/board"

// Profile is a generic struct used to get information about the user.
type Profile struct {
	ID       string `json:"id" validate:"required"`
	Picture  string `json:"picture" validate:"required"`
	Username string `json:"username" validate:"required"`
	Platform string `json:"platform" validate:"required"`
}

func (p Profile) Valid() bool {
	return len(p.ID) > 0 && len(p.Picture) > 0 && len(p.Username) > 0
}

// Possible is a struct indicating which moves are legal.
type Possible struct {
	ID     int8          `json:"id,omitempty" validate:"required"` // [C]
	Points *board.Points `json:"points,omitempty"`                 // [U]
}

// Watchable is a game that could be spectated
type Watchable struct {
	P1  Profile      `json:"p1"`
	P2  Profile      `json:"p2"`
	Brd *board.Board `json:"brd"`
}
