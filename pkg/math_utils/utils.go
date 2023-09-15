package math_utils

import (
	"github.com/pkg/errors"
	"math/big"
)

// Finds a nonnegative integer x < p such that (m * x) % p == n.
func DivMod(n *big.Int, m *big.Int, p *big.Int) (*big.Int, error) {
	a := new(big.Int)
	gcd := new(big.Int)
	gcd.GCD(a, nil, m, p)

	if gcd.Cmp(big.NewInt(1)) != 0 {
		return nil, errors.Errorf("gcd(%s, %s) != 1", m, p)
	}

	return n.Mul(n, a).Mod(n, p), nil
}

func ISqrt(x *big.Int) (*big.Int, error) {
	if x.Sign() == -1 {
		return nil, errors.Errorf("Expected x: %s to be non-negative", x)
	}
	res := new(big.Int)
	return res.Sqrt(x), nil
}
