package main

import (
	"fmt"
	"math/big"
)

func main() {
	private_key, _ := new(big.Int).SetString("0xb5b87bc90b94db1a435a109b4c8c560e4ba0ad9966b5b9cc738ca66295168e0e", 0)
	pub_key := PublickKey(private_key)

	fmt.Println(pub_key)
}
