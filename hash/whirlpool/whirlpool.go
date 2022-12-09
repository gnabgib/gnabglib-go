// Copyright 2022 gnabgib
// This Source Code Form is subject to the terms of the Mozilla Public License v2.0

package whirlpool

import (
	//	"fmt"
	"hash"
	"math/bits"
)

//https://en.wikipedia.org/wiki/Whirlpool_(hash_function)
//https://www2.seas.gwu.edu/~poorvi/Classes/CS381_2007/Whirlpool.pdf

//Limitation: Whirlpool can run on data of length 2^256 bits or 2^253 bytes, but
// we only track up to 2^64 bytes.. 2^67 bits because of size being tracked in bytes with a uint64

const (
	blockSizeBytes  = 64
	blockSizeU64    = blockSizeBytes >> 3 // /8
	digestSizeBytes = blockSizeBytes
	digestSizeU64   = digestSizeBytes >> 3 ///8
	lengthBytes     = 32                   //256bits!
	rounds          = 10
	gfPoly          = 0x11d //GF(2^8) polynomial=x8+x4+x3+x2+1
	//sboxPoly        = 0x13  //x4+x+1 10011
)

var substitutionBox = [...]uint8{
	//block 0
	0x18, 0x23, 0xc6, 0xe8, 0x87, 0xb8, 0x01, 0x4f, 0x36, 0xa6, 0xd2, 0xf5, 0x79, 0x6f, 0x91, 0x52,
	0x60, 0xbc, 0x9b, 0x8e, 0xa3, 0x0c, 0x7b, 0x35, 0x1d, 0xe0, 0xd7, 0xc2, 0x2e, 0x4b, 0xfe, 0x57,
	0x15, 0x77, 0x37, 0xe5, 0x9f, 0xf0, 0x4a, 0xda, 0x58, 0xc9, 0x29, 0x0a, 0xb1, 0xa0, 0x6b, 0x85,
	0xbd, 0x5d, 0x10, 0xf4, 0xcb, 0x3e, 0x05, 0x67, 0xe4, 0x27, 0x41, 0x8b, 0xa7, 0x7d, 0x95, 0xd8,
	//block 1
	0xfb, 0xee, 0x7c, 0x66, 0xdd, 0x17, 0x47, 0x9e, 0xca, 0x2d, 0xbf, 0x07, 0xad, 0x5a, 0x83, 0x33,
	0x63, 0x02, 0xaa, 0x71, 0xc8, 0x19, 0x49, 0xd9, 0xf2, 0xe3, 0x5b, 0x88, 0x9a, 0x26, 0x32, 0xb0,
	0xe9, 0x0f, 0xd5, 0x80, 0xbe, 0xcd, 0x34, 0x48, 0xff, 0x7a, 0x90, 0x5f, 0x20, 0x68, 0x1a, 0xae,
	0xb4, 0x54, 0x93, 0x22, 0x64, 0xf1, 0x73, 0x12, 0x40, 0x08, 0xc3, 0xec, 0xdb, 0xa1, 0x8d, 0x3d,
	//block 2
	0x97, 0x00, 0xcf, 0x2b, 0x76, 0x82, 0xd6, 0x1b, 0xb5, 0xaf, 0x6a, 0x50, 0x45, 0xf3, 0x30, 0xef,
	0x3f, 0x55, 0xa2, 0xea, 0x65, 0xba, 0x2f, 0xc0, 0xde, 0x1c, 0xfd, 0x4d, 0x92, 0x75, 0x06, 0x8a,
	0xb2, 0xe6, 0x0e, 0x1f, 0x62, 0xd4, 0xa8, 0x96, 0xf9, 0xc5, 0x25, 0x59, 0x84, 0x72, 0x39, 0x4c,
	0x5e, 0x78, 0x38, 0x8c, 0xd1, 0xa5, 0xe2, 0x61, 0xb3, 0x21, 0x9c, 0x1e, 0x43, 0xc7, 0xfc, 0x04,
	//block 3
	0x51, 0x99, 0x6d, 0x0d, 0xfa, 0xdf, 0x7e, 0x24, 0x3b, 0xab, 0xce, 0x11, 0x8f, 0x4e, 0xb7, 0xeb,
	0x3c, 0x81, 0x94, 0xf7, 0xb9, 0x13, 0x2c, 0xd3, 0xe7, 0x6e, 0xc4, 0x03, 0x56, 0x44, 0x7f, 0xa9,
	0x2a, 0xbb, 0xc1, 0x53, 0xdc, 0x0b, 0x9d, 0x6c, 0x31, 0x74, 0xf6, 0x46, 0xac, 0x89, 0x14, 0xe1,
	0x16, 0x3a, 0x69, 0x09, 0x70, 0xb6, 0xd0, 0xed, 0xcc, 0x42, 0x98, 0xa4, 0x28, 0x5c, 0xf8, 0x86,
}

