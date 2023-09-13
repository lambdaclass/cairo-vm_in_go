package memory_test

import (
	"reflect"
	"testing"

	"github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
	"github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
)

func TestMaybeRelocatableIsEqual(t *testing.T) {
	rel_felt_one := memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(24))
	rel_felt_two := memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(24))

	rel_rel_one := memory.NewMaybeRelocatableRelocatable(memory.Relocatable{2, 4})
	rel_rel_two := memory.NewMaybeRelocatableRelocatable(memory.Relocatable{2, 4})

	rel_felts_are_equal := rel_felt_one.IsEqual(rel_felt_two)
	rel_rel_are_equal := rel_rel_one.IsEqual(rel_rel_two)

	if !rel_felts_are_equal {
		t.Errorf("%s and %s are not equal", rel_felt_one, rel_felt_two)
	}
	if !rel_rel_are_equal {
		t.Errorf("%s and %s are not equal", rel_rel_one, rel_rel_two)
	}

}

func TestMaybeRelocatableIsNotEqual(t *testing.T) {
	rel_felt_one := memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(2))
	rel_felt_two := memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(4))

	rel_rel_one := memory.NewMaybeRelocatableRelocatable(memory.Relocatable{2, 4})
	rel_rel_two := memory.NewMaybeRelocatableRelocatable(memory.Relocatable{4, 2})

	rel_felts_are_equal := rel_felt_one.IsEqual(rel_felt_two)
	rel_rel_are_equal := rel_rel_one.IsEqual(rel_rel_two)
	rel_rel_and_rel_felt_are_equal := rel_felt_one.IsEqual(rel_rel_one)

	if rel_felts_are_equal {
		t.Errorf("%s and %s are equal", rel_felt_one, rel_felt_two)
	}
	if rel_rel_are_equal {
		t.Errorf("%s and %s are equal", rel_rel_one, rel_rel_two)
	}
	if rel_rel_and_rel_felt_are_equal {
		t.Errorf("%s and %s are equal", rel_rel_one, rel_felt_two)
	}
}

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
	if res != lambdaworks.FeltFromUint64(2) {
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
	res, err := a.Sub(b)
	if err != nil {
		t.Errorf("Relocatable.Sub failed with error: %s", err)
	}
	if res != lambdaworks.FeltFromDecString("-2") {
		t.Errorf("Got wrong value from Relocatable.Sub")
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

func TestRelocatableAddIntPositive(t *testing.T) {
	rel := memory.Relocatable{2, 4}
	res, err := rel.AddInt(24)
	expected := memory.Relocatable{2, 28}
	if err != nil {
		t.Errorf("Relocatable.AddInt failed with error: %s", err)
	}
	if res != expected {
		t.Errorf("got wrong value from Relocatable.AddInt, expected: %v, got: %v", expected, res)
	}
}

func TestRelocatableAddIntNegative(t *testing.T) {
	rel := memory.Relocatable{2, 24}
	res, err := rel.AddInt(-4)
	expected := memory.Relocatable{2, 20}
	if err != nil {
		t.Errorf("Relocatable.AddInt failed with error: %s", err)
	}
	if res != expected {
		t.Errorf("got wrong value from Relocatable.AddInt, expected: %v, got: %v", expected, res)
	}
}
