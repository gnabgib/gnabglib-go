package net

import (
	"net"
	"testing"
)

type maskNumPair struct {
	mask []byte
	num  uint32
}
type cidrMergeSet struct {
	cidr1, cidr2, expected string
}
type rangeMergeSet struct {
	ip1Start, ip1End, ip2Start, ip2End, expectedCidr string
}

var maskNumPairs []maskNumPair
var cidrMergeSets []cidrMergeSet
var rangeMergeSets []rangeMergeSet

func init() {
	maskNumPairs = []maskNumPair{
		{[]byte{255, 255, 255, 255}, 32},
		{[]byte{255, 255, 255, 0b11111110}, 31},
		{[]byte{255, 255, 0b11111111, 0}, 24},
		{[]byte{255, 255, 0b11111110, 0}, 23},
		{[]byte{255, 255, 0b11111100, 0}, 22},
		{[]byte{255, 255, 0b11111000, 0}, 21},
		{[]byte{255, 255, 0b11110000, 0}, 20},
		{[]byte{255, 255, 0b11100000, 0}, 19},
		{[]byte{255, 255, 0b11000000, 0}, 18},
		{[]byte{255, 255, 0b10000000, 0}, 17},
		{[]byte{255, 255, 0b00000000, 0}, 16},
		{[]byte{255, 0, 0, 0}, 8},
		{[]byte{0x80, 0, 0, 0}, 1},
		{[]byte{0, 0, 0, 0}, 0},
	}
	cidrMergeSets = []cidrMergeSet{
		//Subset
		{"192.168.0.0/24", "192.168.0.0/25", "192.168.0.0/24"},
		//Superset
		{"192.168.0.0/25", "192.168.0.0/24", "192.168.0.0/24"},
		//Sequential
		{"192.168.0.0/24", "192.168.1.0/24", "192.168.0.0/23"},
		{"127.1.2.0/24", "127.1.3.0/24", "127.1.2.0/23"},
		{"192.168.0.0/25", "192.168.0.128/25", "192.168.0.0/24"},
		//Overlap isn't possible with Cidr notation
	}
	rangeMergeSets = []rangeMergeSet{
		//Subset
		{"192.168.0.0", "192.168.0.255", "192.168.0.10", "192.168.0.25", "192.168.0.0/24"},
		//Superset
		{"192.168.0.10", "192.168.0.25", "192.168.0.0", "192.168.0.255", "192.168.0.0/24"},
		//Sequential
		{"192.168.0.0", "192.168.0.255", "192.168.1.0", "192.168.1.255", "192.168.0.0/23"},
		{"192.168.0.0", "192.168.0.127", "192.168.0.128", "192.168.0.255", "192.168.0.0/24"},
		{"192.168.0.0","192.168.0.0","192.168.0.1","192.168.0.1","192.168.0.0/31"},
		{"0.0.0.0","0.0.0.0","0.0.0.1","0.0.0.3","0.0.0.0/30"},
		//Overlap:
		{"192.168.0.0", "192.168.0.100", "192.168.0.64", "192.168.1.255", "192.168.0.0/23"},
	}
}

//We need some merge algo, but it's not important to the testing, so
// just chose the first item in the list
func pickFirst(a, b interface{}) interface{} {
	return a
}

func TestEmpty(t *testing.T) {
	tree := New(pickFirst)
	list := tree.ListCidr()
	if len(list) != 0 {
		t.Fatalf("Too many cidr")
	}
}

func TestOneIp(t *testing.T) {
	tree := New(pickFirst)
	ip, expected, _ := net.ParseCIDR("135.101.67.33/32")
	tree.AddIp(ip, "")
	list := tree.ListCidr()
	if len(list) != 1 {
		t.Fatalf("Expecting a list with one CIDR got %v", len(list))
	}
	if !CidrEqual(list[0].Cidr, expected) {
		t.Fatalf("Expecting CIDR %v got %v", expected, list[0])
	}
}

