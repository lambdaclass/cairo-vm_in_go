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

	builtin := builtins.NewBitwiseBuiltinRunner(true)

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

	builtin := builtins.NewBitwiseBuiltinRunner(true)

	address := memory.NewRelocatable(0, 8)

	result, err := builtin.DeduceMemoryCell(address, &mem.Memory)

	if err != nil {
		t.Errorf("Test failed with error: %v", err)
	}
	expected := memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(6))

	if *result != *expected {
		t.Errorf("TestDeduceMemoryCellBitwiseForPresetMemoryValidAnd failed, expected %v, got %v", expected, result)
	}

}

func TestDeduceMemoryCellBitwiseForPresetMemoryValidOr(t *testing.T) {
	mem := memory.NewMemorySegmentManager()
	mem.AddSegment()
	mem.AddSegment()
	mem.Memory.Insert(memory.NewRelocatable(0, 5), memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(10)))
	mem.Memory.Insert(memory.NewRelocatable(0, 6), memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(12)))

	builtin := builtins.NewBitwiseBuiltinRunner(true)

	address := memory.NewRelocatable(0, 9)

	result, err := builtin.DeduceMemoryCell(address, &mem.Memory)

	if err != nil {
		t.Errorf("Test failed with error: %v", err)
	}
	expected := memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(14))

	if *result != *expected {
		t.Errorf("TestDeduceMemoryCellBitwiseForPresetMemoryValidAnd failed, expected %v, got %v", expected, result)
	}

}

func TestDeduceMemoryCellBitwiseForPresetMemoryIncorrectOffset(t *testing.T) {
	mem := memory.NewMemorySegmentManager()
	mem.AddSegment()
	rel1 := memory.NewRelocatable(0, 3)
	rel2 := memory.NewRelocatable(0, 4)
	rel3 := memory.NewRelocatable(0, 5)

	m1 := memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(10))
	m2 := memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(12))
	m3 := memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(0))

	mem.Memory.Insert(rel1, m1)
	mem.Memory.Insert(rel2, m2)
	mem.Memory.Insert(rel3, m3)

	builtin := builtins.NewBitwiseBuiltinRunner(true)

	address := memory.NewRelocatable(0, 5)

	result, err := builtin.DeduceMemoryCell(address, &mem.Memory)

	if err != nil {
		t.Errorf("%v", err)
	}

	if result != nil {
		t.Errorf("TestDeduceMemoryCellBitwiseForPresetMemoryValidAnd failed, expected nil, got %v", result)
	}

}

func TestDeduceMemoryCellBitwiseForPresetMemoryNoValuesToOperate(t *testing.T) {
	mem := memory.NewMemorySegmentManager()
	mem.AddSegment()
	rel1 := memory.NewRelocatable(0, 5)
	rel2 := memory.NewRelocatable(0, 7)

	m1 := memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(10))
	m2 := memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(0))

	mem.Memory.Insert(rel1, m1)
	mem.Memory.Insert(rel2, m2)

	builtin := builtins.NewBitwiseBuiltinRunner(true)

	address := memory.NewRelocatable(0, 5)

	result, err := builtin.DeduceMemoryCell(address, &mem.Memory)

	if err != nil {
		t.Errorf("%v", err)
	}

	if result != nil {
		t.Errorf("TestDeduceMemoryCellBitwiseForPresetMemoryValidAnd failed, expected nil, got %v", result)
	}

}

func TestIntegrationBitwise(t *testing.T) {
	t.Helper()
	_, err := cairo_run.CairoRun("../../cairo_programs/bitwise_builtin_test.json")
	if err != nil {
		t.Errorf("Test failed with error: %v", err)
	}
}
