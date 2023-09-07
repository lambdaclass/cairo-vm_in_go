package lambdaworks_test

import (
	"reflect"
	"testing"

	"github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
)

func TestFromHex(t *testing.T) {
	var h_one = "1a"
	expected := lambdaworks.FeltFromUint64(26)

	result := lambdaworks.FeltFromHex(h_one)
	if result != expected {
		t.Errorf("TestFromHex failed. Expected: %v, Got: %v", expected, result)
	}

}

func TestToHex(t *testing.T) {
	var expected = "0x1a"
	felt := lambdaworks.FeltFromUint64(26)

	result := felt.ToHexString()
	if result != expected {
		t.Errorf("TestFromHex failed. Expected: %v, Got: %v", expected, result)
	}

}

func TestFromDecString(t *testing.T) {
	var s_one = "435"
	expected := lambdaworks.FeltFromUint64(435)

	result := lambdaworks.FeltFromDecString(s_one)
	if result != expected {
		t.Errorf("TestFromDecString failed. Expected: %v, Got: %v", expected, result)
	}
}

func TestFromNegDecString(t *testing.T) {
	var s_one = "-1"
	expected := lambdaworks.FeltFromHex("800000000000011000000000000000000000000000000000000000000000000")

	result := lambdaworks.FeltFromDecString(s_one)
	if result != expected {
		t.Errorf("TestFromNegDecString failed. Expected: %v, Got: %v", expected, result)
	}
}

func TestToLeBytes(t *testing.T) {
	expected := [32]uint8{
		1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	}
	actual := *lambdaworks.FeltOne().ToLeBytes()

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("TestToLeBytes failed. Expected: %v, Got: %v", expected, actual)
	}
}

func TestToBeBytes(t *testing.T) {
	expected := [32]uint8{
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1,
	}
	actual := *lambdaworks.FeltOne().ToBeBytes()

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("TestToBeBytes failed. Expected: %v, Got: %v", expected, actual)
	}
}

func TestFromLeBytes(t *testing.T) {
	bytes := [32]uint8{
		1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	}
	felt_from_bytes := lambdaworks.FeltFromLeBytes(&bytes)

	if !reflect.DeepEqual(felt_from_bytes, lambdaworks.FeltOne()) {
		t.Errorf("TestFromLeBytes failed. Expected 1, Got: %v", felt_from_bytes)
	}
}

func TestFromBeBytes(t *testing.T) {
	bytes := [32]uint8{
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1,
	}
	felt_from_bytes := lambdaworks.FeltFromBeBytes(&bytes)

	if !reflect.DeepEqual(felt_from_bytes, lambdaworks.FeltOne()) {
		t.Errorf("TestToFromBeBytes failed. Expected 1, Got: %v", felt_from_bytes)
	}
}

func TestFeltSub(t *testing.T) {
	f_one := lambdaworks.FeltOne()
	expected := lambdaworks.FeltZero()

	result := f_one.Sub(f_one)
	if result != expected {
		t.Errorf("TestFeltSub failed. Expected: %v, Got: %v", expected, result)
	}
}

func TestFeltAdd(t *testing.T) {
	f_zero := lambdaworks.FeltZero()
	f_one := lambdaworks.FeltOne()
	expected := lambdaworks.FeltOne()

	result := f_zero.Add(f_one)
	if result != expected {
		t.Errorf("TestFeltAdd failed. Expected: %v, Got: %v", expected, result)
	}
}

func TestFeltMul1(t *testing.T) {
	f_one := lambdaworks.FeltOne()
	expected := lambdaworks.FeltOne()

	result := f_one.Mul(f_one)
	if result != expected {
		t.Errorf("TestFeltMul1 failed. Expected: %v, Got: %v", expected, result)
	}
}

func TestFeltMul0(t *testing.T) {
	f_one := lambdaworks.FeltOne()
	f_zero := lambdaworks.FeltZero()
	expected := lambdaworks.FeltZero()

	result := f_zero.Mul(f_one)
	if result != expected {
		t.Errorf("TestFeltMul0 failed. Expected: %v, Got: %v", expected, result)
	}
}

