package rest

import "github.com/toms1441/chess/internal/game"

type User struct {
	Token string `validate:"required" json:"token"`
}

var cls = map[string]*game.Client{}

func GetClient(u User) *game.Client {
	return cls[u.Token]
}
