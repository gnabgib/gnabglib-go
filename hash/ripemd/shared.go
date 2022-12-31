// Copyright 2022 gnabgib
// This Source Code Form is subject to the terms of the Mozilla Public License v2.0

package ripemd

import (
	"unsafe"
)

//https://en.wikipedia.org/wiki/RIPEMD
//https://homes.esat.kuleuven.be/~bosselae/ripemd160.html (1996)

// Constants __ __ __ __ __ __ __ __ __ __ __ __ __ __ __ __ __ __

var f = [...]func(x, y, z uint32) uint32{
	func(x, y, z uint32) uint32 { return x ^ y ^ z },         //Same as MD4-r3
	func(x, y, z uint32) uint32 { return z ^ (x & (y ^ z)) }, // like MD4-r1, optimize from (x&y)|(~x&z)
	func(x, y, z uint32) uint32 { return (x | ^y) ^ z },
	func(x, y, z uint32) uint32 { return y ^ (z & (x ^ y)) }, // like MD4-r1, optimize from (x&z)|(y&~z)
	func(x, y, z uint32) uint32 { return x ^ (y | ^z) },
}

// 0,int(2**30 x sqrt(2)), int(2**30 x sqrt(3)),int(2**30 x sqrt(5)),int(2**30 x sqrt(7))
var k = [...]uint32{0x00000000, 0x5a827999, 0x6ed9eba1, 0x8f1bbcdc, 0xa953fd4e}

// int(2**30 x cbrt(2)),int(2**30 x cbrt(3)),int(2**30 x cbrt(5)),int(2**30 x cbrt(7)),0
var kk = [...]uint32{0x50a28be6, 0x5c4dd124, 0x6d703ef3, 0x7a6d76e9, 0x00000000}

// In 128/256 the last constant of the parallel set is zeroed, but otherwise notice these are the same as @see kk
var kk128 = [...]uint32{0x50a28be6, 0x5c4dd124, 0x6d703ef3, 0x00000000}

var r = [...]int{
	0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, /*r 0..15 -- -- -- -- -- -- -- -- -- -- */
	7, 4, 13, 1, 10, 6, 15, 3, 12, 0, 9, 5, 2, 14, 11, 8, /*r 16..31 -- -- -- -- -- -- -- -- -- --*/
	3, 10, 14, 4, 9, 15, 8, 1, 2, 7, 0, 6, 13, 11, 5, 12, /*r 32..47 -- -- -- -- -- -- -- -- -- --*/
	1, 9, 11, 10, 0, 8, 12, 4, 13, 3, 7, 15, 14, 5, 6, 2, /*r 48..63 -- -- -- -- -- -- -- -- -- --*/
	4, 0, 5, 9, 7, 12, 2, 10, 14, 1, 3, 8, 11, 6, 15, 13, /*r 64..79 -- -- -- -- -- -- -- -- -- --*/
}
var rr = [...]int{
	5, 14, 7, 0, 9, 2, 11, 4, 13, 6, 15, 8, 1, 10, 3, 12, /*r' 0..15  -- -- -- -- -- -- -- -- -- -*/
	6, 11, 3, 7, 0, 13, 5, 10, 14, 15, 8, 12, 4, 9, 1, 2, /*r' 16..31 -- -- -- -- -- -- -- -- -- -*/
	15, 5, 1, 3, 7, 14, 6, 9, 11, 8, 12, 2, 10, 0, 4, 13, /*r' 32..47 -- -- -- -- -- -- -- -- -- -*/
	8, 6, 4, 1, 3, 11, 15, 0, 5, 12, 2, 13, 9, 7, 10, 14, /*r' 48..63 -- -- -- -- -- -- -- -- -- -*/
	12, 15, 10, 4, 1, 5, 8, 7, 6, 2, 13, 14, 0, 3, 9, 11, /*r' 64..79 -- -- -- -- -- -- -- -- -- -*/
}
var s = [...]int{
	11, 14, 15, 12, 5, 8, 7, 9, 11, 13, 14, 15, 6, 7, 9, 8, /*s 0..15 -- -- -- -- -- -- -- -- -- -*/
	7, 6, 8, 13, 11, 9, 7, 15, 7, 12, 15, 9, 11, 7, 13, 12, /*s 16..31 -- -- -- -- -- -- -- -- -- */
	11, 13, 6, 7, 14, 9, 13, 15, 14, 8, 13, 6, 5, 12, 7, 5, /*s 32..47 -- -- -- -- -- -- -- -- -- */
	11, 12, 14, 15, 14, 15, 9, 8, 9, 14, 5, 6, 8, 6, 5, 12, /*s 48..63 -- -- -- -- -- -- -- -- -- */
	9, 15, 5, 11, 6, 8, 13, 12, 5, 12, 13, 14, 11, 8, 5, 6, /*s 64..79 -- -- -- -- -- -- -- -- -- */
}
var ss = [...]int{
	8, 9, 9, 11, 13, 15, 15, 5, 7, 7, 8, 11, 14, 14, 12, 6, /*s 0..15 -- -- -- -- -- -- -- -- -- -*/
	9, 13, 15, 7, 12, 8, 9, 11, 7, 7, 12, 7, 6, 15, 13, 11, /*s 16..31 -- -- -- -- -- -- -- -- -- */
	9, 7, 15, 11, 8, 6, 6, 14, 12, 13, 5, 14, 13, 13, 7, 5, /*s 32..47 -- -- -- -- -- -- -- -- -- */
	15, 5, 8, 11, 14, 14, 6, 14, 6, 9, 12, 9, 12, 5, 15, 8, /*s 48..63 -- -- -- -- -- -- -- -- -- */
	8, 5, 12, 9, 12, 5, 14, 6, 8, 13, 6, 5, 15, 13, 11, 11, /*s 64..79 -- -- -- -- -- -- -- -- -- */
}
var iv = [...]uint32{0x67452301, 0xefcdab89, 0x98badcfe, 0x10325476, 0xc3d2e1f0}
var iv2 = [...]uint32{0x76543210, 0xfedcba98, 0x89abcdef, 0x01234567, 0x3c2d1e0f}

