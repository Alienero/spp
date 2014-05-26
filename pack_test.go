package spp

import (
	"fmt"
	"net"
	"testing"
)

func TestPack(t *testing.T) {
	// Listen at port 8080
	l, _ := net.Listen("tcp", ":8080")
	go func() {
		tc, _ := net.Dial("tcp", "127.0.0.1:8080")
		ttc := tc.(*net.TCPConn)
		cc := NewConn(ttc)
		// pack, err := cc.SetDefaultPack(21, []byte("goodhello"))
		pack, err := cc.SetTempPack(4, 2, 3, []byte("小鸟"))
		if err != nil {
			t.Error(err)
		}
		fmt.Println("Pack body is ", string(pack.Body))
		err = cc.WritePack(pack)
		if err != nil {
			t.Error(err)
		}
	}()
	conn, _ := l.Accept()
	con := conn.(*net.TCPConn)
	c := NewConn(con)
	pack, err := c.ReadPack()
	if err != nil {
		t.Error(err)
	}
	// Print the pack
	fmt.Println("Pack size is ", pack.Size, "Pack type is ", pack.Typ, "Pack body is ", string(pack.Body))
}
