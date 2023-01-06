package net

import (
	"net"
	"testing"
)

var unequalCidrTests = []struct {
	a, b string
}{
	{"192.168.1.0/24", "192.168.1.0/25"},
	{"192.168.1.0/24", "192.168.1.0/23"},
	{"192.168.1.0/24", "192.168.0.0/24"},
}

func TestNotEqual(t *testing.T) {
	for _, rec := range unequalCidrTests {
		_, c1, _ := net.ParseCIDR(rec.a)
		_, c2, _ := net.ParseCIDR(rec.b)
		if CidrEqual(c1, c2) {
			t.Errorf("Expecting CIDRs to be unequal %v==%v", c1, c2)
		}
	}
}

var cidrFirstLastIpTests = []struct {
	cidr, first, last string
}{
	{"192.168.0.0/24", "192.168.0.0", "192.168.0.255"},
	{"192.168.0.0/23", "192.168.0.0", "192.168.1.255"},
	{"1.2.3.4/30", "1.2.3.4", "1.2.3.7"},
}

func TestFirstIpv4(t *testing.T) {
	for _, rec := range cidrFirstLastIpTests {
		_, c, _ := net.ParseCIDR(rec.cidr)
		expect := net.ParseIP(rec.first)
		found := FirstIpv4(c)

		if !net.IP.Equal(expect,found) {
			t.Errorf("Expecting first ip of %v to be %v, found %v", rec.cidr, expect, found)
		}
	}
}

func TestLastIpv4(t *testing.T) {
	for _, rec := range cidrFirstLastIpTests {
		_, c, _ := net.ParseCIDR(rec.cidr)
		expect := net.ParseIP(rec.last)
		found := LastIpv4(c)

		if !net.IP.Equal(expect,found) {
			t.Errorf("Expecting last ip of %v to be %v, found %v", c, expect, found)
		}
	}
}

var cidrContainsTests = []struct {
	cidr, ip string
	expect bool
}{
	{"192.168.0.0/24", "192.168.0.0", true},
	{"192.168.0.0/24", "192.168.0.255", true},
	{"192.168.0.0/24", "192.168.0.128", true},
	{"192.168.0.0/24", "192.168.1.0", false},
	{"192.168.0.0/24", "0.0.0.0", false},
	{"192.168.0.0/24", "2001:db8::68", false},
	//IPv4 mapped Ipv6
	{"192.168.0.0/24", "::FFFF:192.168.0.10", true},
}

func TestCidrContains(t *testing.T) {
	for _, rec := range cidrContainsTests {
		_, c, _ := net.ParseCIDR(rec.cidr)
		ip := net.ParseIP(rec.ip)
		found := c.Contains(ip)

		if rec.expect != found {
			t.Errorf("Expecting %v contains %v to be %v", c, ip, rec.expect)
		}
	}
}

func TestEqual(t *testing.T) {
	_, c1, _ := net.ParseCIDR("192.168.1.0/24")
	_, c2, _ := net.ParseCIDR("192.168.1.0/24")

	if !CidrEqual(c1, c2) {
		t.Errorf("Expecting CIDRs to be found equal")
	}
}