const u32Size = int(unsafe.Sizeof(uint32(0)))
const blockSizeBytes = 64 //512 bits
const blockSizeU32 = blockSizeBytes / u32Size
const sizeSpace = blockSizeBytes - 2*u32Size //64bit uint representing size

// Shared Context/Algo__ __ __ __ __ __ __ __ __ __ __ __ __ __ __

type ripeCtx struct {
	hash     func(*ripeCtx)       //Hashing func for when block is full
	state    [10]uint32           //Runtime state of hash
	stateLen int                  //Part of the state used (variants based)
	len      uint64               //Number of bytes added to state (in total)
	block    [blockSizeBytes]byte //Temp processing block
	bPos     int                  //Position of data written to block
}

func (c *ripeCtx) getX() []uint32 {
	ret := make([]uint32, blockSizeU32)
	for i := 0; i < blockSizeU32; i++ {
		j := i * 4
		//Little Endian conversion
		ret[i] = uint32(c.block[j]) | uint32(c.block[j+1])<<8 | uint32(c.block[j+2])<<16 | uint32(c.block[j+3])<<24
	}
	return ret
}

func (c *ripeCtx) Reset() {
	n := c.stateLen
	if n > 5 {
		//Deal with 256,320 loading iv2 into second half of space
		n /= 2
		for i := 0; i < n; i++ {
			c.state[i] = iv[i]
			c.state[i+n] = iv2[i]
		}
	} else {
		for i := 0; i < n; i++ {
			c.state[i] = iv[i]
		}
	}
	c.len = 0
	c.bPos = 0
}

func (c *ripeCtx) Write(p []byte) (n int, err error) {
	n = len(p)
	c.len += uint64(n)

	nToWrite := n
	space := blockSizeBytes - c.bPos
	for nToWrite > 0 {
		if space > nToWrite {
			//If there's more space than data, copy the data in
			for i := 0; i < nToWrite; i++ {
				c.block[c.bPos+i] = p[i]
			}
			c.bPos += nToWrite
			//And we're done
			return
		}
		//Otherwise write to the end of the space
		for i := 0; i < space; i++ {
			c.block[c.bPos+i] = p[i]
		}
		c.bPos += space
		//Process the block
		c.hash(c)
		//And repeat
		p = p[space:]
		nToWrite -= space
		space = blockSizeBytes
	}
	return
}

func (c *ripeCtx) Sum(in []byte) []byte {
	//Since sum isn't supposed to mutate the hash so far, make a copy
	t := *c
	h := &t

	//Because of the way write works, we must always have at least one
	// byte free (if there was zero, it would be hashed and there'd be 64)
	h.block[h.bPos] = 0x80
	h.bPos++
	//If we don't have enough space for the size, add zeros
	if h.bPos > sizeSpace {
		for h.bPos < blockSizeBytes {
			h.block[h.bPos] = 0
			h.bPos += 1
		}
		h.hash(h)
	}

	//Now add zeros until there's space for the size
	for h.bPos < sizeSpace {
		h.block[h.bPos] = 0
		h.bPos += 1
	}

	//Write the size.. in bits (it's stored in bytes *8 = <<3)
	h.block[h.bPos] = byte(h.len << 3)
	h.block[h.bPos+1] = byte(h.len >> 5)
	h.block[h.bPos+2] = byte(h.len >> 13)
	h.block[h.bPos+3] = byte(h.len >> 21)
	h.block[h.bPos+4] = byte(h.len >> 29)
	h.block[h.bPos+5] = byte(h.len >> 37)
	h.block[h.bPos+6] = byte(h.len >> 45)
	h.block[h.bPos+7] = byte(h.len >> 53)
	(*h).bPos += 8

	h.hash(h)
	//Append the state (which is the hash) to the input
	var out = make([]byte, h.stateLen*4)
	for i := 0; i < h.stateLen; i++ {
		j := i * 4
		out[j] = byte(h.state[i])
		out[j+1] = byte(h.state[i] >> 8)
		out[j+2] = byte(h.state[i] >> 16)
		out[j+3] = byte(h.state[i] >> 24)
	}
	return append(in, out...) //Shake it all about
}

func (c *ripeCtx) BlockSize() int { return blockSizeBytes }

func (c *ripeCtx) Size() int { return c.stateLen * 4 }
