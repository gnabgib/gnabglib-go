package test

import (
	"hash"
	"testing"

	"github.com/gnabgib/gnablib-go/encoding/hex"
)

func HashTest(t *testing.T, h hash.Hash, in []byte, expectHex string) {
	h.Write(in)
	found := h.Sum([]byte{})
	foundHex := hex.FromBytes(found)

	if foundHex != expectHex {
		t.Errorf("Hash(%s)\nExpect: %v\nFound : %v", Abbr(string(in)), expectHex, foundHex)
	}
}

func HashHexTest(t *testing.T, h hash.Hash, inHex string, expectHex string) { 
	b:=hex.ToBytesFast(inHex)
	h.Write(b)
	found := h.Sum([]byte{})
	foundHex := hex.FromBytes(found)

	if foundHex != expectHex {
		t.Errorf("Hash(%s)\nExpect: %v\nFound : %v", "x"+Abbr(inHex), expectHex, foundHex)
	}
}