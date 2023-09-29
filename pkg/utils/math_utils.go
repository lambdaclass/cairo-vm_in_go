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
