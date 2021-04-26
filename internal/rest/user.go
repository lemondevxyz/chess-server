package rest

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/kjk/betterguid"
	"github.com/toms1441/chess-server/internal/game"
	"github.com/toms1441/chess-server/internal/model"
)

type User struct {
	model.CredentialsOrder
	mtx    sync.Mutex
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

func AddClient(profile model.Profile, wc io.WriteCloser) (*User, error) {
	if wc == nil || !profile.Valid() {
		return nil, fmt.Errorf("one of the parameters is nil")
	}

	id := betterguid.New()
	us := &User{
		invite: map[string]*User{},
		cl: &game.Client{
			W: wc,
		},
	}

	us.Token = id
	us.Profile = profile

	usermtx.Lock()
	users[id] = us
	usermtx.Unlock()
	go func() {
		chanuser <- us
	}()

	return us, nil
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
func GetAvaliableUsersHandler(w http.ResponseWriter, r *http.Request) {
	u, err := GetUser(r)
	if err != nil {
		RespondError(w, http.StatusUnauthorized, err)
		return
	}

	ids := []json.RawMessage{}
	for _, v := range users {
		// lamo you can't invite yourself
		if v.Profile.ID == u.Profile.ID && v.Profile.Platform == u.Profile.Platform {
			continue
		}

		if v.Valid() {
			if v.Client().Game() == nil {
				body, err := json.Marshal(v.Profile)
				if err != nil {
					RespondError(w, http.StatusInternalServerError, err)
					return
				}
				ids = append(ids, body)
			}
		}
	}

	RespondJSON(w, http.StatusOK, ids)
}

func (u *User) Client() *game.Client {
	return u.cl
}

func (u *User) Delete() {
	id := u.Token

	u.Profile = model.Profile{}
	u.Token = ""
	if u.cl.Game() != nil {
		u.cl.LeaveGame()
	}
	u.cl = nil
	u.invite = nil

	usermtx.Lock()
	delete(users, id)
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
