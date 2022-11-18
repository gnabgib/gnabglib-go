package net

import (
	"net"
	"testing"
)

type cidrPair struct {
	cidr1,cidr2 string
}

type cidrFirstLast struct {
	cidr,first,last string
}

type cidrContain struct {
	cidr,ip string
	found bool
}

var unEqualCidrPairs []cidrPair
var ranges []cidrFirstLast
var cidrContains []cidrContain

func init() {
	unEqualCidrPairs=[]cidrPair {
		{"192.168.1.0/24","192.168.1.0/25"},
		{"192.168.1.0/24","192.168.1.0/23"},
		{"192.168.1.0/24","192.168.0.0/24"},
	}
	ranges=[]cidrFirstLast {
		{"192.168.0.0/24","192.168.0.0","192.168.0.255"},
		{"192.168.0.0/23","192.168.0.0","192.168.1.255"},
		{"1.2.3.4/30","1.2.3.4","1.2.3.7"},
	}
	cidrContains=[]cidrContain {
		{"192.168.0.0/24","192.168.0.0",true},
		{"192.168.0.0/24","192.168.0.255",true},
		{"192.168.0.0/24","192.168.0.128",true},
		{"192.168.0.0/24","192.168.1.0",false},
		{"192.168.0.0/24","0.0.0.0",false},
		{"192.168.0.0/24","2001:db8::68",false},
		//IPv4 mapped Ipv6
		{"192.168.0.0/24","::FFFF:192.168.0.10",true},
	}
}

func TestEqual(t *testing.T) {
	_,c1,_:=net.ParseCIDR("192.168.1.0/24")
	_,c2,_:=net.ParseCIDR("192.168.1.0/24")

	if !CidrEqual(c1,c2) {
		t.Fatalf("Expecting CIDRs to be found equal")
	}
}

func TestNotEqual(t *testing.T) {
	for i:=0;i<len(unEqualCidrPairs);i++ {
		_,c1,_:=net.ParseCIDR(unEqualCidrPairs[i].cidr1)
		_,c2,_:=net.ParseCIDR(unEqualCidrPairs[i].cidr2)
		if CidrEqual(c1,c2) {
			t.Fatalf("Expecting CIDRs to be unequal %v==%v",c1,c2)
		}
	}
}

func TestFirstIpv4(t *testing.T) {
	for i:=0;i<len(ranges);i++ {
		_,c,_ := net.ParseCIDR(ranges[i].cidr)
		expected := net.ParseIP(ranges[i].first)
		found := FirstIpv4(c)

		if !Ipv4Equal(expected,found) {
			t.Fatalf("Expecting first ip of %v to be %v, found %v",ranges[i].cidr,expected,found)
		}
	}
}

func TestLastIpv4(t *testing.T) {
	for i:=0;i<len(ranges);i++ {
		_,c,_ := net.ParseCIDR(ranges[i].cidr)
		expected := net.ParseIP(ranges[i].last)
		found := LastIpv4(c)

		if !Ipv4Equal(expected,found) {
			t.Fatalf("Expecting last ip of %v to be %v, found %v",c,expected,found)
		}
	}
}

func TestCidrContains(t * testing.T) {
	for i:=0;i<len(cidrContains);i++ {
		_,c,_:=net.ParseCIDR(cidrContains[i].cidr)
		ip := net.ParseIP(cidrContains[i].ip)
		expected := cidrContains[i].found
		found := c.Contains(ip)
		//found := ContainsIpv4(c,ip)

		if expected!=found {
			t.Fatalf("Expecting %v contains %v to be %v",c,ip,expected)
		}
	}
}
//func TestFirstIpv4