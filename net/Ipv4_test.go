package net

import (
	"net"
	"testing"
)

var ipUintTests = []struct {
	ip  string
	num uint32
}{
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

func TestIpv4FromUint(t *testing.T) {
	for _, rec := range ipUintTests {
		//Note ParseIp may return a 16byte array (ipv6)
		expect := net.ParseIP(rec.ip).To4()
		found := Ipv4FromUint(rec.num)
		if !net.IP.Equal(expect,found) {
			t.Errorf("Expecting %v, got %v", expect, found)
		}
	}
}

func TestIpv4ToUint(t *testing.T) {
	for _, rec := range ipUintTests {
		expect := rec.num
		ip := net.ParseIP(rec.ip)
		found := Ipv4ToUint(ip)

		if expect != found {
			t.Errorf("Expecting %v, got %v", expect, found)
		}
	}
}


func TestIpv4ToUint_withv6(t *testing.T) {
	//When IP is v6, a zero is returned.. while 0 is in some ways a
	// valid IP we can't do anything else without mutating the return
	// (using uint64, using uin32,error etc)
	ip := net.ParseIP("2001:db8::68")
	found := Ipv4ToUint(ip)
	if found != 0 {
		t.Errorf("Expecting 0, got %v", found)
	}
}
