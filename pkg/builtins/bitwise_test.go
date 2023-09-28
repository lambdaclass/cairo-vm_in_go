package builtins_test

import (
	"testing"

	"github.com/lambdaclass/cairo-vm.go/pkg/builtins"
	"github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
	"github.com/lambdaclass/cairo-vm.go/pkg/vm"
	"github.com/lambdaclass/cairo-vm.go/pkg/vm/cairo_run"
	"github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
)

func TestDeduceMemoryCellBitwiseForPresetMemoryValidAnd(t *testing.T) {
	mem := memory.NewMemorySegmentManager()
	mem.AddSegment()
	mem.Memory.Insert(memory.NewRelocatable(0, 5), memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(10)))
	mem.Memory.Insert(memory.NewRelocatable(0, 6), memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(12)))
	mem.Memory.Insert(memory.NewRelocatable(0, 7), memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(0)))

	builtin := builtins.DefaultBitwiseBuiltinRunner()

	address := memory.NewRelocatable(0, 7)

	result, err := builtin.DeduceMemoryCell(address, &mem.Memory)
	if err != nil {
		t.Errorf("TestDeduceMemoryCellBitwiseForPresetMemoryValidAnd failed with error:\n %v", err)
	}
	expected := memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(8))

	if *result != *expected {
		t.Errorf("TestDeduceMemoryCellBitwiseForPresetMemoryValidAnd failed, expected %v, got %v", expected, result)
	}

}

func TestDeduceMemoryCellBitwiseForPresetMemoryValidXor(t *testing.T) {
	mem := memory.NewMemorySegmentManager()
	mem.AddSegment()
	mem.Memory.Insert(memory.NewRelocatable(0, 5), memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(10)))
	mem.Memory.Insert(memory.NewRelocatable(0, 6), memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(12)))
	mem.Memory.Insert(memory.NewRelocatable(0, 8), memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(0)))

	builtin := builtins.DefaultBitwiseBuiltinRunner()

	address := memory.NewRelocatable(0, 8)

	result, err := builtin.DeduceMemoryCell(address, &mem.Memory)

	if err != nil {
		t.Errorf("TestDeduceMemoryCellBitwiseForPresetMemoryValidXor failed with error:\n %v", err)
	}
	expected := memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(6))

	if *result != *expected {
		t.Errorf("TestDeduceMemoryCellBitwiseForPresetMemoryValidXor failed, expected %v, got %v", expected, result)
	}

}

func TestDeduceMemoryCellBitwiseForPresetMemoryValidOr(t *testing.T) {
	mem := memory.NewMemorySegmentManager()
	mem.AddSegment()
	mem.AddSegment()
	mem.Memory.Insert(memory.NewRelocatable(0, 5), memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(10)))
	mem.Memory.Insert(memory.NewRelocatable(0, 6), memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(12)))
	mem.Memory.Insert(memory.NewRelocatable(0, 9), memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(0)))

	builtin := builtins.DefaultBitwiseBuiltinRunner()

	address := memory.NewRelocatable(0, 9)

	result, err := builtin.DeduceMemoryCell(address, &mem.Memory)

	if err != nil {
		t.Errorf("TestDeduceMemoryCellBitwiseForPresetMemoryValidOr failed with error:\n %v", err)
	}
	expected := memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(14))

	if *result != *expected {
		t.Errorf("TestDeduceMemoryCellBitwiseForPresetMemoryValidOr failed, expected %v, got %v", expected, result)
	}

}

func TestDeduceMemoryCellBitwiseForPresetMemoryIncorrectOffset(t *testing.T) {
	mem := memory.NewMemorySegmentManager()
	mem.AddSegment()
	mem.Memory.Insert(memory.NewRelocatable(0, 3), memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(10)))
	mem.Memory.Insert(memory.NewRelocatable(0, 4), memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(12)))
	mem.Memory.Insert(memory.NewRelocatable(0, 5), memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(0)))

	builtin := builtins.DefaultBitwiseBuiltinRunner()

	address := memory.NewRelocatable(0, 5)

	result, err := builtin.DeduceMemoryCell(address, &mem.Memory)

	if err != nil {
		t.Errorf("TestDeduceMemoryCellBitwiseForPresetMemoryIncorrectOffset failed with error:\n %v", err)
	}

	if result != nil {
		t.Errorf("TestDeduceMemoryCellBitwiseForPresetMemoryIncorrectOffset failed, expected nil, got %v", result)
	}

}

