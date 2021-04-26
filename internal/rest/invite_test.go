package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
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

	_, err := us1.Invite(model.InviteOrder{
		ID:       us2.Profile.ID,
		Platform: us2.Profile.Platform,
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
	err := us2.AcceptInvite("")
	if err == nil {
		t.Fatalf("us.AcceptInvite: invalid token does not return error")
	}
	// because net.Pipe is synchronous
	ch := make(chan error)

	go func() {
		_, err = us1.Invite(model.InviteOrder{
			ID:       us2.Profile.ID,
			Platform: us2.Profile.Platform,
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
		err := us2.AcceptInvite(inv.ID)
		if err != nil {
			ch <- err
			return
		}

		if us1.Client().Game() == nil || us2.Client().Game() == nil {
			ch <- fmt.Errorf("does not start a new game!")
		}
	}()

	x := <-read(rd2)
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

	if gm1.P1 == gm2.P1 {
		t.Log(gm1.P1, gm2.P1)
		t.Fatalf("player numbers are the same")
	}

}

func TestInviteHandler(t *testing.T) {
	go us1.cl.LeaveGame()
	<-read(rd2)
	//us2.cl.LeaveGame()

	x := model.InviteOrder{
		ID:       us2.Profile.ID,
		Platform: us2.Profile.Platform,
	}
	body, err := json.Marshal(x)
	if err != nil {
		t.Fatalf("json.Marshal: %s", err.Error())
	}

	resp := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "/", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", us1.Token))

	status := http.StatusOK
	handle := http.HandlerFunc(InviteHandler)

	ch := make(chan error)
	go func() {
		body := make([]byte, 1024)
		n, err := rd2.Read(body)
		if err != nil {
			ch <- fmt.Errorf("rd2.Read: %s", err.Error())
			return
		}
		body = body[:n]

		upd := model.Order{}
		err = json.Unmarshal(body, &upd)
		if err != nil {
			ch <- fmt.Errorf("json.Unmarshal: %s", err.Error())
			return
		}

		mod := model.InviteOrder{}
		err = json.Unmarshal(upd.Data, &mod)
		if err != nil {
			ch <- fmt.Errorf("json.Unmarshal: %s", err.Error())
			return
		}

		_inviteCode = mod.ID
		close(ch)
	}()

	handle.ServeHTTP(resp, req)
	rs := resp.Result()

	body, err = ioutil.ReadAll(rs.Body)
	if err != nil {
		t.Fatalf("ioutil.ReadAll: %s", err.Error())
	}

	if rs.StatusCode != status {
		t.Fatalf("%s %d", string(body), rs.StatusCode)
	}

	err = <-ch
	if err != nil {
		t.Fatalf(err.Error())
	}

}

func TestAcceptInviteHandler(t *testing.T) {
	i := model.InviteOrder{
		ID:       _inviteCode,
		Platform: ".",
	}

	body, err := json.Marshal(i)
	if err != nil {
		t.Fatalf("json.Marshal: %s", err.Error())
	}

	resp := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "/", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", us2.Token))

	status := http.StatusOK
	handle := http.HandlerFunc(AcceptInviteHandler)

	done := make(chan struct{})
	go func() {
		handle.ServeHTTP(resp, req)
		close(done)
	}()

	<-read(rd2)
	<-read(rd1)
	<-read(rd2)
	<-read(rd1)

	<-done

	rs := resp.Result()
	body, err = ioutil.ReadAll(rs.Body)
	if err != nil {
		t.Fatalf("ioutil.ReadAll: %s", err.Error())
	}
	if rs.StatusCode != status {
		t.Fatalf("%s %d", string(body), rs.StatusCode)
	}

	if us1.Client().Game() == nil || us2.Client().Game() == nil {
		t.Fatalf("AcceptInvite does not start a new game")
	}

	go func() {
		<-read(rd2)
		<-read(rd1)
	}()
	us1.Client().LeaveGame()

}
