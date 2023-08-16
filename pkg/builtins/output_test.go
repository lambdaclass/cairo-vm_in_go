package builtins_test

import (
	"reflect"
	"testing"

	"github.com/lambdaclass/cairo-vm.go/pkg/builtins"
	"github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
)

func TestOutputInitializeSegments(t *testing.T) {
	mem_manager := memory.NewMemorySegmentManager()
	output := builtins.NewOutputBuiltinRunner(true)
	output.InitializeSegments(&mem_manager)

	if mem_manager.Memory.NumSegments() != 1 {
		t.Errorf("Wrong number of segments after InitializeSegments")
	}

	if !reflect.DeepEqual(output.Base(), memory.NewRelocatable(0, 0)) {
		t.Errorf("Wrong builtin base after InitializeSegments")
	}
}

func TestOutputInitialStackIncluded(t *testing.T) {
	output := builtins.NewOutputBuiltinRunner(true)
	initial_stack := output.InitialStack()
	expected_stack := []memory.MaybeRelocatable{*memory.NewMaybeRelocatableRelocatable(output.Base())}
	if !reflect.DeepEqual(initial_stack, expected_stack) {
		t.Errorf("Wrong initial stack")
	}
}

func TestOutputInitialStackNotIncluded(t *testing.T) {
	output := builtins.NewOutputBuiltinRunner(false)
	if len(output.InitialStack()) != 0 {
		t.Errorf("Initial stack should be empty if not included")
	}
}

func TestOutputAddValidationRule(t *testing.T) {
	empty_mem := memory.NewMemory()
	mem := memory.NewMemory()
	output := builtins.NewOutputBuiltinRunner(true)
	output.AddValidationRule(mem)
	// Check that the memory is equal to a newly created one to check that
	// no validation rules were added
	if !reflect.DeepEqual(mem, empty_mem) {
		t.Errorf("AddValidationRule should do nothing")
	}
}
