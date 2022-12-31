// Copyright 2022 gnabgib
// This Source Code Form is subject to the terms of the Mozilla Public License v2.0

//go:build !tiny
// +build !tiny

package hex

// Decode a hex value using a 256byte lookup table.. for embedded (tiny) applications
// this might be too much memory space to use

const asciiDecode = "" +
	//0-31 NULL-Unit sep
	"~~~~~~~~~~~~~~~~" + "~~~~~~~~~~~~~~~~" +
	//32-63 Space-?
	"~~~~~~~~~~~~~~~~" + "\x00\x01\x02\x03\x04\x05\x06\x07\x08\x09~~~~~~" +
	//64-95 @-_
	"~\x0A\x0B\x0C\x0D\x0E\x0F~~~~~~~~~" + "~~~~~~~~~~~~~~~~" +
	//96-127 `-DEL
	"~\x0A\x0B\x0C\x0D\x0E\x0F~~~~~~~~~" + "~~~~~~~~~~~~~~~~" +
	//128+ (invalid ascii)
	"~~~~~~~~~~~~~~~~" + "~~~~~~~~~~~~~~~~" +
	"~~~~~~~~~~~~~~~~" + "~~~~~~~~~~~~~~~~" +
	"~~~~~~~~~~~~~~~~" + "~~~~~~~~~~~~~~~~" +
	"~~~~~~~~~~~~~~~~" + "~~~~~~~~~~~~~~~~"

//Convert a hex-char byte (0-9,a-f,A-F) into a byte (0-15, >15 invalid)
func decode(b byte) byte {
	return asciiDecode[b]
}
