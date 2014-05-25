package spp

import (
	"bufio"
	"net"
	"time"
)

// The first num is the type's size,second num is the pack len's
// size . They all support byte, uint16, uint32, uint64.
var sizeType *[16][2]int = &[16][2]int{
	{1, 1}, {1, 2}, {1, 3}, {1, 4},
	{2, 1}, {2, 2}, {2, 3}, {2, 4},
	{3, 1}, {3, 2}, {3, 3}, {3, 4},
	{4, 1}, {4, 2}, {4, 3}, {4, 4},
}

type SppError struct {
	err string
}

func (err *SppError) Error() string {
	return err.err
}

type Conn struct {
	// It's the first byte, set the type field and pack
	// len field size .
	sizeByte byte

	readDeadline  time.Duration
	writeDeadline time.Duration

	conn *net.TCPConn
	w    *bufio.Writer
	r    *bufio.Reader

	// Read and Write size temple cache
	rSize []byte
	wSize []byte
}

func NewConn(c *net.TCPConn) *Conn {
	conn := &Conn{
		conn:  c,
		w:     bufio.NewWriter(c),
		r:     bufio.NewReader(c),
		rSize: make([]byte, 8),
		wSize: make([]byte, 8),
	}
	return conn
}
func (c *Conn) getSize(size byte) *[2]int {
	return &sizeType[int(size)]
}
func (c *Conn) getSizeByte(length int, typ int) byte {
	return byte((length-1)*4 + (typ - 1))
}

func (c *Conn) SetSizeByte(length int, typ int) (byte, error) {
	if length > 0 && length < 5 {
		return 0, &SppError{"length is out of range!"}
	}
	if typ > 0 && typ < 5 {
		return 0, &SppError{"typ is out of range!"}
	}
	c.sizeByte = c.getSizeByte(length, typ)
	return c.sizeByte, nil
}
