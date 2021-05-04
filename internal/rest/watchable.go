package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/kjk/betterguid"
	"github.com/toms1441/chess-server/internal/game"
	"github.com/toms1441/chess-server/internal/model"
)

var cacheDuration = time.Minute

type watchableModel struct {
	p1 model.Profile
	p2 model.Profile
	gm *game.Game
}

func (w *watchableModel) MarshalJSON() ([]byte, error) {
	return json.Marshal(model.Watchable{
		P1:  w.p1,
		P2:  w.p2,
		Brd: w.gm.Board(),
	})
}

type cacheWatchable struct {
	mtx   sync.Mutex
	slice map[string]*watchableModel
	cache json.RawMessage
	last  time.Time
}

func (c *cacheWatchable) Add(m watchableModel) string {
	id := betterguid.New()

	watchable.mtx.Lock()
	c.slice[id] = &m
	watchable.mtx.Unlock()

	go c.Rebuild(true)

	return id
}

func (c *cacheWatchable) Rm(id string) {
	watchable.mtx.Lock()
	delete(c.slice, id)
	watchable.mtx.Unlock()

	go c.Rebuild(true)
}

func (c *cacheWatchable) Rebuild(force bool) {
	if force || c.ShouldRebuild() {
		body, _ := json.Marshal(watchable.slice)

		watchable.mtx.Lock()
		watchable.cache = body
		watchable.last = time.Now().UTC()
		watchable.mtx.Unlock()
	}
}

func (c *cacheWatchable) ShouldRebuild() bool {
	if time.Now().UTC().Sub(watchable.last) >= cacheDuration {
		return true
	}

	return false
}

var watchable = cacheWatchable{
	slice: map[string]*watchableModel{},
	cache: json.RawMessage{},
	last:  time.Now().UTC().Add(cacheDuration * -1),
}

func WatchableHandler(w http.ResponseWriter, r *http.Request) {
	_, err := GetUser(r)
	if err != nil {
		RespondError(w, http.StatusUnauthorized, fmt.Errorf("you must be logged in to view watchable games"))
		return
	}

	watchable.Rebuild(false)

	w.WriteHeader(http.StatusOK)
	w.Write(watchable.cache)
}

func WatchHandler(w http.ResponseWriter, r *http.Request) {
	u, err := GetUser(r)
	if err != nil {
		RespondError(w, http.StatusUnauthorized, fmt.Errorf("you must be logged in to view watchable games"))
		return
	}

	generic := model.Generic{}

	err = BindJSON(r, &generic)
	if err != nil {
		RespondError(w, http.StatusBadRequest, err)
		return
	}

	sl, ok := watchable.slice[generic.ID]
	if !ok {
		RespondError(w, http.StatusNotFound, fmt.Errorf("no watchable game has that id"))
		return
	}

	g := sl.gm

	g.AddSpectator(u.Client())

	w.WriteHeader(http.StatusOK)
	w.Write(nil)
}

func LeaveHandler(w http.ResponseWriter, r *http.Request) {
	u, err := GetUser(r)
	if err != nil {
		RespondError(w, http.StatusUnauthorized, fmt.Errorf("you must be logged in to view watchable games"))
		return
	}

	cl := u.Client()
	if cl == nil {
		RespondError(w, http.StatusUnauthorized, game.ErrClientNil)
		return
	}

	g := u.Client().Game()
	if g == nil {
		RespondError(w, http.StatusUnauthorized, game.ErrGameNil)
		return
	}

	g.RmSpectator(u.Client())

	w.WriteHeader(http.StatusOK)
	w.Write(nil)
}
