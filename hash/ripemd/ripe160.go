// Copyright 2022 gnabgib
// This Source Code Form is subject to the terms of the Mozilla Public License v2.0

package ripemd

import (
	"hash"
	"math/bits"
)

const (
	roundCount160 = 5
	hashSize160u32 = 5
)

func hash160(ctx *ripeCtx) {
	a := ctx.state[0]
	b := ctx.state[1]
	c := ctx.state[2]
	d := ctx.state[3]
	e := ctx.state[4]
	aa := ctx.state[0]
	bb := ctx.state[1]
	cc := ctx.state[2]
	dd := ctx.state[3]
	ee := ctx.state[4]

	x := ctx.getX()

	for round := 0; round < roundCount160; round++ {
		n := (round + 1) * 16
		for j := round * 16; j < n; j++ {
			t := e + bits.RotateLeft32(a+f[round](b, c, d)+x[r[j]]+k[round], s[j])
			a, e, d, c, b = e, d, bits.RotateLeft32(c, 10), b, t
			t = ee +
				bits.RotateLeft32(aa+f[4-round](bb, cc, dd)+x[rr[j]]+kk[round], ss[j])
			aa, ee, dd, cc, bb = ee, dd, bits.RotateLeft32(cc, 10), bb, t
		}
	}

	t := ctx.state[1] + c + dd
	ctx.state[1] = ctx.state[2] + d + ee
	ctx.state[2] = ctx.state[3] + e + aa
	ctx.state[3] = ctx.state[4] + a + bb
	ctx.state[4] = ctx.state[0] + b + cc
	ctx.state[0] = t
	ctx.bPos = 0
}

// A new hash for computing RipeMd160
func New160() hash.Hash {
	c := &ripeCtx{
		hash:     hash160,
		stateLen: hashSize160u32}
	c.Reset()
	return c
}
