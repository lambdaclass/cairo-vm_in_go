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
