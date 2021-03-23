package headless

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"time"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/toms1441/chess-server/internal/order"
)

type Client struct {
	conn   net.Conn
	closed bool
	//listen [] chanorder.Order
}

func NewClient(addr string) (*Client, error) {
	ctx, _ := context.WithDeadline(context.Background(), time.Now().Add(time.Millisecond*250))
	conn, _, _, err := ws.Dial(ctx, addr)
	if err != nil {
		return nil, err
	}

	return &Client{
		conn: conn,
	}, nil
}

func (c *Client) Close() error {
	if c.closed {
		return nil
	}

	c.conn.Close()
	return nil
}

func (c *Client) WriteCommand(id uint8, obj interface{}) error {
	data, err := json.Marshal(obj)
	if err != nil {
		return fmt.Errorf("json.Marshal: %s", err.Error())
	}

	x, err := json.Marshal(order.Order{
		ID:   id,
		Data: data,
	})
	if err != nil {
		return fmt.Errorf("json.Marshal(order.Order): %s", err.Error())
	}

	c.conn.SetWriteDeadline(time.Now().Add(time.Millisecond * 100))
	return wsutil.WriteClientText(c.conn, x)
}
