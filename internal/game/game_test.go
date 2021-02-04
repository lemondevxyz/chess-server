package game

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/toms1441/chess/internal/board"
)

var gGame, _ = NewGame(cl1, cl2)

func TestTurns(t *testing.T) {
	resetPipe()
	go func() {
		gGame, _ = NewGame(cl1, cl2)
	}()

	select {
	case <-time.After(time.Millisecond * 100):
		t.Fatalf("timeout")
	case body := <-clientRead(rd1):
		t.Logf("\n%s", string(body))

		x := &ModelUpdateTurn{}
		u := &Update{}

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
		body, err := json.Marshal(ModelCmdPiece{
			Src: board.Point{1, 1},
			Dst: board.Point{2, 1},
		})

		if err != nil {
			cherr <- fmt.Errorf("json.Marshal: %s", err.Error())
			return
		}

		c := Command{
			ID:   CmdPiece,
			Data: body,
		}

		err = cl1.Do(c)
		if err != nil {
			cherr <- err
			return
		}

		cherr <- nil
	}()

	select {
	case <-time.After(time.Millisecond * 100):
		t.Fatalf("timeout")
	case <-cherr:
		select {
		case <-time.After(time.Millisecond * 100):
			t.Fatalf("timeout")
		case body := <-clientRead(rd1):
			t.Logf("\n%s", string(body))

			x := &ModelUpdateTurn{}
			u := &Update{}

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
		}
	}
}