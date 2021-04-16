package game

import (
	"encoding/json"
	"fmt"
	"io"
	"testing"
	"time"

	"github.com/toms1441/chess-server/internal/board"
	"github.com/toms1441/chess-server/internal/order"
)

var kingmap = map[int8]int8{
	7: 6,
	0: 2,
}
var rookmap = map[int8]int8{
	7: 5,
	0: 3,
}

func TestCommandMove(t *testing.T) {
	defer resetPipe()

	cl1.g, cl2.g = nil, nil
	gGame, _ = NewGame(cl1, cl2)

	go func() {
		<-clientRead(rd1)
		<-clientRead(rd2)
	}()
	gGame.SwitchTurn()

	do := func(cl *Client, rd *io.PipeReader, id int8, dst board.Point) error {
		body, err := json.Marshal(order.MoveModel{
			ID:  &id,
			Dst: &dst,
		})

		if err != nil {
			return fmt.Errorf("json.Marshal: %s", err.Error())
		}

		c := order.Order{
			ID:   order.Move,
			Data: body,
		}

		cherr := make(chan error)
		go func() {
			cherr <- cl.Do(c)
		}()

		// turn message
		<-clientRead(rd1)
		<-clientRead(rd2)

		select {
		case <-time.After(time.Millisecond * 10):
			return fmt.Errorf("timeout")
		case body := <-clientRead(rd1):
			x := &order.MoveModel{}
			u := order.Order{}
			err = json.Unmarshal(body, &u)
			if err != nil {
				return fmt.Errorf("json.Unmarshal: %s", err.Error())
			}
			err = json.Unmarshal(u.Data, x)
			if err != nil {
				return fmt.Errorf("json.Unmarshal: %s", err.Error())
			}

			t.Logf("\n%v", x)
		}

		_, pec, err := cl.g.b.Get(dst)
		if err != nil || !pec.Valid() {
			return fmt.Errorf("piece is nil")
		}

		<-clientRead(rd2)

		err = <-cherr
		if err != nil {
			return fmt.Errorf("cl.Do: %s", err.Error())
		}

		return nil
	}

	err := do(cl1, rd1, 17, board.Point{1, 4})
	if err != nil {
		t.Fatalf("client 1 : %v", err)
	}

	resetPipe()

	err = do(cl2, rd2, 9, board.Point{1, 3})
	if err != nil {
		t.Fatalf("client 2 : %v", err)
	}

}

func TestCommandPromotion(t *testing.T) {
	const id = 19

	resetPipe()
	defer resetPipe()

	cl1.g, cl2.g = nil, nil
	gGame, _ = NewGame(cl1, cl2)

	pec, err := gGame.b.GetByIndex(id)
	if err != nil {
		t.Fatalf("board.Get: %s", err)
	}
	pos := pec.Pos

	ch := make(chan error)
	go func() {
		pec, _ := gGame.b.GetByIndex(id)
		pos.Y -= 2
		if !gGame.b.Move(id, pos) {
			//fmt.Println("error 1")
			ch <- fmt.Errorf("here cannot move from %v to %v", pec.Pos, pos)
			return
		}

		gGame.b.Set(3, board.Point{-1, -1})
		gGame.b.Set(11, board.Point{-1, -1})

		for i := 0; i < 4; i++ {
			pec, _ := gGame.b.GetByIndex(id)
			pos.Y -= 1
			if !gGame.b.Move(id, pos) {
				//fmt.Println("error 2")
				ch <- fmt.Errorf("cannot move from %v to %v", pec.Pos, pos)
				return
			}
		}

		close(ch)
	}()

	select {
	case <-time.After(time.Millisecond * 100):
		t.Fatalf("timeout")
	case body := <-clientRead(rd1):
		upd := &order.Order{}
		if err := json.Unmarshal(body, &upd); err != nil {
			t.Fatalf("json.Unmarshal: %s", err.Error())
		}
		promote := &order.PromoteModel{}
		if err := json.Unmarshal(upd.Data, &promote); err != nil {
			t.Fatalf("json.Unmarshal: %s", err.Error())
		}
		if promote.ID != id {
			t.Fatalf("promote.ID != ID : %d != %d", promote.ID, id)
		}

		promote.Kind = board.Queen
		body, err := json.Marshal(promote)
		if err != nil {
			t.Fatalf("json.Marshal: %s", err.Error())
		}

		cmd := &order.Order{
			ID:   order.Promote,
			Data: body,
		}

		go func() {
			// turn
			<-clientRead(rd1)
			<-clientRead(rd2)
			<-clientRead(rd1)
			<-clientRead(rd2)
		}()

		err = cl1.Do(*cmd)
		if err != nil {
			t.Fatalf("cl1.Do: %s", err.Error())
		}
	}

	err = <-ch
	if err != nil {
		t.Fatalf("%v", err)
	}
}

