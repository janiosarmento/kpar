package converter

import (
	"fmt"
	"os"
	"strings"
	"unicode"

	"golang.org/x/text/unicode/norm"
)

// letterReplacements covers letters that do not decompose under NFD and would
// otherwise be lost as separators.
var letterReplacements = map[rune]string{
	'ß': "ss", 'æ': "ae", 'œ': "oe", 'ø': "o", 'đ': "d", 'ł': "l", 'þ': "th", 'ð': "d",
}

// permalink turns a file name (without extension) into a URL-friendly slug:
// lowercase ASCII letters and digits, words joined by single hyphens, no accents.
// Spaces, underscores, dots and any other symbol act as word separators.
func permalink(name string) string {
	var b strings.Builder
	pendingSep := false

	// write appends s, first emitting a hyphen if a separator is pending.
	write := func(s string) {
		if pendingSep && b.Len() > 0 {
			b.WriteByte('-')
		}
		pendingSep = false
		b.WriteString(s)
	}

	for _, r := range norm.NFD.String(name) {
		if unicode.Is(unicode.Mn, r) {
			continue // combining accent mark: drop it, keep the base letter
		}
		r = unicode.ToLower(r)

		if repl, ok := letterReplacements[r]; ok {
			write(repl)
		} else if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			write(string(r))
		} else {
			pendingSep = true
		}
	}

	if b.Len() == 0 {
		return "image"
	}
	return b.String()
}

// uniquePath returns base.ext, or base-2.ext, base-3.ext, ... if that path is
// already taken, so an existing file is never overwritten.
func uniquePath(base, ext string) string {
	candidate := base + "." + ext
	for n := 2; ; n++ {
		if _, err := os.Lstat(candidate); err != nil {
			return candidate
		}
		candidate = fmt.Sprintf("%s-%d.%s", base, n, ext)
	}
}
