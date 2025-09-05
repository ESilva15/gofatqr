package gofatqr

import (
	"fmt"
	"testing"
)

func TestIsValidNIF(t *testing.T) {
	testCases := []struct {
		name   string
		expect bool
		input  string
	}{
		{
			name:   "Normal NIF",
			expect: true,
			input:  "198435678",
		},
		{
			name:   "String too short",
			expect: false,
			input:  "19845678",
		},
		{
			name:   "String too long",
			expect: false,
			input:  "1985678654",
		},
		{
			name:   "Not digits in the string",
			expect: false,
			input:  "1AC43567A",
		},
	}

	for _, tc := range testCases {
		fmt.Printf("Test: %v\n", tc.name)
		res := isValidNIF(tc.input)
		if res != tc.expect {
			t.Errorf("Expected %t for input %s, got %t", tc.expect, tc.input, res)
		}
	}
}
