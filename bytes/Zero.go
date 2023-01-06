// Copyright 2023 gnabgib
// This Source Code Form is subject to the terms of the Mozilla Public License v2.0

package bytes

const z8="\x00\x00\x00\x00\x00\x00\x00\x00"

// Zero the content of a slice
func Zero(target []byte) {
	//Todo optimize with arch specific commands to speed up (AVX2)
	for len(target)>8 {
		copy(target,z8)
		target=target[8:]
	}
	//One last copy to get the 0-7 other elements
	copy(target,z8)
}