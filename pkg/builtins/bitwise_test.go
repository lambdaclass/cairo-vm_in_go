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
