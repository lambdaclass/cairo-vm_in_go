package memory_test

import (
	"reflect"
	"testing"

	"github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
)

func TestMaybeRelocatableIsZeroInt(t *testing.T) {
	zero := memory.NewMaybeRelocatableInt(0)
	if !zero.IsZero() {
		t.Errorf("MaybeRelocatable(0) should be zero")
	}
	not_zero := memory.NewMaybeRelocatableInt(1)
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
	felt := memory.Int{Felt: 5}
	rel := memory.Relocatable{}
	res, err := rel.AddFelt(felt)
	if err != nil {
		t.Errorf("AddFelt failed with error: %s", err)
	}
	if !reflect.DeepEqual(res, memory.Relocatable{0, 5}) {
		t.Errorf("Got wrong value from Relocatable.AddFelt")
	}
}

func TestMaybeRelocatableAddMaybeRelocatableInt(t *testing.T) {
	mr := memory.NewMaybeRelocatableInt(5)
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
