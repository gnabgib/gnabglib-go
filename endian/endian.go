// Copyright 2022 gnabgib
// This Source Code Form is subject to the terms of the Mozilla Public License v2.0

package endian

import (
	"encoding/binary"
	"unsafe"
)

// The native endianness of execution platform
var Platform binary.ByteOrder

// The endianness of source code (always big)
var SourceCode binary.ByteOrder

// The endianness of networks (big by convention)
var Network binary.ByteOrder

func init() {
	//We always write source code in big-endian
	SourceCode = binary.BigEndian
	//Network order is also always big-endian
	Network = binary.BigEndian

	//While the majority of platforms are little-endian, let's detect with a pointer cast
	b := [2]byte{}
	*(*uint16)(unsafe.Pointer(&b[0])) = uint16(0xABCD)
	switch b[0] {
	case 0xCD:
		Platform = binary.LittleEndian
	case 0xAB:
		Platform = binary.BigEndian
	default:
		//We could default to LE, but this should also never fail and
		// we want bug reports if it does
		panic("Endian detection problem")
	}
}
