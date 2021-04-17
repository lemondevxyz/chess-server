package game

import (
	"encoding/json"
	"io"
	"testing"

	"github.com/toms1441/chess-server/internal/board"
	"github.com/toms1441/chess-server/internal/order"
)

// Note: Future contributors, beware of io.Pipe freezing the entire test, wrapping Write operations around a goroutine would make the tests unpredictable.
// Just watch out for different order messages, try mixing the clientRead. Or increase the amount, or decrease it.
var (
	rd1, wr1 = io.Pipe()
	rd2, wr2 = io.Pipe()

	cl1 = &Client{
		W:  wr1,
		p1: true,
		id: "a",
	}

	cl2 = &Client{
		W:  wr2,
		p1: false,
		id: "b",
	}
)

func clientRead(rd *io.PipeReader) chan []byte {

	ch := make(chan []byte)

	go func() {
		body := make([]byte, 128)
		n, err := rd.Read(body)
		if err != nil {
			ch <- nil
		} else {
			ch <- body[:n]
		}
	}()

	return ch
}

func resetPipe() {
	rd1, wr1 = io.Pipe()
	rd2, wr2 = io.Pipe()

	cl1.W = wr1
	cl2.W = wr2
}

func TestInPromotion(t *testing.T) {
	resetPipe()

	go func() {
		<-clientRead(rd2)
	}()
	cl1.LeaveGame()
	//t.Log(string(x))

	var err error
	gGame, err = NewGame(cl1, cl2)
	if err != nil {
		t.Fatalf("NewGame: %s", err.Error())
	}

	go gGame.SwitchTurn()
	<-clientRead(rd1)
	<-clientRead(rd2)

	list := []order.MoveModel{
		{17, board.Point{1, 4}},
		{15, board.Point{7, 3}},
		{17, board.Point{1, 3}},
		{15, board.Point{7, 4}},
		{17, board.Point{1, 2}},
		{15, board.Point{7, 5}},
		{17, board.Point{0, 1}},
		{15, board.Point{6, 6}},
		{17, board.Point{1, 0}},
	}

	p1 := true
	for k, v := range list {
		body, err := json.Marshal(v)
		if err != nil {
			t.Fatalf("json.Marshal: %s", err.Error())
		}

		ord := order.Order{
			ID:   order.Move,
			Data: body,
		}

		go func(index int) {
			<-clientRead(rd1)
			<-clientRead(rd2)
			<-clientRead(rd1)
			<-clientRead(rd2)
		}(k)

		if k+1 == len(list) {
			go clientRead(rd1)
		}

		t.Log("bef", k)
		which := ""
		if p1 {
			which = "cl1.Do"
			err = cl1.Do(ord)
		} else {
			which = "cl2.Do"
			err = cl2.Do(ord)
		}
		t.Log("aft", k)

		if err != nil {
			t.Fatalf("%s: %s", which, err.Error())
		}
		p1 = !p1
	}

	if cl1.inPromotion() == false {
		t.Fatalf("in promotion does not work")
	}

	body, err := json.Marshal(order.PromoteModel{
		ID:   17,
		Kind: board.Queen,
	})
	if err != nil {
		t.Fatalf("json.Marshal: %s", err.Error())
	}

	go clientRead(rd1)
	go clientRead(rd2)

	err = cl1.Do(order.Order{
		ID:   order.Promote,
		Data: body,
	})
	t.Log(err)
	if err != nil {
		t.Fatalf("cl1.Do: %s", err.Error())
	}

	pec, err := cl1.g.b.GetByIndex(17)
	t.Log(pec, err)
	if cl1.inPromotion() == true {
		t.Fatalf("client is in promotion, after promotion.")
	}
}
