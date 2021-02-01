package game

import (
	"encoding/json"
	"fmt"
	"io"
	"testing"
	"time"

	"github.com/toms1441/chess/internal/board"
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
		ID:   CmdSendMessage,
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
	go func() {
		p := gGame.b.Get(board.Point{X: 3, Y: 1})
		p.X = 7
		p.Y = 1
		gGame.b.Set(p)
	}()

	select {
	case <-time.After(time.Millisecond * 10):
		t.Fatalf("timeout")
	}
}
