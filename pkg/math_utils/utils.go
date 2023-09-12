package math_utils

import (
	"math/big"

	"github.com/pkg/errors"
)

// Finds a nonnegative integer x < p such that (m * x) % p == n.
func DivMod(n *big.Int, m *big.Int, p *big.Int) (*big.Int, error) {
	if m.BitLen() == 0 {
		return nil, errors.Errorf("m in div_mod(n, m, p) can't be zero")
	}
	inv_m := new(big.Int).ModInverse(m, p)
	res := new(big.Int).Mul(inv_m, n)
	res.Mod(res, p)
	return res, nil
}
