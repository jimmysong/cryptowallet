// Copyright (C) 2014-15 Michail Kargakis
// This source code is subject to the terms
// of the MIT License
package prime_test

import (
	"log"
	"testing"

	"github.com/kargakis/prime"
)

func TestIsCunninghamChain(t *testing.T) {
	chainTests := []struct {
		name     string
		chain    []int
		expected bool
	}{
		{
			name:     "1st complete 1CC",
			chain:    []int{2, 5, 11, 23, 47},
			expected: true,
		},
		{
			name:     "2nd complete 1CC",
			chain:    []int{89, 179, 359, 719, 1439, 2879},
			expected: true,
		},
		{
			name:     "1st complete 2CC",
			chain:    []int{151, 301, 601, 1201},
			expected: true,
		},
		{
			name:     "1st failed 2CC",
			chain:    []int{19, 37, 73, 145}, // 145 is not a prime
			expected: false,
		},
		{
			name:     "1st failed 1CC",
			chain:    []int{41, 83, 167, 335}, // 335 is not a prime
			expected: false,
		},
	}

	for i, test := range chainTests {
		candidate, err := prime.NewCandidateChain(test.chain)
		if err != nil {
			log.Println(err)
		}
		cunningham := candidate.IsCunninghamChain()
		if !(cunningham == chainTests[i].expected) {
			t.Errorf("IsCunninghamChain #%d (%s) wrong result\n"+
				"got: %v\nwant: %v", i, test.name, cunningham,
				test.expected)
		}
	}
}
