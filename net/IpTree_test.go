package net

import (
	"net"
	"testing"
)

var cidrMergeTests = []struct {
	a, b, expected string
}{
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

func TestCidrMerge(t *testing.T) {
	for _, rec := range cidrMergeTests {
		tree := New(pickFirst)
		_, c1, _ := net.ParseCIDR(rec.a)
		_, c2, _ := net.ParseCIDR(rec.b)
		_, exp, _ := net.ParseCIDR(rec.expected)
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

var rangeMergeTests = []struct {
	ip1Start, ip1End, ip2Start, ip2End, expectedCidr string
}{
	//Subset
	{"192.168.0.0", "192.168.0.255", "192.168.0.10", "192.168.0.25", "192.168.0.0/24"},
	//Superset
	{"192.168.0.10", "192.168.0.25", "192.168.0.0", "192.168.0.255", "192.168.0.0/24"},
	//Sequential
	{"192.168.0.0", "192.168.0.255", "192.168.1.0", "192.168.1.255", "192.168.0.0/23"},
	{"192.168.0.0", "192.168.0.127", "192.168.0.128", "192.168.0.255", "192.168.0.0/24"},
	{"192.168.0.0", "192.168.0.0", "192.168.0.1", "192.168.0.1", "192.168.0.0/31"},
	{"0.0.0.0", "0.0.0.0", "0.0.0.1", "0.0.0.3", "0.0.0.0/30"},
	//Overlap:
	{"192.168.0.0", "192.168.0.100", "192.168.0.64", "192.168.1.255", "192.168.0.0/23"},
}

func TestRangeMerge(t *testing.T) {
	for _, rec := range rangeMergeTests {
		tree := New(pickFirst)
		start1 := net.ParseIP(rec.ip1Start)
		end1 := net.ParseIP(rec.ip1End)
		start2 := net.ParseIP(rec.ip2Start)
		end2 := net.ParseIP(rec.ip2End)
		_, exp, _ := net.ParseCIDR(rec.expectedCidr)

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

// We need some merge algo, but it's not important to the testing, so
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
