// Copyright 2022 gnabgib
// This Source Code Form is subject to the terms of the Mozilla Public License v2.0

package whirlpool

import (
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
	digestSizeU64   = digestSizeBytes >> 3 // /8
	lengthBytes     = 32                   //256bits
	rounds          = 10
	gfPoly          = 0x11d //GF(2^8) polynomial=x8+x4+x3+x2+1
	sBox            = "" +
		"\x18\x23\xc6\xe8\x87\xb8\x01\x4f\x36\xa6\xd2\xf5\x79\x6f\x91\x52" +
		"\x60\xbc\x9b\x8e\xa3\x0c\x7b\x35\x1d\xe0\xd7\xc2\x2e\x4b\xfe\x57" +
		"\x15\x77\x37\xe5\x9f\xf0\x4a\xda\x58\xc9\x29\x0a\xb1\xa0\x6b\x85" +
		"\xbd\x5d\x10\xf4\xcb\x3e\x05\x67\xe4\x27\x41\x8b\xa7\x7d\x95\xd8" +
		//
		"\xfb\xee\x7c\x66\xdd\x17\x47\x9e\xca\x2d\xbf\x07\xad\x5a\x83\x33" +
		"\x63\x02\xaa\x71\xc8\x19\x49\xd9\xf2\xe3\x5b\x88\x9a\x26\x32\xb0" +
		"\xe9\x0f\xd5\x80\xbe\xcd\x34\x48\xff\x7a\x90\x5f\x20\x68\x1a\xae" +
		"\xb4\x54\x93\x22\x64\xf1\x73\x12\x40\x08\xc3\xec\xdb\xa1\x8d\x3d" +
		//
		"\x97\x00\xcf\x2b\x76\x82\xd6\x1b\xb5\xaf\x6a\x50\x45\xf3\x30\xef" +
		"\x3f\x55\xa2\xea\x65\xba\x2f\xc0\xde\x1c\xfd\x4d\x92\x75\x06\x8a" +
		"\xb2\xe6\x0e\x1f\x62\xd4\xa8\x96\xf9\xc5\x25\x59\x84\x72\x39\x4c" +
		"\x5e\x78\x38\x8c\xd1\xa5\xe2\x61\xb3\x21\x9c\x1e\x43\xc7\xfc\x04" +
		//
		"\x51\x99\x6d\x0d\xfa\xdf\x7e\x24\x3b\xab\xce\x11\x8f\x4e\xb7\xeb" +
		"\x3c\x81\x94\xf7\xb9\x13\x2c\xd3\xe7\x6e\xc4\x03\x56\x44\x7f\xa9" +
		"\x2a\xbb\xc1\x53\xdc\x0b\x9d\x6c\x31\x74\xf6\x46\xac\x89\x14\xe1" +
		"\x16\x3a\x69\x09\x70\xb6\xd0\xed\xcc\x42\x98\xa4\x28\x5c\xf8\x86" // substitutionBox

)

var (
	ct = [8 * 256]uint64{} // circulantTable
	rc = [rounds]uint64{}  // roundConstants
)

type ctx struct {
	state    [blockSizeU64]uint64 //Runtime state of hash
	lenBytes uint64               //Number of bytes added to state (in total)
	block    [blockSizeBytes]byte //Temp processing block
	bPos     int                  //Position of data written to block
}

func init() {
	buildCirculantTable()
	buildRoundConstants()
	//Zero pad to 16 digits, the uppercase hex version of a value:
	//fmt.Println(strings.ToUpper(fmt.Sprintf("%016x",circulantTable[q])))
}

func buildCirculantTable() {
	for x := 0; x < 256; x++ {
		v1 := uint64(sBox[x])
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
		ct[x] = (v1 << 56) | (v1 << 48) |
			(v4 << 40) | (v1 << 32) |
			(v8 << 24) | (v5 << 16) |
			(v2 << 8) | v9
		for t := 1; t < 8; t++ {
			ct[(t<<8)|x] = bits.RotateLeft64(ct[((t-1)<<8)|x], -8)
		}
	}
}

func buildRoundConstants() {
	for r := 0; r < rounds; r++ {
		r8 := r << 3
		rc[r] = (ct[r8] & 0xff00000000000000) ^
			(ct[(1<<8)|(r8+1)] & 0x00ff000000000000) ^
			(ct[(2<<8)|(r8+2)] & 0x0000ff0000000000) ^
			(ct[(3<<8)|(r8+3)] & 0x000000ff00000000) ^
			(ct[(4<<8)|(r8+4)] & 0x00000000ff000000) ^
			(ct[(5<<8)|(r8+5)] & 0x0000000000ff0000) ^
			(ct[(6<<8)|(r8+6)] & 0x000000000000ff00) ^
			(ct[(7<<8)|(r8+7)] & 0x00000000000000ff)
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
			L[i] = ct[K[i]>>56&0xff] ^
				ct[256|(K[(i-1)&7]>>48&0xff)] ^
				ct[512|(K[(i-2)&7]>>40&0xff)] ^
				ct[768|(K[(i-3)&7]>>32&0xff)] ^
				ct[1024|(K[(i-4)&7]>>24&0xff)] ^
				ct[1280|(K[(i-5)&7]>>16&0xff)] ^
				ct[1536|(K[(i-6)&7]>>8&0xff)] ^
				ct[1792|(K[(i-7)&7]&0xff)]
		}
		L[0] ^= rc[r]
		for i := 0; i < 8; i++ {
			//Update K for next round
			K[i] = L[i]
			L[i] ^= ct[state[i]>>56&0xff] ^
				ct[256|(state[(i-1)&7]>>48&0xff)] ^
				ct[512|(state[(i-2)&7]>>40&0xff)] ^
				ct[768|(state[(i-3)&7]>>32&0xff)] ^
				ct[1024|(state[(i-4)&7]>>24&0xff)] ^
				ct[1280|(state[(i-5)&7]>>16&0xff)] ^
				ct[1536|(state[(i-6)&7]>>8&0xff)] ^
				ct[1792|(state[(i-7)&7]&0xff)]
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
	c.lenBytes += uint64(n)

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
	h.block[blockSizeBytes-1] = byte(c.lenBytes << 3)
	h.block[blockSizeBytes-2] = byte(c.lenBytes >> 5)
	h.block[blockSizeBytes-3] = byte(c.lenBytes >> 13)
	h.block[blockSizeBytes-4] = byte(c.lenBytes >> 21)
	h.block[blockSizeBytes-5] = byte(c.lenBytes >> 29)
	h.block[blockSizeBytes-6] = byte(c.lenBytes >> 37)
	h.block[blockSizeBytes-7] = byte(c.lenBytes >> 45)
	h.block[blockSizeBytes-8] = byte(c.lenBytes >> 53)
	h.block[blockSizeBytes-9] = byte(c.lenBytes >> 61)
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
	c.lenBytes = 0
	c.bPos = 0
}

func (c *ctx) Size() int { return digestSizeBytes }

func (c *ctx) BlockSize() int { return blockSizeBytes }
