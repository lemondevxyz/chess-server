package rest

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/toms1441/chess-server/internal/game"
	"github.com/toms1441/chess-server/internal/model"
	"github.com/toms1441/chess-server/internal/rest/auth"
)

type WsClient struct {
	net.Conn
	W        chan []byte
	r        []chan []byte
	c        []chan bool
	pingPong chan struct{}
	u        *User
	closed   bool
}

const (
	writeWait  = time.Second * 10
	pongWait   = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10
)

func (cl *WsClient) ReadBytes() chan []byte {
	x := make(chan []byte)
	cl.r = append(cl.r, x)

	return x
}

// ClosedChannel returns a channel that gets closed when the client is closed
func (cl *WsClient) ClosedChannel() chan bool {
	x := make(chan bool)
	cl.c = append(cl.c, x)

	return x
}

// Close closes the underlying websocket connection. Sudden determines if the websocket connection closed due to an unreceived pong frame, or if it wasn't sudden
func (cl *WsClient) Close(status ws.StatusCode, reason string) error {
	if cl.closed {
		return nil
	}

	cl.closed = true

	load := ws.NewCloseFrameBody(status, reason)
	wsutil.WriteServerMessage(cl.Conn, ws.OpClose, load)

	go func() {
		for _, v := range cl.r {
			close(v)
		}
		for _, v := range cl.c {
			close(v)
		}
	}()

	cl.u.Delete()
	cl.Conn.Close()

	return nil
}

func (cl *WsClient) Write(b []byte) (n int, err error) {
	cl.W <- b

	return 0, nil
}

func UpgradeConn(conn net.Conn) (*WsClient, error) {
	if conn == nil {
		return nil, fmt.Errorf("conn is nil")
	}

	cl := &WsClient{
		Conn: conn,
		W:    make(chan []byte, 8),
		c:    []chan bool{},
		r:    []chan []byte{},
	}

	// read any close messages
	go func() {
		for !cl.closed {
			header, err := ws.ReadHeader(conn)
			if err != nil {
				cl.Close(ws.StatusProtocolError, "")
				return
			}

			if header.OpCode == ws.OpClose {
				cl.Close(ws.StatusGoingAway, "")
				return
			}

			cl.Conn.SetReadDeadline(time.Now().Add(pongWait))
			if header.OpCode == ws.OpPong {
				continue
			}
		}
	}()

	go func() {
		defer func() {
			cl.Close(ws.StatusAbnormalClosure, "cannot send data")
		}()

		ticker := time.NewTicker(pingPeriod)
		for !cl.closed {
			select {
			case message := <-cl.W:
				cl.Conn.SetWriteDeadline(time.Now().Add(writeWait))

				writer := wsutil.NewWriter(conn, ws.StateServerSide, ws.OpText)
				_, err := writer.Write(message)
				if err != nil {
					return
				}
				err = writer.Flush()
				if err != nil {
					return
				}

				for i := 0; i < len(cl.W); i++ {
					_, err = writer.Write(<-cl.W)
					if err != nil {
						return
					}
					if writer.Flush() != nil {
						return
					}
				}
			case <-ticker.C:
				cl.Conn.SetWriteDeadline(time.Now().Add(writeWait))
				if err := wsutil.WriteServerMessage(cl.Conn, ws.OpPing, nil); err != nil {
					return
				}
			}
		}
	}()

	gc := &game.Client{
		W: cl,
	}

	u := AddClient(gc)
	cl.u = u

	// send token to the client
	ch := make(chan error)
	go func(u *User, ch chan error) {
		body, err := json.Marshal(u)
		if err != nil {
			ch <- err
		}

		upd := model.Order{
			ID:   model.OrCredentials,
			Data: body,
		}
		body, err = json.Marshal(upd)
		if err != nil {
			ch <- err
		}

		cl.Write(body)
		ch <- nil
		close(ch)
	}(u, ch)

	err := <-ch
	if err != nil {
		return nil, err
	}

	return cl, nil
}

func WebsocketHandler(w http.ResponseWriter, r *http.Request) {
	authuser := auth.Identify(r)
	fmt.Println(authuser)

	conn, _, _, err := ws.UpgradeHTTP(r, w)
	if err != nil {
		RespondError(w, http.StatusBadRequest, err)
		return
	}

	UpgradeConn(conn)
}

func WebsocketServe(ln net.Listener) {
	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}

		_, err = ws.DefaultUpgrader.Upgrade(conn)
		if err != nil {
			continue
		}

		UpgradeConn(conn)
	}
}
