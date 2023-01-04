package checksum

//Different package space because of some of the checksums depending on the checksum
// namespace, and that leads to a circular import reference on the test files (if same namespace)

import (
	"hash/adler32"
	"testing"
)

var adlerTests = []struct {
	s string
	c uint32
}{
	//Wiki
	{"Wikipedia", 0x11E60398},
	//Others
	{"abcde", 0x05C801F0},
	{"abcdef", 0x081E0256},
	{"abcdefgh", 0x0E000325},
	{"", 1},
	{"\x00", 0x10001},
	{"\x00\x00", 0x20001},
	{"\x00\x00\x00", 0x30001},
	{"\x00\x00\x00\x00", 0x40001},
	{"\x01", 0x20002},
	{"\x01\x02", 0x60004},
	{"\x01\x02\x03", 0xD0007},
	{"a", 0x00620062},
	{"a\x00", 0xC40062},
	{"a\x00\x00", 0x01260062},
	{"a\x00\x00\x00", 0x01880062},
	{"ab", 0x012600C4},
	{"abc", 0x024D0127},
	{"abcd", 0x03D8018B},
	{"f", 0x00670067},
	{"fo", 0x013D00D6},
	{"foo", 0x02820145},
	{"foob", 0x042901A7},
	{"fooba", 0x06310208},
	{"foobar", 0x08AB027A},
	{"123456789", 0x091E01DE},
	{"foo bar baz٪☃🍣", 0x5c010a36},
	{"gnabgib", 0x0B4202CB},
}

func TestAdler(t *testing.T) {
	for _, rec := range adlerTests {
		found := adler32.Checksum([]byte(rec.s))
		if found != rec.c {
			t.Fatalf("Hashing %v, expecting %v, got %v", rec.s, rec.c, found)
		}
	}
}
