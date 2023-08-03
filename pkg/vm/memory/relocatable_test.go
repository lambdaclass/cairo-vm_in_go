package memory_test

import (
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
