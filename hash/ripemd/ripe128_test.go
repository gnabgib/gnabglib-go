package ripemd

import (
	"strings"
	"testing"

	"github.com/gnabgib/gnablib-go/test"
)

type ripePair struct {
	in  string
	hex string
}

var ripe128pairs = []ripePair{
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

func TestRipe128(t *testing.T) {
	for _, rec := range ripe128pairs {
		d := New128()
		test.HashTest(t, d, []byte(rec.in), rec.hex)
	}
}

func TestDoubleWrite128Sum(t *testing.T) {
	d := New128()
	test.HashTest(t, d, []byte("a"), "86BE7AFA339D0FC7CFC785E72F578D33")
	test.HashTest(t, d, []byte("bc"), "C14A12199C66E4BA84636B0F69144C77")
}
