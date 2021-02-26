package rest

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/kjk/betterguid"
	"github.com/thanhpk/randstr"
	"github.com/toms1441/chess-server/internal/game"
)

type User struct {
	Token  string `validate:"required" json:"string"`
	invite map[string]*User
	cl     *game.Client
}

var users = map[string]*User{}

const (
	InviteLifespan = time.Second * 10
)

func AddClient(c *game.Client) *User {
	id := betterguid.New()
	us := &User{
		Token:  id,
		invite: map[string]*User{},
		cl:     c,
	}

	users[id] = us
	return us
}

func GetUser(r *http.Request) (*User, error) {
	str := r.Header.Get("Authorization")
	str = strings.ReplaceAll(str, "Bearer ", "")

	cl, ok := users[str]
	if !ok {
		return nil, game.ErrClientNil
	}

	return cl, nil
}

func (u *User) Client() *game.Client {
	return u.cl
}

func (u *User) Delete() {
	u.cl = nil
	u.invite = nil
	delete(users, u.Token)
}

func (u *User) Invite(tok string, lifespan time.Duration) error {
	// make sure panic don't happen
	if u.Client() == nil {
		return game.ErrClientNil
	}
	if u.Client().Game() != nil {
		return game.ErrGameIsNotNil
	}

	id := randstr.String(4)
	param := game.ModelUpdateInvite{
		ID: id,
	}
	body, err := json.Marshal(param)
	if err != nil {
		return err
	}

	vs := users[tok]
	if vs == nil || vs.Client() == nil {
		return game.ErrClientNil
	}
	if vs.Client().Game() != nil {
		return game.ErrGameIsNotNil
	}
	// u invited vs
	vs.invite[id] = u
	gu := game.Update{
		ID:   game.UpdateInvite,
		Data: body,
	}
	body, err = json.Marshal(gu)
	if err != nil {
		return err
	}

	// delete after X amount of time
	go func(u *User, id string) {
		<-time.After(lifespan)
		delete(u.invite, id)
	}(vs, id)

	vs.Client().W.Write(body)

	return nil
}

func (u User) AcceptInvite(tok string) error {
	x := u.invite[tok]
	if x == nil || x.Client() == nil {
		return game.ErrClientNil
	}
	if x.Client().Game() != nil || u.Client().Game() != nil {
		return game.ErrGameIsNotNil
	}

	_, err := game.NewGame(u.Client(), x.Client())
	if err != nil {
		return err
	}

	return nil
}
