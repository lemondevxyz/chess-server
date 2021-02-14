package rest

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/toms1441/chess/serv/internal/game"
)

type WsClient struct {
	net.Conn
	W      chan []byte
	r      []chan []byte
	c      []chan bool
	u      User
	closed bool
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

func (cl *WsClient) Close() error {
	if cl.closed {
		return nil
	}

	cl.closed = true
	wsutil.WriteServerMessage(cl.Conn, ws.OpClose, nil)

	for _, v := range cl.c {
		close(v)
	}

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
	}

	// read any close messages
	go func() {
		defer cl.Close()

		for {
			body, opcode, err := wsutil.ReadClientData(cl.Conn)
			//fmt.Println(string(body), opcode, err)
			if err != nil || opcode == ws.OpClose {
				return
			} else if opcode == ws.OpText {
				for _, v := range cl.r {
					v <- body
				}
			}
		}
	}()

	go func() {
		defer func() {
			cl.Close()
		}()

		ticker := time.NewTicker(pingPeriod)
		for {
			select {
			case message := <-cl.W:
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
				cl.Conn.SetReadDeadline(time.Now().Add(writeWait))
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
	go func() {
		cl.Write([]byte(u.Token))
	}()

	return cl, nil
}

func WebsocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, _, _, err := ws.UpgradeHTTP(r, w)
	if err != nil {
		respondError(w, http.StatusBadRequest, err)
		return
	}

	UpgradeConn(conn)
}
