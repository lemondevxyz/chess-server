package rest

import (
	"testing"
	"time"

	"github.com/toms1441/chess-server/internal/board"
	"github.com/toms1441/chess-server/internal/model"
	"github.com/toms1441/chess-server/internal/model/local"
)

var _watchableId = ""

func TestWatchableAdd(t *testing.T) {
	oldcache := string(watchable.cache)
	_watchableId = watchable.Add(model.Watchable{
		P1:  local.NewUser(),
		P2:  local.NewUser(),
		Brd: *board.NewBoard(),
	})

	if len(_watchableId) == 0 {
		t.Fatalf("watchable.Add: empty id")
	}

	time.Sleep(time.Millisecond * 10)
	newcache := string(watchable.cache)

	if oldcache == newcache {
		t.Fatalf("does not rebuild after adding a new model")
	}
}

func TestWatchableRebuild(t *testing.T) {
	oldtime := watchable.last
	watchable.Rebuild(true)
	newtime := watchable.last

	if oldtime.Equal(newtime) {
		t.Fatalf("last cache time does not change when rebuilding")
	}
}

func TestWatchableShouldRebuild(t *testing.T) {
	cacheDuration = time.Millisecond

	time.Sleep(time.Millisecond * 5)

	if !watchable.ShouldRebuild() {
		t.Fatalf("should rebuild has an error with time")
	}
}

func TestWatchableRemove(t *testing.T) {

	oldcache := string(watchable.cache)
	watchable.Rm(_watchableId)
	// cause goroutine
	time.Sleep(time.Millisecond * 10)
	newcache := string(watchable.cache)

	if oldcache == newcache {
		t.Fatalf("cache does not change when rm")
	}
}
