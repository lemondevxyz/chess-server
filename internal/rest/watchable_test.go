package rest

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/toms1441/chess-server/internal/model"
	"github.com/toms1441/chess-server/internal/model/local"
)

var _watchableId = ""

var (
	specR, specW = io.Pipe()
	specCl, _    = AddClient(local.NewUser(), specW)
)

func TestWatchableAdd(t *testing.T) {

	go func() {
		<-read(rd2)
	}()
	us1.cl.LeaveGame()

	t.Log(len(watchable.slice))

	oldcache := string(watchable.cache)
	_watchableId = watchable.Add(watchableModel{
		p1: us1.Profile,
		p2: us2.Profile,
		gm: gGame,
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

func TestWatchableWatchHandler(t *testing.T) {
	go func() {
		<-read(rd2)
	}()

	err := us1.Invite(model.InviteOrder{
		Profile: us2.Profile,
	}, InviteLifespan)
	if err != nil {
		t.Fatalf("us1.Invite: %s", err.Error())
	}

	go func() {
		<-read(rd2)
		<-read(rd1)
		<-read(rd2)
		<-read(rd1)
	}()

	err = us2.AcceptInvite(us1.Profile.ID + "_" + us1.Profile.Platform)
	if err != nil {
		t.Fatalf("us2.Invite: %s", err.Error())
	}

	var id string
	for k := range watchable.slice {
		id = k
		break
	}
	if len(id) == 0 {
		t.Fatalf("empty id")
	}

	body, err := json.Marshal(model.Generic{
		ID: id,
	})
	if err != nil {
		t.Fatalf("json.Marshal: %s", err.Error())
	}

	resp := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "/", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("http.NewRequest: %s", err.Error())
	}

	hd := req.Header
	hd.Add("Authorization", "Bearer "+specCl.Token)

	req.Header = hd

	handle := http.HandlerFunc(WatchableJoinHandler)

	go func() {
		<-read(specR)
		<-read(specR)
	}()
	handle.ServeHTTP(resp, req)

	if resp.Result().StatusCode != 200 {
		t.Fatalf("%d: %s", resp.Result().StatusCode, resp.Body.String())
	}
}

func TestWatchableUpdate(t *testing.T) {
	/* just to make sure */
	go us1.cl.Game().SwitchTurn()

	<-read(rd2)
	<-read(rd1)

	x := model.Order{}
	err := json.Unmarshal(<-read(specR), &x)
	if err != nil {
		t.Fatalf("json.Unmarshal: %s", err.Error())
	}
}

func TestWatchableLeaveHandler(t *testing.T) {
	resp := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "/", nil)
	if err != nil {
		t.Fatalf("http.NewRequest: %s", err.Error())
	}

	hd := req.Header
	hd.Add("Authorization", "Bearer "+specCl.Token)

	req.Header = hd

	handle := http.HandlerFunc(WatchableLeaveHandler)
	handle.ServeHTTP(resp, req)

	if resp.Result().StatusCode != 200 {
		t.Fatalf("%d: %s", resp.Result().StatusCode, resp.Body.String())
	}

	go us1.cl.Game().SwitchTurn()

	<-read(rd2)
	<-read(rd1)

	select {
	case <-time.After(time.Millisecond * 25):
		break
	case <-read(specR):
		t.Fatalf("LeaveHandler does not work propely.")
	}
}
