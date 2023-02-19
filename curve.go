package main

import (
	"math/big"
)

type curve struct {
	P       *big.Int
	N       *big.Int
	B       *big.Int
	Gx, Gy  *big.Int
	BitSize int
}

func (c *curve) Double(p *Point) *Point {
	// k = (3*x_p**2 + a) * (2*y_p)**(-1) mod p
	x_p_2_ := new(big.Int).Mul(p.X, p.X)
	x_p_2 := new(big.Int).Mul(x_p_2_, big.NewInt(3))
	k_1_ := new(big.Int).Add(x_p_2, big.NewInt(0))
	k_2_ := new(big.Int).Mul(p.Y, big.NewInt(2))
	k := new(big.Int).Mul(k_1_, inverse(k_2_, c.P))

	// x_r = (k**2 - x_p - x_p) mod p
	k_2 := new(big.Int).Mul(k, k)
	x_r_ := new(big.Int).Sub(k_2, p.X)
	x_r := new(big.Int).Sub(x_r_, p.X)
	x_r_mod := new(big.Int).Mod(x_r, c.P)

	// y_r = (k*(x_p - x_r) - y_p) mod p
	x_d := new(big.Int).Sub(p.X, x_r)
	y_r_ := new(big.Int).Mul(k, x_d)
	y_r := new(big.Int).Sub(y_r_, p.Y)
	y_r_mod := new(big.Int).Mod(y_r, c.P)

	return NewPoint(x_r_mod, y_r_mod)
}

func (c *curve) add(p, q *Point) *Point {
	// k = (p_y - q_y) * (p_x - p_y)**(-1) mod p
	k_1_ := new(big.Int).Sub(p.Y, q.Y)
	k_2_ := new(big.Int).Sub(p.X, q.X)
	k := new(big.Int).Mul(k_1_, inverse(k_2_, c.P))

	// x_r = (k**2 - x_p - x_q) mod p
	k_2 := new(big.Int).Mul(k, k)
	x_r_ := new(big.Int).Sub(k_2, p.X)
	x_r := new(big.Int).Sub(x_r_, q.X)
	x_r_mod := new(big.Int).Mod(x_r, c.P)

	// y_r = (k*(x_p - x_r) - y_p) mod p
	x_d := new(big.Int).Sub(p.X, x_r)
	y_r_ := new(big.Int).Mul(k, x_d)
	y_r := new(big.Int).Sub(y_r_, p.Y)
	y_r_mod := new(big.Int).Mod(y_r, c.P)

	return NewPoint(x_r_mod, y_r_mod)
}

func (c *curve) Add(p, q *Point) *Point {
	if p.Zero() {
		return q
	} else if q.Zero() {
		return p
	} else if p.Opposite(q) {
		return DefaulttPoint()
	} else if p.Identical(q) {
		return c.Double(p)
	} else {
		return c.add(p, q)
	}
}

func (c *curve) Mul(n *big.Int) *Point {
	out := DefaulttPoint()
	iter := NewPoint(c.Gx, c.Gy)
	for n.Sign() > 0 {
		if new(big.Int).And(n, big.NewInt(1)).Sign() != 0 {
			out = c.Add(out, iter)
		}
		iter = c.Double(iter)
		n.Rsh(n, 1)
	}

	return out
}
