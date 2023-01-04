// Copyright 2022 gnabgib
// This Source Code Form is subject to the terms of the Mozilla Public License v2.0

package ripemd

import (
	"encoding/binary"
	"hash"
	"math/bits"
)

const hashSize256u32 = 8

func hash256(ctx *ripeCtx) {
	a := ctx.state[0]
	b := ctx.state[1]
	c := ctx.state[2]
	d := ctx.state[3]
	aa := ctx.state[4]
	bb := ctx.state[5]
	cc := ctx.state[6]
	dd := ctx.state[7]

	x := ctx.getX()

	j := 0
	round := 0
	tk := binary.BigEndian.Uint32([]byte(k[round*4 : round*4+4]))
	tkk := binary.BigEndian.Uint32([]byte(kk128[round*4 : round*4+4]))
	var t uint32

	for ; j < 16; j++ {
		t = bits.RotateLeft32(a+f[round](b, c, d)+x[r[j]]+tk, int(s[j]))
		a, d, c, b = d, c, b, t
		t = bits.RotateLeft32(aa+f[3-round](bb, cc, dd)+x[rr[j]]+tkk, int(ss[j]))
		aa, dd, cc, bb = dd, cc, bb, t
	}
	a, aa = aa, a

	round = 1
	tk = binary.BigEndian.Uint32([]byte(k[round*4 : round*4+4]))
	tkk = binary.BigEndian.Uint32([]byte(kk128[round*4 : round*4+4]))
	for ; j < 32; j++ {
		t = bits.RotateLeft32(a+f[round](b, c, d)+x[r[j]]+tk, int(s[j]))
		a, d, c, b = d, c, b, t
		t = bits.RotateLeft32(aa+f[3-round](bb, cc, dd)+x[rr[j]]+tkk, int(ss[j]))
		aa, dd, cc, bb = dd, cc, bb, t
	}
	b, bb = bb, b

	round = 2
	tk = binary.BigEndian.Uint32([]byte(k[round*4 : round*4+4]))
	tkk = binary.BigEndian.Uint32([]byte(kk128[round*4 : round*4+4]))
	for ; j < 48; j++ {
		t = bits.RotateLeft32(a+f[round](b, c, d)+x[r[j]]+tk, int(s[j]))
		a, d, c, b = d, c, b, t
		t = bits.RotateLeft32(aa+f[3-round](bb, cc, dd)+x[rr[j]]+tkk, int(ss[j]))
		aa, dd, cc, bb = dd, cc, bb, t
	}
	c, cc = cc, c

	round = 3
	tk = binary.BigEndian.Uint32([]byte(k[round*4 : round*4+4]))
	tkk = binary.BigEndian.Uint32([]byte(kk128[round*4 : round*4+4]))
	for ; j < 64; j++ {
		t = bits.RotateLeft32(a+f[round](b, c, d)+x[r[j]]+tk, int(s[j]))
		a, d, c, b = d, c, b, t
		t = bits.RotateLeft32(aa+f[3-round](bb, cc, dd)+x[rr[j]]+tkk, int(ss[j]))
		aa, dd, cc, bb = dd, cc, bb, t
	}
	d, dd = dd, d

	ctx.state[0] += a
	ctx.state[1] += b
	ctx.state[2] += c
	ctx.state[3] += d
	ctx.state[4] += aa
	ctx.state[5] += bb
	ctx.state[6] += cc
	ctx.state[7] += dd
	ctx.bPos = 0
}

// A new hash for computing RipeMd160
func New256() hash.Hash {
	c := &ripeCtx{
		hash:     hash256,
		stateLen: hashSize256u32}
	c.Reset()
	return c
}
