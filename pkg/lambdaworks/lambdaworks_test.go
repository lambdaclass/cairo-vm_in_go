package lambdaworks

import (
	"testing"
)

func TestFromHex(t *testing.T) {
	var h_one = "1a"
	expected := From(26)

	result := FeltFromHex(h_one)
	if result != expected {
		t.Errorf("TestFromHex failed. Expected: %v, Got: %v", expected, result)
	}

}

func TestFromDecString(t *testing.T) {
	var s_one = "435"
	expected := From(435)

	result := FeltFromDecString(s_one)
	if result != expected {
		t.Errorf("TestFromDecString failed. Expected: %v, Got: %v", expected, result)
	}
}

func TestFeltSub(t *testing.T) {
	var felt Felt
	f_one := felt.One()
	expected := felt.Zero()

	result := f_one.Sub(f_one)
	if result != expected {
		t.Errorf("TestFeltSub failed. Expected: %v, Got: %v", expected, result)
	}
}

func TestFeltAdd(t *testing.T) {
	var felt Felt
	f_zero := felt.Zero()
	f_one := felt.One()
	expected := felt.One()

	result := f_zero.Add(f_one)
	if result != expected {
		t.Errorf("TestFeltAdd failed. Expected: %v, Got: %v", expected, result)
	}
}

func TestFeltAddFelts(t *testing.T) {
	var felt Felt
	f_zero := felt.Zero()
	felts := [3]Felt{felt.One(), felt.One(), felt.One()}
	expected := From(3)

	result := f_zero.AddFelts(felts[:])
	if result != expected {
		t.Errorf("TestFeltAdd failed. Expected: %v, Got: %v", expected, result)
	}
}

func TestFeltMul1(t *testing.T) {
	var felt Felt
	f_one := felt.One()
	expected := felt.One()

	result := f_one.Mul(f_one)
	if result != expected {
		t.Errorf("TestFeltMul1 failed. Expected: %v, Got: %v", expected, result)
	}
}

func TestFeltMul0(t *testing.T) {
	var felt Felt
	f_one := felt.One()
	f_zero := felt.Zero()
	expected := felt.Zero()

	result := f_zero.Mul(f_one)
	if result != expected {
		t.Errorf("TestFeltMul0 failed. Expected: %v, Got: %v", expected, result)
	}
}

func TestFeltMul9(t *testing.T) {
	f_three := From(3)
	expected := From(9)

	result := f_three.Mul(f_three)
	if result != expected {
		t.Errorf("TestFeltMul9 failed. Expected: %v, Got: %v", expected, result)
	}
}

func TestFeltDiv3(t *testing.T) {
	f_three := From(3)
	expected := From(1)

	result := f_three.Div(f_three)
	if result != expected {
		t.Errorf("TestFeltDiv3 failed. Expected: %v, Got: %v", expected, result)
	}
}

func TestFeltDiv4(t *testing.T) {
	f_four := From(4)
	f_two := From(2)

	expected := From(2)

	result := f_four.Div(f_two)
	if result != expected {
		t.Errorf("TestFeltDiv4 failed. Expected: %v, Got: %v", expected, result)
	}
}

func TestFeltDiv4Error(t *testing.T) {
	f_four := From(4)
	f_one := From(1)

	expected := From(45)

	result := f_four.Div(f_one)
	if result == expected {
		t.Errorf("TestFeltDiv4Error failed. Expected: %v, Got: %v", expected, result)
	}
}
