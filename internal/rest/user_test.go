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
