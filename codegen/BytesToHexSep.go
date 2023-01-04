// Copyright 2022 gnabgib
// This Source Code Form is subject to the terms of the Mozilla Public License v2.0

package codegen

import (
	"github.com/gnabgib/gnablib-go/encoding/hex"
)

// Format for use as an array or slice, using hex representation of the numbers
// bytesPerSection should be >1 and indicate how many bytes to include before
// adding a break.  There are 5 characters per byte ("0x" + two hex + ",").
func BytesToHexSep(in []byte, bytesPerSection int) string {
	if bytesPerSection < 1 {
		panic("bytesPerSection must be >=1")
	}

	h := []byte(hex.FromBytes(in))
	n := len(h)
	//For each byte we need a "0x" prefix and ", " suffix and two hex chars.. so 6 chars
	retLen := len(in) * 6
	hexPerSection := bytesPerSection * 2
	//For each break we need 1 char "\n"
	retLen += ((n / hexPerSection) - 1)
	//And we need open/close bracket +2, but we won't have a trailing ", "
	//retLen+=2-2

	ret := make([]byte, retLen)
	ret[0] = '{'
	ptr := 1
	for i, b := range h {
		if i&1 == 0 {
			if ptr > 1 {
				ret[ptr] = ','
				ret[ptr+1] = ' '
				ptr += 2
				if i%hexPerSection == 0 {
					ret[ptr] = '\n'
					ptr += 1
				}
			}
			ret[ptr] = '0'
			ret[ptr+1] = 'x'
			ptr += 2
		}
		ret[ptr] = b
		ptr += 1
	}
	ret[ptr] = '}'
	return string(ret)
}
