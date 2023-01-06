// Copyright 2023 gnabgib
// This Source Code Form is subject to the terms of the Mozilla Public License v2.0

package bytes

import (
	"encoding/binary"

	"github.com/gnabgib/gnablib-go/endian"
)

const (
	bytes16  = 2
	bytes32  = 4
	bytes64 = 8
	bytes128 = 16
	bytes256 = 32
	bytes512 = 64
)

// Add the first 2 bytes of two byte arrays as if they're 16bit integers
// in platform endian format
func Add16(sum, a, b []byte) (carry byte, err error) {
	if endian.Platform == binary.LittleEndian {
		err = padLE(sum, a, bytes16)
		carry = addEq16LE(sum, b, 0)
	} else {
		err = padBE(sum, a, bytes16)
		carry = addEq16BE(sum, b, 0)
	}
	return
}

// Add the first 4 bytes of two byte arrays as if they're 32bit integers
// in platform endian format
func Add32(sum, a, b []byte) (carry byte, err error) {
	if endian.Platform == binary.LittleEndian {
		err = padLE(sum, a, bytes32)
		carry = addEq32LE(sum, b, 0)
	} else {
		err = padBE(sum, a, bytes32)
		carry = addEq32BE(sum, b, 0)
	}
	return
}

// Add the first 8 bytes of two byte arrays as if they're 64bit integers
// in platform endian format
func Add64(sum, a, b []byte) (carry byte, err error) {
	
	if endian.Platform == binary.LittleEndian {
		carry,err = Add64LE(sum,a,b)
	} else {
		carry,err = Add64BE(sum,a,b)
	}
	return
}

// Add the first 16 bytes of two byte arrays as if they're 128bit integers
// in platform endian format
func Add128(sum, a, b []byte) (carry byte, err error) {
	if endian.Platform == binary.LittleEndian {
		carry,err = Add128LE(sum,a,b)
	} else {
		carry,err = Add128BE(sum,a,b)
	}
	return
}

// Add the first 32 bytes of two byte arrays as if they're 256bit integers
// in platform endian format
func Add256(sum, a, b []byte) (carry byte, err error) {
	if endian.Platform == binary.LittleEndian {
		carry,err = Add256LE(sum,a,b)
	} else {
		carry,err = Add256BE(sum,a,b)
	}
	return
}

// Add the first 64 bytes of two byte arrays as if they're 512bit integers
// in platform endian format
func Add512(sum, a, b []byte) (carry byte, err error) {
	if endian.Platform == binary.LittleEndian {
		carry,err = Add512LE(sum,a,b)
	} else {
		carry,err = Add512BE(sum,a,b)
	}
	return
}

// Update sum to add bytes from b in integer format
// - `sum` MUST have at least 1 byte of content
// - `b` can have any number of bytes, 0-1 will be used
// - Only 0/1 will be read from `carryBit`
// - `carry` can only contain 0/1
func addEq8(sum, b []byte, carryBit byte) (carry byte) {
	a16 := uint16(sum[0])
	b16 := uint16(0)
	if len(b) > 0 {
		b16 = uint16(b[0])
	}
	c16 := a16 + b16 + uint16(carryBit&1)
	sum[0] = byte(c16)
	carry = byte(c16 >> 8)
	return
}
