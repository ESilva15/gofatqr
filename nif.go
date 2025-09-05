package gofatqr

import (
	"unicode"
)

// TODO
// Add a way to check with an API or whatever like nif.pt
// Also need to be able to indentify other countries NIFs I guess
// some validation at the end

func isValidNIF(n string) bool {
	if len(n) != 9 {
		return false
	}

	for _, c := range n {
		if !unicode.IsDigit(c) {
			return false
		}
	}

	return true
}
