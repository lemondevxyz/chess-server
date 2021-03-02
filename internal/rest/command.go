package rest

import (
	"net/http"

	"github.com/toms1441/chess-server/internal/game"
	"github.com/toms1441/chess-server/internal/order"
)

func CmdHandler(w http.ResponseWriter, r *http.Request) {
	u, err := GetUser(r)
	if err != nil {
		RespondError(w, http.StatusUnauthorized, err)
		return
	}
	cl := u.Client()

	cmd := order.Order{}
	err = BindJSON(r, &cmd)
	if err != nil {
		RespondError(w, http.StatusBadRequest, err)
		return
	}

	g := cl.Game()
	if g == nil {
		RespondError(w, http.StatusNotFound, game.ErrGameNil)
		return
	}

	err = cl.Do(cmd)
	if err != nil {
		RespondError(w, http.StatusBadRequest, err)
		return
	}

	RespondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
	})
	return
}
