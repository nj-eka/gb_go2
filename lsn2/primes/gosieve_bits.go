// Package primes implements simple functions for searching prime numbers and some comments to test godoc :)
package primes

import (
	"github.com/yourbasic/bit"
	"math"
)

// SieveWithBits - the fastest algo due to bits structure
// BenchmarkSieveWithBits-4           23870            121063 ns/op           13440 B/op          2 allocs/op
func SieveWithBits(top int) *bit.Set {
	sieve := bit.New().AddRange(2, top + 1)
	sqrtN := int(math.Sqrt(float64(top)))
	for p := 2; p <= sqrtN; p = sieve.Next(p) {
		for k := p * p; k <= top; k += p {
			sieve.Delete(k)
		}
	}
	return sieve
}
