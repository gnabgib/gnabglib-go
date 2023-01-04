// Copyright 2022 gnabgib
// This Source Code Form is subject to the terms of the Mozilla Public License v2.0

package ripemd

import (
	"encoding/binary"
	"unsafe"
)

//https://en.wikipedia.org/wiki/RIPEMD
//https://homes.esat.kuleuven.be/~bosselae/ripemd/rmd128.txt
//https://homes.esat.kuleuven.be/~bosselae/ripemd160.html (1996)

// Constants __ __ __ __ __ __ __ __ __ __ __ __ __ __ __ __ __ __

func f0(x, y, z uint32) uint32 { return x ^ y ^ z }         //Same as MD4-r3
func f1(x, y, z uint32) uint32 { return z ^ (x & (y ^ z)) } // like MD4-r1, optimize from (x&y)|(~x&z)
func f2(x, y, z uint32) uint32 { return (x | ^y) ^ z }
func f3(x, y, z uint32) uint32 { return y ^ (z & (x ^ y)) } // like MD4-r1, optimize from (x&z)|(y&~z)
func f4(x, y, z uint32) uint32 { return x ^ (y | ^z) }

const (
	r = "" +
		//r 0..15
		"\x00\x01\x02\x03\x04\x05\x06\x07\x08\x09\x0A\x0B\x0C\x0D\x0E\x0F" +
		//r 16..31
		"\x07\x04\x0D\x01\x0A\x06\x0F\x03\x0C\x00\x09\x05\x02\x0E\x0B\x08" +
		//r 32..47
		"\x03\x0A\x0E\x04\x09\x0F\x08\x01\x02\x07\x00\x06\x0D\x0B\x05\x0C" +
		//r 48..63
		"\x01\x09\x0B\x0A\x00\x08\x0C\x04\x0D\x03\x07\x0F\x0E\x05\x06\x02" +
		//r 64..79
		"\x04\x00\x05\x09\x07\x0C\x02\x0A\x0E\x01\x03\x08\x0B\x06\x0F\x0D"
	rr = "" +
		//r' 0..15
		"\x05\x0E\x07\x00\x09\x02\x0B\x04\x0D\x06\x0F\x08\x01\x0A\x03\x0C" +
		//r' 16..31
		"\x06\x0B\x03\x07\x00\x0D\x05\x0A\x0E\x0F\x08\x0C\x04\x09\x01\x02" +
		//r' 32..47
		"\x0F\x05\x01\x03\x07\x0E\x06\x09\x0B\x08\x0C\x02\x0A\x00\x04\x0D" +
		//r' 48..63
		"\x08\x06\x04\x01\x03\x0B\x0F\x00\x05\x0C\x02\x0D\x09\x07\x0A\x0E" +
		//r' 64..79
		"\x0C\x0F\x0A\x04\x01\x05\x08\x07\x06\x02\x0D\x0E\x00\x03\x09\x0B"
	s = "" +
		//s 0..15
		"\x0B\x0E\x0F\x0C\x05\x08\x07\x09\x0B\x0D\x0E\x0F\x06\x07\x09\x08" +
		//s 16..31
		"\x07\x06\x08\x0D\x0B\x09\x07\x0F\x07\x0C\x0F\x09\x0B\x07\x0D\x0C" +
		//s 32..47
		"\x0B\x0D\x06\x07\x0E\x09\x0D\x0F\x0E\x08\x0D\x06\x05\x0C\x07\x05" +
		//s 48..63
		"\x0B\x0C\x0E\x0F\x0E\x0F\x09\x08\x09\x0E\x05\x06\x08\x06\x05\x0C" +
		//s 64..79
		"\x09\x0F\x05\x0B\x06\x08\x0D\x0C\x05\x0C\x0D\x0E\x0B\x08\x05\x06"
	ss = "" +
		//s' 0..15
		"\x08\x09\x09\x0B\x0D\x0F\x0F\x05\x07\x07\x08\x0B\x0E\x0E\x0C\x06" +
		//s' 16..31
		"\x09\x0D\x0F\x07\x0C\x08\x09\x0B\x07\x07\x0C\x07\x06\x0F\x0D\x0B" +
		//s' 32..47
		"\x09\x07\x0F\x0B\x08\x06\x06\x0E\x0C\x0D\x05\x0E\x0D\x0D\x07\x05" +
		//s' 48..63
		"\x0F\x05\x08\x0B\x0E\x0E\x06\x0E\x06\x09\x0C\x09\x0C\x05\x0F\x08" +
		//s' 64..79
		"\x08\x05\x0C\x09\x0C\x05\x0E\x06\x08\x0D\x06\x05\x0F\x0D\x0B\x0B"
	//0,int(2**30 x sqrt(2)), int(2**30 x sqrt(3)),int(2**30 x sqrt(5)),int(2**30 x sqrt(7))
	k = "\x00\x00\x00\x00" +
		"\x5a\x82\x79\x99" +
		"\x6e\xd9\xeb\xa1" +
		"\x8f\x1b\xbc\xdc" +
		"\xa9\x53\xfd\x4e"
	// int(2**30 x cbrt(2)),int(2**30 x cbrt(3)),int(2**30 x cbrt(5)),int(2**30 x cbrt(7)),0
	kk = "\x50\xa2\x8b\xe6" +
		"\x5c\x4d\xd1\x24" +
		"\x6d\x70\x3e\xf3" +
		"\x7a\x6d\x76\xe9" +
		"\x00\x00\x00\x00"
	// In 128/256 the last constant of the parallel set is zeroed, but otherwise notice these are the same as @see kk
	kk128 = "\x50\xa2\x8b\xe6" +
		"\x5c\x4d\xd1\x24" +
		"\x6d\x70\x3e\xf3" +
		"\x00\x00\x00\x00"
	iv = "\x67\x45\x23\x01" +
		"\xef\xcd\xab\x89" +
		"\x98\xba\xdc\xfe" +
		"\x10\x32\x54\x76" +
		"\xc3\xd2\xe1\xf0"
	iv2 = "\x76\x54\x32\x10" +
		"\xfe\xdc\xba\x98" +
		"\x89\xab\xcd\xef" +
		"\x01\x23\x45\x67" +
		"\x3c\x2d\x1e\x0f"
	u32Size        = int(unsafe.Sizeof(uint32(0)))
	blockSizeBytes = 64 //512 bits
	blockSizeU32   = blockSizeBytes / u32Size
	sizeSpace      = blockSizeBytes - 2*u32Size //64bit uint representing size
)

var (
	f = [...]func(x, y, z uint32) uint32{f0, f1, f2, f3, f4}
)

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
			c.state[i] = binary.BigEndian.Uint32([]byte(iv[i*4 : i*4+4]))
			c.state[i+n] = binary.BigEndian.Uint32([]byte(iv2[i*4 : i*4+4]))
		}
	} else {
		for i := 0; i < n; i++ {
			c.state[i] = binary.BigEndian.Uint32([]byte(iv[i*4 : i*4+4]))
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
