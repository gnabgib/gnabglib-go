// Copyright 2022 gnabgib
// This Source Code Form is subject to the terms of the Mozilla Public License v2.0

package lrc

//[Longitudinal redundancy check](https://en.wikipedia.org/wiki/Longitudinal_redundancy_check)

import "github.com/gnabgib/gnablib-go/checksum"

type digest uint8

// A new Hash8 for computing the Longitudinal redundancy check checksum
func New() checksum.Hash8 {
	d := new(digest)
	d.Reset()
	return d
}

func (d *digest) update(p []byte) {
	dig:=uint8(*d)
	for i:=0;i<len(p);i++ {
		dig+=p[i]
	}
	*d=digest(^dig+1)
}

func (d *digest) Write(p []byte) (n int, err error) {
	d.update(p)
	return len(p), nil
}

func (d *digest) Sum(in []byte) []byte { return append(in, byte(*d)) }

func (d *digest) Reset() { *d = 0 }

func (d *digest) Size() int { return 1 }

func (d *digest) BlockSize() int { return 1 }

func (d *digest) Sum8() uint8 { return byte(*d) }
