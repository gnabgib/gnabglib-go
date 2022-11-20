package checksum_test

import (
	"hash/crc32"
	"testing"
)

type stringCrc32 struct {
	s string
	h uint32
}

var stringCrc32s []stringCrc32

func init() {
	stringCrc32s = []stringCrc32{
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
}

func TestCrc32(t *testing.T) {
	for i := 0; i < len(stringCrc32s); i++ {
		found := crc32.ChecksumIEEE([]byte(stringCrc32s[i].s))
		if found != stringCrc32s[i].h {
			t.Fatalf("Hashing %v, expecting %v, got %v", stringCrc32s[i].s, stringCrc32s[i].h, found)
		}
	}
}
