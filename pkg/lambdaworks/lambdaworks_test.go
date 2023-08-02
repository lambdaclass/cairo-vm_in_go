package lambdaworks

import (
	"testing"
)

func TestFromHex(t *testing.T) {
	var h_one = "1a"
	expected := From(26)

	result := FromHex(h_one)
	if result != expected {
		t.Errorf("TestFromHex failed. Expected: %v, Got: %v", expected, result)
	}

}

func TestFromString(t *testing.T) {
	var s_one = "9"
	expected := From(9)

	result, err := FromString(s_one)
	if result != expected {
		t.Errorf("TestFromString failed. Expected: %v, Got: %v", expected, result)
		t.Errorf("Error is: %s", err)
	}
}

func TestFeltSub(t *testing.T) {
	var felt Felt
	f_one := felt.One()
	expected := felt.Zero()

	result := Sub(f_one, f_one)
	if result != expected {
		t.Errorf("TestFeltSub failed. Expected: %v, Got: %v", expected, result)
	}
}

func TestFeltAdd(t *testing.T) {
	var felt Felt
	f_zero := felt.Zero()
	f_one := felt.One()
	expected := felt.One()

	result := Add(f_one, f_zero)
	if result != expected {
		t.Errorf("TestFeltAdd failed. Expected: %v, Got: %v", expected, result)
	}
}

func TestFeltMul1(t *testing.T) {
	var felt Felt
	f_one := felt.One()
	expected := felt.One()

	result := Mul(f_one, f_one)
	if result != expected {
		t.Errorf("TestFeltMul1 failed. Expected: %v, Got: %v", expected, result)
	}
}

func TestFeltMul0(t *testing.T) {
	var felt Felt
	f_one := felt.One()
	f_zero := felt.Zero()
	expected := felt.Zero()

	result := Mul(f_one, f_zero)
	if result != expected {
		t.Errorf("TestFeltMul0 failed. Expected: %v, Got: %v", expected, result)
	}
}

func TestFeltMul9(t *testing.T) {
	f_three := From(3)
	expected := From(9)

	result := Mul(f_three, f_three)
	if result != expected {
		t.Errorf("TestFeltMul9 failed. Expected: %v, Got: %v", expected, result)
	}
}

func TestFeltDiv3(t *testing.T) {
	f_three := From(3)
	expected := From(1)

	result := Div(f_three, f_three)
	if result != expected {
		t.Errorf("TestFeltDiv3 failed. Expected: %v, Got: %v", expected, result)
	}
}

func TestFeltDiv4(t *testing.T) {
	f_four := From(4)
	f_two := From(2)

	expected := From(2)

	result := Div(f_four, f_two)
	if result != expected {
		t.Errorf("TestFeltDiv4 failed. Expected: %v, Got: %v", expected, result)
	}
}

func TestFeltDiv4Error(t *testing.T) {
	f_four := From(4)
	f_one := From(1)

	expected := From(45)

	result := Div(f_four, f_one)
	if result == expected {
		t.Errorf("TestFeltDiv4Error failed. Expected: %v, Got: %v", expected, result)
	}
}
