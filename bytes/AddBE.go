// Copyright 2023 gnabgib
// This Source Code Form is subject to the terms of the Mozilla Public License v2.0

package bytes

import (
	"encoding/binary"
)

// Add the first 64 bytes of two slices together into sum (which must have space),
// as if they're 512 bit big endian integers.
// - `a` or `b` can have less than 64 bytes of content and will be zero padded (at the end).
// - `carry` will only contain a 0 or 1
func Add512BE(sum, a, b []byte) (carry byte, err error) {
	carry, err = addAnyBE(sum, a, b, bytes512)
	return
}

// Add the first 32 bytes of two slices together into sum (which must have space),
// as if they're 256 bit big endian integers.
// - `a` or `b` can have less than 32 bytes of content and will be zero padded (at the end).
// - `carry` will only contain a 0 or 1
func Add256BE(sum, a, b []byte) (carry byte, err error) {
	carry, err = addAnyBE(sum, a, b, bytes256)
	return
}

// Add the first 16 bytes of two slices together into sum (which must have space),
// as if they're 128 bit big endian integers.
// - `a` or `b` can have less than 16 bytes of content and will be zero padded (at the end).
// - `carry` will only contain a 0 or 1
func Add128BE(sum, a, b []byte) (carry byte, err error) {
	carry, err = addAnyBE(sum, a, b, bytes128)
	return
}

// Add the first 8 bytes of two slices together into sum (which must have space),
// as if they're 64 bit big endian integers.
// - `a` or `b` can have less than 8 bytes of content and will be zero padded (at the end).
// - `carry` will only contain a 0 or 1
func Add64BE(sum, a, b []byte) (carry byte, err error) {
	carry, err = addAnyBE(sum, a, b, bytes64)
	return
}

// Add the first 4 bytes of two slices together into sum (which must have space),
// as if they're 32 bit big endian integers.
// - `a` or `b` can have less than 4 bytes of content and will be zero padded (at the end)
// - `carry` will only contain a 0 or 1
func Add32BE(sum, a, b []byte) (carry byte, err error) {
	err = padBE(sum, a, bytes32)
	if err != nil {
		return
	}
	carry = addEq32BE(sum, b, 0)
	return
}

// Add the first 2 bytes of two slices together into sum (which must have space),
// as if they're 16 bit big endian integers.
// - `a` or `b` can have less than 2  bytes of content and will be zero padded (at the end)
// - `carry` will only contain a 0 or 1
func Add16BE(sum, a, b []byte) (carry byte, err error) {
	err = padBE(sum, a, bytes16)
	if err != nil {
		return
	}
	carry = addEq16BE(sum, b, 0)
	return
}

// Add the first `byteSize` bytes of two byte arrays if they're byte*8
// bit big endian unsigned integers.  `carry` will only contain a 0 or 1
func addAnyBE(sum, a, b []byte, byteSize int) (carry byte, err error) {
	err = padBE(sum, a, byteSize)
	empty := []byte{}
	if err != nil {
		return
	}

	in32 := (byteSize / bytes32)-1
	in8 := byteSize % bytes32

	for ; in32 >= 0 && len(b) > bytes32; in32-- {
		carry = addEq32BE(sum[in32*bytes32:], b[len(b)-bytes32:], carry)
		b = b[:len(b)-bytes32]
	}
	//Add a partial B if there is one
	if len(b)>0 && in32>=0 {
		carry = addEq32BE(sum[in32*bytes32:], b, carry)
		in32--
		b=b[:0]
	}
	//Zero the rest of the adds (we still propagate carry)
	for ; in32 >= 0; in32-- {
		carry = addEq32BE(sum[in32*bytes32:], empty, carry)
	}

	//Add a 16 bit if needed
	if in8 > 1 {
		in8-=2
		carry = addEq16BE(sum[in8:], b, carry)
		//Update b
		if len(b) == 1 {
			b = b[:0]
		} else if len(b) >= bytes16 {
			b = b[:len(b)-bytes16]
		}
	}
	//Add 8 bit if needed (only one)
	if in8 > 0 {
		carry = addEq8(sum, b, carry)
	}
	return
}

// Update sum to add bytes from b in big endian integer format
// - `sum` MUST have at least 4 bytes of content
// - `b` can have any number of bytes, 0-4 will be used
// - Only 0/1 will be read from `carryBit`
// - `carry` can only contain 0/1
func addEq32BE(sum, b []byte, carryBit byte) (carry byte) {
	a64 := uint64(binary.BigEndian.Uint32(sum))
	b64 := uint64(0)
	n := bytes32
	if len(b) < n {
		n = len(b)
	}
	for i := 0; i < n; i++ {
		b64 = b64<<8 | uint64(b[i])
	}

	c64 := a64 + b64 + uint64(carryBit&1)
	binary.BigEndian.PutUint32(sum, uint32(c64))
	carry = byte(c64 >> 32)
	return
}

// Update sum to add bytes from b in big endian integer format
// - `sum` MUST have at least 2 bytes of content
// - `b` can have any number of bytes, 0-2 will be used
// - Only 0/1 will be read from `carryBit`
// - `carry` can only contain 0/1
func addEq16BE(sum, b []byte, carryBit byte) (carry byte) {
	a32 := uint32(binary.BigEndian.Uint16(sum))
	b32 := uint32(0)
	i := 0
	n := bytes16
	if len(b) < n {
		n = len(b)
	}
	for ; i < n; i++ {
		b32 = b32<<8 | uint32(b[i])
	}

	c32 := a32 + b32 + uint32(carryBit&1)
	binary.BigEndian.PutUint16(sum, uint16(c32))
	carry = byte(c32 >> 16)
	return
}
