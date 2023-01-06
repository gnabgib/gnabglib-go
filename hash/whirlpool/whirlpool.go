// Copyright 2022 gnabgib
// This Source Code Form is subject to the terms of the Mozilla Public License v2.0

package whirlpool

import (
	"encoding/binary"
	"hash"

	"github.com/gnabgib/gnablib-go/bytes"
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
)

type ctx struct {
	state    [blockSizeU64]uint64 //Runtime state of hash
	lenBytes uint64               //Number of bytes added to state (in total)
	block    [blockSizeBytes]byte //Temp processing block
	bPos     int                  //Position of data written to block
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
			copy(c.block[c.bPos:], p)
			//Update block pos and return
			c.bPos += nToWrite
			return
		}
		//Otherwise write to the end of the space
		copy(c.block[c.bPos:], p[0:space])
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

	//There can never be 0 free bytes so write the separator
	h.block[h.bPos] = 0x80
	h.bPos++

	//If there's not enough space for the (enormous) size, zero out and hash
	if h.bPos > sizeSpace {
		bytes.Zero(h.block[h.bPos:])
		h.hash()
	}
	//Zero rest of bytes (over zero, but useful for size)
	bytes.Zero(h.block[h.bPos : blockSizeBytes-9])

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
		binary.BigEndian.PutUint64(out[i*8:], h.state[i])
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
