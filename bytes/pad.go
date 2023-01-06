// Copyright 2023 gnabgib
// This Source Code Form is subject to the terms of the Mozilla Public License v2.0

package bytes

import "errors"

// Pad 0-n bytes input in to n bytes (zero end), if len(in)>n only n bytes will be used.
// Always returns a copy
func padLE(dst,src []byte, n int) error {
	//Todo optimize with arch specific commands to speed up (AVX2)
	if len(dst)<n {
		return errors.New("not enough space")
	}
	sn:=n
	if len(src)<sn {
		sn=len(src)
	}
	//When src/dst are pointing to the same memory
	if sn>0 && &dst[0]==&src[0] {
		return nil
	}
	copy(dst,src[:sn])
	//Zero any remaining space
	if sn<n {
		Zero(dst[sn:n])
	}
	return nil
}

// Pad 0-n bytes input in to n bytes (zero start), if len(in)>n only n bytes will be used
// Always returns a copy
func padBE(dst,src []byte, n int) error {
	if len(dst)<n {
		return errors.New("not enough space")
	}
	sn:=n
	if len(src)<sn {
		sn=len(src)
	}
	//Zero starting space
	if sn<n {
		Zero(dst[0:n-sn])
	}
	copy(dst[n-sn:],src)
	return nil
}
