package checksum_test

import (
	"testing"

	"github.com/gnabgib/gnablib-go/checksum/fletcher"
	"github.com/gnabgib/gnablib-go/endian"
)

type stringF16 struct {
	s string
	h uint16
}

type stringF32 struct {
	s string
	h uint32
}

type stringF64 struct {
	s string
	h uint64
}

var stringF16s []stringF16
var stringF32s []stringF32
var stringF64s []stringF64

func init() {
	stringF16s = []stringF16{
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
	stringF32s = []stringF32{
		//Wiki
		{"abcde", 0xf04fc729},
		{"abcdef", 0x56502d2a},
		{"abcdefgh", 0xebe19591},
		// Others
		{"", 0},
		{"\x00", 0},
		{"\x00\x00", 0},
		{"\x00\x00\x00", 0},
		{"\x00\x00\x00\x00", 0},
		{"\x01", 0x010001},
		{"\x01\x02", 0x02010201},
		{"\x01\x02\x03", 0x04050204},
		{"a", 0x00610061},
		{"a\x00", 0x00610061},
		{"a\x00\x00", 0x00c20061}, //0061+0000 = 0061 | 0061+0061 = 00C2
		{"a\x00\x00\x00", 0x00c20061},
		{"ab", 0x62616261},
		{"abc", 0xc52562c4},  //6261+0063 = 62C4 | 6261+62C4 = C525
		{"abcd", 0x2926c6c4}, //6261+6463 = C6C4 | 6261+C6C4 = 2926 (mod)
		{"f", 0x660066},
		{"fo", 0x6f666f66},     //6f66
		{"foo", 0xdf3b6fd5},    //6f66+006f = 6FD5 | 6f66+6FD5 = DF3B
		{"foob", 0x413cd1d5},   //6F66+626F = D1D5 | 6F66 + D1D5 = 413C (mod)
		{"fooba", 0x1373d236},  //6F66+626F = D1D5 +0061=D236 | 6F66+D1D5= 413C +D236=1373(mod)
		{"foobar", 0x85734437}, //6F66+626F = D1D5 +7261=4437 | 6F66+D1D5= 413C +4437=8573
		{"123456789", 0xdf09d509},
		{"foo bar baz٪☃🍣", 0xecb2f648},
		{"gnabgib", 0xb3f23a92},
	}
	stringF64s = []stringF64{
		//Wiki
		{"abcde", 0xC8C6C527646362C6},
		{"abcdef", 0xC8C72B276463C8C6},
		{"abcdefgh", 0x312E2B28CCCAC8C6},
		// Others
		{"", 0},
		{"\x00", 0},
		{"\x00\x00", 0},
		{"\x00\x00\x00", 0},
		{"\x00\x00\x00\x00", 0},
		{"\x01", 0x0000000100000001},
		{"\x01\x02", 0x0000020100000201},
		{"\x01\x02\x03", 0x0003020100030201},
		{"a", 0x0000006100000061},
		{"a\x00", 0x0000006100000061},
		{"a\x00\x00", 0x0000006100000061},
		{"a\x00\x00\x00", 0x0000006100000061},
		{"ab", 0x0000626100006261},
		{"abc", 0x0063626100636261},
		{"abcd", 0x6463626164636261},
		{"f", 0x0000006600000066},
		{"fo", 0x00006F6600006F66},
		{"foo", 0x006F6F66006F6F66},
		{"foob", 0x626F6F66626F6F66},
		{"fooba", 0xC4DEDF2D626F6FC7},  //626F6F66+00000061=  626F6FC7 | 626F6FC7 + 626F6F66 = C4DE DF2D
		{"foobar", 0xC4DF512D626FE1C7}, //626F6F66+00007261=  626F E1C7 | 626FE1C7 + 626F6F66 = C4DF 512D
		{"123456789", 0x0D0803376C6A689F},
		{"foo bar baz٪☃🍣", 0x5B253BF54182B4C6},
		{"gnabgib", 0xC525463562C3D7CE},
	}
}

func TestFletcher16(t *testing.T) {
	for i := 0; i < len(stringF16s); i++ {
		d := fletcher.New16()
		d.Write([]byte(stringF16s[i].s))
		found := d.Sum16()
		if found != stringF16s[i].h {
			t.Fatalf("Hashing %v, expecting %v, got %v", stringF16s[i].s, stringF16s[i].h, found)
		}
		found2 := endian.SourceCode.Uint16(d.Sum([]byte{}))
		if found2 != stringF16s[i].h {
			t.Fatalf("Hashing %v, expecting %v, got %v", stringF16s[i].s, stringF16s[i].h, found2)
		}
	}
}

func TestFletcher32(t *testing.T) {
	for i := 0; i < len(stringF32s); i++ {
		d := fletcher.New32()
		d.Write([]byte(stringF32s[i].s))
		found := d.Sum32()
		if found != stringF32s[i].h {
			t.Fatalf("Hashing %v, expecting %v, got %v", stringF32s[i].s, stringF32s[i].h, found)
		}
		found2 := endian.SourceCode.Uint32(d.Sum([]byte{}))
		if found2 != stringF32s[i].h {
			t.Fatalf("Hashing %v, expecting %v, got %v", stringF32s[i].s, stringF32s[i].h, found2)
		}
	}
}

func TestFletcher64(t *testing.T) {
	for i := 0; i < len(stringF64s); i++ {
		d := fletcher.New64()
		d.Write([]byte(stringF64s[i].s))
		found := d.Sum64()
		if found != stringF64s[i].h {
			t.Fatalf("Hashing %v, expecting %v, got %v", stringF64s[i].s, stringF64s[i].h, found)
		}
		found2 := endian.SourceCode.Uint64(d.Sum([]byte{}))
		if found2 != stringF64s[i].h {
			t.Fatalf("Hashing %v, expecting %v, got %v", stringF64s[i].s, stringF64s[i].h, found2)
		}
	}
}
