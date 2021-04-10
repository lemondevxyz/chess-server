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

	"github.com/toms1441/chess-server/internal/board"
	"github.com/toms1441/chess-server/internal/game"
	"github.com/toms1441/chess-server/internal/order"
)

var (
	rd1, wr1 = io.Pipe()
	rd2, wr2 = io.Pipe()

	cl1 = &game.Client{
		W: wr1,
	}

	cl2 = &game.Client{
		W: wr2,
	}

	us1 = AddClient(cl1)
	us2 = AddClient(cl2)

	gGame, _ = game.NewGame(cl1, cl2)
)

func read(r io.Reader) chan []byte {
	ch := make(chan []byte)

	go func() {
		x := make([]byte, 1024)
		n, _ := r.Read(x)
		x = x[:n]

		ch <- x
	}()

	return ch
}

func TestCommandRequest(t *testing.T) {

	go gGame.SwitchTurn()
	<-read(rd1)
	<-read(rd2)

	x := order.MoveModel{
		ID:  17, // {1, 6}
		Dst: board.Point{1, 4},
	}

	byt, err := json.Marshal(x)
	if err != nil {
		t.Fatalf("json.Marshal: %s", err.Error())
	}

	cmd := order.Order{
		ID:   order.Move,
		Data: byt,
	}

	byt, err = json.Marshal(cmd)
	if err != nil {
		t.Fatalf("json.Marshal: %s", err.Error())
	}

	resp := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "/", bytes.NewBuffer(byt))

	if err != nil {
		t.Fatal(err)
	}
	handle := http.HandlerFunc(CmdHandler)

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", us1.Token))

	done := make(chan struct{})
	go func() {
		handle.ServeHTTP(resp, req)
		done <- struct{}{}
	}()

	<-read(rd1)
	<-read(rd2)
	<-read(rd1)
	<-read(rd2)

	<-done
	t.Log(resp.Body.String())
	hh := resp.Header()
	if hh.Get("Content-Type") != "application/json" {
		t.Fatalf("bad content type")
	}

	if p, err := ioutil.ReadAll(resp.Body); err != nil {
		t.Fatalf("ioutil.ReadAll: %s", err.Error())
	} else {
		obj := map[string]interface{}{}

		err := json.Unmarshal(p, &obj)
		if err != nil {
			t.Fatalf("json.Unmarshal: %s", err.Error())
		} else {
			err, ok := obj["error"]
			if ok {
				t.Fatalf("err: %s", err.(string))
			} else {
				t.Logf("%v", obj)
			}
		}
	}

}

func TestPossibleRequest(t *testing.T) {

	id := int8(15)
	possib := order.PossibleModel{
		ID: &id,
	}

	byt, err := json.Marshal(possib)
	if err != nil {
		t.Fatalf("json.Marshal: %s", err.Error())
	}

	resp := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "/", bytes.NewBuffer(byt))

	if err != nil {
		t.Fatal(err)
	}

	handle := http.HandlerFunc(PossibHandler)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", us1.Token))

	handle.ServeHTTP(resp, req)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("ioutil.ReadAll: %s", err.Error())
	}

	err = json.Unmarshal(body, &possib)
	if err != nil {
		t.Fatalf("json.Unmarshal: %s", err.Error())
	}

	t.Log(possib)
	t.Log(string(body))
}
