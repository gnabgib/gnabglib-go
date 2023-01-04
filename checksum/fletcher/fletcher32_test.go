package fletcher

import (
	"testing"

	"github.com/gnabgib/gnablib-go/endian"
)

var fletcher32tests = []struct {
	s string
	c uint32
}{
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

func TestFletcher32(t *testing.T) {
	for _, rec := range fletcher32tests {
		d := New32()
		d.Write([]byte(rec.s))
		found := d.Sum32()
		if found != rec.c {
			t.Fatalf("Hashing %v, expecting %v, got %v", rec.s, rec.c, found)
		}
		found2 := endian.SourceCode.Uint32(d.Sum([]byte{}))
		if found2 != rec.c {
			t.Fatalf("Hashing %v, expecting %v, got %v", rec.s, rec.c, found2)
		}
	}
}
