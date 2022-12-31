// Copyright 2022 gnabgib
// This Source Code Form is subject to the terms of the Mozilla Public License v2.0

package ripemd

import (
	"hash"
	"math/bits"
)

const (
	roundCount128  = 4
	hashSize128u32 = 4
)

func hash128(ctx *ripeCtx) {
	a := ctx.state[0]
	b := ctx.state[1]
	c := ctx.state[2]
	d := ctx.state[3]
	aa := ctx.state[0]
	bb := ctx.state[1]
	cc := ctx.state[2]
	dd := ctx.state[3]

	x := ctx.getX()

	for round := 0; round < roundCount128; round++ {
		n := (round + 1) * 16
		for j := round * 16; j < n; j++ {
			t := bits.RotateLeft32(a+f[round](b, c, d)+x[r[j]]+k[round], int(s[j]))
			a, d, c, b = d, c, b, t
			t = bits.RotateLeft32(aa+f[3-round](bb, cc, dd)+x[rr[j]]+kk128[round], int(ss[j]))
			aa, dd, cc, bb = dd, cc, bb, t
		}
	}
	t := ctx.state[1] + c + dd
	ctx.state[1] = ctx.state[2] + d + aa
	ctx.state[2] = ctx.state[3] + a + bb
	ctx.state[3] = ctx.state[0] + b + cc
	ctx.state[0] = t
	ctx.bPos = 0
}

// A new hash for computing RipeMd160
func New128() hash.Hash {
	c := &ripeCtx{
		hash:     hash128,
		stateLen: hashSize128u32}
	c.Reset()
	return c
}