func TestDeduceMemoryCellBitwiseForPresetMemoryNoValuesToOperate(t *testing.T) {
	mem := memory.NewMemorySegmentManager()
	mem.AddSegment()
	mem.Memory.Insert(memory.NewRelocatable(0, 5), memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(10)))
	mem.Memory.Insert(memory.NewRelocatable(0, 7), memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(0)))

	builtin := builtins.DefaultBitwiseBuiltinRunner()

	address := memory.NewRelocatable(0, 5)

	result, err := builtin.DeduceMemoryCell(address, &mem.Memory)

	if err != nil {
		t.Errorf("TestDeduceMemoryCellBitwiseForPresetMemoryNoValuesToOperate failed with error:\n %v", err)
	}

	if result != nil {
		t.Errorf("TestDeduceMemoryCellBitwiseForPresetMemoryNoValuesToOperate failed, expected nil, got %v", result)
	}

}

func TestGetAllocatedMemoryUnitsBitwise(t *testing.T) {
	bitwise := builtins.DefaultBitwiseBuiltinRunner()
	bitwise.Include(true)

	vm := vm.NewVirtualMachine()
	vm.CurrentStep = 256
	mem_units, err := bitwise.GetAllocatedMemoryUnits(&vm.Segments, vm.CurrentStep)

	if err != nil {
		t.Error("test failed with error: ", err)
	}

	if mem_units != 5 {
		t.Errorf("expected memory units to be 5, got: %d", mem_units)
	}
}

func TestIntegrationBitwise(t *testing.T) {
	t.Helper()
	cairoRunConfig := cairo_run.CairoRunConfig{DisableTracePadding: false, Layout: "all_cairo", ProofMode: false}
	_, err := cairo_run.CairoRun("../../cairo_programs/bitwise_builtin_test.json", cairoRunConfig)
	if err != nil {
		t.Errorf("TestIntegrationBitwise failed with error:\n %v", err)
	}
}

func TestGetUsedDilutedCheckUnitsA(t *testing.T) {
	builtin := builtins.NewBitwiseBuiltinRunner(256)

	result := builtin.GetUsedDilutedCheckUnits(12, 2)
	var expected uint = 535
	if result != expected {
		t.Errorf("Wrong Value for GetUsedDilutedChecks, should be %d, got %d", expected, result)
	}
}

func TestGetUsedDilutedCheckUnitsB(t *testing.T) {
	builtin := builtins.NewBitwiseBuiltinRunner(256)

	result := builtin.GetUsedDilutedCheckUnits(30, 56)
	var expected uint = 150
	if result != expected {
		t.Errorf("Wrong Value for GetUsedDilutedChecks, should be %d, got %d", expected, result)
	}
}

func TestGetUsedDilutedCheckUnitsC(t *testing.T) {
	builtin := builtins.NewBitwiseBuiltinRunner(256)

	result := builtin.GetUsedDilutedCheckUnits(50, 25)
	var expected uint = 250
	if result != expected {
		t.Errorf("Wrong Value for GetUsedDilutedChecks, should be %d, got %d", expected, result)
	}
}

func TestRunSecurityChecksEmptyMemory(t *testing.T) {
	builtin := builtins.NewBitwiseBuiltinRunner(256)
	segments := memory.NewMemorySegmentManager()
	err := builtin.RunSecurityChecks(&segments)
	if err != nil {
		t.Errorf("RunSecurityChecks failed with error: %s", err.Error())
	}
}

func TestRunSecurityChecksMissingMemoryCells(t *testing.T) {
	builtin := builtins.NewBitwiseBuiltinRunner(256)
	segments := memory.NewMemorySegmentManager()

	builtin.InitializeSegments(&segments)
	builtinBase := builtin.Base()
	// A bitwise cell consists of 5 elements: 2 input cells & 3 output cells
	// In this test we insert cells 4-5 and leave the first input cell empty
	// This will fail the security checks, as the memory cell with offset 0 will be missing
	builtinSegment := []memory.MaybeRelocatable{
		*memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint(1)),
		*memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint(2)),
		*memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint(3)),
		*memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint(4)),
	}
	segments.LoadData(builtinBase.AddUint(1), &builtinSegment)

	err := builtin.RunSecurityChecks(&segments)
	if err == nil {
		t.Errorf("RunSecurityChecks should have failed")
	}
}

func TestRunSecurityChecksMissingMemoryCellsNCheck(t *testing.T) {
	builtin := builtins.NewBitwiseBuiltinRunner(256)
	segments := memory.NewMemorySegmentManager()

	builtin.InitializeSegments(&segments)
	builtinBase := builtin.Base()
	// n = max(offsets) // cellsPerInstance + 1
	// n = max[(0]) // 5 + 1 = 0 // 5 + 1 = 1
	// len(offsets) // inputCells = 1 // 2
	// This will fail the security checks, as n > len(offsets) // inputCells
	builtinSegment := []memory.MaybeRelocatable{
		*memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint(1)),
	}
	segments.LoadData(builtinBase, &builtinSegment)

	err := builtin.RunSecurityChecks(&segments)
	if err == nil {
		t.Errorf("RunSecurityChecks should have failed")
	}
}
