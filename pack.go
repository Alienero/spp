package spp

import (
	"time"
)

type Pack struct {
	Size byte
	Typ  int
	Body []byte
}

func (c *Conn) ReadPack() (*Pack, error) {
	err := c.conn.SetReadDeadline(time.Now().Add(c.readDeadline))
	if err != nil {
		return nil, err
	}
	ts, err := c.r.ReadByte()
	if err != nil {
		return nil, err
	}

	pack := new(Pack)

	if ts >= 20 {
		if ts < 0 || ts > 15 {
			return nil, &SppError{"size byte out of range!"}
		}
		size := c.getSize(ts)
		// read pack length
		err = c.readAll(c.rSize, size[0])
		if err != nil {
			return nil, err
		}
		// convert
	} else {
		ts = 1
	}

	return pack, nil
}
func (c *Conn) readAll(data []byte, size int) error {
	hasRead := 0
	read := 0
	for {
		read, err = c.r.Read(data[hasRead:size])
		if err != nil {
			return nil, err
		}
		hasRead += read
		if hasRead == size {
			break
		}
	}
}
func (c *Conn) getInt(data []byte, length int) {
	switch length {

	}
}