//1^F=E
//E=1
//1^1=0
//F^1=E

var circulantTable = [8 * 256]uint64{}
var roundConstants = [rounds]uint64{}

type ctx struct {
	state [blockSizeU64]uint64 //Runtime state of hash
	len   uint64               //Number of bytes added to state (in total)
	block [blockSizeBytes]byte //Temp processing block
	bPos  int                  //Position of data written to block
}

func init() {
	buildCirculantTable()
	buildRoundConstants()

	// es := [...]uint8{0x1, 0xb, 0x9, 0xc, 0xd, 0x6, 0xf, 0x3, 0xe, 0x8, 0x7, 0x4, 0xa, 0x2, 0x5, 0x0}
	// eps := [...]uint8{0xf, 0x0, 0xd, 0x7, 0xb, 0xe, 0x5, 0xa, 0x9, 0x2, 0xc, 0x1, 0x3, 0x4, 0x8, 0x6}
	// rs := [...]uint8{0x7, 0xc, 0xb, 0xd, 0xe, 0x4, 0x9, 0xf, 0x6, 0x3, 0x8, 0xa, 0x2, 0x5, 0x1, 0x0}

	// sb := [256]uint8{}
	// for i := 0; i < 16; i++ {
	// 	for j := 0; j < 16; j++ {
	// 		e := es[i] ^ sboxPoly
	// 		ep := eps[j] ^ sboxPoly
	// 		ex := e ^ ep
	// 		r := rs[ex]
	// 		sb[i*16+j] = ((e ^ r) << 4) | (ep ^ r)
	// 	}
	// }
	// fmt.Println(sb)

	//Zero pad to 16 digits, the uppercase hex version of a value:
	//fmt.Println(strings.ToUpper(fmt.Sprintf("%016x",circulantTable[q])))
}

func buildCirculantTable() {
	for x := 0; x < 256; x++ {
		v1 := uint64(substitutionBox[x])
		v2 := v1 << 1
		if v2 >= 0x100 {
			v2 ^= gfPoly
		}
		v4 := v2 << 1
		if v4 >= 0x100 {
			v4 ^= gfPoly
		}
		v5 := v4 ^ v1
		v8 := v4 << 1
		if v8 >= 0x100 {
			v8 ^= gfPoly
		}
		v9 := v8 ^ v1
		circulantTable[x] = (v1 << 56) | (v1 << 48) |
			(v4 << 40) | (v1 << 32) |
			(v8 << 24) | (v5 << 16) |
			(v2 << 8) | v9
		for t := 1; t < 8; t++ {
			circulantTable[(t<<8)|x] = bits.RotateLeft64(circulantTable[((t-1)<<8)|x], -8)
		}
	}
}

func buildRoundConstants() {
	for r := 0; r < rounds; r++ {
		r8 := r << 3
		roundConstants[r] = (circulantTable[r8] & 0xff00000000000000) ^
			(circulantTable[(1<<8)|(r8+1)] & 0x00ff000000000000) ^
			(circulantTable[(2<<8)|(r8+2)] & 0x0000ff0000000000) ^
			(circulantTable[(3<<8)|(r8+3)] & 0x000000ff00000000) ^
			(circulantTable[(4<<8)|(r8+4)] & 0x00000000ff000000) ^
			(circulantTable[(5<<8)|(r8+5)] & 0x0000000000ff0000) ^
			(circulantTable[(6<<8)|(r8+6)] & 0x000000000000ff00) ^
			(circulantTable[(7<<8)|(r8+7)] & 0x00000000000000ff)
	}
}

// A new hash for computing Whirlpool
func New() hash.Hash {
	c := &ctx{}
	c.Reset()
	return c
}

