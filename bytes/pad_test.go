package bytes

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/gnabgib/gnablib-go/encoding/hex"
	"github.com/gnabgib/gnablib-go/test"
)


func padEqual(t *testing.T, expect, found []byte, size int) {
	if !bytes.Equal(found, expect) {
		test.StringMatchTitle(
			t,
			fmt.Sprintf("Padding to %d",size),
			"",
			hex.FromBytes(expect),
			hex.FromBytes(found))
	}
}

func TestPadLE16(t *testing.T) {
	var tests = []struct {
		in, out []byte
	}{
		{[]byte{}, []byte{0, 0}},
		{[]byte{0}, []byte{0, 0}},
		{[]byte{0x10}, []byte{0x10, 0}},
		{[]byte{0, 0x10}, []byte{0, 0x10}},
		{[]byte{0x01, 0x00}, []byte{0x01, 0}},
		{[]byte{0x10, 0x00}, []byte{0x10, 0}},
		{[]byte{1}, []byte{1, 0}},
		{[]byte{1, 2}, []byte{1, 2}},
		{[]byte{1, 2, 3}, []byte{1, 2}},
		{[]byte{1, 2, 3, 4}, []byte{1, 2}},
	}
	const size=2
	dst:=make([]byte,size)
	for _, rec := range tests {
		padLE(dst,rec.in,size)
		padEqual(t,rec.out,dst,size)
	}
}

func TestPadBE16(t *testing.T) {
	var tests = []struct {
		in, out []byte
	}{
		{[]byte{}, []byte{0, 0}},
		{[]byte{0}, []byte{0, 0}},
		{[]byte{0x10}, []byte{0, 0x10}},
		{[]byte{0, 0x10}, []byte{0, 0x10}},
		{[]byte{0x01, 0}, []byte{0x01, 0}},
		{[]byte{0x10, 0}, []byte{0x10, 0}},
		{[]byte{1}, []byte{0, 1}},
		{[]byte{1, 2}, []byte{1, 2}},
		{[]byte{1, 2, 3}, []byte{1, 2}},
		{[]byte{1, 2, 3, 4}, []byte{1, 2}},
	}
	const size=2
	dst:=make([]byte,size)
	for _, rec := range tests {
		padBE(dst,rec.in,size)
		padEqual(t,rec.out,dst,size)
	}
}

func TestPadLE32(t *testing.T) {
	var tests = []struct {
		in, out []byte
	}{
		{[]byte{}, []byte{0, 0, 0, 0}},
		{[]byte{0}, []byte{0, 0, 0, 0}},
		{[]byte{1}, []byte{1, 0, 0, 0}},
		{[]byte{0x10}, []byte{0x10, 0, 0, 0}},
		{[]byte{0, 0x10}, []byte{0, 0x10, 0, 0}},
		{[]byte{0x01, 0}, []byte{0x01, 0, 0, 0}},
		{[]byte{0x10, 0}, []byte{0x10, 0, 0, 0}},
		{[]byte{1}, []byte{1, 0, 0, 0}},
		{[]byte{1, 2}, []byte{1, 2, 0, 0}},
		{[]byte{1, 2, 3}, []byte{1, 2, 3, 0}},
		{[]byte{1, 2, 3, 4}, []byte{1, 2, 3, 4}},
		{[]byte{1, 2, 3, 4, 5}, []byte{1, 2, 3, 4}},
		{[]byte{1, 2, 3, 4, 5, 6}, []byte{1, 2, 3, 4}},
		{[]byte{1, 2, 3, 4, 5, 6, 7}, []byte{1, 2, 3, 4}},
		{[]byte{1, 2, 3, 4, 5, 6, 7, 8}, []byte{1, 2, 3, 4}},
		{[]byte{1, 0}, []byte{1, 0, 0, 0}},
		{[]byte{1, 0, 0}, []byte{1, 0, 0, 0}},
		{[]byte{1, 0, 0, 0}, []byte{1, 0, 0, 0}},
		{[]byte{1, 0, 0, 0, 0}, []byte{1, 0, 0, 0}},
		{[]byte{1, 0, 0, 0, 0, 0}, []byte{1, 0, 0, 0}},
		{[]byte{1, 0, 0, 0, 0, 0, 0}, []byte{1, 0, 0, 0}},
		{[]byte{1, 0, 0, 0, 0, 0, 0, 0}, []byte{1, 0, 0, 0}},
	}
	const size=4
	dst:=make([]byte,size)
	for _, rec := range tests {
		padLE(dst,rec.in,size)
		padEqual(t,rec.out,dst,size)
	}
}

func TestPadBE32(t *testing.T) {
	var tests = []struct {
		in, out []byte
	}{
		{[]byte{}, []byte{0, 0, 0, 0}},
		{[]byte{0}, []byte{0, 0, 0, 0}},
		{[]byte{1}, []byte{0, 0, 0, 1}},
		{[]byte{0x10}, []byte{0, 0, 0, 0x10}},
		{[]byte{0, 0x10}, []byte{0, 0, 0, 0x10}},
		{[]byte{0x01, 0}, []byte{0, 0, 1, 0}},
		{[]byte{0x10, 0}, []byte{0, 0, 0x10, 0}},
		{[]byte{1}, []byte{0, 0, 0, 1}},
		{[]byte{1, 2}, []byte{0, 0, 1, 2}},
		{[]byte{1, 2, 3}, []byte{0, 1, 2, 3}},
		{[]byte{1, 2, 3, 4}, []byte{1, 2, 3, 4}},
		{[]byte{1, 2, 3, 4, 5}, []byte{1, 2, 3, 4}},
		{[]byte{1, 2, 3, 4, 5, 6}, []byte{1, 2, 3, 4}},
		{[]byte{1, 2, 3, 4, 5, 6, 7}, []byte{1, 2, 3, 4}},
		{[]byte{1, 2, 3, 4, 5, 6, 7, 8}, []byte{1, 2, 3, 4}},
		{[]byte{1, 0}, []byte{0, 0, 1, 0}},
		{[]byte{1, 0, 0}, []byte{0, 1, 0, 0}},
		{[]byte{1, 0, 0, 0}, []byte{1, 0, 0, 0}},
		{[]byte{1, 0, 0, 0, 0}, []byte{1, 0, 0, 0}},
		{[]byte{1, 0, 0, 0, 0, 0}, []byte{1, 0, 0, 0}},
		{[]byte{1, 0, 0, 0, 0, 0, 0}, []byte{1, 0, 0, 0}},
		{[]byte{1, 0, 0, 0, 0, 0, 0, 0}, []byte{1, 0, 0, 0}},
	}
	const size=4
	dst:=make([]byte,size)
	for _, rec := range tests {
		padBE(dst,rec.in,size)
		padEqual(t,rec.out,dst,size)
	}
}
