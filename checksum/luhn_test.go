package checksum_test

import (
	"testing"

	"github.com/gnabgib/gnablib-go/checksum/luhn"
)

type uintLuhn struct {
	n uint64
	h uint8
}

var uintLuhns []uintLuhn

func init() {
	uintLuhns=[]uintLuhn {
		{7992739871, 3},
		{79927398710, 4},
		{1234567, 4},
		{12345678, 2},
		{123456789, 7},
		{1234567890, 3},
		{4992739871, 6},
		{123456781234567, 0},
		{411111111111111, 1},
		{123456781234567, 0},
		{987654321,7},
		{1,8},
		{10,9},
		{100,8},
		{1000,9},
		{10000,8},
		{2,6},
		{12,5},
		{212,1},
		{1212,0},
		{21212,6},
		{18,2},
		{182,6},
		{1826,7},
		{18267,5},
	}
}

func TestLuhn(t *testing.T) {
	for i:=0;i<len(uintLuhns);i++ {
		rec := uintLuhns[i]
		found:=luhn.Checksum(rec.n)
		if found != rec.h {
			t.Fatalf("Hashing %v, expecting %v, got %v",rec.n,rec.h,found)	
		}
	}
}