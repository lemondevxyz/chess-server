package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/thanhpk/randstr"
	"github.com/toms1441/chess-server/internal/game"
	"github.com/toms1441/chess-server/internal/model"
)

func (u *User) Invite(inv model.InviteOrder, lifespan time.Duration) (string, error) {
	// make sure panic don't happen
	if !u.Valid() {
		return "", game.ErrClientNil
	}
	if u.Client().Game() != nil {
		return "", game.ErrGameIsNotNil
	}

	var vs *User
	for _, v := range users {
		pro := v.Profile
		if pro.ID == inv.ID && pro.Platform == inv.Platform {
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
		ID:      id,
		Profile: vs.Profile,
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

	u.invite = map[string]*User{}

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
	if err != nil || len(inv.Platform) == 0 {
		RespondError(w, http.StatusBadRequest, err)
		return
	}

	_, err = u.Invite(inv, InviteLifespan)
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

	err = u.AcceptInvite(inv.ID)
	if err != nil {
		RespondError(w, http.StatusBadRequest, err)
		return
	}

	RespondJSON(w, http.StatusOK, nil)
}
