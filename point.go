package main

import "math/big"

type Point struct {
	X *big.Int
	Y *big.Int
}

func DefaulttPoint() *Point {
	return NewPoint(
		big.NewInt(0),
		big.NewInt(0))
}

func NewPoint(x, y *big.Int) *Point {
	return &Point{
		X: x,
		Y: y,
	}
}

func (p *Point) Zero() bool {
	if p.X.Sign() == 0 && p.Y.Sign() == 0 {
		return true
	}
	return false
}

func (p *Point) Opposite(q *Point) bool {
	if p.X.Cmp(q.X) == 0 && p.Y.Cmp(q.Y) != 0 && p.Y.CmpAbs(q.Y) == 0 {
		return true
	}
	return false
}

func (p *Point) Identical(q *Point) bool {
	if p.X.Cmp(q.X) == 0 && p.Y.Cmp(q.Y) == 0 {
		return true
	}
	return false
}
