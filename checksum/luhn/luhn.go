// Copyright 2022 gnabgib
// This Source Code Form is subject to the terms of the Mozilla Public License v2.0

package luhn

// https://www.dcode.fr/luhn-algorithm


func Checksum(i uint64) uint8 {
	var ret= uint64(0)
	var mul= uint64(2)
	for i>0 {
		v := i%10
		i=(i-v)/10
		v*=mul
        ret+=(v%10)+(v/10)
        mul=1+mul%2
	}
	return uint8((10-ret%10)%10)
}

// import "github.com/gnabgib/gnablib-go/checksum"

// type digest uint8

// // A new Hash8 for computing the Longitudinal redundancy check checksum
// func New() checksum.Hash8 {
// 	d := new(digest)
// 	d.Reset()
// 	return d
// }

// func (d *digest) update(p []byte) {
// 	dig:=uint8(*d)
// 	for i:=0;i<len(p);i++ {
// 		dig+=p[i]
// 	}
// 	*d=digest(^dig+1)
// }

// func (d *digest) Write(p []byte) (n int, err error) {
// 	d.update(p)
// 	return len(p), nil
// }

// func (d *digest) Sum(in []byte) []byte { return append(in, byte(*d)) }

// func (d *digest) Reset() { *d = 0 }

// func (d *digest) Size() int { return 1 }

// func (d *digest) BlockSize() int { return 1 }

// func (d *digest) Sum8() uint8 { return byte(*d) }
