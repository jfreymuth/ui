package text

import (
	"unicode"
	"unicode/utf8"
)

func FindWord(t string, c int) (int, int) {
	s := c
	r, size := utf8.DecodeRuneInString(t[c:])
	for unicode.IsLetter(r) || unicode.IsNumber(r) {
		c += size
		r, size = utf8.DecodeRuneInString(t[c:])
	}
	r, size = utf8.DecodeLastRuneInString(t[:s])
	for unicode.IsLetter(r) || unicode.IsNumber(r) {
		s -= size
		r, size = utf8.DecodeLastRuneInString(t[:s])
	}
	if c == s && c < len(t) {
		c++
	}
	return s, c
}

func NextWord(t string, c int) int {
	if c == len(t) {
		return c
	}
	r, size := utf8.DecodeRuneInString(t[c:])
	for unicode.IsSpace(r) {
		c += size
		r, size = utf8.DecodeRuneInString(t[c:])
	}
	if unicode.IsLetter(r) || unicode.IsNumber(r) {
		for unicode.IsLetter(r) || unicode.IsNumber(r) {
			c += size
			r, size = utf8.DecodeRuneInString(t[c:])
		}
		return c
	} else {
		for r != utf8.RuneError && !unicode.IsLetter(r) && !unicode.IsNumber(r) && !unicode.IsSpace(r) {
			c += size
			r, size = utf8.DecodeRuneInString(t[c:])
		}
		return c
	}
}

func PreviousWord(t string, c int) int {
	if c == 0 {
		return c
	}
	r, size := utf8.DecodeLastRuneInString(t[:c])
	for unicode.IsSpace(r) {
		c -= size
		r, size = utf8.DecodeLastRuneInString(t[:c])
	}
	if unicode.IsLetter(r) || unicode.IsNumber(r) {
		for unicode.IsLetter(r) || unicode.IsNumber(r) {
			c -= size
			r, size = utf8.DecodeLastRuneInString(t[:c])
		}
		return c
	} else {
		for r != utf8.RuneError && !unicode.IsLetter(r) && !unicode.IsNumber(r) && !unicode.IsSpace(r) {
			c -= size
			r, size = utf8.DecodeLastRuneInString(t[:c])
		}
		return c
	}
}
