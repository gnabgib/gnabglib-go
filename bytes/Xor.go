// Copyright 2022 gnabgib
// This Source Code Form is subject to the terms of the Mozilla Public License v2.0

package bytes

// XORs elements from `other` into `dst`.  Returns the number of elements
// XORed, which will be the minimum of len(dst) len(other)
// dst ^= other
func XorEq(dst,other []byte) (n int){
	n=len(dst)
	//If other is smaller, shrink xor count
	if len(other)<n {
		n=len(other)
	}
	for i:=0;i<n;i++ {
		dst[i]^=other[i]
	}
	return
}