func TestCommandCastling(t *testing.T) {

	defer resetPipe()

	id := order.Castling
	do := func(other bool, rook int, king int, cl *Client) {
		resetPipe()
		cl1.g, cl2.g = nil, nil
		gGame, _ = NewGame(cl1, cl2)

		go func() {
			gGame.SwitchTurn()
			if other {
				gGame.SwitchTurn()
			}
		}()

		<-clientRead(rd1)
		<-clientRead(rd2)
		if other {
			<-clientRead(rd1)
			<-clientRead(rd2)
		}

		row := board.GetRangeStart(cl.p1)
		for _, v := range row {

			pec, _ := gGame.b.GetByIndex(v)
			if pec.Kind != board.Rook && pec.Kind != board.King {
				gGame.b.Set(v, board.Point{-1, -1})
			}
		}

		// thanks golang
		p1, p2 := &rook, &king
		cast := order.CastlingModel{
			Src: p1,
			Dst: p2,
		}

		body, err := json.Marshal(cast)
		if err != nil {
			t.Fatalf("json.Marshal: %s", err.Error())
		}

		x1, x2 := make(chan []byte), make(chan []byte)
		go func() {
			x1 = clientRead(rd1)
			clientRead(rd2)
			clientRead(rd1)
			x2 = clientRead(rd2)
		}()

		pecking, err := gGame.b.GetByIndex(king)
		if err != nil {
			t.Fatalf("board.GetByIndex(king): %s", err.Error())
		}
		pecrook, err := gGame.b.GetByIndex(rook)
		if err != nil {
			t.Fatalf("board.GetByIndex(king): %s", err.Error())
		}

		err = cl.Do(order.Order{
			ID:   id,
			Data: body,
		})
		if err != nil {
			t.Fatalf("cl.Do: %s", err.Error())
		}

		b1, b2 := <-x1, <-x2
		t.Log(string(b1), string(b2))

		y := pecking.Pos.Y
		rookx := pecrook.Pos.X
		if rookx == 4 {
			rookx = pecking.Pos.X
		}

		pecking, err = gGame.b.GetByIndex(king)
		if err != nil {
			t.Fatalf("board.GetByIndex(king): %s", err.Error())
		}
		pecrook, err = gGame.b.GetByIndex(rook)
		if err != nil {
			t.Fatalf("board.GetByIndex(king): %s", err.Error())
		}

		if pecrook.Kind == board.King && pecking.Kind == board.Rook {
			pecrook, pecking = pecking, pecrook
		}

		want := board.Point{rookmap[rookx], y}
		if !pecrook.Pos.Equal(want) {
			t.Logf("\n%s", gGame.b.String())
			t.Fatalf("rook's position hasn't changed. want: %s | have: %s", want, pecrook.Pos)
		}
		want = board.Point{kingmap[rookx], y}
		if !pecking.Pos.Equal(want) {
			t.Logf("\n%s", gGame.b.String())
			t.Fatalf("king's position hasn't changed. want: %s | have: %s", want, pecking.Pos)
		}
	}

	king := board.GetKing(cl1.p1)
	rks := board.GetRooks(cl1.p1)

	do(false, rks[1], king, cl1)
	t.Log("aft 1")
	do(false, rks[0], king, cl1)
	t.Log("aft 2")
	do(false, king, rks[0], cl1)
	t.Log("aft 3")
	do(false, king, rks[1], cl1)
	t.Log("aft 4")

	king = board.GetKing(cl2.p1)
	rks = board.GetRooks(cl2.p1)

	do(true, rks[1], king, cl2)
	t.Log("aft 5")
	do(true, rks[0], king, cl2)
	t.Log("aft 6")
	do(true, king, rks[0], cl2)
	t.Log("aft 7")
	do(true, king, rks[1], cl2)
	t.Log("aft 8")

}

func TestCommandDone(t *testing.T) {

	defer resetPipe()

	go cl1.LeaveGame()
	x := <-clientRead(rd2)
	t.Log(string(x))

	var err error
	gGame, err = NewGame(cl1, cl2)
	if err != nil {
		t.Fatalf("NewGame: %s", err.Error())
	}

	done := make(chan map[string]interface{})
	go func() {
		b := <-clientRead(rd1)
		<-clientRead(rd2)
		pam := map[string]interface{}{}

		json.Unmarshal(b, &pam)

		done <- pam
	}()
	err = cl1.Do(order.Order{
		ID:   order.Done,
		Data: nil,
	})
	if err != nil {
		t.Fatalf("cl.Do: %s", err.Error())
	}

	pam := <-done
	data := pam["data"].(map[string]interface{})

	won := data["p1"].(bool)

	if cl1.P1() == won {
		t.Fatalf("cl1 should be the one who's losing, not winning")
	}
	if cl2.P1() != won {
		t.Fatalf("cl2 should be the one who won, not losing..")
	}

}
