package codegen

import (
	"testing"

	"github.com/gnabgib/gnablib-go/test"
)

type bytesHexPair struct {
	bytes []byte
	hex   string
}

var bytesStrPairs = []bytesHexPair{
	{[]byte{}, "\"\""},
	{[]byte{0}, "\"\\x00\""},
	{[]byte{0xde, 0xad, 0xbe, 0xef}, "\"\\xDE\\xAD\\xBE\\xEF\""},
}

var bytesStrSep4 = []bytesHexPair{
	{[]byte{0xf0, 0xf1, 0xf2, 0xf3, 0xf4, 0xf5, 0xf6, 0xf7},
		"\"\\xF0\\xF1\\xF2\\xF3\"+\"\\xF4\\xF5\\xF6\\xF7\""},
	{[]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15},
		"\"\\x00\\x01\\x02\\x03\"+\"\\x04\\x05\\x06\\x07\"+\"\\x08\\x09\\x0A\\x0B\"+\"\\x0C\\x0D\\x0E\\x0F\""},
	{[]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		"\"\\x00\\x00\\x00\\x00\"+\"\\x00\\x00\\x00\\x00\"+\"\\x00\\x00\\x00\\x00\"+\"\\x00\\x00\\x00\\x00\"+\"\\x00\\x00\\x00\\x00\"+\"\\x00\\x00\\x00\\x00\"+\"\\x00\\x00\\x00\\x00\"+\"\\x00\\x00\\x00\\x00\""},
}

func TestGenString(t *testing.T) {
	for _, rec := range bytesStrPairs {
		found := BytesToString(rec.bytes)
		test.StringMatch(t, found, rec.hex)
	}
}

func TestGenStringSep4(t *testing.T) {
	for _, rec := range bytesStrSep4 {
		found := BytesToStringSep(rec.bytes, 4)
		test.StringMatch(t, found, rec.hex)
	}
}
