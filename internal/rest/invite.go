package rest

import (
	"net/http"

	"github.com/toms1441/chess-server/internal/model"
)

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
