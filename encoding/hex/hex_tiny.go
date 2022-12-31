// Copyright 2022 gnabgib
// This Source Code Form is subject to the terms of the Mozilla Public License v2.0

//go:build tiny
// +build tiny

package hex

// Instead of a 256 byte table, use a mask and some bit shifting to convert (~half speed)

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

// Using lookup tables for sub/fold/max/add totalling 16 bytes, it takes ~1100ns
// We divide ascii into quarters
// const (
// 	sub4  = "0A  "          //16b
// 	fold4 = "\xFF_\xFF\xFF" //Only fold 2nd col making a->A using 0b1011111=_
// 	max4  = "\x0A\x06\x00\x00"
// 	add4  = "\x00\x0A  "
// )

// func decode_v4(b byte) byte {
// 	//~1100ns, 16b const
// 	//Divide in 4 (6 bits each)
// 	u2 := b >> 6
// 	b = b&fold4[u2] - sub4[u2]
// 	if b >= max4[u2] {
// 		return 20
// 	}
// 	return b + add4[u2]
// }

// Using lookup tables for sub/max/add totally 20 bytes, it takes ~1100ns
// const (
// 	sub3 = " 0Aa    " //20b
// 	max3 = "\x00\x0A\x06\x06\x00\x00\x00\x00"
// 	add3 = " \x00\x0A\x0A"
// )

// func decode_v3(b byte) byte {
// 	//~1100ns, 20b const
// 	//Divide into 8 (5 bits each)
// 	u3 := b >> 5
// 	b -= sub3[u3]
// 	if b >= max3[u3] {
// 		return 20
// 	}
// 	return b + add3[u3]
// }

// Using lookup tables for  whether a value is less than 10, 6 (halving first)
// and confirming the right "page" of ascii (3=0-9,4=A-F,6=a-f)
// Takes ~2600ns and uses 6b of storage
// const (
// 	lt10lookup = byte(0b00011111)
// 	lt6lookup  = byte(0b00000111)
// 	u4_num     = uint16(0b0000000000001000) //u4=3
// 	u4_alpha   = uint16(0b0000000001010000) //u4=4,6
// )

// func decode_v2(b byte) byte {
// 	//~2600ns 10b const
// 	//We can ignore x upper bit because of invalid calc below
// 	// u4=011 l4<10
// 	// u4=1*0 l4-1<6
// 	u4 := b >> 4
// 	l4 := b & 0xf
// 	l4lt10 := 1 & byte(u4_num>>u4) & (lt10lookup >> (l4 >> 1))
// 	l4lt6 := 1 & byte(u4_alpha>>u4) & (lt6lookup >> ((l4 - 1) >> 1))
// 	invalid := 1 - (l4lt6 | l4lt10)
// 	return l4 +
// 		l4lt6*9 +
// 		invalid*20
// }

// Using lookup tables for whether value is less than 10,6 (halving first)
// ASCII page is determined with bit logic (expensive)
// Takes ~3200ns, 2b of storage
// const (
// 	lt10lookup = byte(0b00011111)
// 	lt6lookup  = byte(0b00000111)
// )

// func decode_v1(b byte) byte {
// 	//~3200ns, 2b const
// 	//Stock is ~600ns (x64 arch)
// 	//We can ignore x upper bit because of invalid calc below
// 	//Valid x=*0, y=11, l4<10
// 	//Valid x=*1, y=*0, l4-1<6
// 	// by=011 l4<10
// 	// by=1*0 l4-1<6
// 	x := b >> 6
// 	y := (b >> 4) & 3
// 	l4 := b & 0xf

// 	l4lt10 := (1 - x&1) & (y & 1 & (y >> 1)) & ((lt10lookup >> (l4 >> 1)) & 1)
// 	l4lt6 := (x & 1) & (1 - y&1) & ((lt6lookup >> ((l4 - 1) >> 1)) & 1)

// 	invalid := b>>7 | (1 - (l4lt10 | l4lt6))

// 	return (b & 0xf) + (b >> 6) + ((b >> 6) << 3) + invalid*20
// }
