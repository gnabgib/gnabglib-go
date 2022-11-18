package net

import (
	"net"
	"testing"
)

type ipUintPair struct {
	ip  string
	num uint32
}
type ipEqTest struct {
	a, b string
	eq   bool
}

var ipUintPairs []ipUintPair
var ipEqTests []ipEqTest

func init() {
	ipUintPairs = []ipUintPair{
		{"0.0.0.0", 0},
		{"0.1.2.3", 66051},
		{"0.1.2.30", 66078},
		{"1.1.2.0", 16843264},
		{"1.1.1.1", 16843009},
		{"1.2.3.4", 16909060},
		{"4.3.2.1", 67305985},
		{"8.7.6.5", 134678021},
		{"8.8.8.8", 134744072},
		{"100.200.150.250", 1690867450},
		{"127.0.0.1", 2130706433},
		{"127.255.255.255", 2147483647}, //Max (signed) int32
		{"192.168.1.1", 3232235777},
		{"255.3.2.1", 4278387201},
		{"255.255.255.255", 4294967295},
		{"135.101.67.33", 2271560481}, //87 65 43 21, each segment has dif value
	}
	ipEqTests = []ipEqTest{
		{"192.168.1.0", "192.168.1.0", true},
		{"192.168.1.0", "192.168.0.1", false},
		{"192.168.0.1", "192.168.1.0", false},
		{"1.2.3.4", "4.3.2.1", false},
		{"2001:db8::68", "192.168.0.1", false},
		{"192.168.0.1", "2001:db8::68", false},
		{"2001:db8::68", "2001:db8::68", false}, //Note because neither are v4 this is still false
	}
}


func TestIpv4FromUint(t *testing.T) {
	for i := 0; i < len(ipUintPairs); i++ {
		//Note ParseIp may return a 16byte array (ipv6)
		expect := net.ParseIP(ipUintPairs[i].ip).To4()
		found := Ipv4FromUint(ipUintPairs[i].num)
		if !Ipv4Equal(expect, found) {
			t.Fatalf("Expecting %v, got %v", expect, found)
		}
	}
}

func TestIpv4ToUint(t *testing.T) {
	for i := 0; i < len(ipUintPairs); i++ {
		expect := ipUintPairs[i].num
		ip := net.ParseIP(ipUintPairs[i].ip)
		found := Ipv4ToUint(ip)

		if expect != found {
			t.Fatalf("Expecting %v, got %v", expect, found)
		}
	}
}

func TestIpv4ToUint_withv6(t *testing.T) {
	//When IP is v6, a zero is returned.. while 0 is in some ways a 
	// valid IP we can't do anything else without mutating the return 
	// (using uint64, using uin32,error etc)
	ip := net.ParseIP("2001:db8::68")
	found := Ipv4ToUint(ip)
	if found!=0 {
		t.Fatalf("Expecting 0, got %v",found)
	}
}

func TestIpv4Equal(t *testing.T) {
	for i := 0; i < len(ipEqTests); i++ {
		a := net.ParseIP(ipEqTests[i].a)
		b := net.ParseIP(ipEqTests[i].b)

		if Ipv4Equal(a, b) != ipEqTests[i].eq {
			t.Fatalf("Expected %v==%v to be %v", a, b, ipEqTests[i].eq)
		}
	}
}
