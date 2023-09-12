package builtins_test

import (
	"reflect"
	"testing"

	"github.com/lambdaclass/cairo-vm.go/pkg/builtins"
	"github.com/lambdaclass/cairo-vm.go/pkg/vm"
	"github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
)

func TestOutputDeduceMemoryCell(t *testing.T) {
	output := builtins.NewOutputBuiltinRunner()
	a, b := output.DeduceMemoryCell(memory.NewRelocatable(0, 0), memory.NewMemory())
	if a != nil || b != nil {
		t.Errorf("DeduceMemoryCell should do nothing")
	}
}

func TestOutputInitializeSegments(t *testing.T) {
	mem_manager := memory.NewMemorySegmentManager()
	output := builtins.NewOutputBuiltinRunner()
	output.InitializeSegments(&mem_manager)

	if mem_manager.Memory.NumSegments() != 1 {
		t.Errorf("Wrong number of segments after InitializeSegments")
	}

	if !reflect.DeepEqual(output.Base(), memory.NewRelocatable(0, 0)) {
		t.Errorf("Wrong builtin base after InitializeSegments")
	}
}

func TestOutputInitialStackIncluded(t *testing.T) {
	output := builtins.NewOutputBuiltinRunner()
	output.Include(true)
	initial_stack := output.InitialStack()
	expected_stack := []memory.MaybeRelocatable{*memory.NewMaybeRelocatableRelocatable(output.Base())}
	if !reflect.DeepEqual(initial_stack, expected_stack) {
		t.Errorf("Wrong initial stack")
	}
}

func TestOutputInitialStackNotIncluded(t *testing.T) {
	output := builtins.NewOutputBuiltinRunner()
	if len(output.InitialStack()) != 0 {
		t.Errorf("Initial stack should be empty if not included")
	}
}

func TestOutputAddValidationRule(t *testing.T) {
	empty_mem := memory.NewMemory()
	mem := memory.NewMemory()
	output := builtins.NewOutputBuiltinRunner()
	output.AddValidationRule(mem)
	// Check that the memory is equal to a newly created one to check that
	// no validation rules were added
	if !reflect.DeepEqual(mem, empty_mem) {
		t.Errorf("AddValidationRule should do nothing")
	}
}

func TestGetAllocatedMemoryUnitsOutput(t *testing.T) {
	output := builtins.NewOutputBuiltinRunner()
	vm := vm.NewVirtualMachine()

	mem_units, err := output.GetAllocatedMemoryUnits(&vm.Segments, vm.CurrentStep)
	if err != nil {
		t.Error("test failed with error: ", err)
	}
	if mem_units != 0 {
		t.Errorf("expected memory units to be 5, got: %d", mem_units)
	}
}
