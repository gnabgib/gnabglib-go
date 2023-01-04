// Copyright 2022 gnabgib
// This Source Code Form is subject to the terms of the Mozilla Public License v2.0

package net

import (
	"net"
)

type iTreeNode interface {
	getValue() interface{}
	each(position uint32, mask byte, output func(position uint32, mask byte, value interface{}))
}

// All values below are represented
type all struct {
	value interface{}
}

func (n all) getValue() interface{} {
	return n.value
}
func (n all) each(position uint32, mask byte, output func(position uint32, mask byte, value interface{})) {
	output(position, 31-mask, n.value)
}

func newAll(value interface{}) all {
	return all{value}
}

// No values below are represented
type none struct{}

func (n none) getValue() interface{} {
	return nil
}
func (n none) each(position uint32, mask byte, output func(position uint32, mask byte, value interface{})) {
	//nop
}

// A node/fork that represents two paths
type treeNode struct {
	l, r iTreeNode
}

func newTreeNode() treeNode {
	r := treeNode{&none{}, &none{}}
	return r
}

func (n treeNode) getValue() interface{} {
	return nil
}
func (n treeNode) each(position uint32, mask byte, output func(position uint32, mask byte, value interface{})) {
	n.l.each(position, mask-1, output)
	n.r.each(position|(1<<mask), mask-1, output)
}
func add(parent iTreeNode, pos uint32, mask, end byte, value interface{}, merge func(a, b interface{}) interface{}) iTreeNode {
	//If at an all node, no further processing needed
	if _, isAll := parent.(all); isAll {
		return parent
	}
	//If we've descended as far as we need to, make it an all node and be done
	if mask <= end {
		return newAll(value)
	}

	//If this is a none node (ie not a TreeNode), switch (all has already been filtered)
	ret, isTree := parent.(treeNode)
	if !isTree {
		ret = newTreeNode()
	}

	//Now descend
	mask -= 1
	odd := (pos >> mask) & 1
	if odd == 1 {
		ret.r = add(ret.r, pos, mask, end, value, merge)
	} else {
		ret.l = add(ret.l, pos, mask, end, value, merge)
	}

	//If children are all, switch to all node
	lAll, lIsAll := ret.l.(all)
	rAll, rIsAll := ret.r.(all)
	if lIsAll && rIsAll {
		return newAll(merge(lAll.value, rAll.value))
	}

	//Since you can't switch to none-nodes, and you only switch from if there are
	// children.. we don't need to check for children=none
	return ret
}

type Tree struct {
	root  iTreeNode
	merge func(a, b interface{}) interface{}
}

func New(merge func(a, b interface{}) interface{}) *Tree {
	tree := &Tree{none{}, merge}
	return tree
}

// Add a single IP to the tree
func (t *Tree) AddIp(ip net.IP, value interface{}) {
	u := Ipv4ToUint(ip)
	t.root = add(t.root, u, 32, 0, value, t.merge)
}

// Add a range of IPs to the tree (inclusive)
func (t *Tree) AddRange(start, end net.IP, value interface{}) {
	uStart := Ipv4ToUint(start)
	uEnd := Ipv4ToUint(end)
	for uStart <= uEnd {
		m := (uStart - 1) & ^uStart
		for uStart+m > uEnd {
			m >>= 1
		}
		var bit byte
		bit = 0
		temp := m
		for temp != 0 {
			bit += 1
			temp >>= 1
		}
		t.root = add(t.root, uStart, 32, bit, value, t.merge)
		uStart += m + 1
	}
}

type CidrValue struct {
	Cidr  *net.IPNet
	Value interface{}
}

// Add a CIDR to the tree
func (t *Tree) AddCidr(cidr *net.IPNet, value interface{}) {
	u := Ipv4ToUint(cidr.IP)
	b := MaskToByte(cidr.Mask)
	t.root = add(t.root, u, 32, 32-b, value, t.merge)
}

// List all CIDR that describe the contents of the tree
func (t *Tree) ListCidr() []CidrValue {
	ret := make([]CidrValue, 0)
	t.root.each(0, 31, func(position uint32, bit byte, value interface{}) {
		ip := Ipv4FromUint(position)
		mask := MaskFromByte(bit)
		cidr := &net.IPNet{IP: ip, Mask: mask}
		ret = append(ret, CidrValue{cidr, value})
	})
	return ret
}
