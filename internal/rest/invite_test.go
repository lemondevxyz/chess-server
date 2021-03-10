package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/toms1441/chess-server/internal/order"
)

var _inviteCode = ""

func TestInviteHandler(t *testing.T) {
	us1.cl.LeaveGame()
	us2.cl.LeaveGame()

	x := order.InviteModel{
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

		upd := order.Order{}
		err = json.Unmarshal(body, &upd)
		if err != nil {
			ch <- fmt.Errorf("json.Unmarshal: %s", err.Error())
			return
		}

		mod := order.InviteModel{}
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
	i := order.InviteModel{
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

	// game update
	x := make([]byte, 1024)
	rd2.Read(x)

	x = make([]byte, 1024)
	rd1.Read(x)

	x = make([]byte, 1024)
	rd2.Read(x)

	x = make([]byte, 1024)
	rd1.Read(x)

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

	us1.Client().LeaveGame()
	us2.Client().LeaveGame()

}
