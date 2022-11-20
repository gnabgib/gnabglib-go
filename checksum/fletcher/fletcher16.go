// Copyright 2022 gnabgib
// This Source Code Form is subject to the terms of the Mozilla Public License v2.0

package fletcher

//http://www.zlib.net/maxino06_fletcher-adler.pdf -> Lower cpu and Adler and mostly more effective (in their tests)
//https://datatracker.ietf.org/doc/html/rfc1146 (Appendix I)

import (
	"github.com/gnabgib/gnablib-go/checksum"
)

const size16=2

// Fletcher 16 -- --
type digest16 struct {
	a, b uint8
}

// A new Hash16 for computing the Fletcher 16 checksum
func New16() checksum.Hash16 {
	d := new(digest16)
	d.Reset()
	return d
}

func (d *digest16) update(p []byte) {
	//Using the Anastase Nakassis optimization
	//c1 overflow solve to 5802 (from eq: n > 0 and n * (n+1) / 2 * (2^8-1) < (2^32-1))
	const mod = 0xff
	const c1Overflow = 5803
	var c0, c1 uint32
	c0 = uint32((*d).a)
	c1 = uint32((*d).b)

	for i := 0; i < len(p); {
		len := len(p)
		if len > c1Overflow {
			len = c1Overflow
		}
		for j := 0; j < len; j++ {
			//We only use j for the batch-count (incrementing i each loop)
			c0 += uint32(p[i])
			c1 += c0
			i++
		}
		c0 %= mod
		c1 %= mod
	}
	(*d).a = byte(c0)
	(*d).b = byte(c1)
}

func (d *digest16) Write(p []byte) (n int, err error) {
	d.update(p)
	return len(p), nil
}

func (d *digest16) Sum(in []byte) []byte {
	return append(in, (*d).b, (*d).a)
}

func (d *digest16) Reset() {
	(*d).a = 0
	(*d).b = 0
}

func (d *digest16) Size() int { return size16 }

func (d *digest16) BlockSize() int { return size16/2 }

func (d *digest16) Sum16() uint16 { return uint16((*d).b)<<8 | uint16((*d).a) }