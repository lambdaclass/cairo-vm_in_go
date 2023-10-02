package utils_test

import (
	"math/big"
	"testing"

	. "github.com/lambdaclass/cairo-vm.go/pkg/utils"
)

func TestSafeDivBigOkSmallNumbers(t *testing.T) {
	a := big.NewInt(10)
	b := big.NewInt(5)
	expected := big.NewInt(2)
	val, err := SafeDivBig(a, b)
	if err != nil || val.Cmp(expected) != 0 {
		t.Error("Wrong value returned by SafeDivBig")
	}
}

func TestSafeDivBigOkBigNumbers(t *testing.T) {
	a, _ := new(big.Int).SetString("1354671610274991796869769298862800192014", 10)
	b, _ := new(big.Int).SetString("37853847583", 10)
	expected, _ := new(big.Int).SetString("35786893453952854863476753458", 10)
	val, err := SafeDivBig(a, b)
	if err != nil || val.Cmp(expected) != 0 {
		t.Error("Wrong value returned by SafeDivBig")
	}
}

func TestSafeDivBigErrNotDivisible(t *testing.T) {
	a := big.NewInt(10)
	b := big.NewInt(7)
	_, err := SafeDivBig(a, b)
	if err == nil {
		t.Error("SafeDivBig should have failed")
	}
}

func TestSafeDivBigErrZeroDivison(t *testing.T) {
	a := big.NewInt(10)
	b := big.NewInt(0)
	_, err := SafeDivBig(a, b)
	if err == nil {
		t.Error("SafeDivBig should have failed")
	}
}

func TestIgcdex11(t *testing.T) {
	a := big.NewInt(1)
	b := big.NewInt(1)
	expectedX, expectedY, expectedZ := big.NewInt(0), big.NewInt(1), big.NewInt(1)
	x, y, z := Igcdex(a, b)
	if x.Cmp(expectedX) != 0 || y.Cmp(expectedY) != 0 || z.Cmp(expectedZ) != 0 {
		t.Error("Wrong values returned by Igcdex")
	}
}

func TestIgcdex00(t *testing.T) {
	a := big.NewInt(0)
	b := big.NewInt(0)
	expectedX, expectedY, expectedZ := big.NewInt(0), big.NewInt(1), big.NewInt(0)
	x, y, z := Igcdex(a, b)
	if x.Cmp(expectedX) != 0 || y.Cmp(expectedY) != 0 || z.Cmp(expectedZ) != 0 {
		t.Error("Wrong values returned by Igcdex")
	}
}

func TestIgcdex10(t *testing.T) {
	a := big.NewInt(1)
	b := big.NewInt(0)
	expectedX, expectedY, expectedZ := big.NewInt(1), big.NewInt(0), big.NewInt(1)
	x, y, z := Igcdex(a, b)
	if x.Cmp(expectedX) != 0 || y.Cmp(expectedY) != 0 || z.Cmp(expectedZ) != 0 {
		t.Error("Wrong values returned by Igcdex")
	}
}

func TestIgcdex46(t *testing.T) {
	a := big.NewInt(4)
	b := big.NewInt(6)
	expectedX, expectedY, expectedZ := big.NewInt(-1), big.NewInt(1), big.NewInt(2)
	x, y, z := Igcdex(a, b)
	if x.Cmp(expectedX) != 0 || y.Cmp(expectedY) != 0 || z.Cmp(expectedZ) != 0 {
		t.Error("Wrong values returned by Igcdex")
	}
}

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
