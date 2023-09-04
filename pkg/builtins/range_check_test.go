package builtins_test

import (
	"testing"

	"github.com/lambdaclass/cairo-vm.go/pkg/builtins"
	"github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
)

func TestBaseRangeCheck(t *testing.T) {
	check_range := builtins.NewRangeCheckBuiltinRunner(true)
	if check_range.Base() != memory.NewRelocatable(0, 0) {
		t.Errorf("Wrong base value in %s builtin", check_range.Name())
	}
}

func TestInitializeSegmentsForRangeCheck(t *testing.T) {
	check_range := builtins.NewRangeCheckBuiltinRunner(true)
	segments := memory.NewMemorySegmentManager()
	check_range.InitializeSegments(&segments)
	if check_range.Base() != memory.NewRelocatable(0, 0) {
		t.Errorf("Builtin %s base is not 0", check_range.Name())
	}
}

func TestGetInitialStackForRangeCheckWithBase(t *testing.T) {
	check_range := builtins.NewRangeCheckBuiltinRunner(true)
	initialStack := check_range.InitialStack()
	stackValue := initialStack[0]
	expectedValue := memory.NewMaybeRelocatableRelocatable(check_range.Base())
	if !stackValue.IsEqual(expectedValue) {
		t.Errorf("Wrong stack value in %s builtin", check_range.Name())
	}
}

func TestDeduceMemoryCellRangeCheck(t *testing.T) {
	check_range := builtins.NewRangeCheckBuiltinRunner(true)
	a, b := check_range.DeduceMemoryCell(memory.NewRelocatable(0, 0), memory.NewMemory())
	if a != nil || b != nil {
		t.Errorf("Deduce memory cell on %s builtin should return (nil, nil)", check_range.Name())
	}
}
