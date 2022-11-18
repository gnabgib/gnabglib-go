// Copyright 2022 gnabgib
// This Source Code Form is subject to the terms of the Mozilla Public License v2.0

package net

import (
	"encoding/binary"
	"net"
)

// Convert a uint into an IPMask (from CIDR notation to go-internal)
func MaskFromByte(b byte) net.IPMask {
	var mask uint32
	mask = 0xffffffff << (32 - b)
	ret := make([]byte, 4)
	binary.BigEndian.PutUint32(ret, mask)
	return ret
}

// Convert an IPMask into a byte, much like in CIDR notation (0-32)
func MaskToByte(m net.IPMask) byte {
	ret, _ := m.Size()
	return byte(ret)
}