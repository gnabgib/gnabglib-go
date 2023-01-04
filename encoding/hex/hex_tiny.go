// Copyright 2022 gnabgib
// This Source Code Form is subject to the terms of the Mozilla Public License v2.0

//go:build tiny
// +build tiny

package hex

// Instead of a 256 byte table, use a mask and some bit shifting to convert (~50% slower)

const mask6 = 0b11011111

//Convert a hex-char byte (0-9,a-f,A-F) into a byte (0-15, >15 invalid)
func decode(b byte) byte {
	// Zero out th 6th bit, subtract 48 and a further 7 when value >17
	// Takes ~950ns, uses 1 byte

	//Surprisingly a branch here is swifter than bit shift and invert (-200ns)
	// mask:=^(((b+6)>>1)&0b00100000)
	// b=(b&mask)-48
	if b > 57 {
		//Fold a-f -> A-F,
		//Push :-? into invalid territory
		//Make sure max numbers are 32 lower
		b &= mask6
	}
	b -= 48

	//0-9 are now 0-9            - Good
	//:-? are out of range       - Good
	//a-f are converted into A-F - Good
	//A-F are 17-22              - FIX
	//Instead of a divide (/17) or a branch (>=17), let's +15 and make sure it's
	// a multiple of 32 before subtracting 7
	return b - 7*((b+15)>>5)
}
