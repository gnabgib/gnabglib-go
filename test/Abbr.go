// Copyright 2022 gnabgib
// This Source Code Form is subject to the terms of the Mozilla Public License v2.0

package test

//Abbreviate a string if there's more than 20 characters
func Abbr(s string) string {
	n := len(s)
	if n > 20 {
		return s[0:19] + "…"
	}
	return s
}

