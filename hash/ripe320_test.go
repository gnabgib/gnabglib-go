package hash_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/gnabgib/gnablib-go/hash/ripemd"
)

type ripe320pair struct {
	s   string
	hex string
}

var ripe320pairs []ripe320pair

func init() {
	ripe320pairs = []ripe320pair{
		//Source: https://homes.esat.kuleuven.be/~bosselae/ripemd160.html
		{"", "22D65D5661536CDC75C1FDF5C6DE7B41B9F27325EBC61E8557177D705A0EC880151C3A32A00899B8"},
		{"a", "CE78850638F92658A5A585097579926DDA667A5716562CFCF6FBE77F63542F99B04705D6970DFF5D"},
		{"abc", "DE4C01B3054F8930A79D09AE738E92301E5A17085BEFFDC1B8D116713E74F82FA942D64CDBC4682D"},
		{"message digest", "3A8E28502ED45D422F68844F9DD316E7B98533FA3F2A91D29F84D425C88D6B4EFF727DF66A7C0197"},
		{"abcdefghijklmnopqrstuvwxyz", "CABDB1810B92470A2093AA6BCE05952C28348CF43FF60841975166BB40ED234004B8824463E6B009"},
		{"abcdbcdecdefdefgefghfghighijhijkijkljklmklmnlmnomnopnopq",
			"D034A7950CF722021BA4B84DF769A5DE2060E259DF4C9BB4A4268C0E935BBC7470A969C9D072A1AC"},
		{"ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789",
			"ED544940C86D67F250D232C30B7B3E5770E0C60C8CB9A4CAFE3B11388AF9920E1B99230B843C86A4"},
		{strings.Repeat("1234567890", 8), "557888AF5F6D8ED62AB66945C6D2A0A47ECD5341E915EB8FEA1D0524955F825DC717E4A008AB2D42"},
		{strings.Repeat("a", 1000000), "BDEE37F4371E20646B8B0D862DDA16292AE36F40965E8C8509E63D1DBDDECC503E2B63EB9245BB66"},
		
		//Other using https://md5calc.com/hash/ripemd320
		{"The quick brown fox jumps over the lazy dog",
			"E7660E67549435C62141E51C9AB1DCC3B1EE9F65C0B3E561AE8F58C5DBA3D21997781CD1CC6FBC34"},
		{"The quick brown fox jumps over the lazy cog",
			"393E0DF728C4CE3D79E7DCFD357D5C26F5C6D64C6D652DC53B6547B214EA9183E4F61C477EBF5CB0"},
		{"The quick brown fox jumps over the lazy dog.",
			"4B743C0A2262F904097FDA33F0B8E03819BC012C5FC39643F17049C566EEE0B1961D1BD7B25A3E4D"},
		{"gnabgib", "EFB9A476DAC2762191AA9890089ECCE637DAA14C22A2F8B6D024DF1AAB2DDDFAE159071D96196EFF"},
	}
}

func TestRipe320(t *testing.T) {
	for i := 0; i < len(ripe320pairs); i++ {
		rec := ripe320pairs[i]
		d := ripemd.New320()
		d.Write([]byte(rec.s))
		found := d.Sum([]byte{})
		foundHex := strings.ToUpper(fmt.Sprintf("%x", found))

		if foundHex != rec.hex {
			t.Fatalf("Ripe320 %v, expecting %v, got %v", abbr(rec.s), rec.hex, foundHex)
		}
	}
}
