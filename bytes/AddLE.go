// Copyright 2023 gnabgib
// This Source Code Form is subject to the terms of the Mozilla Public License v2.0

package bytes

import (
	"encoding/binary"
)

// Add the first 64 bytes of two slices together into sum (which must have space),
// as if they're 512 bit little endian integers.
// - `a` or `b` can have less than 64 bytes of content and will be zero padded (at the end).
// - `carry` will only contain a 0 or 1
func Add512LE(sum, a, b []byte) (carry byte, err error) {
	carry, err = addAnyLE(sum, a, b, bytes512)
	return
}

// Add the first 32 bytes of two slices together into sum (which must have space),
// as if they're 256 bit little endian integers.
// - `a` or `b` can have less than 32 bytes of content and will be zero padded (at the end).
// - `carry` will only contain a 0 or 1
func Add256LE(sum, a, b []byte) (carry byte, err error) {
	carry, err = addAnyLE(sum, a, b, bytes256)
	return
}

// Add the first 16 bytes of two slices together into sum (which must have space),
// as if they're 128 bit little endian integers.
// - `a` or `b` can have less than 16 bytes of content and will be zero padded (at the end).
// - `carry` will only contain a 0 or 1
func Add128LE(sum, a, b []byte) (carry byte, err error) {
	carry, err = addAnyLE(sum, a, b, bytes128)
	return
}

// Add the first 8 bytes of two slices together into sum (which must have space),
// as if they're 64 bit little endian integers.
// - `a` or `b` can have less than 8 bytes of content and will be zero padded (at the end).
// - `carry` will only contain a 0 or 1
func Add64LE(sum, a, b []byte) (carry byte, err error) {
	carry, err = addAnyLE(sum, a, b, bytes64)
	return
}

// Add the first 4 bytes of two slices together into sum (which must have space),
// as if they're 32 bit little endian integers.
// - `a` or `b` can have less than 4 bytes of content and will be zero padded (at the end)
// - `carry` will only contain a 0 or 1
func Add32LE(sum, a, b []byte) (carry byte, err error) {
	err = padLE(sum, a, bytes32)
	if err != nil {
		return
	}
	carry = addEq32LE(sum, b, 0)
	return
}

// Add the first 2 bytes of two slices together into sum (which must have space),
// as if they're 16 bit little endian integers.
// - `a` or `b` can have less than 2  bytes of content and will be zero padded (at the end)
// - `carry` will only contain a 0 or 1
func Add16LE(sum, a, b []byte) (carry byte, err error) {
	err = padLE(sum, a, bytes16)
	if err != nil {
		return
	}
	carry = addEq16LE(sum, b, 0)
	return
}

// Add the first `byteSize` bytes of two byte arrays if they're byte*8
// bit little endian unsigned integers.  `carry` will only contain a 0 or 1
func addAnyLE(sum, a, b []byte, byteSize int) (carry byte, err error) {
	err = padLE(sum, a, byteSize)
	empty := []byte{}
	if err != nil {
		return
	}

	in32 := byteSize / bytes32
	in8 := byteSize % bytes32

	numU32 := 0
	numU16 := 0
	for ; numU32 < in32 && len(b) > bytes32; numU32++ {
		carry = addEq32LE(sum[numU32*bytes32:], b, carry)
		b=b[bytes32:]
	}
	//Add a partial B if there is one
	if len(b)>0 && numU32<in32 {
		carry = addEq32LE(sum[numU32*bytes32:], b, carry)
		numU32++
		b=b[len(b):]
	}
	//Zero the rest of the adds (we still propagate carry)
	for ; numU32 < in32; numU32++ {
		carry = addEq32LE(sum[numU32*bytes32:], empty, carry)
	}

	//Add a 16 bit if needed
	if in8 > 1 {
		carry = addEq16LE(sum[numU32*bytes32:], b, carry)
		//Update b
		if len(b) == 1 {
			b = b[1:]
		} else if len(b) >= bytes16 {
			b = b[bytes16:]
		}
		numU16++
		in8 -= 2
	}
	//Add 8 bit if needed (only one)
	if in8 > 0 {
		carry = addEq8(sum[numU32*bytes32+numU16*bytes16:], b, carry)
	}
	return
}

// Update sum to add bytes from b in little endian integer format
// - `sum` MUST have at least 4 bytes of content
// - `b` can have any number of bytes, 0-4 will be used
// - Only 0/1 will be read from `carryBit`
// - `carry` can only contain 0/1
func addEq32LE(sum, b []byte, carryBit byte) (carry byte) {
	a64 := uint64(binary.LittleEndian.Uint32(sum))
	b64 := uint64(0)
	n := bytes32
	if len(b) < n {
		n = len(b)
	}
	for i := 0; i < n; i++ {
		b64 += uint64(b[i]) << (8 * i)
	}

	c64 := a64 + b64 + uint64(carryBit&1)
	binary.LittleEndian.PutUint32(sum, uint32(c64))
	carry = byte(c64 >> 32)
	return
}

// Update sum to add bytes from b in little endian integer format
// - `sum` MUST have at least 2 bytes of content
// - `b` can have any number of bytes, 0-2 will be used
// - Only 0/1 will be read from `carryBit`
// - `carry` can only contain 0/1
func addEq16LE(sum, b []byte, carryBit byte) (carry byte) {
	a32 := uint32(binary.LittleEndian.Uint16(sum))
	b32 := uint32(0)
	i := 0
	n := bytes16
	if len(b) < n {
		n = len(b)
	}
	for ; i < n; i++ {
		b32 += uint32(b[i]) << (8 * i)
	}

	c32 := a32 + b32 + uint32(carryBit&1)
	binary.LittleEndian.PutUint16(sum, uint16(c32))
	carry = byte(c32 >> 16)
	return
}