// aka transform
func (c *ctx) hash() {
	x := [blockSizeU64]uint64{}     // data
	K := [blockSizeU64]uint64{}     // round key
	state := [blockSizeU64]uint64{} // cipher state
	L := [blockSizeU64]uint64{}

	//Initialize K and state
	for i := 0; i < blockSizeU64; i++ {
		j := i * 8
		x[i] = uint64(c.block[j])<<56 |
			uint64(c.block[j+1])<<48 |
			uint64(c.block[j+2])<<40 |
			uint64(c.block[j+3])<<32 |
			uint64(c.block[j+4])<<24 |
			uint64(c.block[j+5])<<16 |
			uint64(c.block[j+6])<<8 |
			uint64(c.block[j+7])
		K[i] = c.state[i]
		state[i] = x[i] ^ K[i]
	}

	for r := 0; r < rounds; r++ {
		for i := 0; i < 8; i++ {
			L[i] = circulantTable[K[i]>>56&0xff] ^
				circulantTable[256|(K[(i-1)&7]>>48&0xff)] ^
				circulantTable[512|(K[(i-2)&7]>>40&0xff)] ^
				circulantTable[768|(K[(i-3)&7]>>32&0xff)] ^
				circulantTable[1024|(K[(i-4)&7]>>24&0xff)] ^
				circulantTable[1280|(K[(i-5)&7]>>16&0xff)] ^
				circulantTable[1536|(K[(i-6)&7]>>8&0xff)] ^
				circulantTable[1792|(K[(i-7)&7]&0xff)]
		}
		L[0] ^= roundConstants[r]
		for i := 0; i < 8; i++ {
			//Update K for next round
			K[i] = L[i]
			L[i] = L[i] ^
				circulantTable[state[i]>>56&0xff] ^
				circulantTable[256|(state[(i-1)&7]>>48&0xff)] ^
				circulantTable[512|(state[(i-2)&7]>>40&0xff)] ^
				circulantTable[768|(state[(i-3)&7]>>32&0xff)] ^
				circulantTable[1024|(state[(i-4)&7]>>24&0xff)] ^
				circulantTable[1280|(state[(i-5)&7]>>16&0xff)] ^
				circulantTable[1536|(state[(i-6)&7]>>8&0xff)] ^
				circulantTable[1792|(state[(i-7)&7]&0xff)]
		}
		//L references state (in prior loop) so this needs to be separate
		for i := 0; i < 8; i++ {
			state[i] = L[i]
		}
	}

	//Miyaguchi-Preneel compression
	for i := 0; i < blockSizeU64; i++ {
		c.state[i] ^= state[i] ^ x[i]
	}
	c.bPos = 0
}

func (c *ctx) Write(p []byte) (n int, err error) {
	n = len(p)
	c.len += uint64(n)

	nToWrite := n
	space := blockSizeBytes - c.bPos
	for nToWrite > 0 {
		if space > nToWrite {
			//More space than data, copy the data in
			for i := 0; i < nToWrite; i++ {
				c.block[c.bPos+i] = p[i]
			}
			//Update block pos and return
			c.bPos += nToWrite
			return
		}
		//Otherwise write to the end of the space
		for i := 0; i < space; i++ {
			c.block[c.bPos+i] = p[i]
		}
		c.bPos += space
		c.hash()      //Process the block
		p = p[space:] //Re-slice for the next section
		nToWrite -= space
		space = blockSizeBytes //Max space from now on
	}
	return
}

func (c *ctx) Sum(in []byte) []byte {
	//Since sum isn't supposed to mutate the hash so far, make a copy
	h := *c

	const sizeSpace = blockSizeBytes - lengthBytes

	//Zero array for when we need more space
	zeros := make([]byte, blockSizeBytes)

	//There can never be 0 free bytes so write the separator
	h.block[h.bPos] = 0x80
	h.bPos++

	//If there's not enough space for the (enormous) size, zero out and hash
	if h.bPos > sizeSpace {
		copy(h.block[h.bPos:], zeros)
		h.hash()
		h.bPos = 0
	}
	//Zero rest of bytes (over zero, but useful for size)
	copy(h.block[h.bPos:blockSizeBytes-9], zeros)
	//Big endian size
	// Because we have length in bytes, not bits the shifts are all -3 (left 3)
	h.block[blockSizeBytes-1] = byte(c.len << 3)
	h.block[blockSizeBytes-2] = byte(c.len >> 5)
	h.block[blockSizeBytes-3] = byte(c.len >> 13)
	h.block[blockSizeBytes-4] = byte(c.len >> 21)
	h.block[blockSizeBytes-5] = byte(c.len >> 29)
	h.block[blockSizeBytes-6] = byte(c.len >> 37)
	h.block[blockSizeBytes-7] = byte(c.len >> 45)
	h.block[blockSizeBytes-8] = byte(c.len >> 53)
	h.block[blockSizeBytes-9] = byte(c.len >> 61)
	h.hash()

	out := make([]byte, digestSizeBytes)
	for i := 0; i < digestSizeU64; i++ {
		j := i * 8
		out[j] = byte(h.state[i] >> 56)
		out[j+1] = byte(h.state[i] >> 48)
		out[j+2] = byte(h.state[i] >> 40)
		out[j+3] = byte(h.state[i] >> 32)
		out[j+4] = byte(h.state[i] >> 24)
		out[j+5] = byte(h.state[i] >> 16)
		out[j+6] = byte(h.state[i] >> 8)
		out[j+7] = byte(h.state[i])
	}
	return append(in, out...) //Shake it all about
}

func (c *ctx) Reset() {
	for i := 0; i < blockSizeU64; i++ {
		c.state[i] = 0
	}
	c.len = 0
	c.bPos = 0
}

func (c *ctx) Size() int { return digestSizeBytes }

func (c *ctx) BlockSize() int { return blockSizeBytes }
