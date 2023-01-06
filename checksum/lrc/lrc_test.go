package lrc

import (
	"testing"
)

var lrcTests = []struct {
	s string
	c byte
}{
	//The problem with LRC - size doesn't change the checksum
	{"", 0},
	{"\x00", 0},
	{"\x00\x00", 0},
	{"\x00\x00\x00", 0},
	{"\x00\x00\x00\x00", 0},
	{"\x01\x02", 0xfd},
	{"\x01\x02\x03", 0xfa},
	{"Wikipedia", 0x69},
	{"a", 0x9f},
	{"ab", 0x3d},
	{"abc", 0xda},
	{"abcd", 0x76},
	{"abcde", 0x11},
	{"abcdef", 0xab},
	{"abcdefg", 0x44},
	{"abcdefgh", 0xdc},
	{"gnabgib", 0x36},
	{"\xFF\xEE\xDD", 0x36},
	{"f", 0x9a},
	{"fo", 0x2b},
	{"foo", 0xbc},
	{"foob", 0x5a},
	{"fooba", 0xf9},
	{"foobar", 0x87},
	{"foo bar baz٪☃🍣", 0xcb},
}

func TestLrc(t *testing.T) {
	for _, rec := range lrcTests {
		d := New()
		d.Write([]byte(rec.s))
		found := d.Sum8()
		if found != rec.c {
			t.Errorf("Hashing %v, expecting %v, got %v", rec.s, rec.c, found)
		}
	}
}
