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
	id := u.Profile.ID + "_" + u.Profile.Platform

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

	// lamo you though you were slick
	if vs == u {
		return game.ErrClientNil
	}

	param := model.InviteOrder{
		Profile: vs.Profile,
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

	jsu, err := json.Marshal(model.GameOrder{
		P1:      u.Client().P1(),
		Profile: vs.Profile,
		Brd:     b,
	})
	if err != nil {
		return cancel(err)
	}
	jsv, err := json.Marshal(model.GameOrder{
		P1:  vs.Client().P1(),
		Brd: b,
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

	id := watchable.Add(model.Watchable{
		P1:  u.Profile,
		P2:  vs.Profile,
		Brd: g.Board(),
	})

	go func() {
		<-g.ListenForDone()
		watchable.Rm(id)
	}()

	u.mtx.Lock()
	u.invite = map[string]*User{}
	u.mtx.Unlock()

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

	err = u.AcceptInvite(inv.Profile.ID + "_" + inv.Profile.Platform)
	if err != nil {
		RespondError(w, http.StatusBadRequest, err)
		return
	}

	RespondJSON(w, http.StatusOK, nil)
}
