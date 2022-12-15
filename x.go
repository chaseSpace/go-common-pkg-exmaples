package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common/math"
	"math/big"
)

func extended_euclid_gcd(a, b *big.Int) (*big.Int, *big.Int, *big.Int) {
	s, oldS := big.NewInt(0), big.NewInt(1)
	t, oldT := big.NewInt(1), big.NewInt(0)
	r, oldR := b, a

	for r.Cmp(big.NewInt(0)) != 0 {
		quotient := oldR.Div(oldR, r)
		oldR, r = r, oldR.Sub(oldR, quotient.Mul(quotient, r))
		oldS, s = s, oldS.Sub(oldS, quotient.Mul(quotient, s))
		oldT, t = t, oldT.Sub(oldT, quotient.Mul(quotient, t))
	}
	return oldR, oldS, oldT
}

func inverse_mod(a, n *big.Int) *big.Int {
	_, x, _ := extended_euclid_gcd(a, n)
	if x.Cmp(big.NewInt(0)) == -1 {
		x = x.Add(x, n)
	}
	return x
}

func derivate_privkey(p, r, s1, s2, z1, z2 *big.Int) (*big.Int, *big.Int) {
	z := z1.Sub(z1, z2)
	s := s1.Sub(s1, s2)
	rInv := inverse_mod(r, p)
	sInv := inverse_mod(s, p)

	zxsInv := z.Mul(z, sInv)
	k := zxsInv.Mod(zxsInv, p)

	s1xk := s1.Mul(s1, k)
	_left := rInv.Mul(rInv, s1xk.Sub(s1xk, z1))
	d := _left.Mod(_left, p)
	return d, k
}

func main() {
	z1, _ := math.ParseBig256("0x4f6a8370a435a27724bbc163419042d71b6dcbeb61c060cc6816cda93f57860c")
	s1, _ := math.ParseBig256("0x2bbd9c2a6285c2b43e728b17bda36a81653dd5f4612a2e0aefdb48043c5108de")
	r, _ := math.ParseBig256("0x69a726edfb4b802cbf267d5fd1dabcea39d3d7b4bf62b9eeaeba387606167166")
	z2, _ := math.ParseBig256("0x350f3ee8007d817fbd7349c477507f923c4682b3e69bd1df5fbb93b39beb1e04")
	s2, _ := math.ParseBig256("0x7724cedeb923f374bef4e05c97426a918123cc4fec7b07903839f12517e1b3c8")
	p, _ := math.ParseBig256("0xFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEBAAEDCE6AF48A03BBFD25E8CD0364141")

	r1, r2 := derivate_privkey(p, r, s1, s2, z1, z2)

	fmt.Printf("key:%s\n k:%s %s", r1, r2, z1)
}
