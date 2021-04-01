package game

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/toms1441/chess-server/internal/board"
	"github.com/toms1441/chess-server/internal/order"
)

var gGame, _ = NewGame(cl1, cl2)

func TestTurns(t *testing.T) {
	resetPipe()

	cl1.g, cl2.g = nil, nil
	gGame, _ = NewGame(cl1, cl2)
	go gGame.SwitchTurn()

	time.Sleep(time.Millisecond * 10)

	select {
	case <-time.After(time.Millisecond * 100):
		t.Fatalf("timeout")
	case body := <-clientRead(rd1):
		t.Logf("\n%s", string(body))

		x := &order.TurnModel{}
		u := &order.Order{}

		err := json.Unmarshal(body, u)
		if err != nil {
			t.Fatalf("json.Unmarshal: %s", err.Error())
		}

		err = json.Unmarshal(u.Data, x)
		if err != nil {
			t.Fatalf("json.Unmarshal: %s", err.Error())
		}

		if x.Player != 1 {
			t.Fatalf("player is not one: %d", x.Player)
		}
	}

	resetPipe()

	cherr := make(chan error)
	go func() {
		body, err := json.Marshal(order.MoveModel{
			Src: board.Point{6, 1},
			Dst: board.Point{5, 1},
		})

		if err != nil {
			cherr <- fmt.Errorf("json.Marshal: %s", err.Error())
			return
		}

		c := order.Order{
			ID:   order.Move,
			Data: body,
		}

		err = cl2.Do(c)
		if err == nil {
			cherr <- fmt.Errorf("it's not 2 turn, yet it works .. :(")
			return
		}

		err = cl1.Do(c)
		if err != nil {
			cherr <- err
			return
		}
		//cherr <- nil
		close(cherr)
	}()

	<-clientRead(rd1)
	<-clientRead(rd2)
	select {
	case <-time.After(time.Millisecond * 100):
		t.Fatalf("timeout")
	case body := <-clientRead(rd1):
		t.Logf("\n%s", string(body))

		x := &order.TurnModel{}
		u := &order.Order{}

		err := json.Unmarshal(body, u)
		if err != nil {
			t.Fatalf("json.Unmarshal: %s", err.Error())
		}

		err = json.Unmarshal(u.Data, x)
		if err != nil {
			t.Fatalf("json.Unmarshal: %s", err.Error())
		}

		t.Log(x)
		if x.Player != 2 {
			t.Fatalf("player is not two: %d", x.Player)
		}

		<-clientRead(rd2)

		select {
		case <-time.After(time.Millisecond * 100):
			t.Fatalf("timeout")
		case err := <-cherr:
			if err != nil {
				t.Fatalf("err: %s", err.Error())
			}
		}
	}
}

func TestGameDone(t *testing.T) {

	resetPipe()

	cl1.g = nil
	cl2.g = nil

	gGame, _ = NewGame(cl1, cl2)

	go gGame.SwitchTurn()

	by1 := <-clientRead(rd1)
	by2 := <-clientRead(rd2)
	t.Log(string(by1), string(by2))

	done := false

	doMove := func(src, dst board.Point) {
		cl := gGame.cs[gGame.turn-1]

		x, err := json.Marshal(order.MoveModel{
			Src: src,
			Dst: dst,
		})
		if err != nil {
			t.Fatalf("json.Marshal: %s", err)
		}

		if !done {
			go func() {
				<-clientRead(rd1)
				<-clientRead(rd2)
				<-clientRead(rd1)
				<-clientRead(rd2)
			}()
		}

		err = cl.Do(order.Order{
			ID:   order.Move,
			Data: x,
		})
		if err != nil {
			t.Fatalf("cl.Do: %s", err.Error())
		}

	}
	t.Log("first move")
	doMove(board.Point{6, 5}, board.Point{5, 5})
	t.Log("second move")
	doMove(board.Point{1, 4}, board.Point{3, 4})
	t.Log("third move")
	doMove(board.Point{6, 6}, board.Point{4, 6})
	done = true
	resetPipe()
	go func() {
		<-clientRead(rd1)
		<-clientRead(rd2)
		<-clientRead(rd2)
		<-clientRead(rd1)
	}()
	t.Log("fourth move")
	doMove(board.Point{0, 3}, board.Point{4, 7})
	t.Log("after fourth move")

	lc1 := gGame.cs[0]
	lc2 := gGame.cs[1]
	if lc1 != nil || lc2 != nil {
		t.Fatalf("gGame cs: %v | %v", lc1, lc2)
	}

	t.Logf("\n%s", gGame.b.String())

}
