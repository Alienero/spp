package spp

import (
	"bufio"
	"net"
	"time"
)

// The first num is the pack len's size , the second
// num is the type's size,second  . They all support byte, uint16, uint32, uint64.
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
	sizeByte    byte
	packLenSize int
	typLenSize  int

	readDeadline  time.Duration
	writeDeadline time.Duration

	conn *net.TCPConn
	w    *bufio.Writer
	r    *bufio.Reader

	// Read and Write size temple cache
	rSize []byte
	wSize []byte

	pakcLenLimit int
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
	c.packLenSize = c.getSize(c.sizeByte)[0]
	c.typLenSize = c.getSize(c.sizeByte)[1]
	return c.sizeByte, nil
}

func (c *Conn) SetReadDeadline(n time.Duration) {
	c.readDeadline = n
}
func (c *Conn) SetWeadDeadline(n time.Duration) {
	c.writeDeadline = n
}
