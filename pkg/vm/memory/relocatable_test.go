package memory_test

import (
	"reflect"
	"testing"

	"github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
	"github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
)

func TestMaybeRelocatableIsZeroInt(t *testing.T) {
	zero := memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(0))
	if !zero.IsZero() {
		t.Errorf("MaybeRelocatable(0) should be zero")
	}
	not_zero := memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(1))
	if not_zero.IsZero() {
		t.Errorf("MaybeRelocatable(1) should not be zero")
	}
}

func TestMaybeRelocatableIsZeroRelocatable(t *testing.T) {
	ptr_zero_zero := memory.NewMaybeRelocatableRelocatable(memory.Relocatable{})
	if ptr_zero_zero.IsZero() {
		t.Errorf("MaybeRelocatable(0:0) should not be zero")
	}

	ptr_one_one := memory.NewMaybeRelocatableRelocatable(memory.Relocatable{1, 1})
	if ptr_one_one.IsZero() {
		t.Errorf("MaybeRelocatable(1:1) should not be zero")
	}
}

func TestMaybeRelocatableAddFelt(t *testing.T) {
	felt := lambdaworks.FeltFromUint64(5)
	rel := memory.Relocatable{}
	res, err := rel.AddFelt(felt)
	if err != nil {
		t.Errorf("AddFelt failed with error: %s", err)
	}
	if !reflect.DeepEqual(res, memory.Relocatable{0, 5}) {
		t.Errorf("Got wrong value from Relocatable.AddFelt")
	}
}

func TestRelocatableIsEqual(t *testing.T) {
	a := memory.Relocatable{2, 4}
	b := memory.Relocatable{2, 4}

	is_equal := a.IsEqual(&b)
	if !is_equal {
		t.Errorf("TestRelocatableIsEqual failed epected true, got %v", is_equal)
	}

}

func TestRelocatableIsNotEqual(t *testing.T) {
	a := memory.Relocatable{2, 4}
	b := memory.Relocatable{4, 2}

	is_equal := a.IsEqual(&b)
	if is_equal {
		t.Errorf("TestRelocatableIsNotEqual failed epected false, got %v", is_equal)
	}
}

func TestRelocatableAddUint(t *testing.T) {
	rel := memory.Relocatable{2, 4}
	res := rel.AddUint(24)
	expected := memory.Relocatable{2, 28}
	if res != expected {
		t.Errorf("got wrong value from Relocatable.AddUint, expected: %v, got: %v", expected, res)
	}
}

func TestRelocatableSubOk(t *testing.T) {
	a := memory.Relocatable{1, 7}
	b := memory.Relocatable{1, 5}
	res, err := a.Sub(b)
	if err != nil {
		t.Errorf("Relocatable.Sub failed with error: %s", err)
	}
	if res != 2 {
		t.Errorf("Got wrong value from Relocatable.Sub")
	}
}

func TestRelocatableSubDiffIndex(t *testing.T) {
	a := memory.Relocatable{1, 7}
	b := memory.Relocatable{2, 5}
	_, err := a.Sub(b)
	if err == nil {
		t.Errorf("Relocatable.Sub should have failed")
	}
}

func TestRelocatableSubNegativeDifference(t *testing.T) {
	a := memory.Relocatable{1, 7}
	b := memory.Relocatable{1, 9}
	_, err := a.Sub(b)
	if err == nil {
		t.Errorf("Relocatable.Sub should have failed")
	}
}

func TestMaybeRelocatableSubFelt(t *testing.T) {
	felt := lambdaworks.FeltFromUint64(5)
	rel := memory.Relocatable{1, 7}
	res, err := rel.SubFelt(felt)
	if err != nil {
		t.Errorf("SubFelt failed with error: %s", err)
	}
	if !reflect.DeepEqual(res, memory.Relocatable{1, 2}) {
		t.Errorf("Got wrong value from Relocatable.SubFelt")
	}
}

func TestMaybeRelocatableSubFeltOutOfRange(t *testing.T) {
	felt := lambdaworks.FeltFromUint64(5)
	rel := memory.Relocatable{}
	_, err := rel.SubFelt(felt)
	if err == nil {
		t.Errorf("SubFelt should have failed")
	}
}

func TestMaybeRelocatableAddMaybeRelocatableInt(t *testing.T) {
	mr := memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(5))
	rel := memory.Relocatable{}
	res, err := rel.AddMaybeRelocatable(*mr)
	if err != nil {
		t.Errorf("AddMaybeRelocatable failed with error: %s", err)
	}
	if !reflect.DeepEqual(res, memory.Relocatable{0, 5}) {
		t.Errorf("Got wrong value from Relocatable.AddMaybeRelocatable")
	}
}

func TestMaybeRelocatableAddMaybeRelocatableRelocatable(t *testing.T) {
	mr := memory.NewMaybeRelocatableRelocatable(memory.Relocatable{})
	rel := memory.Relocatable{}
	_, err := rel.AddMaybeRelocatable(*mr)
	if err == nil {
		t.Errorf("Addition between relocatable values should fail")
	}
}

func TestMaybeRelocatableSubBothFelts(t *testing.T) {
	a := memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(7))
	b := memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(5))
	res, err := a.Sub(*b)
	if err != nil {
		t.Errorf("MaybeRelocatable.Sub failed with error: %s", err)
	}
	if !reflect.DeepEqual(res, *memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(2))) {
		t.Errorf("Got wrong value from Relocatable.MaybeRelocatable.Sub")
	}
}

func TestMaybeRelocatableSubBothRelocatable(t *testing.T) {
	a := memory.NewMaybeRelocatableRelocatable(memory.Relocatable{1, 7})
	b := memory.NewMaybeRelocatableRelocatable(memory.Relocatable{1, 5})
	res, err := a.Sub(*b)
	if err != nil {
		t.Errorf("MaybeRelocatable.Sub failed with error: %s", err)
	}
	if !reflect.DeepEqual(res, *memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(2))) {
		t.Errorf("Got wrong value from Relocatable.MaybeRelocatable.Sub")
	}
}

func TestMaybeRelocatableSubFeltFromRelocatable(t *testing.T) {
	a := memory.NewMaybeRelocatableRelocatable(memory.Relocatable{1, 7})
	b := memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(5))
	res, err := a.Sub(*b)
	if err != nil {
		t.Errorf("MaybeRelocatable.Sub failed with error: %s", err)
	}
	if !reflect.DeepEqual(res, *memory.NewMaybeRelocatableRelocatable(memory.Relocatable{1, 2})) {
		t.Errorf("Got wrong value from Relocatable.MaybeRelocatable.Sub")
	}
}

func TestMaybeRelocatableSubRelFromFelt(t *testing.T) {
	a := memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(7))
	b := memory.NewMaybeRelocatableRelocatable(memory.Relocatable{})
	_, err := a.Sub(*b)
	if err == nil {
		t.Errorf("Subtraction of relocatable from felt should fail")
	}
}
