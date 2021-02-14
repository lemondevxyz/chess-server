package game

import (
	"encoding/json"
	"fmt"
	"io"
	"testing"
	"time"

	"github.com/toms1441/chess-server/internal/board"
)

func TestCommandSendMessage(t *testing.T) {
	defer resetPipe()

	body, err := json.Marshal(ModelCmdMessage{
		Message: "test",
	})

	if err != nil {
		t.Fatalf("json.Marshal: %s", err.Error())
	}

	c := Command{
		ID:   CmdMessage,
		Data: body,
	}

	cherr := make(chan error)
	go func() {
		cherr <- cl1.Do(c)
	}()

	err = <-cherr
	if err != nil {
		t.Fatalf("cl.Do: %s", err.Error())
	}

	// first message is the turn message
	<-clientRead(rd1)

	select {
	case <-time.After(time.Millisecond * 10):
		t.Fatalf("timeout")
	case body := <-clientRead(rd1):

		x := ModelCmdMessage{}

		u := Update{}

		err := json.Unmarshal(body, &u)
		if err != nil {
			t.Fatalf("json.Unmarshal: %s", err.Error())
		}

		t.Logf("%s", string(body))

		err = json.Unmarshal(u.Data, &x)
		if err != nil {
			t.Fatalf("json.Unmarshal: %s", err.Error())
		}

		if x.Message != "test" {
			t.Fatalf("json.Unmarshal: unwanted result")
		}
	}
}

func TestCommandPiece(t *testing.T) {
	defer resetPipe()

	do := func(cl *Client, rd *io.PipeReader, src, dst board.Point) error {
		body, err := json.Marshal(ModelCmdPiece{
			Src: src,
			Dst: dst,
		})

		if err != nil {
			return fmt.Errorf("json.Marshal: %s", err.Error())
		}

		c := Command{
			ID:   CmdPiece,
			Data: body,
		}

		cherr := make(chan error)
		go func() {
			cherr <- cl.Do(c)
		}()

		err = <-cherr
		if err != nil {
			return fmt.Errorf("cl.Do: %s", err.Error())
		}

		// turn message
		<-clientRead(rd1)
		select {
		case <-time.After(time.Millisecond * 10):
			return fmt.Errorf("timeout")
		case body := <-clientRead(rd1):
			x := &board.Board{}
			u := Update{}
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

		return nil
	}

	err := do(cl1, rd1, board.Point{X: 1, Y: 1}, board.Point{X: 3, Y: 1})
	if err != nil {
		t.Fatalf("client 1 : %v", err)
	}

	resetPipe()

	err = do(cl2, rd2, board.Point{X: 6, Y: 1}, board.Point{X: 4, Y: 1})
	if err != nil {
		t.Fatalf("client 2 : %v", err)
	}

}

func TestCommandPromotion(t *testing.T) {
	defer resetPipe()

	p := gGame.b.Get(board.Point{X: 1, Y: 5})
	go func() {

		gGame.b.Move(p, board.Point{X: p.X + 2, Y: p.Y})

		gGame.b.Set(&board.Piece{T: board.Empty, X: 0, Y: 4})
		gGame.b.Set(&board.Piece{T: board.Empty, X: 6, Y: 5})
		gGame.b.Set(&board.Piece{T: board.Empty, X: 7, Y: 5})

		gGame.b.Move(p, board.Point{X: p.X + 1, Y: p.Y})
		gGame.b.Move(p, board.Point{X: p.X + 1, Y: p.Y})
		gGame.b.Move(p, board.Point{X: p.X + 1, Y: p.Y})
		gGame.b.Move(p, board.Point{X: p.X + 1, Y: p.Y})
	}()

	select {
	case <-time.After(time.Millisecond * 100):
		t.Fatalf("timeout")
	case body := <-clientRead(rd1):
		t.Logf("\n%s", gGame.b.String())
		t.Log(string(body))

		u := Update{}
		err := json.Unmarshal(body, &u)
		if err != nil {
			t.Fatalf("json.Unmarshal: %s", err.Error())
		}

		parameter := ModelUpdatePromotion{}
		err = json.Unmarshal(u.Data, &parameter)
		if err != nil {
			t.Fatalf("json.Unmarshal: %s", err.Error())
		}

		x := ModelCmdPromotion{
			Src:  parameter.Dst,
			Type: board.Queen,
		}

		data, err := json.Marshal(x)
		if err != nil {
			t.Fatalf("json.Marshal: %s", err.Error())
		}

		err = cl1.Do(Command{
			ID:   CmdPromotion,
			Data: data,
		})
		if err != nil {
			t.Fatalf("cl.Do: %s", err.Error())
		}

		v := gGame.b.Get(parameter.Dst)
		if p.T != board.Queen || (v != nil && v.T != board.Queen) {
			t.Fatalf("promotion dont work")
		}
		t.Logf("\n%s", gGame.b.String())
	}
}
