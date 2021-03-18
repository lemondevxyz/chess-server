package rest

import (
	"fmt"
	"net/http"

	"github.com/toms1441/chess-server/internal/order"
)

func InviteHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("user")
	u, err := GetUser(r)
	if err != nil {
		RespondError(w, http.StatusUnauthorized, err)
		return
	}
	fmt.Println("done user")

	inv := order.InviteModel{}
	err = BindJSON(r, &inv)
	if err != nil {
		RespondError(w, http.StatusBadRequest, err)
		return
	}

	fmt.Println("invite")
	_, err = u.Invite(inv.ID, InviteLifespan)
	if err != nil {
		RespondError(w, http.StatusBadRequest, err)
		return
	}
	fmt.Println("done invite")

	RespondJSON(w, http.StatusOK, nil)
}

func AcceptInviteHandler(w http.ResponseWriter, r *http.Request) {
	u, err := GetUser(r)
	if err != nil {
		RespondError(w, http.StatusUnauthorized, err)
		return
	}

	inv := order.InviteModel{}
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
