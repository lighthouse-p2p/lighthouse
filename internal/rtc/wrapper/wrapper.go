package wrapper

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"time"

	"github.com/pion/webrtc/v2"
)

// NilAddr is an empty address
type NilAddr struct {
	ID string
}

// Network needs to be exported to be compatible with net.Addr
func (a *NilAddr) Network() string {
	return "WebRTC"
}

func (a *NilAddr) String() string {
	return ""
}

// Conn is used to wrap the connection in a golang net.Conn
type Conn struct {
	dc         *webrtc.DataChannel
	laddr      net.Addr
	raddr      net.Addr
	p          *io.PipeReader
	isClosed   bool
	closeMutex sync.Mutex
}

// WrapConn wraps the datachannel in a net.Conn
func WrapConn(dc *webrtc.DataChannel, laddr, raddr net.Addr) (net.Conn, error) {
	r, w := io.Pipe()

	conn := &Conn{
		dc:    dc,
		laddr: laddr,
		raddr: raddr,
		p:     r,
	}

	go func() {
		dc.OnMessage(func(msg webrtc.DataChannelMessage) {
			if !msg.IsString {
				w.Write(msg.Data)
			}
		})
	}()

	return conn, nil
}

// Read reads data from the underlying the data channel
func (c *Conn) Read(b []byte) (int, error) {
	if c.isClosed {
		return 0, errors.New("read on closed conn")
	}
	i, err := c.p.Read(b)
	return i, err
}

// Write writes the data to the underlying data channel
func (c *Conn) Write(b []byte) (int, error) {
	if c.isClosed {
		return 0, errors.New("write on closed conn")
	}
	err := c.dc.Send(b)
	if err != nil {
		return 0, err
	}
	return len(b), nil
}

// Close closes the datachannel and peerconnection
func (c *Conn) Close() error {
	if c.isClosed {
		return errors.New("close on closed conn")
	}

	// Prevent concurrent closing of the datachannel
	c.closeMutex.Lock()
	c.isClosed = true

	defer c.closeMutex.Unlock()

	// Unblock readers
	err := c.p.Close()
	if err != nil {
		fmt.Println("failed to close pipe:", err)
		return err
	}

	c.dc.Close()
	return nil
}

// LocalAddr TODO
func (c *Conn) LocalAddr() net.Addr {
	return c.laddr
}

// RemoteAddr TODO
func (c *Conn) RemoteAddr() net.Addr {
	return c.raddr
}

// SetDeadline TODO
func (c *Conn) SetDeadline(t time.Time) error {
	panic("TODO")
}

// SetReadDeadline TODO
func (c *Conn) SetReadDeadline(t time.Time) error {
	panic("TODO")

}

// SetWriteDeadline TODO
func (c *Conn) SetWriteDeadline(t time.Time) error {
	panic("TODO")
}

// JoinStreams can proxy data from stream 1 to stream 2 and vice-versa
// statsCallback will give the amount of data copied into the first stream
func JoinStreams(c1, c2 net.Conn, statsCallback func(stats int64)) {
	defer c1.Close()
	defer c2.Close()

	errc := make(chan error, 2)
	statsc := make(chan int64)

	go func() {
		stats, err := io.Copy(c1, c2)
		statsc <- stats - 150
		errc <- err
	}()
	go func() {
		_, err := io.Copy(c2, c1)
		errc <- err
	}()

	statsCallback(<-statsc)

	err := <-errc

	if err != nil && !errors.Is(err, io.EOF) && !errors.Is(err, context.Canceled) {
		log.Printf("Stream err %s\n", err)
	}
}
