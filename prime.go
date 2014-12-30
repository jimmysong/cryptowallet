// Copyright (C) 2014-15 Michail Kargakis
// This source code is subject to the terms
// of the MIT License
package prime

import (
	"errors"
	"fmt"
	"math/big"
	"sort"
)

// TODO
// * AKS primality test (https://en.wikipedia.org/wiki/AKS_primality_test)
// * Prime gap (https://en.wikipedia.org/wiki/Prime_gap)
// * Bi-twin chains (https://en.wikipedia.org/wiki/Bi-twin_chain)
// * Twin primes (https://en.wikipedia.org/wiki/Twin_prime)
// * 100% testing

// Kind defines the kind of a Cunningham list.
type Kind uint8

const (
	// No Kind defines that a given candidate Cunningham chain
	// is not actually any kind of Cunningham chain.
	NoKind Kind = iota
	// First Kind chain of length n is a sequence
	// of prime numbers (p1, ..., pn) such that for all
	// 1 ≤ i < n, pi+1 = 2pi + 1.
	FirstKind
	// Second Kind chain of length n is a sequence
	// of prime numbers (p1, ..., pn) such that for all
	// 1 ≤ i < n, pi+1 = 2pi + 1.
	SecondKind
)

// CandidateChain is a candidate chain for being a
// Cunningham chain.
type CandidateChain struct {
	actualChain []*big.Int
	kind        Kind
	checked     bool
}

// NewCandidateChain turns an integer chain into a CandidateChain.
func NewCandidateChain(chain []int) (*CandidateChain, error) {
	// We need at least two numbers to form a chain.
	if len(chain) < 2 {
		return nil, errors.New("The passed integer chain does not contain enough numbers.")
	}
	candidateCh := &CandidateChain{
		actualChain: make([]*big.Int, len(chain)),
		// The rest are initialized to their zero values
		// (false and NoKind respectively) thanks to Go.
	}
	// Sort chain if it's not already sorted.
	if !sort.IntsAreSorted(chain) {
		sort.Ints(chain)
	}
	// Copy the passed integers chain into our CandidateChain struct.
	for i, num := range chain {
		if num < 1 {
			// This number cannot be a prime number
			return nil, errors.New("The passed integer chain contains non-prime numbers.")
		}
		candidateCh.actualChain[i] = big.NewInt(int64(num))
	}
	return candidateCh, nil
}

// KindOf returns the kind of the CandidateChain.
func (c *CandidateChain) KindOf() Kind {
	if c.checked == false {
		fmt.Println("The kind of this chain is unknown")
	}
	// TODO: Probably add a mutex here
	return c.kind
}

// IsCunninghamChain returns wheter a CandidateChain is a Cunningham
// chain or not.
func (c *CandidateChain) IsCunninghamChain() bool {
	return checkForCunningham(c)
}

// checkForCunningham checks whether the passed chain
// of numbers is a Cunningham chain, a certain sequence
// of prime numbers.
func checkForCunningham(chain *CandidateChain) bool {
	chain.checked = true
	// A Sophie Germain prime is a prime p if (p-1)/2
	// is also a prime number.
	sophieGermain := chain.actualChain[0]
	if !sophieGermain.ProbablyPrime(1000) { // TODO: I have to work on this prime check
		fmt.Printf("This number is not a prime: %d\n", sophieGermain)
		return false
	}
	// A safe prime is a prime p if 2p+1 is also a prime number.
	safePrime := chain.actualChain[1]
	if !safePrime.ProbablyPrime(1000) { // TODO: I have to work on this prime check
		fmt.Printf("This number is not a prime: %d\n", safePrime)
		return false
	}
	// Check what kind of a candidate Cunningham chain are we on.
	if safePrime.Cmp(safePrime1CC(sophieGermain)) == 0 {
		// First kind of Cunningham chain if all this goes well
		chain.kind = FirstKind
	} else if safePrime.Cmp(safePrime2CC(sophieGermain)) == 0 {
		// Second kind
		chain.kind = SecondKind
	} else {
		// No kind, this is not a Cunningham chain
		return false
	}
	sophieGermain = safePrime
	for _, safePrime := range chain.actualChain[2:] {
		if !safePrime.ProbablyPrime(1000) { // TODO: I have to work on this prime check
			fmt.Printf("This number is not a prime: %d\n", safePrime)
			chain.kind = NoKind
			return false
		}
		// Continue confirming that we are on the appropriate kind of chain.
		switch chain.kind {
		case FirstKind:
			if !(safePrime.Cmp(safePrime1CC(sophieGermain)) == 0) {
				chain.kind = NoKind
				return false
			}
		case SecondKind:
			if !(safePrime.Cmp(safePrime2CC(sophieGermain)) == 0) {
				chain.kind = NoKind
				return false
			}
		}
		sophieGermain = safePrime
	}
	// All the given chain has to be consisted of Sophie Germain and safe primes.
	// Cunningham subchains will be ignored for the time being.
	return true
}

// safePrime1CC returns the next safe prime number of the 1st kind
func safePrime1CC(sophieGermain *big.Int) *big.Int {
	sf := big.NewInt(0)
	return sf.Add(sophieGermain, sophieGermain).Add(sf, big.NewInt(1))
}

// safePrime2CC returns the next safe prime number of the 2nd kind
func safePrime2CC(sophieGermain *big.Int) *big.Int {
	sf := big.NewInt(0)
	return sf.Add(sophieGermain, sophieGermain).Sub(sf, big.NewInt(1))
}
