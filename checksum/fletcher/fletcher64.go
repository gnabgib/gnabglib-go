// Copyright 2022 gnabgib
// This Source Code Form is subject to the terms of the Mozilla Public License v2.0

package fletcher

//http://www.zlib.net/maxino06_fletcher-adler.pdf -> Lower cpu and Adler and mostly more effective (in their tests)
//https://datatracker.ietf.org/doc/html/rfc1146 (Appendix I)

import (
	"hash"
)

const size64 = 8

type digest64 struct {
	a, b uint32
}

// A new Hash64 for computing the Fletcher 64 checksum
func New64() hash.Hash64 {
	d := new(digest64)
	d.Reset()
	return d
}

func (d *digest64) update(p []byte) {
	//Using the Anastase Nakassis optimization, requires 64bit processor/os
	const mod = 0xffffffff
	const c1Overflow = 92681
	var c0, c1 uint64
	c0 = uint64((*d).a)
	c1 = uint64((*d).b)

	n := len(p)
	for i := 0; i < n; {
		len := n - i
		if len > c1Overflow {
			len = c1Overflow
		}
		for j := 0; j < len; j += 4 {
			//We only use j for the batch-count (incrementing i each loop)

			//Fletcher is little endian
			add := uint64(p[i])
			i++
			if i < n {
				add |= uint64(p[i]) << 8
				i++
				if i < n {
					add |= uint64(p[i]) << 16
					i++
					if i < n {
						add |= uint64(p[i]) << 24
						i++
					}
				}
			}
			c0 += add
			c1 += c0
		}
		c0 %= mod
		c1 %= mod
	}
	(*d).a = uint32(c0)
	(*d).b = uint32(c1)
}

func (d *digest64) Write(p []byte) (n int, err error) {
	d.update(p)
	return len(p), nil
}

func (d *digest64) Sum(in []byte) []byte {
	return append(in, byte((*d).b>>24), byte((*d).b>>16), byte((*d).b>>8), byte((*d).b),
		byte((*d).a>>24), byte((*d).a>>16), byte((*d).a>>8), byte((*d).a))
}

func (d *digest64) Reset() {
	(*d).a = 0
	(*d).b = 0
}

func (d *digest64) Size() int { return size64 }

func (d *digest64) BlockSize() int { return size64 / 2 }

func (d *digest64) Sum64() uint64 { return uint64((*d).b)<<32 | uint64((*d).a) }