func TestFeltMul9(t *testing.T) {
	f_three := lambdaworks.FeltFromUint64(3)
	expected := lambdaworks.FeltFromUint64(9)

	result := f_three.Mul(f_three)
	if result != expected {
		t.Errorf("TestFeltMul9 failed. Expected: %v, Got: %v", expected, result)
	}
}

func TestFeltDiv3(t *testing.T) {
	f_three := lambdaworks.FeltFromUint64(3)
	expected := lambdaworks.FeltFromUint64(1)

	result := f_three.Div(f_three)
	if result != expected {
		t.Errorf("TestFeltDiv3 failed. Expected: %v, Got: %v", expected, result)
	}
}

func TestFeltDiv4(t *testing.T) {
	f_four := lambdaworks.FeltFromUint64(4)
	f_two := lambdaworks.FeltFromUint64(2)

	expected := lambdaworks.FeltFromUint64(2)

	result := f_four.Div(f_two)
	if result != expected {
		t.Errorf("TestFeltDiv4 failed. Expected: %v, Got: %v", expected, result)
	}
}

func TestFeltDiv4Error(t *testing.T) {
	f_four := lambdaworks.FeltFromUint64(4)
	f_one := lambdaworks.FeltFromUint64(1)

	expected := lambdaworks.FeltFromUint64(45)

	result := f_four.Div(f_one)
	if result == expected {
		t.Errorf("TestFeltDiv4Error failed. Expected: %v, Got: %v", expected, result)
	}
}

func TestBits(t *testing.T) {
	f_zero := lambdaworks.FeltZero()
	if f_zero.Bits() != 0 {
		t.Errorf("TestBits failed. Expected: %d, Got: %d", 1, f_zero.Bits())
	}
	f_one := lambdaworks.FeltOne()
	if f_one.Bits() != 1 {
		t.Errorf("TestBits failed. Expected: %d, Got: %d", 1, f_one.Bits())
	}
	f_eight := lambdaworks.FeltFromUint64(8)
	if f_eight.Bits() != 4 {
		t.Errorf("TestBits failed. Expected: %d, Got: %d", 4, f_eight.Bits())
	}
	f_fifteen := lambdaworks.FeltFromUint64(15)
	if f_fifteen.Bits() != 4 {
		t.Errorf("TestBits failed. Expected: %d, Got: %d", 4, f_fifteen.Bits())
	}

	f_neg_one := lambdaworks.FeltFromDecString("-1")
	if f_neg_one.Bits() != 252 {
		t.Errorf("TestBits failed. Expected: %d, Got: %d", 252, f_neg_one.Bits())
	}
}

func TestToU641(t *testing.T) {
	felt := lambdaworks.FeltOne()
	result, err := felt.ToU64()

	var expected uint64 = 1

	if expected != result {
		t.Errorf("Error in conversion expected: %v, got %v with err: %v", expected, result, err)
	}

}

func TestToU6410230(t *testing.T) {
	felt := lambdaworks.FeltFromUint64(10230)
	result, err := felt.ToU64()

	var expected uint64 = 10230

	if expected != result {
		t.Errorf("Error in conversion expected: %v, got %v with err: %v", expected, result, err)
	}
}

func TestToU64Fail(t *testing.T) {
	felt := lambdaworks.FeltFromDecString("9999999999999999999999999")

	_, err := felt.ToU64()
	expected_err := lambdaworks.ConversionError(felt, "u64")

	if err.Error() != expected_err.Error() {
		t.Errorf("Conversion test should fail with error: %v", expected_err)
	}
}
func TestFeltIsZero(t *testing.T) {
	f_zero := lambdaworks.FeltZero()

	is_zero := f_zero.IsZero()

	if !is_zero {
		t.Errorf("TestFeltIsZero failed. Expected true, Got: %v", is_zero)
	}
}

func TestFeltIsNotZero(t *testing.T) {
	f_one := lambdaworks.FeltOne()

	is_zero := f_one.IsZero()

	if is_zero {
		t.Errorf("TestFeltIsNotZero failed. Expected false, Got: %v", is_zero)
	}
}
