package hash_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/gnabgib/gnablib-go/hash/ripemd"
)

type ripe128pair struct {
	s   string
	hex string
}

var ripe128pairs []ripe128pair

func init() {
	ripe128pairs = []ripe128pair{
		//Source: https://homes.esat.kuleuven.be/~bosselae/ripemd160.html
		{"", "CDF26213A150DC3ECB610F18F6B38B46"},
		{"a", "86BE7AFA339D0FC7CFC785E72F578D33"},
		{"abc", "C14A12199C66E4BA84636B0F69144C77"},
		{"message digest", "9E327B3D6E523062AFC1132D7DF9D1B8"},
		{"abcdefghijklmnopqrstuvwxyz", "FD2AA607F71DC8F510714922B371834E"},
		{"abcdbcdecdefdefgefghfghighijhijkijkljklmklmnlmnomnopnopq", "A1AA0689D0FAFA2DDC22E88B49133A06"},
		{"ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789",
			"D1E959EB179C911FAEA4624C60C5C702"},
		{strings.Repeat("1234567890", 8), "3F45EF194732C2DBB2C4A2C769795FA3"},
		{"ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789",
			"D1E959EB179C911FAEA4624C60C5C702"},

		//Other
		{"The quick brown fox jumps over the lazy dog", "3FA9B57F053C053FBE2735B2380DB596"},
		{"The quick brown fox jumps over the lazy cog", "3807AAAEC58FE336733FA55ED13259D9"},
		{"The quick brown fox jumps over the lazy dog.", "F39288A385E5996C6FB2C01F7F8FBF2F"},
		{"gnabgib", "1B6B23F7CFBA4BBF4209757466C1561B"},
		{strings.Repeat("a", 1000000), "4A7F5723F954EBA1216C9D8F6320431F"},
	}
}

func abbr(s string) string {
	n := len(s)
	if n > 10 {
		return s[0:7] + "..."
	}
	return s
}

func TestRipe128(t *testing.T) {
	for i := 0; i < len(ripe128pairs); i++ {
		rec := ripe128pairs[i]
		d := ripemd.New128()
		d.Write([]byte(rec.s))
		found := d.Sum([]byte{})
		foundHex := strings.ToUpper(fmt.Sprintf("%x", found))

		if foundHex != rec.hex {
			t.Fatalf("Ripe128 %v, expecting %v, got %v", abbr(rec.s), rec.hex, foundHex)
		}
	}
}

func TestDoubleWrite128Sum(t *testing.T) {
	d := ripemd.New128()
	d.Write([]byte("a"))
	sum1:=d.Sum([]byte{})
	sum1hex:=strings.ToUpper(fmt.Sprintf("%x",sum1))
	d.Write([]byte("bc"))
	sum2:=d.Sum([]byte{})
	sum2hex:=strings.ToUpper(fmt.Sprintf("%x",sum2))

	if sum1hex!="86BE7AFA339D0FC7CFC785E72F578D33" {
		t.Fatalf("Ripe128 first hash mismatch, got %v", sum1hex)
	}
	if sum2hex!="C14A12199C66E4BA84636B0F69144C77" {
		t.Fatalf("Ripe128 second hash mismatch, got %v", sum1hex)
	}

}