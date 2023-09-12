package builtins_test

import (
	"testing"

	"github.com/lambdaclass/cairo-vm.go/pkg/builtins"
	"github.com/lambdaclass/cairo-vm.go/pkg/vm"
	"github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
)

func TestBaseRangeCheck(t *testing.T) {
	check_range := builtins.DefaultBitwiseBuiltinRunner()
	if check_range.Base() != memory.NewRelocatable(0, 0) {
		t.Errorf("Wrong base value in %s builtin", check_range.Name())
	}
}

func TestInitializeSegmentsForRangeCheck(t *testing.T) {
	check_range := builtins.DefaultBitwiseBuiltinRunner()
	segments := memory.NewMemorySegmentManager()
	check_range.InitializeSegments(&segments)
	if check_range.Base() != memory.NewRelocatable(0, 0) {
		t.Errorf("Builtin %s base is not 0", check_range.Name())
	}
}

func TestGetInitialStackForRangeCheckWithBase(t *testing.T) {
	check_range := builtins.DefaultBitwiseBuiltinRunner()
	check_range.Include(true)
	initialStack := check_range.InitialStack()
	stackValue := initialStack[0]
	expectedValue := memory.NewMaybeRelocatableRelocatable(check_range.Base())
	if !stackValue.IsEqual(expectedValue) {
		t.Errorf("Wrong stack value in %s builtin", check_range.Name())
	}
}

func TestDeduceMemoryCellRangeCheck(t *testing.T) {
	check_range := builtins.DefaultBitwiseBuiltinRunner()
	a, b := check_range.DeduceMemoryCell(memory.NewRelocatable(0, 0), memory.NewMemory())
	if a != nil || b != nil {
		t.Errorf("Deduce memory cell on %s builtin should return (nil, nil)", check_range.Name())
	}
}

// #[test]
// #[cfg_attr(target_arch = "wasm32", wasm_bindgen_test)]
// fn get_allocated_memory_units_range_check() {
// 	let builtin = BuiltinRunner::RangeCheck(RangeCheckBuiltinRunner::new(Some(8), 8, true));
// 	let mut vm = vm!();
// 	vm.current_step = 8;
// 	assert_eq!(builtin.get_allocated_memory_units(&vm), Ok(1));
// }

func TestGetAllocatedMemoryUnitsRangeCheck(t *testing.T) {
	range_check := builtins.DefaultRangeCheckBuiltinRunner()
	vm := vm.NewVirtualMachine()
	vm.CurrentStep = 8
	mem_units, err := range_check.GetAllocatedMemoryUnits(&vm.Segments, vm.CurrentStep)
	if err != nil {
		t.Error("test failed with error: ", err)
	}
	if mem_units != 1 {
		t.Errorf("expected memory units to be 1, got: %d", mem_units)
	}
}
