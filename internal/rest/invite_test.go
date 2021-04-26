package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/toms1441/chess-server/internal/model"
)

var _inviteCode = ""

func TestUserInvite(t *testing.T) {
	go read(rd2)
	us1.cl.LeaveGame()

	rd1, wr1 = io.Pipe()
	rd2, wr2 = io.Pipe()

	us1.Client().W = wr1
	us2.Client().W = wr2

	const lifespan = time.Millisecond * 10
	// because net.Pipe is synchronous, we have to read through it. otherwise it would hang forever...
	go read(rd2)

	err := us1.Invite(model.InviteOrder{
		Profile: us2.Profile,
	}, lifespan)
	if err != nil {
		t.Fatalf("us.Invite: %s", err.Error())
	}
	if len(us2.invite) == 0 {
		t.Fatalf("vs invite map is empty")
	}

	<-time.After(lifespan * 2)
	if len(us2.invite) == 1 {
		t.Fatalf("vs lifespan does not work")
	}
}

func TestUserAcceptInvite(t *testing.T) {
	oldlen := len(watchable.slice)

	err := us2.AcceptInvite("")
	if err == nil {
		t.Fatalf("us.AcceptInvite: invalid token does not return error")
	}
	// because net.Pipe is synchronous
	ch := make(chan error)

	go func() {
		err = us1.Invite(model.InviteOrder{
			Profile: us2.Profile,
		}, InviteLifespan)
		if err != nil {
			ch <- fmt.Errorf("us.Invite: %s", err.Error())
			return
		}
		close(ch)
	}()

	update := model.Order{}
	err = json.Unmarshal(<-read(rd2), &update)
	if err != nil {
		t.Fatalf("json.Unmarshal: %s", err.Error())
	}

	err = <-ch
	if err != nil {
		t.Fatalf("rd2.Read: %s", err.Error())
	}

	inv := model.InviteOrder{}
	json.Unmarshal(update.Data, &inv)
	if err != nil {
		t.Fatalf("json.Unmarshal: %s", err.Error())
	}

	ch = make(chan error)
	go func() {
		err := us2.AcceptInvite(us1.Profile.ID + "_" + us1.Profile.Platform)
		if err != nil {
			ch <- err
			return
		}

		if us1.Client().Game() == nil || us2.Client().Game() == nil {
			ch <- fmt.Errorf("does not start a new game!")
			return
		}
		close(ch)
	}()

	x := []byte{}
	x = <-read(rd2)
	o := model.Order{}
	err = json.Unmarshal(x, &o)
	if err != nil {
		t.Log(string(x))
		t.Fatalf("json.Unmarshal: %s", err.Error())
	}
	gm2 := model.GameOrder{}
	err = json.Unmarshal(o.Data, &gm2)
	if err != nil {
		t.Fatalf("json.Unmarshal: %s", err.Error())
	}

	x = <-read(rd1)
	o = model.Order{}
	err = json.Unmarshal(x, &o)
	if err != nil {
		t.Fatalf("json.Unmarshal: %s", err.Error())
	}
	gm1 := model.GameOrder{}
	err = json.Unmarshal(o.Data, &gm1)
	if err != nil {
		t.Fatalf("json.Unmarshal: %s", err.Error())
	}

	<-read(rd2)
	<-read(rd1)

	err = <-ch
	if err != nil {
		t.Fatalf(err.Error())
	}

	if gm1.P1 == gm2.P1 {
		t.Log(gm1.P1, gm2.P1)
		t.Fatalf("player numbers are the same")
	}

	if oldlen == len(watchable.slice) {
		t.Fatalf("AcceptInvite does not add model.Watchable")
	}
}

func TestInviteHandler(t *testing.T) {
	go read(rd2)
	us1.cl.LeaveGame()

	marshal, _ := json.Marshal(model.InviteOrder{
		Profile: us2.Profile,
	})

	resp := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "/", bytes.NewBuffer(marshal))

	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", us1.Token))

	handle := http.HandlerFunc(InviteHandler)
	done := make(chan struct{})
	go func() {
		handle.ServeHTTP(resp, req)
		close(done)
	}()

	<-read(rd2)
	<-done

	res := resp.Result()
	if res.StatusCode != 200 {
		t.Log(resp.Body.String())
		t.Fatalf("http.StatusCode: %d", res.StatusCode)
	}
}

func TestAcceptInviteHandler(t *testing.T) {
	marshal, _ := json.Marshal(model.InviteOrder{
		Profile: us1.Profile,
	})

	resp := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "/", bytes.NewBuffer(marshal))

	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", us2.Token))

	handle := http.HandlerFunc(AcceptInviteHandler)
	done := make(chan struct{})
	go func() {
		handle.ServeHTTP(resp, req)
		close(done)
	}()

	//x := []byte{}
	<-read(rd2)
	//t.Log(string(x))
	<-read(rd1)
	//t.Log(string(x))
	<-read(rd2)
	//t.Log(string(x))
	<-read(rd1)
	//t.Log(string(x))

	<-done

	res := resp.Result()
	if res.StatusCode != 200 {
		t.Log(resp.Body.String())
		t.Fatalf("http.StatusCode: %d", res.StatusCode)
	}

}
