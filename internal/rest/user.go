package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/kjk/betterguid"
	"github.com/thanhpk/randstr"
	"github.com/toms1441/chess-server/internal/game"
	"github.com/toms1441/chess-server/internal/model"
)

type User struct {
	model.CredentialsOrder
	invite map[string]*User
	cl     *game.Client
}

var users = map[string]*User{}
var usermtx = sync.Mutex{}

var chanuser = make(chan *User)

const (
	InviteLifespan = time.Second * 30
)

func AllClients() map[string]*User {
	return users
}

func ClientChannel() chan *User {
	return chanuser
}

func AddClient(profile model.Profile, c *game.Client) *User {
	usermtx.Lock()

	id := betterguid.New()
	us := &User{
		invite: map[string]*User{},
		cl:     c,
	}

	us.Token = id
	us.CredentialsOrder.Profile = profile

	users[id] = us
	go func() {
		chanuser <- us
	}()

	usermtx.Unlock()

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

// GetAvaliableUsersHandler returns a list of public ids that are looking to play.
// TODO: exclude this user's public id from the returned value
func GetAvaliableUsersHandler(w http.ResponseWriter, r *http.Request) {
	_, err := GetUser(r)
	if err != nil {
		RespondError(w, http.StatusUnauthorized, err)
		return
	}

	ids := []string{}
	for _, v := range users {
		// lamo you can't invite yourself
		if v.Profile.GetPublicID() == v.Profile.GetPublicID() && v.Profile.GetPlatform() == v.Profile.GetPlatform() {
			continue
		}

		if v.Valid() {
			if v.Client().Game() == nil {
				ids = append(ids, v.Profile.GetPublicID())
			}
		}
	}

	RespondJSON(w, http.StatusOK, ids)
}

func (u *User) Client() *game.Client {
	return u.cl
}

func (u *User) Delete() {
	u.Profile = nil
	u.Token = ""
	if u.cl.Game() != nil {
		u.cl.LeaveGame()
	}
	u.cl = nil
	u.invite = nil
	usermtx.Lock()
	delete(users, u.Token)
	usermtx.Unlock()
}

func (u *User) Valid() bool {
	if u.Client() == nil {
		return false
	}
	x, ok := users[u.Token]
	if !ok {
		return false
	}
	if x != u {
		u.Delete()
		return false
	}

	return true
}

func (u *User) Invite(pubid string, lifespan time.Duration) (string, error) {
	// make sure panic don't happen
	if !u.Valid() {
		return "", game.ErrClientNil
	}
	if u.Client().Game() != nil {
		return "", game.ErrGameIsNotNil
	}

	var vs *User
	for _, v := range users {
		if v.Profile.GetPublicID() == pubid {
			vs = v
			break
		}
	}

	if vs == nil || !vs.Valid() {
		return "", game.ErrClientNil
	}
	if vs.Client().Game() != nil {
		return "", game.ErrGameIsNotNil
	}

	// lamo you though you were slick
	if vs == u {
		return "", game.ErrClientNil
	}

	id := randstr.String(4)
	param := model.InviteOrder{
		ID: id,
	}

	body, err := json.Marshal(param)
	if err != nil {
		return "", err
	}

	// u invited vs
	vs.invite[id] = u
	gu := model.Order{
		ID:   model.OrInvite,
		Data: body,
	}
	send, err := json.Marshal(gu)
	if err != nil {
		return "", err
	}

	// delete after X amount of time
	go func(vs *User, id string) {
		<-time.After(lifespan)
		delete(vs.invite, id)
	}(vs, id)

	vs.Client().W.Write(send)

	return id, nil
}

// AcceptInvite accepts the invite from the user.
func (u *User) AcceptInvite(tok string) error {
	vs, ok := u.invite[tok]
	if !ok {
		return ErrInvalidInvite
	}

	if vs == nil || vs.Client() == nil {
		return game.ErrClientNil
	}

	g, err := game.NewGame(u.Client(), vs.Client())
	if err != nil {
		return err
	}

	b := g.Board()

	cancel := func(err error) error {
		u.Client().LeaveGame()
		vs.Client().LeaveGame()

		return fmt.Errorf("%s | %w", err.Error(), ErrInternal)
	}

	jsu, err := json.Marshal(model.GameOrder{
		P1:    u.Client().P1(),
		Board: b,
	})
	if err != nil {
		return cancel(err)
	}
	jsv, err := json.Marshal(model.GameOrder{
		P1:    vs.Client().P1(),
		Board: b,
	})
	if err != nil {
		return cancel(err)
	}
	data, err := json.Marshal(model.Order{
		ID:   model.OrGame,
		Data: jsu,
	})
	if err != nil {
		return cancel(err)
	}
	u.Client().W.Write(data)
	data, err = json.Marshal(model.Order{
		ID:   model.OrGame,
		Data: jsv,
	})
	if err != nil {
		return cancel(err)
	}
	vs.Client().W.Write(data)

	g.SwitchTurn()

	u.invite = map[string]*User{}

	return nil
}
