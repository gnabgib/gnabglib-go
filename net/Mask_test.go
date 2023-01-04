package net

import (
	"bytes"
	"net"
	"testing"
)

var maskTests = []struct {
	mask []byte
	b    byte
}{
	{[]byte{0x00, 0, 0, 0}, 0},
	{[]byte{0x80, 0, 0, 0}, 1},
	{[]byte{0xC0, 0, 0, 0}, 2},
	{[]byte{0xE0, 0, 0, 0}, 3},
	{[]byte{0xF0, 0, 0, 0}, 4},
	{[]byte{0xF8, 0, 0, 0}, 5},
	{[]byte{0xFC, 0, 0, 0}, 6},
	{[]byte{0xFE, 0, 0, 0}, 7},
	{[]byte{0xFF, 0, 0, 0}, 8},
	{[]byte{255, 0x80, 0, 0}, 9},
	{[]byte{255, 0xC0, 0, 0}, 10},
	{[]byte{255, 0xE0, 0, 0}, 11},
	{[]byte{255, 0xF0, 0, 0}, 12},
	{[]byte{255, 0xF8, 0, 0}, 13},
	{[]byte{255, 0xFC, 0, 0}, 14},
	{[]byte{255, 0xFE, 0, 0}, 15},
	{[]byte{255, 0xFF, 0, 0}, 16},
	{[]byte{255, 255, 0x80, 0}, 17},
	{[]byte{255, 255, 0xC0, 0}, 18},
	{[]byte{255, 255, 0xE0, 0}, 19},
	{[]byte{255, 255, 0xF0, 0}, 20},
	{[]byte{255, 255, 0xF8, 0}, 21},
	{[]byte{255, 255, 0xFC, 0}, 22},
	{[]byte{255, 255, 0xFE, 0}, 23},
	{[]byte{255, 255, 0xFF, 0}, 24},
	{[]byte{255, 255, 255, 0x80}, 25},
	{[]byte{255, 255, 255, 0xC0}, 26},
	{[]byte{255, 255, 255, 0xE0}, 27},
	{[]byte{255, 255, 255, 0xF0}, 28},
	{[]byte{255, 255, 255, 0xF8}, 29},
	{[]byte{255, 255, 255, 0xFC}, 30},
	{[]byte{255, 255, 255, 0xFE}, 31},
	{[]byte{255, 255, 255, 0xFF}, 32},
}

func TestToByte(t *testing.T) {
	for _, rec := range maskTests {
		mask := net.IPMask(rec.mask)
		found := MaskToByte(mask)
		expected := rec.b

		if found != expected {
			t.Fatalf("Expecting byte form of mask %v to be %v, got %v", mask, expected, found)
		}
	}
}

func TestFromByte(t *testing.T) {
	for _, rec := range maskTests {

		expected := net.IPMask(rec.mask)
		found := MaskFromByte(rec.b)

		if !bytes.Equal(found, expected) {
			t.Fatalf("Expecting mask from byte %v to be %v, got %v", rec.b, expected, found)
		}
	}
}
