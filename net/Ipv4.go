// Copyright 2022 gnabgib
// This Source Code Form is subject to the terms of the Mozilla Public License v2.0

package net

import (
	"bytes"
	"encoding/binary"
	"net"
)

// Turn an IP address into a 32bit uint
func Ipv4ToUint(ip net.IP) uint32 {
	i4 := ip.To4()
	if i4==nil {
		return 0
	}
	return binary.BigEndian.Uint32(i4)
}

// Get an IP address from a 32bit uint
func Ipv4FromUint(u uint32) net.IP {
	ret := make([]byte, 4)
	binary.BigEndian.PutUint32(ret, u)
	return ret
}

// Whether two Ipv4 addresses are equal
func Ipv4Equal(a, b net.IP) bool {
	//Note we cast down to v4 because a net.IP can also be v6
	// (will be nil in v4 form)
	a4 := a.To4()
	b4 := b.To4()
	//We're only looking for v4 equality so if either (or both)
	// are not v4, they're not equal (note this means two equal v6
	// addresses will be reported as not-equal)
	if a4 == nil || b4 == nil {
		return false
	}
	return bytes.Equal(a4,b4)
}
