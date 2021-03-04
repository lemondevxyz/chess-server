package game

import (
	"io"
)

var (
	rd1, wr1 = io.Pipe()
	rd2, wr2 = io.Pipe()

	cl1 = &Client{
		W:   wr1,
		num: 1,
		id:  "a",
	}

	cl2 = &Client{
		W:   wr2,
		num: 2,
		id:  "b",
	}
)

func clientRead(rd *io.PipeReader) chan []byte {

	ch := make(chan []byte)

	go func() {
		body := make([]byte, 128)
		n, err := rd.Read(body)
		if err != nil {
			ch <- nil
		} else {
			ch <- body[:n]
		}
	}()

	return ch
}

func resetPipe() {
	rd1, wr1 = io.Pipe()
	rd2, wr2 = io.Pipe()

	cl1.W = wr1
	cl2.W = wr2
}
