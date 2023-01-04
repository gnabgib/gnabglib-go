// Copyright 2022 gnabgib
// This Source Code Form is subject to the terms of the Mozilla Public License v2.0

package codegen

import (
	"github.com/gnabgib/gnablib-go/encoding/hex"
)

// Generate a (likely invalid UTF8) string holding bytes, which can be used as a constant in golang
func BytesToString(in []byte) string {
	h := []byte(hex.FromBytes(in))
	n := len(h)
	//For each pair of characters we need a "\x" prefix.. so 4 characters per
	// hex pair (that is unless we switch to ASCII for 32-126)
	ret := make([]byte, n*2+2)
	ret[0] = '"'
	ptr := 1
	for i, b := range h {
		if i&1 == 0 {
			ret[ptr] = '\\'
			ret[ptr+1] = 'x'
			ptr += 2
		}
		ret[ptr] = b
		ptr += 1
	}
	ret[ptr] = '"'
	return string(ret)
}

// Generate a (likely invalid UTF8) string holding bytes, which can be used as a constant
// in golang.  bytesPerSection should be >1 and indicate how many bytes to include before
// adding a break.  There are 4 characters per byte ("\x" plus two hex).
func BytesToStringSep(in []byte, bytesPerSection int) string {
	if bytesPerSection < 1 {
		panic("bytesPerSection must be >=1")
	}

	h := []byte(hex.FromBytes(in))
	n := len(h)
	//For each pair of characters we need a "\x" prefix.. so 4 characters per
	// hex pair (that is unless we switch to ASCII for 32-126)
	retLen := n * 2
	hexPerSection := bytesPerSection * 2
	//For each break we need 3 characters "\"+\"" (we don't include new lines and whitespace)
	retLen += 3 * ((n / hexPerSection) - 1)
	//And finally we need the starting and ending quote
	retLen += 2

	ret := make([]byte, retLen)
	ret[0] = '"'
	ptr := 1
	for i, b := range h {
		if i&1 == 0 {
			if i%hexPerSection == 0 && ptr > 1 {
				ret[ptr] = '"'
				ret[ptr+1] = '+'
				ret[ptr+2] = '"'
				ptr += 3
			}
			ret[ptr] = '\\'
			ret[ptr+1] = 'x'
			ptr += 2
		}
		ret[ptr] = b
		ptr += 1
	}
	ret[ptr] = '"'
	return string(ret)
}
