package rest

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/toms1441/chess-server/internal/game"
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
	err := us1.Invite(us2.PublicID, lifespan)
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
	x := make([]byte, 64)

	go func() {
		n, err := rd2.Read(x)
		ch <- err
		x = x[:n]
	}()

	err = us1.Invite(us2.PublicID, InviteLifespan)
	if err != nil {
		t.Fatalf("us.Invite: %s", err.Error())
	}
	err = <-ch
	if err != nil {
		t.Fatalf("rd2.Read: %s", err.Error())
	}

	update := game.Update{}
	err = json.Unmarshal(x, &update)
	if err != nil {
		t.Fatalf("json.Unmarshal: %s", err.Error())
	}

	inv := game.ModelUpdateInvite{}
	json.Unmarshal(update.Data, &inv)
	if err != nil {
		t.Fatalf("json.Unmarshal: %s", err.Error())
	}

	err = us2.AcceptInvite(inv.ID)
	if err != nil {
		t.Fatalf("us.AcceptInvite: %s", err.Error())
	}

	if us1.Client().Game() == nil || us2.Client().Game() == nil {
		t.Fatalf("us.AcceptInvite: does not start a new game!")
	}
}
