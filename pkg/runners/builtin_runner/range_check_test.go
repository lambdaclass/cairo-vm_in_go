package builtinrunner

import (
	"testing"

	"github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
)

func TestBaseRangeCheck(t *testing.T) {
	builtin := NewRangeCheckBuiltinRunner(true)
	if builtin.base != memory.NewRelocatable(0, 0) {
		t.Errorf("Wrong base value in %s builtin", builtin.Name())
	}
}

func TestInitializeSegmentsForRangeCheck(t *testing.T) {
	builtin := NewRangeCheckBuiltinRunner(true)
	segments := memory.NewMemorySegmentManager()
	builtin.InitializeSegments(&segments)
	if builtin.base != memory.NewRelocatable(0, 0) {
		t.Errorf("Builtin %s base is not 0", builtin.Name())
	}
}

func TestGetInitialStackForRangeCheckWithBase(t *testing.T) {
	builtin := NewRangeCheckBuiltinRunner(true)
	builtin.base = memory.NewRelocatable(1, 0)
	initialStack := builtin.InitialStack()
	stackValue := initialStack[0]
	expectedValue := memory.NewMaybeRelocatableRelocatable(builtin.base)
	if !stackValue.IsEqual(expectedValue) {
		t.Errorf("Wrong stack value in %s builtin", builtin.Name())
	}
}

func TestDeduceMemoryCellRangeCheck(t *testing.T) {
	builtin := NewRangeCheckBuiltinRunner(true)
	a, b := builtin.DeduceMemoryCell(memory.NewRelocatable(0, 0), memory.NewMemory())
	if a != nil || b != nil {
		t.Errorf("Deduce memory cell on %s builtin should return (nil, nil)", builtin.Name())
	}
}
