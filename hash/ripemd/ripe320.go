// Copyright 2022 gnabgib
// This Source Code Form is subject to the terms of the Mozilla Public License v2.0

package ripemd

import (
	"hash"
	"math/bits"
)

const hashSize320u32 = 10

func hash320(ctx *ripeCtx) {
	a := ctx.state[0]
	b := ctx.state[1]
	c := ctx.state[2]
	d := ctx.state[3]
	e := ctx.state[4]
	aa := ctx.state[5]
	bb := ctx.state[6]
	cc := ctx.state[7]
	dd := ctx.state[8]
	ee := ctx.state[9]

	x := ctx.getX()

	j := 0
	round := 0
	var t uint32

	for ; j < 16; j++ {
		t = e + bits.RotateLeft32(a+f[round](b, c, d)+x[r[j]]+k[round], s[j])
		a, e, d, c, b = e, d, bits.RotateLeft32(c, 10), b, t
		t = ee +
			bits.RotateLeft32(aa+f[4-round](bb, cc, dd)+x[rr[j]]+kk[round], ss[j])
		aa, ee, dd, cc, bb = ee, dd, bits.RotateLeft32(cc, 10), bb, t
	}
	b, bb = bb, b

	round = 1
	for ; j < 32; j++ {
		t = e + bits.RotateLeft32(a+f[round](b, c, d)+x[r[j]]+k[round], s[j])
		a, e, d, c, b = e, d, bits.RotateLeft32(c, 10), b, t
		t = ee +
			bits.RotateLeft32(aa+f[4-round](bb, cc, dd)+x[rr[j]]+kk[round], ss[j])
		aa, ee, dd, cc, bb = ee, dd, bits.RotateLeft32(cc, 10), bb, t
	}
	d, dd = dd, d

	round = 2
	for ; j < 48; j++ {
		t = e + bits.RotateLeft32(a+f[round](b, c, d)+x[r[j]]+k[round], s[j])
		a, e, d, c, b = e, d, bits.RotateLeft32(c, 10), b, t
		t = ee +
			bits.RotateLeft32(aa+f[4-round](bb, cc, dd)+x[rr[j]]+kk[round], ss[j])
		aa, ee, dd, cc, bb = ee, dd, bits.RotateLeft32(cc, 10), bb, t
	}
	a, aa = aa, a

	round = 3
	for ; j < 64; j++ {
		t = e + bits.RotateLeft32(a+f[round](b, c, d)+x[r[j]]+k[round], s[j])
		a, e, d, c, b = e, d, bits.RotateLeft32(c, 10), b, t
		t = ee +
			bits.RotateLeft32(aa+f[4-round](bb, cc, dd)+x[rr[j]]+kk[round], ss[j])
		aa, ee, dd, cc, bb = ee, dd, bits.RotateLeft32(cc, 10), bb, t
	}
	c, cc = cc, c

	round = 4
	for ; j < 80; j++ {
		t = e + bits.RotateLeft32(a+f[round](b, c, d)+x[r[j]]+k[round], s[j])
		a, e, d, c, b = e, d, bits.RotateLeft32(c, 10), b, t
		t = ee +
			bits.RotateLeft32(aa+f[4-round](bb, cc, dd)+x[rr[j]]+kk[round], ss[j])
		aa, ee, dd, cc, bb = ee, dd, bits.RotateLeft32(cc, 10), bb, t
	}
	e, ee = ee, e

	ctx.state[0] += a
	ctx.state[1] += b
	ctx.state[2] += c
	ctx.state[3] += d
	ctx.state[4] += e
	ctx.state[5] += aa
	ctx.state[6] += bb
	ctx.state[7] += cc
	ctx.state[8] += dd
	ctx.state[9] += ee
	ctx.bPos = 0
}

// A new hash for computing RipeMd160
func New320() hash.Hash {
	c := &ripeCtx{
		hash:     hash320,
		stateLen: hashSize320u32}
	c.Reset()
	return c
}
