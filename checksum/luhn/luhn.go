// Copyright 2022 gnabgib
// This Source Code Form is subject to the terms of the Mozilla Public License v2.0

package luhn

// https://www.dcode.fr/luhn-algorithm


// Generate a single digit (0-9) checksum for the given number
func Checksum(i uint64) uint8 {
	var ret= uint64(0)
	var mul= uint64(2)
	for i>0 {
		v := i%10
		i=(i-v)/10
		v*=mul
        ret+=(v%10)+(v/10)
        mul=1+mul%2
	}
	return uint8((10-ret%10)%10)
}