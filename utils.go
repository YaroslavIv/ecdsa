package main

import (
	"math/big"

	"golang.org/x/crypto/sha3"
)

func hash(message string) *big.Int {
	keccak256 := sha3.NewLegacyKeccak256()
	keccak256.Write([]byte(message))
	return new(big.Int).SetBytes(keccak256.Sum(nil))
}