func TestOneCidr(t *testing.T) {
	tree := New(pickFirst)
	_, expected, _ := net.ParseCIDR("1.2.3.4/16") //Note this parses to 1.2.0.0/16
	tree.AddCidr(expected, "")
	//fmt.Println(tree)
	list := tree.ListCidr()
	if len(list) != 1 {
		t.Fatalf("Expecting a list with one CIDR got %v", len(list))
	}
	if !CidrEqual(list[0].Cidr, expected) {
		t.Fatalf("Expecting CIDR %v got %v", expected, list[0])
	}
}

func TestOneRange(t *testing.T) {
	tree := New(pickFirst)
	start := net.ParseIP("192.168.1.0")
	end := net.ParseIP("192.168.1.63")
	_, expected, _ := net.ParseCIDR("192.168.1.0/26")

	tree.AddRange(start, end, "")
	list := tree.ListCidr()
	if len(list) != 1 {
		t.Fatalf("Expecting a list with one CIDR got %v", len(list))
	}
	if !CidrEqual(list[0].Cidr, expected) {
		t.Fatalf("Expecting CIDR %v got %v", expected, list[0])
	}
}

func TestCidrMerge(t *testing.T) {
	for i := 0; i < len(cidrMergeSets); i++ {
		tree := New(pickFirst)
		_, c1, _ := net.ParseCIDR(cidrMergeSets[i].cidr1)
		_, c2, _ := net.ParseCIDR(cidrMergeSets[i].cidr2)
		_, exp, _ := net.ParseCIDR(cidrMergeSets[i].expected)
		tree.AddCidr(c1, "")
		tree.AddCidr(c2, "")
		list := tree.ListCidr()
		if len(list) != 1 {
			t.Fatalf("With %v + %v, expecting a list with one CIDR got %v", c1, c2, len(list))
		}
		if !CidrEqual(list[0].Cidr, exp) {
			t.Fatalf("With %v + %v, expecting CIDR %v got %v", c1, c2, exp, list[0])
		}
	}
}

func TestRangeMerge(t *testing.T) {
	for i := 0; i < len(rangeMergeSets); i++ {
		tree := New(pickFirst)
		start1 := net.ParseIP(rangeMergeSets[i].ip1Start)
		end1 := net.ParseIP(rangeMergeSets[i].ip1End)
		start2 := net.ParseIP(rangeMergeSets[i].ip2Start)
		end2 := net.ParseIP(rangeMergeSets[i].ip2End)
		_, exp, _ := net.ParseCIDR(rangeMergeSets[i].expectedCidr)

		tree.AddRange(start1, end1, "")
		tree.AddRange(start2, end2, "")
		list := tree.ListCidr()
		if len(list) != 1 {
			t.Fatalf("Expecting a list with one CIDR got %v", len(list))
		}
		if !CidrEqual(list[0].Cidr, exp) {
			t.Fatalf("Expecting CIDR %v got %v", exp, list[0])
		}
	}
}

func TestIpRangeCidr(t *testing.T) {
	tree := New(pickFirst)
	ip := net.ParseIP("192.168.1.0")
	rangeStart := net.ParseIP("192.168.1.1")
	rangeEnd := net.ParseIP("192.168.1.129") //Note slight overlap with cidr
	_, cidr, _ := net.ParseCIDR("192.168.1.128/25")
	_, expected, _ := net.ParseCIDR("192.168.1.0/24")

	tree.AddIp(ip, "")
	tree.AddRange(rangeStart, rangeEnd, "")
	tree.AddCidr(cidr, "")

	list := tree.ListCidr()
	if len(list) != 1 {
		t.Fatalf("Expecting a list with one CIDR got %v", len(list))
	}
	if !CidrEqual(list[0].Cidr, expected) {
		t.Fatalf("Expecting CIDR %v got %v", expected, list[0])
	}

}