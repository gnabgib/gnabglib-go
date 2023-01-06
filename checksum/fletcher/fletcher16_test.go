package fletcher

import (
	"testing"

	"github.com/gnabgib/gnablib-go/endian"
)

var fletcher16tests = []struct {
	s string
	c uint16
}{
	//Wiki
	{"abcde", 0xc8f0},
	{"abcdef", 0x2057},
	{"abcdefgh", 0x0627},
	// Others
	{"", 0},
	{"\x00", 0},
	{"\x00\x00", 0},
	{"\x00\x00\x00", 0},
	{"\x00\x00\x00\x00", 0},
	{"\x01", 0x0101},
	{"\x01\x02", 0x0403},
	{"\x01\x02\x03", 0xa06},
	{"a", 0x6161},
	{"ab", 0x25c3},
	{"abc", 0x4c27},
	{"abcd", 0xd78b},
	{"f", 0x6666},
	{"fo", 0x3cd5},
	{"foo", 0x8145},
	{"foob", 0x29a7},
	{"fooba", 0x3209},
	{"foobar", 0xad7b},
	{"123456789", 0x1ede},
	{"foo bar baz٪☃🍣", 0x493f},
	{"gnabgib", 0x46cc},
	{"Z", 0x5a5a},
}

func TestFletcher16(t *testing.T) {
	for _, rec := range fletcher16tests {
		d := New16()
		d.Write([]byte(rec.s))
		found := d.Sum16()
		if found != rec.c {
			t.Errorf("Hashing %v, expecting %v, got %v", rec.s, rec.c, found)
		}
		found2 := endian.SourceCode.Uint16(d.Sum([]byte{}))
		if found2 != rec.c {
			t.Errorf("Hashing %v, expecting %v, got %v", rec.s, rec.c, found2)
		}
	}
}
