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

var kingmap = map[int]int{
	7: 6,
	0: 2,
}
var rookmap = map[int]int{
	7: 5,
	0: 3,
}

func TestCommandSendMessage(t *testing.T) {
	defer resetPipe()

	body, err := json.Marshal(order.MessageModel{
		Message: "test",
	})

	if err != nil {
		t.Fatalf("json.Marshal: %s", err.Error())
	}

	c := order.Order{
		ID:   order.Message,
		Data: body,
	}

	cherr := make(chan error)
	go func() {
		cherr <- cl1.Do(c)
	}()

	body = <-clientRead(rd1)

	x := order.MessageModel{}

	u := order.Order{}

	err = json.Unmarshal(body, &u)
	if err != nil {
		t.Fatalf("json.Unmarshal: %s", err.Error())
	}

	t.Logf("%s", string(body))

	err = json.Unmarshal(u.Data, &x)
	if err != nil {
		t.Fatalf("json.Unmarshal: %s", err.Error())
	}

	if x.Message != "[Player 1]: test" {
		t.Fatalf("json.Unmarshal: unwanted result")
	}

	<-clientRead(rd2)

	err = <-cherr
	if err != nil {
		t.Fatalf("cl.Do: %s", err.Error())
	}
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

	do := func(cl *Client, rd *io.PipeReader, src, dst board.Point) error {
		body, err := json.Marshal(order.MoveModel{
			Src: src,
			Dst: dst,
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

		p := cl.g.b.Get(dst)
		if p == nil {
			return fmt.Errorf("piece is nil")
		}

		<-clientRead(rd2)

		err = <-cherr
		if err != nil {
			return fmt.Errorf("cl.Do: %s", err.Error())
		}

		return nil
	}

	err := do(cl1, rd1, board.Point{X: 6, Y: 1}, board.Point{X: 4, Y: 1})
	if err != nil {
		t.Fatalf("client 1 : %v", err)
	}

	resetPipe()

	err = do(cl2, rd2, board.Point{X: 1, Y: 1}, board.Point{X: 3, Y: 1})
	if err != nil {
		t.Fatalf("client 2 : %v", err)
	}

}

func TestCommandPromotion(t *testing.T) {
	resetPipe()
	defer resetPipe()

	cl1.g, cl2.g = nil, nil
	gGame, _ = NewGame(cl1, cl2)

	pec := gGame.b.Get(board.Point{6, 3})
	pos := pec.Pos

	ch := make(chan error)
	go func() {
		pos.X -= 2
		if !gGame.b.Move(pec, pos) {
			//fmt.Println("error 1")
			ch <- fmt.Errorf("cannot move from %v to %v", pec.Pos, pos)
			return
		}

		gGame.b.Set(&board.Piece{Pos: board.Point{1, 3}, T: board.Empty})
		gGame.b.Set(&board.Piece{Pos: board.Point{0, 3}, T: board.Empty})

		for i := 0; i < 4; i++ {
			pos.X -= 1
			if !gGame.b.Move(pec, pos) {
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
		if !promote.Src.Equal(pos) {
			t.Fatalf("promote coordinates aren't at x=7|x=1")
		}

		promote.Type = board.Queen
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

	err := <-ch
	if err != nil {
		t.Fatalf("%v", err)
	}
}

func TestCommandCastling(t *testing.T) {

	defer resetPipe()

	id := order.Castling
	do := func(x int, rooky int, kingy int, cl *Client) {
		resetPipe()
		cl1.g, cl2.g = nil, nil
		gGame, _ = NewGame(cl1, cl2)

		go func() {
			gGame.SwitchTurn()
			if x == 0 {
				gGame.SwitchTurn()
			}
		}()

		<-clientRead(rd1)
		<-clientRead(rd2)
		if x == 0 {
			<-clientRead(rd1)
			<-clientRead(rd2)
		}

		for y := 1; y < 7; y++ {
			if y == 4 {
				continue
			}

			gGame.b.Set(&board.Piece{
				Pos: board.Point{x, y},
				T:   board.Empty,
			})
		}

		t.Logf("\n%s", gGame.b.String())

		cast := order.CastlingModel{
			Src: board.Point{x, kingy},
			Dst: board.Point{x, rooky},
		}

		body, err := json.Marshal(cast)
		if err != nil {
			t.Fatalf("json.Marshal: %s", err.Error())
		}

		x1, x2 := make(chan []byte), make(chan []byte)
		go func() {
			<-clientRead(rd1)
			<-clientRead(rd2)
			x1 = clientRead(rd1)
			x2 = clientRead(rd2)
		}()

		t.Logf(string(body))
		err = cl.Do(order.Order{
			ID:   id,
			Data: body,
		})
		if err != nil {
			t.Fatalf("cl.Do: %s", err.Error())
		}

		b1, b2 := <-x1, <-x2
		t.Log(string(b1), string(b2))

		if kingy != 4 {
			rooky, kingy = kingy, rooky
		}

		pecrook, pecking := gGame.b.Get(board.Point{x, rookmap[rooky]}), gGame.b.Get(board.Point{x, kingmap[rooky]})
		if pecrook == nil || pecking == nil || pecrook.T != board.Rook || pecking.T != board.King {
			t.Fatalf("unpredictable results. %s: %s | %s: %s", pecrook.Pos, pecrook, pecking.Pos, pecking)
		}
	}

	do(7, 7, 4, cl1)
	do(7, 0, 4, cl1)
	do(7, 4, 7, cl1)
	do(7, 4, 0, cl1)

	do(0, 7, 4, cl2)
	do(0, 0, 4, cl2)
	do(0, 4, 7, cl2)
	do(0, 4, 7, cl2)

	id = order.Move
	t.Log("henlo")
	do(7, 7, 4, cl1)

}
