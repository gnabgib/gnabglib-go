// Copyright 2022 gnabgib
// This Source Code Form is subject to the terms of the Mozilla Public License v2.0

package bcc

//[Block check character](https://en.wikipedia.org/wiki/Block_check_character)

import "github.com/gnabgib/gnablib-go/checksum"

type digest uint8

func New() checksum.Hash8 {
	d:=new(digest)
	d.Reset()
	return d
}

func (d *digest) update(p []byte) {
	dig:=byte(*d)
	for i:=0;i<len(p);i++ {
		dig^=p[i]
	}
	*d=digest(dig)
}

func (d *digest) Write(p []byte) (n int, err error) {
	d.update(p)
	return len(p), nil
}

func (d *digest) Sum(in []byte) []byte {
	s := byte(*d)
	return append(in,s)
}


func(d *digest) Reset() {
	*d=0
}

func (d *digest) Size() int {
	return 1
}

func (d *digest) BlockSize() int { return 1 }

func (d *digest) Sum8() uint8 { return byte(*d) }