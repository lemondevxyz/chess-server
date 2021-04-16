package rest

import (
	"fmt"
	"net/http"

	"github.com/toms1441/chess-server/internal/board"
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

	g := cl.Game()
	if g == nil {
		RespondError(w, http.StatusNotFound, game.ErrGameNil)
		return
	}

	cmd := order.Order{}
	err = BindJSON(r, &cmd)
	if err != nil {
		RespondError(w, http.StatusBadRequest, err)
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

// since this is a specific handler and not via CmdHandler then there is no need to parse order.Order.
func PossibHandler(w http.ResponseWriter, r *http.Request) {
	u, err := GetUser(r)
	if err != nil {
		RespondError(w, http.StatusUnauthorized, err)
		return
	}

	cl := u.Client()
	if cl == nil {
		RespondError(w, http.StatusUnauthorized, game.ErrClientNil)
		return
	}

	gm := cl.Game()
	if gm == nil {
		RespondError(w, http.StatusUnauthorized, game.ErrGameNil)
		return
	}

	possib := order.PossibleModel{}
	err = BindJSON(r, &possib)
	if err != nil {
		RespondError(w, http.StatusBadRequest, err)
		return
	}

	if !board.BelongsTo(possib.ID, cl.P1()) {
		RespondError(w, http.StatusUnauthorized, fmt.Errorf("piece doesn't belong to you"))
		return
	}

	brd := gm.Board()

	points, _ := brd.Possib(int(possib.ID))

	possib = order.PossibleModel{}
	possib.Points = &points

	RespondJSON(w, http.StatusOK, possib)
}
