package rest

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/kjk/betterguid"
	"github.com/toms1441/chess-server/internal/model"
)

var cacheDuration = time.Minute

type cacheWatchable struct {
	mtx   sync.Mutex
	slice map[string]model.Watchable
	cache json.RawMessage
	last  time.Time
}

func (c *cacheWatchable) Add(m model.Watchable) string {
	id := betterguid.New()

	watchable.mtx.Lock()
	c.slice[id] = m
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
	slice: map[string]model.Watchable{},
	cache: json.RawMessage{},
	last:  time.Now().UTC().Add(cacheDuration * -1),
}

func WatchableHandler(w http.ResponseWriter, r *http.Request) {
	watchable.Rebuild(false)

	BindJSON(r, watchable.cache)
}
