package game

import (
	"io"
)

// Note: Future contributors, beware of io.Pipe freezing the entire test, wrapping Write operations around a goroutine would make the tests unpredictable.
// Just watch out for different order messages, try mixing the clientRead. Or increase the amount, or decrease it.
var (
	rd1, wr1 = io.Pipe()
	rd2, wr2 = io.Pipe()

	cl1 = &Client{
		W:  wr1,
		p1: true,
		id: "a",
	}

	cl2 = &Client{
		W:  wr2,
		p1: false,
		id: "b",
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
