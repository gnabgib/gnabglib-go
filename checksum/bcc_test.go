package checksum_test

import (
	"github.com/gnabgib/gnablib-go/checksum/bcc"
	"testing"
)

type stringBcc struct {
	s string
	h byte
}

var stringBccs []stringBcc

func init() {
	stringBccs = []stringBcc{
		//The problem with BCC - size doesn't change the checksum
		{"", 0},
		{"\x00", 0},
		{"\x00\x00", 0},
		{"\x00\x00\x00", 0},
		{"\x00\x00\x00\x00", 0},
		{"\x01\x02", 3},
		{"\x01\x02\x03", 0},//3 undoes all the hard work of 1, 2
		{"Wikipedia",0x45},
		{"Freddy",0x28},
		{"a", 0x61},
		{"ab", 0x3},
		{"abc", 0x60},
		{"abcd", 0x4},
		{"abcde", 0x61},
		{"abcdef", 0x7},
		{"abcdefg", 0x60},
		{"abcdefgh", 0x8},
		{"gnabgib",0x66},
		{"BCCXOR",0x07},
		{"\xFF\xEE\xDD", 0xcc},
		{"f", 0x66},
		{"fo", 0x9},
		{"foo", 0x66},
		{"foob", 0x4},
		{"fooba", 0x65},
		{"foobar", 0x17},
		{"foo bar baz٪☃🍣", 0xa5},
	}
}

func TestBcc(t *testing.T) {
	for i := 0; i < len(stringBccs); i++ {
		d := bcc.New()
		d.Write([]byte(stringBccs[i].s))
		found := d.Sum8()
		//found := bcc.New()
		//found := bcc.Checksum([]byte(stringBccs[i].s))
		if found != stringBccs[i].h {
			t.Fatalf("Hashing %v, expecting %v, got %v", stringBccs[i].s, stringBccs[i].h, found)
		}
	}
}
