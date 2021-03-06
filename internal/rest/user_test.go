package rest

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/toms1441/chess-server/internal/game"
	"github.com/toms1441/chess-server/internal/order"
)

var (
	us = &User{}
)

func TestNewUser(t *testing.T) {
	us = AddClient(cl1)
}

func TestGetUser(t *testing.T) {
	resp := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", us.Token))

	handle := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := GetUser(r)
		if err != nil {
			RespondError(w, http.StatusBadRequest, err)
			return
		}

		RespondJSON(w, http.StatusOK, "success")
	})

	handle.ServeHTTP(resp, req)
	if resp.Result().StatusCode != http.StatusOK {
		t.Fatalf("header authentication doesnt work")
	}

}

func TestUserClient(t *testing.T) {
	cl := us.Client()
	if cl1 != cl {
		t.Fatalf("not the same pointers")
	}
}

func TestUserDelete(t *testing.T) {
	us.Delete()

	if us.Client() != nil {
		t.Fatalf("delete does not delete")
	}
}

func TestUserInvite(t *testing.T) {
	rd1, wr1 = io.Pipe()
	rd2, wr2 = io.Pipe()

	us1.Client().W = wr1
	us2.Client().W = wr2

	const lifespan = time.Millisecond * 10
	// because net.Pipe is synchronous, we have to read through it. otherwise it would hang forever...
	go func() {
		rd2.Read(make([]byte, 64))
	}()
	_, err := us1.Invite(us2.PublicID, lifespan)
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
	x := make([]byte, 128)
	var n int

	go func() {
		n, err = rd2.Read(x)
		ch <- err
		x = x[:n]
	}()

	_, err = us1.Invite(us2.PublicID, InviteLifespan)
	if err != nil {
		t.Fatalf("us.Invite: %s", err.Error())
	}
	err = <-ch
	if err != nil {
		t.Fatalf("rd2.Read: %s", err.Error())
	}

	update := order.Order{}
	err = json.Unmarshal(x[:n], &update)
	if err != nil {
		t.Fatalf("json.Unmarshal: %s", err.Error())
	}

	inv := order.InviteModel{}
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

	x = make([]byte, 1280)
	n, err = rd2.Read(x)
	if err != nil {
		t.Fatalf("rd2.Read: %s", err.Error())
	}
	o := order.Order{}
	err = json.Unmarshal(x[:n], &o)
	if err != nil {
		t.Fatalf("json.Unmarshal: %s", err.Error())
	}
	gm2 := order.GameModel{}
	err = json.Unmarshal(o.Data, &gm2)
	if err != nil {
		t.Fatalf("json.Unmarshal: %s", err.Error())
	}

	x = make([]byte, 1280)
	n, err = rd1.Read(x)
	if err != nil {
		t.Fatalf("rd1.Read: %s", err.Error())
	}
	o = order.Order{}
	err = json.Unmarshal(x[:n], &o)
	if err != nil {
		t.Fatalf("json.Unmarshal: %s", err.Error())
	}
	gm1 := order.GameModel{}
	err = json.Unmarshal(o.Data, &gm1)
	if err != nil {
		t.Fatalf("json.Unmarshal: %s", err.Error())
	}

	if gm1.Player == gm2.Player {
		t.Fatalf("player numbers are the same")
	}
	t.Log(gm1.Player, gm2.Player)

}

func TestGetAvaliableUsersHandler(t *testing.T) {

	us1.Client().LeaveGame()
	us2.Client().LeaveGame()

	handle := http.HandlerFunc(GetAvaliableUsersHandler)

	resp := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	avali := []string{}

	handle.ServeHTTP(resp, req)
	body, err := ioutil.ReadAll(resp.Body)
	if resp.Result().StatusCode != http.StatusOK {
		t.Fatalf("status is not ok.")
	}

	if err := json.Unmarshal(body, &avali); err != nil {
		t.Fatalf("json.Unmarshal: %s", err.Error())
	}

	if len(avali) != 2 {
		t.Fatalf("len(avali): %d - want: 2", len(avali))
	}

	game.NewGame(cl1, cl2)

	avali = []string{}

	handle.ServeHTTP(resp, req)
	body, err = ioutil.ReadAll(resp.Body)
	if resp.Result().StatusCode != http.StatusOK {
		t.Fatalf("status is not ok.")
	}

	if err := json.Unmarshal(body, &avali); err != nil {
		t.Fatalf("json.Unmarshal: %s", err.Error())
	}

	if len(avali) != 0 {
		t.Fatalf("len(avali): %d - want: 0", len(avali))
	}

	rd1, wr1 = io.Pipe()
	rd2, wr2 = io.Pipe()

	us1.Client().W = wr1
	us2.Client().W = wr2
}
