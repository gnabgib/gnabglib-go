// Copyright 2023 gnabgib
// This Source Code Form is subject to the terms of the Mozilla Public License v2.0

package test

import (
	"strings"
	"testing"
)

// Compares two strings, creates errors if they differ in length, and shows
// the values side by side with style markers to show character mismatch.
// Note this does no UTF8 folding etc, the strings are to be byte-identical
func StringMatch(t *testing.T, expect, found string) {
	StringMatchTitle(t, "", expect, found)
}

func StringMatchTitle(t *testing.T, title, expect, found string) {
	//Simple case.. the byte strings are identical (yay)
	if found == expect {
		return
	}
	if len(found) != len(expect) {
		t.Errorf("Expected %d chars, got %d", len(expect), len(found))
	}

	eRunes := strings.Split(expect, "")
	fRunes := strings.Split(found, "")

	nRunes := len(eRunes)
	nBytes := len(expect)
	if len(found) > nBytes {
		nBytes = len(found)
		nRunes = len(fRunes)
	}

	//Make a marker line, only needs one char/rune (some runes can be >1 byte)
	// and one extra to show the end as being the issue
	marker := make([]byte, nRunes+1)

	first := findFirstMismatch(eRunes, fRunes, marker)
	last := findLastMismatch(eRunes, fRunes, marker, first)
	fullTitle:=""
	if len(title)>0 {
		fullTitle=title+": "
	}
	t.Errorf("%sExpect, found:\n%s␃\n%s␃\n%s",
		fullTitle,
		renderText(expect, nBytes-len(expect), first, last),
		renderText(found, nBytes-len(found), first, last),
		marker)
}

// The first differing rune between the two arrays
func findFirstMismatch(eRunes, fRunes []string, marker []byte) int {
	//Note the runes come in as strings, but we know they're just runes
	// due to the rule of Split on empty delimiter
	i := 0
	for ; i < len(eRunes); i++ {
		if i >= len(fRunes) || eRunes[i] != fRunes[i] {
			break
		}
		switch {
		case eRunes[i] == "\x7f":
			//Duplicate del in marker
			marker[i] = '\x7f'
		case int(eRunes[i][0]) < 32:
			//If the char is control (other than above), copy it to output
			marker[i] = byte(eRunes[i][0])
		default:
			marker[i] = '-'
		}
	}
	marker[i] = '^'
	return i
}

// The last differing rune between two arrays, returned as negative offset from the end
func findLastMismatch(eRunes, fRunes []string, marker []byte, first int) int {
	markerDelta := len(marker) - len(eRunes) - 1
	e := len(eRunes) - 1
	for f := len(fRunes) - 1; e > first && f > first; {
		//Marker is 1 longer than eRunes
		if eRunes[e] != fRunes[f] {
			break
		}
		marker[e+markerDelta] = '-'
		e--
		f--
	}
	//Fill in spaces between markers
	for i := e - 1 + markerDelta; i > first; i-- {
		marker[i] = ' '
	}
	//Mark the mismatch
	marker[e+markerDelta] = '^'
	return e - len(eRunes) + 1
}

// Render a text line, which is equal to length of text + delta (which
// is >0 iff found is longer and equal to the length-diff).  Up until firstDiff
// the line matches expect then there are different/missing char indicators
func renderText(t string, blankChars, markStart, markEndOffset int) string {
	//Return length is at least expected length
	n := len(t)
	//Add 4 characters (each) for start/stop style
	n += 4 + 4
	//Add 2 bytes for each, because the missing char maker (·) is 2 bytes
	n += blankChars * 2

	//Convert offset into position
	markEnd := len(t) + markEndOffset

	ret := make([]byte, n)
	//Copy in starting matching text
	copy(ret, t[0:markStart])
	//Add start change marker (invert fg/bg colours)
	copy(ret[markStart:], "\x1b[7m")
	//Add the diff text
	copy(ret[markStart+4:], t[markStart:markEnd])
	//Add any blanks
	for i := 0; i < blankChars*2; i += 2 {
		ret[markEnd+4+i] = '\xc2'
		ret[markEnd+4+i+1] = '\xb7'
	}
	//Add end change marker
	copy(ret[markEnd+4+blankChars*2:], "\x1b[0m")
	//Add ending content
	copy(ret[markEnd+8+blankChars*2:], t[markEnd:])

	return string(ret)
}
