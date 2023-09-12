package math_utils_test

import (
	"math/big"
	"testing"

	. "github.com/lambdaclass/cairo-vm.go/pkg/math_utils"
)

func TestDivModOk(t *testing.T) {
	a, _ := new(big.Int).SetString("11260647941622813594563746375280766662237311019551239924981511729608487775604310196863705127454617186486639011517352066501847110680463498585797912894788", 10)
	b, _ := new(big.Int).SetString("4020711254448367604954374443741161860304516084891705811279711044808359405970", 10)
	prime, _ := new(big.Int).SetString("800000000000011000000000000000000000000000000000000000000000001", 16)
	expected, _ := new(big.Int).SetString("2904750555256547440469454488220756360634457312540595732507835416669695939476", 10)
	num, err := DivMod(a, b, prime)
	if err != nil {
		t.Errorf("DivMod failed with error: %s", err)
	}
	if num.Cmp(expected) != 0 {
		t.Errorf("Expected result: %s to be equal to %s", num, expected)
	}
}

// func TestDivModFail(t *testing.T) {
// 	b := big.NewInt(1)
// 	a, _ := new(big.Int).SetString("11260647941622813594563746375280766662237311019551239924981511729608487775604310196863705127454617186486639011517352066501847110680463498585797912894788", 10)
// 	prime, _ := new(big.Int).SetString("800000000000011000000000000000000000000000000000000000000000001", 16)

// 	_, err := DivMod(a, b, prime)
// 	if err == nil {
// 		t.Errorf("DivMod expected to fail with division by zero")
// 	}
// }
