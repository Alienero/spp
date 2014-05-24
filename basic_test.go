package spp

import (
	"testing"
)

func TestSize(t *testing.T) {
	c := new(Conn)
	var arr [16][2]int
	count := 0
	for i := 1; i <= 4; i++ {
		for j := 1; j <= 4; j++ {
			arr[count] = *c.getSize(c.getSizeByte(i, j))
			count++
		}
	}
	for i := 0; i < 16; i++ {
		if arr[i][0] != sizeType[i][0] || arr[i][1] != sizeType[i][1] {
			t.Errorf("Error nq! \n")
		}
	}
}
