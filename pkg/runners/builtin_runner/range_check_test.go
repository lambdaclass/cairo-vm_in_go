package builtinrunner

import (
	"testing"

	"github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
)

func TestBaseRangeCheck(t *testing.T) {
	ratio := uint32(8)
	builtin := NewRangeCheckBuiltinRunner(&ratio, 8, true)
	if builtin.base != 0 {
		t.Errorf("Wrong base value in %s builtin", builtin.Name())
	}
}

func TestInitializeSegmentsForRangeCheck(t *testing.T) {
	ratio := uint32(8)
	builtin := NewRangeCheckBuiltinRunner(&ratio, 8, true)
	segments := memory.NewMemorySegmentManager()
	builtin.InitializeSegments(&segments)
	if builtin.base != 0 {
		t.Errorf("Builtin %s base is not 0", builtin.Name())
	}
}

func TestGetInitialStackForRangeCheckWithBase(t *testing.T) {
	ratio := uint32(8)
	builtin := NewRangeCheckBuiltinRunner(&ratio, 8, true)
	builtin.base = 1
	initialStack := builtin.InitialStack()
	stackValue := initialStack[0]
	expectedValue := memory.NewMaybeRelocatableRelocatable(memory.NewRelocatable(builtin.base, 0))
	if !stackValue.IsEqual(expectedValue) {
		t.Errorf("Wrong stack value in %s builtin", builtin.Name())
	}
}

func TestDeduceMemoryCellRangeCheck(t *testing.T) {
	ratio := uint32(8)
	builtin := NewRangeCheckBuiltinRunner(&ratio, 8, true)
	a, b := builtin.DeduceMemoryCell(memory.NewRelocatable(0, 0), memory.NewMemory())
	if a != nil || b != nil {
		t.Errorf("Deduce memory cell on %s builtin should return (nil, nil)", builtin.Name())
	}
}
