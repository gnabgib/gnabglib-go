// Copyright 2022 gnabgib
// This Source Code Form is subject to the terms of the Mozilla Public License v2.0

package net

import (
	"bytes"
	"net"
)

// Whether two Cidr are equal (IP and masks match)
func CidrEqual(a, b *net.IPNet) bool {
	m := bytes.Compare(a.Mask, b.Mask)
	if m != 0 {
		//Masks don't match
		return false
	}
	//Make sure IP compare is only the v4 versions
	i := bytes.Compare(a.IP.To4(), b.IP.To4())
	return i == 0
}

// First Ipv4 address in a CIDR (ignores network address)
func FirstIpv4(n *net.IPNet) net.IP {
	// ones, _ := n.Mask.Size()
	// if ones >= 31 {
	// 	return n.IP
	// }
	// otherwise: next IP
	return n.IP
}

// Last Ipv4 address in a CIDR (ignores broadcast)
func LastIpv4(n *net.IPNet) net.IP {
	ret := make([]byte, 4)
	ret[0]=n.IP[0] | ^n.Mask[0]
	ret[1]=n.IP[1] | ^n.Mask[1]
	ret[2]=n.IP[2] | ^n.Mask[2]
	ret[3]=n.IP[3] | ^n.Mask[3]
	return ret
}

// This is provided by go IPNet.Contains(ip IP) bool
// // Whether the given IP is within the CIDR
// func ContainsIpv4(cidr *net.IPNet,ip net.IP) bool {
// 	start := FirstIpv4(cidr)
// 	//A v6 address cannot be contained
// 	i4:=ip.To4()
// 	if i4==nil {
// 		return false
// 	}

// 	masked := make([]byte, 4)
// 	masked[0]=i4[0] & cidr.Mask[0]
// 	masked[1]=i4[1] & cidr.Mask[1]
// 	masked[2]=i4[2] & cidr.Mask[2]
// 	masked[3]=i4[3] & cidr.Mask[3]

// 	return bytes.Equal(start,masked)
// }