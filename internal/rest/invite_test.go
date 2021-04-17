package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/toms1441/chess-server/internal/model"
)

var _inviteCode = ""

func TestInviteHandler(t *testing.T) {
	go us1.cl.LeaveGame()
	<-read(rd2)
	//us2.cl.LeaveGame()

	x := model.InviteOrder{
		ID: us2.PublicID,
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
		ID: _inviteCode,
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
