package math_utils_test

import (
	"math/big"
	"testing"

	. "github.com/lambdaclass/cairo-vm.go/pkg/math_utils"
)

func TestDivModOk(t *testing.T) {
	a := new(big.Int)
	b := new(big.Int)
	prime := new(big.Int)
	expected := new(big.Int)

	a.SetString("11260647941622813594563746375280766662237311019551239924981511729608487775604310196863705127454617186486639011517352066501847110680463498585797912894788", 10)
	b.SetString("4020711254448367604954374443741161860304516084891705811279711044808359405970", 10)
	prime.SetString("800000000000011000000000000000000000000000000000000000000000001", 16)
	expected.SetString("2904750555256547440469454488220756360634457312540595732507835416669695939476", 10)

	num, err := DivMod(a, b, prime)
	if err != nil {
		t.Errorf("DivMod failed with error: %s", err)
	}
	if num.Cmp(expected) != 0 {
		t.Errorf("Expected result: %s to be equal to %s", num, expected)
	}
}

func TestDivModMZeroFail(t *testing.T) {
	a := new(big.Int)
	b := new(big.Int)
	prime := new(big.Int)

	a.SetString("11260647941622813594563746375280766662237311019551239924981511729608487775604310196863705127454617186486639011517352066501847110680463498585797912894788", 10)
	prime.SetString("800000000000011000000000000000000000000000000000000000000000001", 16)

	_, err := DivMod(a, b, prime)
	if err == nil {
		t.Errorf("DivMod expected to failed with gcd != 1")
	}
}

func TestDivModMEqPFail(t *testing.T) {
	a := new(big.Int)
	b := new(big.Int)
	prime := new(big.Int)

	a.SetString("11260647941622813594563746375280766662237311019551239924981511729608487775604310196863705127454617186486639011517352066501847110680463498585797912894788", 10)
	b.SetString("800000000000011000000000000000000000000000000000000000000000001", 16)
	prime.SetString("800000000000011000000000000000000000000000000000000000000000001", 16)

	_, err := DivMod(a, b, prime)
	if err == nil {
		t.Errorf("DivMod expected to failed with gcd != 1")
	}
}

func TestIsSqrtOk(t *testing.T) {
	x := new(big.Int)
	y := new(big.Int)
	x.SetString("4573659632505831259480", 10)
	y.Mul(x, x)

	sqr_y, err := ISqrt(y)
	if err != nil {
		t.Errorf("ISqrt failed with error: %s", err)
	}
	if x.Cmp(sqr_y) != 0 {
		t.Errorf("Failed to get square root of x^2, x: %s", x)
	}
}

func TestCalculateIsqrtA(t *testing.T) {
	x := new(big.Int)
	x.SetString("81", 10)
	sqrt, err := ISqrt(x)
	if err != nil {
		t.Error("ISqrt failed")
	}

	expected := new(big.Int)
	expected.SetString("9", 10)

	if sqrt.Cmp(expected) != 0 {
		t.Errorf("ISqrt failed, expected %d, got %d", expected, sqrt)
	}
}

func TestCalculateIsqrtB(t *testing.T) {
	x := new(big.Int)
	x.SetString("4573659632505831259480", 10)
	square := new(big.Int)
	square = square.Mul(x, x)

	sqrt, err := ISqrt(square)
	if err != nil {
		t.Error("ISqrt failed")
	}

	if sqrt.Cmp(x) != 0 {
		t.Errorf("ISqrt failed, expected %d, got %d", x, sqrt)
	}
}

func TestCalculateIsqrtC(t *testing.T) {
	x := new(big.Int)
	x.SetString("3618502788666131213697322783095070105623107215331596699973092056135872020481", 10)
	square := new(big.Int)
	square = square.Mul(x, x)

	sqrt, err := ISqrt(square)
	if err != nil {
		t.Error("ISqrt failed")
	}

	if sqrt.Cmp(x) != 0 {
		t.Errorf("ISqrt failed, expected %d, got %d", x, sqrt)
	}
}

func TestIsSqrtFail(t *testing.T) {
	x := big.NewInt(-1)

	_, err := ISqrt(x)
	if err == nil {
		t.Errorf("expected ISqrt to fail")
	}
}
