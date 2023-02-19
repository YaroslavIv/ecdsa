package main

import (
	"math/big"
)

func inverse(a, p *big.Int) *big.Int {
	// a*x + p*y = gcd(a,p)
	// (a*x + p*y) mod b = gcd(a,p) mod p
	// a*x mod p = 1 mod p
	// x = a**(-1) mod p
	return new(big.Int).ModInverse(a, p)
}
