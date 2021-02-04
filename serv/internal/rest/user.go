package rest

import (
	"github.com/kjk/betterguid"
	"github.com/toms1441/chess/serv/internal/game"
)

type User struct {
	Token string `validate:"required" json:"string"`
}

var users = map[string]*game.Client{}

func addClient(c *game.Client) User {
	id := betterguid.New()
	users[id] = c

	return User{
		Token: id,
	}
}

func (u User) Client() *game.Client {
	return users[u.Token]
}

func (u User) Delete() {
	delete(users, u.Token)
}
