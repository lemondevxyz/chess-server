package rest

import (
	"net/http"
)

type InviteModel struct {
	ID string `validate:"required" json:"id"`
}

func InviteHandler(w http.ResponseWriter, r *http.Request) {
	u, err := GetUser(r)
	if err != nil {
		RespondError(w, http.StatusUnauthorized, err)
		return
	}

	inv := InviteModel{}
	err = BindJSON(r, &inv)
	if err != nil {
		RespondError(w, http.StatusBadRequest, err)
		return
	}

	err = u.Invite(inv.ID, InviteLifespan)
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

	inv := InviteModel{}
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
