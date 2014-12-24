// Copyright (C) 2014-15 Michail Kargakis
// This source code is subject to the terms
// of the MIT License
package prime

import (
	"fmt"
	"math/big"
	"sort"
)

// TODO: 100% testing

// Kind defines the kind of a Cunningham list.
type Kind uint8

const (
	NoKind Kind = iota
	// A chain of the 1st kind of length n is a sequence
	// of prime numbers (p1, ..., pn) such that for all
	// 1 ≤ i < n, pi+1 = 2pi + 1.
	FirstKind Kind
	// A chain of the 2nd kind of length n is a sequence
	// of prime numbers (p1, ..., pn) such that for all
	// 1 ≤ i < n, pi+1 = 2pi + 1.
	SecondKind Kind
)

// Candidate is a candidate number for being a prime.
type Candidate struct {
	value   *big.Int
	isPrime bool
}

// CandidateChain is a candidate chain for being a
// Cunningham chain.
type CandidateChain struct {
	actualChain []*Candidate
	kind        Kind
	checked     bool
}

func NewCandidateChain(chain []int) *CandidateChain {
	candidateCh := &CandidateChain{
		actualChain: make([]*Candidate, len(chain)),
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
			fmt.Println("The passed integer chain cannot be a Cunningham chain")
			return nil
		}
		candidateCh.actualChain[i].value = big.NewInt(int64(num))
	}
	return candidateCh
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
	// We need at least two numbers to form a chain.
	if len(chain) < 2 {
		return false
	}
	// A Sophie Germain prime is a prime p if (p-1)/2
	// is also a prime number.
	sophieGermain := chain.actualChain[0].value
	if !sophieGermain.ProbablyPrime(5) { // TODO: I have to work on this prime check
		return false
	}
	chain.actualChain[0].isPrime = true
	// Check what kind of a candidate Cunningham chain are we on.
	next := chain.actualChain[1].value
	if next == 2*sophieGermain+1 {
		chain.kind = FirstKind
	} else if next == 2*sophieGermain-1 {
		chain.kind = SecondKind
	} else {
		// This is not a Cunningham chain
		return false
	}
	sophieGermain = next
	// A safe prime is a prime p if 2p+1 is also a prime number.
	for i, safePrime := range chain.actualChain[2:] {
		sf := safePrime.value
		if !sf.ProbablyPrime(5) { // TODO: I have to work on this prime check
			chain.kind = NoKind
			return false
		}
		switch sf {
		case 2*sophieGermain + 1:
			// First kind of Cunningham chain if all this goes well
		case 2*sophieGermain - 1:
			// Second kind
		default:
			// No kind
			chain.kind = NoKind
			return false
		}
		chain.actualChain[i].isPrime = true
		sophieGermain = sf
	}
	// All the given chain has to be consisted of Sophie Germain and safe primes.
	// Cunningham subchains will be ignored for the time being.
	return true
}
