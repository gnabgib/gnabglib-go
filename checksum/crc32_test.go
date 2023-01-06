package checksum

import (
	"hash/crc32"
	"testing"
)

var crc32tests = []struct {
	s string
	c uint32
}{
	{"", 0},
	{"a", 0xE8B7BE43},
	{"f", 0x76D32BE0},
	{"fo", 0xAF73A217},
	{"foo", 0x8C736521},
	{"foob", 0x3D5B8CC2},
	{"fooba", 0x9DE04653},
	{"foobar", 0x9EF61F95},
	{"abcde", 0x8587D865},
	{"123456789", 0xCBF43926},
	{"foo bar baz٪☃🍣", 0x5B4B18F3},
	{"gnabgib", 0x8EE5AF75},
}

func TestCrc32(t *testing.T) {
	for _, rec := range crc32tests {
		found := crc32.ChecksumIEEE([]byte(rec.s))
		if found != rec.c {
			t.Errorf("Hashing %v, expecting %v, got %v", rec.s, rec.c, found)
		}
	}
}
