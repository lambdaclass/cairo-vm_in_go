package utils

import (
	"math"

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
	return 1 + (x-1)/y
}
