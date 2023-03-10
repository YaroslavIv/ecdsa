package main

import (
	"crypto/rand"
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
	return c.mul(n, Point{
		X: c.Gx,
		Y: c.Gy,
	})
}

func (c *curve) mul(n *big.Int, init_point Point) *Point {
	N := new(big.Int).Set(n)
	out := DefaulttPoint()
	iter := NewPoint(init_point.X, init_point.Y)
	for N.Sign() > 0 {
		if new(big.Int).And(N, big.NewInt(1)).Sign() != 0 {
			out = c.Add(out, iter)
		}
		iter = c.Double(iter)
		N.Rsh(N, 1)
	}

	return out
}

type Signatrue struct {
	R *big.Int
	S *big.Int
}

func (c *curve) sign(private_key *big.Int, message string) Signatrue {
	mes := hash(message)

	for {

		k, err := rand.Int(rand.Reader, c.N)
		if err != nil {
			continue
		}

		// R = k*G mod n
		point := c.Mul(k)
		r := new(big.Int).Mod(point.X, c.N)

		// S = ((mes + r * private_key) * k**(-1)) mod n
		x1 := new(big.Int).Mul(r, private_key)
		x2 := new(big.Int).Add(mes, x1)
		x3 := new(big.Int).Mul(x2, inverse(k, c.N))
		s := new(big.Int).Mod(x3, c.N)

		return Signatrue{
			R: r,
			S: s,
		}
	}
}

func (c *curve) verify(public_key *Point, message string, signature Signatrue) bool {
	mes := hash(message)

	r := signature.R
	s := signature.S

	// S = (mes + r * private_key) * k**(-1)
	// k = (mes + r * private_key) * S**(-1)
	// k = mes * S**(-1) + r * private_key * S**(-1)
	// k*G = mes * S**(-1) * G + r * S**(-1) * private_key * G
	// private_key * G = public_key
	// R = k*G = mes * S**(-1) * G + r * S**(-1) * public_key
	// R.x mod n == r mod n

	w := inverse(s, c.N)

	u1_ := new(big.Int).Mul(mes, w)
	u1 := new(big.Int).Mod(u1_, c.N)

	u2_ := new(big.Int).Mul(r, w)
	u2 := new(big.Int).Mod(u2_, c.N)

	point := c.Add(c.Mul(u1), c.mul(u2, *public_key))

	x := new(big.Int).Set(point.X)
	x_mod := new(big.Int).Mod(x, c.N)

	r_mod := new(big.Int).Mod(r, c.N)

	if x_mod.Cmp(r_mod) == 0 {
		return true
	}

	return false
}
