// Copyright 2022 gnabgib
// This Source Code Form is subject to the terms of the Mozilla Public License v2.0

package net

import (
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