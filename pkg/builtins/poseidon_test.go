package builtins_test

import (
	"reflect"
	"testing"

	"github.com/lambdaclass/cairo-vm.go/pkg/builtins"
	"github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
)

func TestInitializeSegments(t *testing.T) {
	mem_manager := memory.NewMemorySegmentManager()
	poseidon := builtins.NewPoseidonBuiltinRunner(true)
	poseidon.InitializeSegments(&mem_manager)

	if mem_manager.Memory.NumSegments() != 1 {
		t.Errorf("Wrong number of segments after InitializeSegments")
	}
	if !reflect.DeepEqual(poseidon.Base(), memory.Relocatable{SegmentIndex: 0, Offset: 0}) {
		t.Errorf("Wrong builtin base after InitializeSegments")
	}

}

func TestInitialStackIncluded(t *testing.T) {
	poseidon := builtins.NewPoseidonBuiltinRunner(true)
	initial_stack := poseidon.InitialStack()
	expected_stack := []memory.MaybeRelocatable{*memory.NewMaybeRelocatableRelocatable(poseidon.Base())}
	if !reflect.DeepEqual(initial_stack, expected_stack) {
		t.Errorf("Wrong initial stack")
	}
}

func TestInitialStackNotIncluded(t *testing.T) {
	poseidon := builtins.NewPoseidonBuiltinRunner(false)
	if len(poseidon.InitialStack()) != 0 {
		t.Errorf("Initial stack should be empty if not included")
	}
}

func TestAddValidationRule(t *testing.T) {
	empty_mem := memory.NewMemory()
	mem := memory.NewMemory()
	poseidon := builtins.NewPoseidonBuiltinRunner(true)
	poseidon.AddValidationRule(mem)
	// Check that the memory is equal to a newly created one to check that
	// no validation rules were added
	if !reflect.DeepEqual(mem, empty_mem) {
		t.Errorf("AddValidationRule should do nothing")
	}

}
