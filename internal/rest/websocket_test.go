package rest

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

var (
	ln    net.Listener
	nconn net.Conn
	wcl   *WsClient
)

func TestWsDial(t *testing.T) {

	var err error
	ln, err = net.Listen("tcp", "localhost:12345")
	if err != nil {
		t.Fatalf("net.Listen: %v", err)
	}

	x := make(chan error)

	go func(ln net.Listener) {

		y, err := ln.Accept()
		if err != nil {
			return
		}

		_, err = ws.Upgrade(y)
		if err != nil {
			return
		}

		wcl, err = UpgradeConn(y)
		if err != nil {
			x <- err
		} else {
			close(x)
		}
	}(ln)

	nconn, _, _, err = ws.Dial(context.Background(), "ws://"+ln.Addr().String())
	if err != nil {
		t.Fatalf("ws.Dial: %s", err.Error())
	}

	select {
	case <-time.After(time.Millisecond * 100):
		t.Fatalf("timeout")
	case err = <-x:
		if err != nil {
			t.Fatalf("UpgradeConn: %s", err.Error())
		}
	}

}

// test out first write, which is id
func TestWsWrite(t *testing.T) {

	/*
		data := []byte("test")
		go func() {
			wcl.Write(data)
		}()
	*/

	// guarntee it won't freeze
	time.Sleep(time.Millisecond * 10)

	nconn.SetReadDeadline(time.Now().Add(time.Millisecond * 100))

	b, err := wsutil.ReadServerText(nconn)
	if err != nil {
		t.Fatalf("wsutil.ReadServerText: %s", err.Error())
	}

	if string(b) != wcl.u.Token {
		t.Fatalf("ids do not match")
	}
	/*
		y := make(chan error)
		x := make(chan []byte)
		go func() {
			b, err := wsutil.ReadServerText(nconn)

			x <- b
			y <- err
		}()

		select {
		case <-time.After(time.Millisecond * 100):
			t.Fatalf("timeout")
		case err := <-y:
			if err != nil {
				t.Fatalf("wsutil.ReadServerText: %s", err.Error())
			}

			z := <-x
			if string(z) != wcl.u.Token {
				t.Fatalf("ids do not match")
			}
		}
	*/

}

// test out read, but it's not actually necessary
func TestWsRead(t *testing.T) {

	// guarntee it won't freeze
	time.Sleep(time.Millisecond * 10)

	data := []byte("test")

	nconn.SetWriteDeadline(time.Now().Add(time.Millisecond * 100))
	err := wsutil.WriteClientText(nconn, data)
	if err != nil {
		t.Fatalf("wsutil.WriteClientText: %s", err.Error())
	}

	y := wcl.ReadBytes()
	select {
	case <-time.After(time.Millisecond * 10):
	case b := <-y:
		str := string(b)

		if str != string(data) {
			t.Fatalf("data is not the same")
		}
	}

}

func TestWsClose(t *testing.T) {
	x := wcl.ClosedChannel()

	go func() {
		wcl.Close()
	}()

	select {
	case <-time.After(time.Millisecond * 100):
		t.Fatalf("timeout")
	case <-x:
		t.Logf("worked as expected")
	}
}

func TestWsDoesntPanic(t *testing.T) {
	wcl.Write([]byte("asddas"))
}