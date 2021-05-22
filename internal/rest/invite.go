package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/toms1441/chess-server/internal/game"
	"github.com/toms1441/chess-server/internal/model"
)

func (u *User) Invite(inv model.InviteOrder, lifespan time.Duration) error {
	// make sure panic don't happen
	if !u.Valid() {
		return game.ErrClientNil
	}
	if u.Client().Game() != nil {
		return game.ErrGameIsNotNil
	}
	id := u.Profile.GetInviteID()

	var vs *User
	for _, v := range users {
		pro := v.Profile
		if v.Client().Game() == nil {
			if pro == inv.Profile {
				vs = v
				break
			}
		}
	}

	if vs == nil || !vs.Valid() {
		return game.ErrClientNil
	}

	// lamo you thought you were slick
	if vs == u {
		return game.ErrClientNil
	}

	param := model.InviteOrder{
		Profile: u.Profile,
	}

	body, err := json.Marshal(param)
	if err != nil {
		return err
	}

	// u invited vs
	vs.mtx.Lock()
	vs.invite[id] = u
	vs.mtx.Unlock()

	gu := model.Order{
		ID:   model.OrInvite,
		Data: body,
	}
	send, err := json.Marshal(gu)
	if err != nil {
		return err
	}

	// delete after X amount of time
	go func(vs *User, id string) {
		<-time.After(lifespan)
		vs.mtx.Lock()
		delete(vs.invite, id)
		vs.mtx.Unlock()
	}(vs, id)

	vs.Client().W.Write(send)

	return nil
}

// AcceptInvite accepts the invite from the user.
func (u *User) AcceptInvite(tok string) error {
	u.mtx.Lock()
	vs, ok := u.invite[tok]
	u.mtx.Unlock()
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
		go u.Client().LeaveGame()
		go vs.Client().LeaveGame()

		return fmt.Errorf("%s | %w", err.Error(), ErrInternal)
	}

	p1 := u.Client().P1()
	jsu, err := json.Marshal(model.GameOrder{
		P1:      &p1,
		Profile: &vs.Profile,
		Brd:     b,
	})
	if err != nil {
		return cancel(err)
	}
	p1 = vs.Client().P1()
	jsv, err := json.Marshal(model.GameOrder{
		P1:      &p1,
		Profile: &u.Profile,
		Brd:     b,
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

	id := watchable.Add(watchableModel{
		p1: u.Profile,
		p2: u.Profile,
		gm: g,
	})

	go func() {
		<-g.ListenForDone()
		watchable.Rm(id)
	}()

	u.mtx.Lock()
	u.invite = map[string]*User{}
	u.mtx.Unlock()

	vs.mtx.Lock()
	u.invite = map[string]*User{}
	vs.mtx.Unlock()

	return nil
}

func InviteHandler(w http.ResponseWriter, r *http.Request) {
	u, err := GetUser(r)
	if err != nil {
		RespondError(w, http.StatusUnauthorized, err)
		return
	}

	inv := model.InviteOrder{}
	err = BindJSON(r, &inv)
	if err != nil {
		RespondError(w, http.StatusBadRequest, err)
		return
	}

	err = u.Invite(inv, InviteLifespan)
	if err != nil {
		RespondError(w, http.StatusBadRequest, err)
		return
	}

	RespondJSON(w, http.StatusOK, nil)
}

func AcceptInviteHandler(w http.ResponseWriter, r *http.Request) {
	u, err := GetUser(r)
	if err != nil {
		RespondError(w, http.StatusUnauthorized, err)
		return
	}

	inv := model.InviteOrder{}

	err = BindJSON(r, &inv)
	if err != nil {
		RespondError(w, http.StatusBadRequest, err)
		return
	}

	err = u.AcceptInvite(inv.Profile.GetInviteID())
	if err != nil {
		RespondError(w, http.StatusBadRequest, err)
		return
	}

	RespondJSON(w, http.StatusOK, nil)
}
