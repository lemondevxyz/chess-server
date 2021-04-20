package rest

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/toms1441/chess-server/internal/game"
	"github.com/toms1441/chess-server/internal/model"
	"github.com/toms1441/chess-server/internal/model/local"
)

var (
	us       = &User{}
	urd, uwr = io.Pipe()
)

func TestNewUser(t *testing.T) {
	var err error
	us, err = AddClient(local.NewUser(), uwr)

	if err != nil {
		t.Fatalf("AddClient: %s", err.Error())
	}
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

func TestUserDelete(t *testing.T) {
	us.Delete()

	if us.Client() != nil {
		t.Fatalf("delete does not delete")
	}
}

/*
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

	_, err := us1.Invite(model.InviteOrder{
		ID:       us2.Profile.GetPublicID(),
		Platform: us2.Profile.GetPlatform(),
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
			ID:       us2.Profile.GetPublicID(),
			Platform: us2.Profile.GetPlatform(),
		}, InviteLifespan)
		if err != nil {
			ch <- fmt.Errorf("us.Invite: %s", err.Error())
			return
		}
		close(ch)
	}()

	update := model.Order{}
	err = json.Unmarshal(<-read(rd2), &update)
	t.Log()
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
*/

func TestGetAvaliableUsersHandler(t *testing.T) {
	go func() {
		// turn update
		<-read(rd2)
		<-read(rd1)
		// game done update
		<-read(rd2)
		<-read(rd1)
	}()
	us1.Client().LeaveGame()

	handle := http.HandlerFunc(GetAvaliableUsersHandler)

	resp := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Add("Authorization", us1.Token)

	avali := []model.Profile{}

	handle.ServeHTTP(resp, req)
	body, err := ioutil.ReadAll(resp.Body)
	if resp.Result().StatusCode != http.StatusOK {
		t.Fatalf("status is not ok.")
	}

	if err := json.Unmarshal(body, &avali); err != nil {
		t.Fatalf("json.Unmarshal: %s", err.Error())
	}

	if len(avali) != 1 {
		t.Fatalf("len(avali): %d - want: 1", len(avali))
	}

	game.NewGame(cl1, cl2)

	avali = []model.Profile{}

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
