package spp

import (
	"encoding/binary"
	"errors"
	"time"
)

type Pack struct {
	Size byte
	Typ  int
	Body []byte
}

func (c *Conn) setPack(packLenSize int, typ int, typSize int, body []byte) (*Pack, error) {
	pack := new(pack)
	if typSize == 1 && packLenSize == 2 {
		if typ <= 20 {
			return nil, &SppError{"typ<=20"}
		}
		pack.Size = 1
		pack.Typ = typ
	} else {
		pack.Size = c.getSizeByte(packLenSize, typSize)
		pack.Typ = typ
	}
	pack.Body = body
	return pack, nil
}
func (c *Conn) SetDefaultPack(typ int, body []byte) (*Pack, error) {
	return c.setPack(2, typ, 1, body)
}
func (c *Conn) SetStandardPack(typ int, body []byte) (*Pack, error) {
	return c.setPack(c.packLenSize, typ, c.typLenSize, body)
}
func (c *Conn) SetTempPack(packLenSize int, typ int, typSize int, body []byte) (*Pack, error) {
	if (packLenSize < 0 || packLenSize > 4) || (typSize < 0 || typSize > 4) {
		return nil, &SppError{"length is out of range!"}
	}
	return c.setPack(packLenSize, typ, typSize, body)
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

	if ts <= 20 {
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
		packLen, _ := c.getInt(c.rSize[:size[0]])

		if packLen > c.pakcLenLimit {
			return nil, &SppError{"pack length is too long!"}
		}
		// Read the whole pack
		data := make([]byte, packLen)
		err = c.readAll(data, packLen)
		if err != nil {
			return nil, err
		}
		pack.Typ, _ = c.getInt(data[:size[1]])

		pack.Body = data[size[1]:]

	} else {
		// Pack length's size is 2 bytes, not include pack type
		// This byte is type number
		pack.Typ, _ = int(ts) // type number
		ts = 1
		size := c.getSize(ts)
		// Read the pack len
		err = c.readAll(c.rSize, size[1])
		if err != nil {
			return nil, err
		}
		// Read the all pack
		packLen, _ := c.getInt(c.rSize[:size[1]])

		data := make([]byte, packLen)
		err = c.readAll(data, size[1])
		if err != nil {
			return nil, err
		}
		pack.Body = data
	}
	pack.Size = ts

	return pack, nil
}
func (c *Conn) readAll(data []byte, size int) (err error) {
	hasRead := 0
	read := 0
	for {
		read, err = c.r.Read(data[hasRead:size])
		if err != nil {
			return
		}
		hasRead += read
		if hasRead == size {
			break
		}
	}
	return
}

func (c *Conn) WritePack(pack *Pack) error {
	var err error
	err = c.conn.SetWriteDeadline(time.Now().Add(c.writeDeadline))
	if err != nil {
		return err
	}
	// parse pack
	if pack.Size == 1 {
		// Type size is 1 byte
		// PackLength is 2 bytes
		err = c.w.WriteByte(byte(pack.Typ))
		if err != nil {
			return err
		}
		// Write pack length(except type's length)
		err = c.writeAll(c.getBytes(len(pack.Body), 2))
		if err != nil {
			return err
		}

	} else {
		// Write the flag byte
		err = c.w.WriteByte(pack.Size)
		if err != nil {
			return err
		}
		// Write pack length(include pack length and type length)
		err = c.writeAll(c.getBytes(len(pack.Body)+c.getSize(pack.Size)[1], c.getSize(pack.Size[0])))
		if err != nil {
			return err
		}
		// Write pack type
		err = c.writeAll(c.getBytes(pack.Typ, c.getSize(pack.Size)[1]))
		if err != nil {
			return err
		}
	}
	// Write the pack body
	return c.writeAll(pack.Body)

}
func (c *Conn) writeAll(data []byte) (err error) {
	write := 0
	for hasWrite := 0; hasWrite < len(data); {
		write, err = c.w.Write(data[hasWrite:])
		if err != nil {
			return
		}
		hasWrite += write
	}
	err = c.w.Flush()
	return
}

func (c *Conn) getInt(data []byte) (i int, err error) {
	j, n := binary.Varint(data)
	if n <= 0 {
		err = &SppError{"Out of range!"}
	} else {
		i = int(j)
	}
	return
}
func (c *Conn) getBytes(n int, size int) []byte {
	defer func() {
		if err := recover(); err != nil {
		}
	}()
	buf := make([]byte, size)
	binary.PutVarint(buf, n)
	return buf
}
