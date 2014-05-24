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
	if ts < 0 || ts > 15 {
		return nil, &SppError{"size byte out of range!"}
	}
	pack := &Pack{
		Size: ts,
	}
	// size := c.getSize(ts)
	return pack, nil
}
