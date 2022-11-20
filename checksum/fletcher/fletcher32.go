// Copyright 2022 gnabgib
// This Source Code Form is subject to the terms of the Mozilla Public License v2.0

package fletcher

//http://www.zlib.net/maxino06_fletcher-adler.pdf -> Lower cpu and Adler and mostly more effective (in their tests)
//https://datatracker.ietf.org/doc/html/rfc1146 (Appendix I)

import (
	"hash"
)

const size32=4

type digest32 struct {
	a, b uint16
}

// A new Hash32 for computing the Fletcher 32 checksum
func New32() hash.Hash32 {
	d := new(digest32)
	d.Reset()
	return d
}

func (d *digest32) update(p []byte) {
	//Using the Anastase Nakassis optimization, requires 64bit processor/os
	const mod = 0xffff
	const c1Overflow = 23726746
	var c0, c1 uint64
	c0 = uint64((*d).a)
	c1 = uint64((*d).b)

	n := len(p)
	for i := 0; i < n; {
		len := n - i
		if len > c1Overflow {
			len = c1Overflow
		}
		for j := 0; j < len; j += 2 {
			//We only use j for the batch-count (incrementing i each loop)

			//Fletcher is little endian
			add := uint64(p[i])
			i++
			if i < n {
				add |= uint64(p[i]) << 8
				i++
			}
			c0 += add
			c1 += c0
		}
		c0 %= mod
		c1 %= mod
	}
	(*d).a = uint16(c0)
	(*d).b = uint16(c1)
}

func (d *digest32) Write(p []byte) (n int, err error) {
	d.update(p)
	return len(p), nil
}

func (d *digest32) Sum(in []byte) []byte {
	return append(in, byte((*d).b>>8), byte((*d).b), byte((*d).a>>8), byte((*d).a))
}

func (d *digest32) Reset() {
	(*d).a = 0
	(*d).b = 0
}

func (d *digest32) Size() int { return size32 }

func (d *digest32) BlockSize() int { return size32/2 }

func (d *digest32) Sum32() uint32 { return uint32((*d).b)<<16 | uint32((*d).a) }
