// Copyright 2022 gnabgib
// This Source Code Form is subject to the terms of the Mozilla Public License v2.0

// This packages maps gnablib terms to go (to speed up dev) and adheres to
// https://datatracker.ietf.org/doc/html/rfc4648#section-8 which says HEX is
// case-insensitive, but uses upper case letters in all examples
// no need extra features
package hex

import (
	"errors"
	"fmt"
)

const hexTable = "0123456789ABCDEF"

// Length error (1 extra character.. is it the high or low nibble?)
var ErrLength = errors.New("hex string must be composed of pairs of hex values")

// Invalid hex error type
type invalidHexAtError struct {
	Byte byte
	At   int
}

func (e invalidHexAtError) Error() string {
	return fmt.Sprintf("Invalid hex: %s @ %d", string(e.Byte), e.At)
}

// An invalid character(b) found at position(at) in a hex string
func InvalidHexAt(b byte, at int) invalidHexAtError {
	return invalidHexAtError{Byte: b, At: at}
}

// Convert a byte array into a hex string, uses uppercase chars
func FromBytes(src []byte) string {
	n := len(src)
	ret := make([]byte, n*2)
	ptr := 0
	for i := 0; i < n; i++ {
		//High nibble
		ret[ptr] = hexTable[src[i]>>4]
		//Low nibble
		ret[ptr+1] = hexTable[src[i]&0xF]
		ptr += 2
	}
	return string(ret)
}

//func decode(b byte) byte - Defined in _256b and _tiny files

// Covert a case-insensitive hex string into a byte slice, if any of the hex is malformed
// (invalid character) or there's an odd number of chars, then an error is returned (and no slice)
func ToBytes(hex string) ([]byte, error) {
	n := len(hex)
	//Since there are always 2 hex chars per byte, let's allocate the return with that expectation
	ret := make([]byte, n/2)
	i := 1
	ptr := 0

	for ; i < n; i += 2 {
		high := decode(hex[i-1])
		low := decode(hex[i])
		//Making sure the upper bits are zero is faster than >15
		if high&0xF0 != 0 {
			return nil, InvalidHexAt(hex[i-1], i-1)
		}
		if low&0xF0 != 0 {
			return nil, InvalidHexAt(hex[i], i)
		}
		ret[ptr] = high<<4 | low
		ptr++
	}
	//If there's an extra (valid) character, the length is invalid (is it the low or high nibble?)
	if n&1 == 1 {
		//Consistency: report content errors before length errors
		if decode(hex[i-1]) > 15 {
			return nil, InvalidHexAt(hex[i-1], i-1)
		}
		return nil, ErrLength
	}

	return ret, nil
}

// Convert a case-insensitive hex string into a byte slice, if any of the hex is malformed
// (invalid character) or there's an odd number of chars, then nil is returned (which may also be a valid answer if len(hex)=0)
// 3% faster than ToBytes (x64 tests)
func ToBytesFast(hex string) []byte {
	n := len(hex)
	//Since there are always 2 hex chars per byte, let's allocate the return with that expectation
	ret := make([]byte, n/2)
	i := 1
	ptr := 0

	for ; i < n; i += 2 {
		high := decode(hex[i-1])
		low := decode(hex[i])
		//Making sure the upper bits are zero is faster than >15
		if (high|low)&0xF0 != 0 {
			return nil
		}
		ret[ptr] = high<<4 | low
		ptr++
	}
	//If there's an extra (valid) character, the length is invalid (is it the low or high nibble?)
	if n&1 == 1 {
		return nil
	}
	return ret
}