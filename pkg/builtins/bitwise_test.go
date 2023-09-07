package builtins_test

import (
	"testing"

	"github.com/lambdaclass/cairo-vm.go/pkg/builtins"
	"github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
	"github.com/lambdaclass/cairo-vm.go/pkg/vm/cairo_run"
	"github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
)

func TestDeduceMemoryCellBitwiseForPresetMemoryValidAnd(t *testing.T) {
	mem := memory.NewMemorySegmentManager()
	mem.AddSegment()
	mem.Memory.Insert(memory.NewRelocatable(0, 5), memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(10)))
	mem.Memory.Insert(memory.NewRelocatable(0, 6), memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(12)))
	mem.Memory.Insert(memory.NewRelocatable(0, 7), memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(0)))

	builtin := builtins.NewBitwiseBuiltinRunner()

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

	builtin := builtins.NewBitwiseBuiltinRunner()

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

	builtin := builtins.NewBitwiseBuiltinRunner()

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

	builtin := builtins.NewBitwiseBuiltinRunner()

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

	builtin := builtins.NewBitwiseBuiltinRunner()

	address := memory.NewRelocatable(0, 5)

	result, err := builtin.DeduceMemoryCell(address, &mem.Memory)

	if err != nil {
		t.Errorf("TestDeduceMemoryCellBitwiseForPresetMemoryNoValuesToOperate failed with error:\n %v", err)
	}

	if result != nil {
		t.Errorf("TestDeduceMemoryCellBitwiseForPresetMemoryNoValuesToOperate failed, expected nil, got %v", result)
	}

}

func TestIntegrationBitwise(t *testing.T) {
	t.Helper()
	_, err := cairo_run.CairoRun("../../cairo_programs/bitwise_builtin_test.json", "small", false)
	if err != nil {
		t.Errorf("TestIntegrationBitwise failed with error:\n %v", err)
	}
}
