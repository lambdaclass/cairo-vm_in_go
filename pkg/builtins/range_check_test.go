package builtins_test

import (
	"testing"

	"github.com/lambdaclass/cairo-vm.go/pkg/builtins"
	"github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
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

func TestGetRangeCheckUsageSuccessfulA(t *testing.T) {
	var builtin = builtins.DefaultRangeCheckBuiltinRunner()
	builtin.Include(true)
	builtin.SetBase(memory.NewRelocatable(0, 0))

	var memoryManager = memory.NewMemorySegmentManager()
	memoryManager.AddSegment()

	memoryManager.Memory.Insert(memory.NewRelocatable(0, 0), memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(1)))
	memoryManager.Memory.Insert(memory.NewRelocatable(0, 1), memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(2)))
	memoryManager.Memory.Insert(memory.NewRelocatable(0, 2), memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(3)))
	memoryManager.Memory.Insert(memory.NewRelocatable(0, 3), memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(4)))

	var resultMin, resultMax = builtin.GetRangeCheckUsage(&memoryManager.Memory)

	if *resultMin != 0 {
		t.Errorf("rcMin should return 0, got %d", *resultMin)
	}

	if *resultMax != 4 {
		t.Errorf("rcMax should return 4, got %d", *resultMax)
	}
}

func TestGetRangeCheckUsageSuccessfulB(t *testing.T) {
	var builtin = builtins.DefaultRangeCheckBuiltinRunner()
	builtin.Include(true)
	builtin.SetBase(memory.NewRelocatable(0, 0))

	var memoryManager = memory.NewMemorySegmentManager()
	memoryManager.AddSegment()

	memoryManager.Memory.Insert(memory.NewRelocatable(0, 0), memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(1465218365)))
	memoryManager.Memory.Insert(memory.NewRelocatable(0, 1), memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(2134570341)))
	memoryManager.Memory.Insert(memory.NewRelocatable(0, 2), memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(uint64(31349610736))))
	memoryManager.Memory.Insert(memory.NewRelocatable(0, 3), memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(uint64(413468326585859))))

	var resultMin, resultMax = builtin.GetRangeCheckUsage(&memoryManager.Memory)

	if *resultMin != 0 {
		t.Errorf("rcMin should return 0, got %d", *resultMin)
	}

	if *resultMax != 62821 {
		t.Errorf("rcMax should return 62821, got %d", *resultMax)
	}
}

func TestGetRangeCheckUsageSuccessfulC(t *testing.T) {
	var builtin = builtins.DefaultRangeCheckBuiltinRunner()
	builtin.Include(true)
	builtin.SetBase(memory.NewRelocatable(0, 0))

	var memoryManager = memory.NewMemorySegmentManager()
	memoryManager.AddSegment()

	memoryManager.Memory.Insert(memory.NewRelocatable(0, 0), memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(uint64(634834751465218365))))
	memoryManager.Memory.Insert(memory.NewRelocatable(0, 1), memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(uint64(42876922134570341))))
	memoryManager.Memory.Insert(memory.NewRelocatable(0, 2), memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(uint64(uint64(23469831349610736)))))
	memoryManager.Memory.Insert(memory.NewRelocatable(0, 3), memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromHex("145B08DB55D996603")))
	memoryManager.Memory.Insert(memory.NewRelocatable(0, 4), memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromHex("415A2EB28BD916266")))
	memoryManager.Memory.Insert(memory.NewRelocatable(0, 5), memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromHex("44D0E30CF1E224448BCBD11549C")))

	var resultMin, resultMax = builtin.GetRangeCheckUsage(&memoryManager.Memory)

	if *resultMin != 0 {
		t.Errorf("rcMin should return 0, got %d", *resultMin)
	}

	if *resultMax != 61576 {
		t.Errorf("rcMax should return 61576, got %d", *resultMax)
	}
}

func TestGetRangeCheckUsageEmptyMemory(t *testing.T) {
	var builtin = builtins.DefaultRangeCheckBuiltinRunner()
	builtin.Include(true)
	builtin.SetBase(memory.NewRelocatable(0, 0))

	var memoryManager = memory.NewMemorySegmentManager()

	var resultMin, resultMax = builtin.GetRangeCheckUsage(&memoryManager.Memory)

	if resultMin != nil {
		t.Errorf("rcMin should return nil, got %d", *resultMin)
	}

	if resultMax != nil {
		t.Errorf("rcMax should return nil, got %d", *resultMax)
	}
}

// Range check bound is calculated via the constant RANGE_CHECK_N_PARTS.
// If something changes and the bound is set to zero, there could be unexpected errors.
func TestBoundIsNotZero(t *testing.T) {
	rangeCheck := builtins.DefaultRangeCheckBuiltinRunner()

	bound := rangeCheck.Bound()

	if bound.IsZero() {
		t.Error("range check bound should never be zero")
	}
}
