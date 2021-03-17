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

		if x.Player != 2 {
			t.Fatalf("player is not two: %d", x.Player)
		}

		<-clientRead(rd2)
		<-clientRead(rd1)
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
				by1 := <-clientRead(rd1)
				by2 := <-clientRead(rd2)
				fmt.Println("bc", string(by1), "|", string(by2))
				by1 = <-clientRead(rd1)
				by2 = <-clientRead(rd2)
				fmt.Println("yz", string(by1), "|", string(by2))
			}()
		}

		err = cl.Do(order.Order{
			ID:   order.Move,
			Data: x,
		})
		if err != nil {
			t.Fatalf("cl.Do: %s", err)
		}
	}

	doMove(board.Point{6, 5}, board.Point{5, 5})
	doMove(board.Point{1, 4}, board.Point{3, 4})
	doMove(board.Point{6, 6}, board.Point{4, 6})
	done = true
	go func() {
		by2 = <-clientRead(rd2)
		by1 = <-clientRead(rd1)
		fmt.Println("kxxx", string(by1), "|", string(by2))
	}()
	fmt.Println("xxx")
	doMove(board.Point{0, 3}, board.Point{4, 7})
	fmt.Println("adsdas")

	t.Logf("\n%s", gGame.b.String())

}
