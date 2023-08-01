package lambdaworks

import (
	"testing"
)

func TestFeltSub(t *testing.T) {
	var f_one Felt
	f_one.One()
	var expected Felt
	expected.One()

	result := Sub(f_one, f_one)
	if result != expected {
		t.Errorf("TestFeltSub failed. Expected: %v, Got: %v", expected, result)
	}
}

func TestFeltAdd(t *testing.T) {
	var f_one Felt
	var f_zero Felt
	f_zero.Zero()
	f_one.One()

	var expected Felt
	expected.One()

	result := Add(f_one, f_zero)
	if result != expected {
		t.Errorf("TestFeltAdd failed. Expected: %v, Got: %v", expected, result)
	}
}

func TestFeltMul1(t *testing.T) {
	var f_one Felt
	f_one.One()

	var expected Felt
	expected.One()

	result := Mul(f_one, f_one)
	if result != expected {
		t.Errorf("TestFeltMul1 failed. Expected: %v, Got: %v", expected, result)
	}
}

func TestFeltMul0(t *testing.T) {
	var f_one Felt
	var f_zero Felt
	f_one.One()
	f_zero.Zero()

	var expected Felt
	expected.Zero()

	result := Mul(f_one, f_zero)
	if result != expected {
		t.Errorf("TestFeltMul0 failed. Expected: %v, Got: %v", expected, result)
	}
}

func TestFeltMul9(t *testing.T) {
	f_three := From(3)

	var expected Felt = From(9)

	result := Mul(f_three, f_three)
	if result != expected {
		t.Errorf("TestFeltMul9 failed. Expected: %v, Got: %v", expected, result)
	}
}

func TestFeltDiv3(t *testing.T) {
	f_three := From(3)

	var expected Felt = From(1)

	result := Div(f_three, f_three)
	if result != expected {
		t.Errorf("TestFeltDiv3 failed. Expected: %v, Got: %v", expected, result)
	}
}

func TestFeltDiv4(t *testing.T) {
	f_four := From(4)
	f_two := From(2)

	var expected Felt = From(2)

	result := Div(f_four, f_two)
	if result != expected {
		t.Errorf("TestFeltDiv4 failed. Expected: %v, Got: %v", expected, result)
	}
}

func TestFeltDiv4Error(t *testing.T) {
	f_four := From(4)
	f_one := From(1)

	var expected Felt = From(45)

	result := Div(f_four, f_one)
	if result == expected {
		t.Errorf("TestFeltDiv4Error failed. Expected: %v, Got: %v", expected, result)
	}
}
