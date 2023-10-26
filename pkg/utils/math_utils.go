package utils

import (
	"math"
	"math/big"

	"github.com/pkg/errors"
)

func NextPowOf2(n uint) uint {
	var k uint = 1
	for k < n {
		k = k << 1
	}
	return k
}

// Performs integer division between x and y; fails if x is not divisible by y.
func SafeDiv(x uint, y uint) (uint, error) {
	if y == 0 {
		return 0, errors.New("Attempted to divide by zero")
	}
	rem := math.Remainder(float64(x), float64(y))
	if rem != 0 {
		return 0, errors.Errorf("%d is not divisible by %d", x, y)
	}

	return x / y, nil
}

func MinInt(x int, y int) int {
	if x < y {
		return x
	} else {
		return y
	}
}

func MaxInt(x int, y int) int {
	if x > y {
		return x
	} else {
		return y
	}
}

func DivCeil(x uint, y uint) uint {
	q := x / y
	if x%y != 0 {
		q++
	}
	return q
}

// Performs integer division between x and y; fails if x is not divisible by y.
func SafeDivBig(x *big.Int, y *big.Int) (*big.Int, error) {
	if y.Cmp(big.NewInt(0)) == 0 {
		return &big.Int{}, errors.New("SafeDiv: Attempted to divide by zero")
	}
	q, r := x.DivMod(x, y, new(big.Int))
	if r.Cmp(big.NewInt(0)) != 0 {
		return &big.Int{}, errors.Errorf("SafeDiv: %s is not divisible by %s", x.Text(10), y.Text(10))
	}
	return q, nil
}

// Finds a nonnegative integer x < p such that (m * x) % p == n.
func DivMod(n *big.Int, m *big.Int, p *big.Int) (*big.Int, error) {
	a, _, c := Igcdex(m, p)
	if c.Cmp(big.NewInt(1)) != 0 {
		return nil, errors.Errorf("Operation failed: divmod(%s, %s, %s), igcdex(%s, %s) != 1 ", n.Text(10), m.Text(10), p.Text(10), m.Text(10), p.Text(10))
	}
	return new(big.Int).Mod(new(big.Int).Mul(n, a), p), nil
}

func Igcdex(a *big.Int, b *big.Int) (*big.Int, *big.Int, *big.Int) {
	zero := big.NewInt(0)
	one := big.NewInt(1)
	switch true {
	case a.Cmp(zero) == 0 && b.Cmp(zero) == 0:
		return zero, one, zero
	case a.Cmp(zero) == 0:
		return zero, big.NewInt(int64(a.Sign())), new(big.Int).Abs(b)
	case b.Cmp(zero) == 0:
		return big.NewInt(int64(a.Sign())), zero, new(big.Int).Abs(a)
	default:
		xSign := big.NewInt(int64(a.Sign()))
		ySign := big.NewInt(int64(b.Sign()))
		a = new(big.Int).Abs(a)
		b = new(big.Int).Abs(b)
		x, y, r, s := big.NewInt(1), big.NewInt(0), big.NewInt(0), big.NewInt(1)
		for b.Cmp(zero) != 0 {
			q, c := new(big.Int).DivMod(a, b, new(big.Int))
			x = new(big.Int).Sub(x, new(big.Int).Mul(q, r))
			y = new(big.Int).Sub(y, new(big.Int).Mul(q, s))

			a, b, r, s, x, y = b, c, x, y, r, s
		}

		return new(big.Int).Mul(x, xSign), new(big.Int).Mul(y, ySign), a

	}
}

func IsEven(n *big.Int) bool {
	res := new(big.Int).And(n, big.NewInt(1))
	return res.Cmp(big.NewInt(0)) != 0
}

func ISqrt(x *big.Int) (*big.Int, error) {
	if x.Sign() == -1 {
		return nil, errors.Errorf("Expected x: %s to be non-negative", x)
	}
	res := new(big.Int)
	return res.Sqrt(x), nil
}
