package main

import "math/big"

func Secp256k1() *curve {
	secp256k1 := curve{}

	secp256k1.P, _ = new(big.Int).SetString("0xFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEFFFFFC2F", 0)
	secp256k1.N, _ = new(big.Int).SetString("0xFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEBAAEDCE6AF48A03BBFD25E8CD0364141", 0)
	secp256k1.B, _ = new(big.Int).SetString("0x0000000000000000000000000000000000000000000000000000000000000007", 0)
	secp256k1.Gx, _ = new(big.Int).SetString("0x79BE667EF9DCBBAC55A06295CE870B07029BFCDB2DCE28D959F2815B16F81798", 0)
	secp256k1.Gy, _ = new(big.Int).SetString("0x483ADA7726A3C4655DA4FBFC0E1108A8FD17B448A68554199C47D08FFB10D4B8", 0)
	secp256k1.BitSize = 256

	return &secp256k1
}

func PublickKey(private_key *big.Int) string {
	p := Secp256k1().Mul(private_key)

	x := p.X.Text(16)
	for len(x) != 64 {
		x = "0" + x
	}
	y := p.Y.Text(16)
	for len(y) != 64 {
		y = "0" + y
	}

	return "04" + x + y
}
