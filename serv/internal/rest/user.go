package rest

import (
	"net/http"
	"strings"

	"github.com/kjk/betterguid"
	"github.com/toms1441/chess/serv/internal/game"
)

type User struct {
	Token string `validate:"required" json:"string"`
}

var users = map[string]*game.Client{}

func AddClient(c *game.Client) User {
	id := betterguid.New()
	users[id] = c

	return User{
		Token: id,
	}
}

func GetUser(r *http.Request) (*game.Client, error) {
	str := r.Header.Get("Authorization")
	str = strings.ReplaceAll(str, "Bearer ", "")

	cl, ok := users[str]
	if !ok {
		return nil, game.ErrClientNil
	}

	return cl, nil
}

func (u User) Client() *game.Client {
	return users[u.Token]
}

func (u User) Delete() {
	delete(users, u.Token)
}
