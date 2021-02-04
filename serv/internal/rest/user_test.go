package rest

import (
	"io"
	"testing"

	"github.com/toms1441/chess/serv/internal/game"
)

var (
	rd2, wr1 = io.Pipe()
	cl1      = &game.Client{W: wr1}
	us       = User{}
)

func TestNewUser(t *testing.T) {
	us = addClient(cl1)
}

func TestUserClient(t *testing.T) {
	cl := us.Client()
	if cl1 != cl {
		t.Fatalf("not the same pointers")
	}
}

func TestUserDelete(t *testing.T) {
	us.Delete()
	if us.Client() != nil {
		t.Fatalf("delete does not delete")
	}
}
