package test

func Abbr(s string) string {
	n := len(s)
	if n > 10 {
		return s[0:17] + "..."
	}
	return s
}

