package fletcher

import (
	"testing"

	"github.com/gnabgib/gnablib-go/endian"
)

var fletcher64tests = []struct {
	s string
	c uint64
}{
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

func TestFletcher64(t *testing.T) {
	for _, rec := range fletcher64tests {
		d := New64()
		d.Write([]byte(rec.s))
		found := d.Sum64()
		if found != rec.c {
			t.Fatalf("Hashing %v, expecting %v, got %v", rec.s, rec.c, found)
		}
		found2 := endian.SourceCode.Uint64(d.Sum([]byte{}))
		if found2 != rec.c {
			t.Fatalf("Hashing %v, expecting %v, got %v", rec.s, rec.c, found2)
		}
	}
}
