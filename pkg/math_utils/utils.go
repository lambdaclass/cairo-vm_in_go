package math_utils

import (
	"github.com/pkg/errors"
	"math/big"
)

// Finds a nonnegative integer x < p such that (m * x) % p == n.
func DivMod(n *big.Int, m *big.Int, p *big.Int) (*big.Int, error) {
	if m.BitLen() == 0 {
		return nil, errors.Errorf("m in div_mod(n, m, p) can't be zero")
	}
	var inv_m, res *big.Int
	inv_m.ModInverse(m, p)
	res.Mul(inv_m, n)
	res.Mod(res, p)
	return res, nil
}
