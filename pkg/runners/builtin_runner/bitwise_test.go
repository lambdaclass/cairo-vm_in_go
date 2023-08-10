package builtinrunner

import (
	"testing"

	"github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
	"github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
)

func TestDeduceMemoryCellBitwiseForPresetMemoryValidAnd(t *testing.T) {
	mem := memory.NewMemorySegmentManager()
	mem.AddSegment()
	rel1 := memory.NewRelocatable(0, 5)
	rel2 := memory.NewRelocatable(0, 6)
	rel3 := memory.NewRelocatable(0, 7)

	m1 := memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(10))
	m2 := memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(12))
	m3 := memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(0))

	mem.Memory.Insert(rel1, m1)
	mem.Memory.Insert(rel2, m2)
	mem.Memory.Insert(rel3, m3)

	var ratio uint = 256

	builtin_instance := BitwiseInstanceDef{Ratio: &ratio, TotalNBits: 251}
	builtin := NewBitwiseBuiltinRunner(builtin_instance, true)

	address := memory.NewRelocatable(0, 7)

	result, err := builtin.DeduceMemoryCell(address, &mem.Memory)
	expected := memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(8))

	if err != nil {
		t.Errorf("%v", err)
	}

	if *result != *expected {
		t.Errorf("TestDeduceMemoryCellBitwiseForPresetMemoryValidAnd Failed, expected %v, got %v", expected, result)
	}

}

func TestDeduceMemoryCellBitwiseForPresetMemoryValidXor(t *testing.T) {
	mem := memory.NewMemorySegmentManager()
	mem.AddSegment()
	rel1 := memory.NewRelocatable(0, 5)
	rel2 := memory.NewRelocatable(0, 6)
	rel3 := memory.NewRelocatable(0, 8)

	m1 := memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(10))
	m2 := memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(12))
	m3 := memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(0))

	mem.Memory.Insert(rel1, m1)
	mem.Memory.Insert(rel2, m2)
	mem.Memory.Insert(rel3, m3)

	var ratio uint = 256
	builtin_instance := BitwiseInstanceDef{Ratio: &ratio, TotalNBits: 251}
	builtin := NewBitwiseBuiltinRunner(builtin_instance, true)

	address := memory.NewRelocatable(0, 8)

	result, err := builtin.DeduceMemoryCell(address, &mem.Memory)
	expected := memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(6))

	if err != nil {
		t.Errorf("%v", err)
	}

	if *result != *expected {
		t.Errorf("TestDeduceMemoryCellBitwiseForPresetMemoryValidAnd Failed, expected %v, got %v", expected, result)
	}

}

func TestDeduceMemoryCellBitwiseForPresetMemoryValidOr(t *testing.T) {
	mem := memory.NewMemorySegmentManager()
	mem.AddSegment()
	rel1 := memory.NewRelocatable(0, 5)
	rel2 := memory.NewRelocatable(0, 6)
	rel3 := memory.NewRelocatable(0, 9)

	m1 := memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(10))
	m2 := memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(12))
	m3 := memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(0))

	mem.Memory.Insert(rel1, m1)
	mem.Memory.Insert(rel2, m2)
	mem.Memory.Insert(rel3, m3)

	var ratio uint = 256
	builtin_instance := BitwiseInstanceDef{Ratio: &ratio, TotalNBits: 251}
	builtin := NewBitwiseBuiltinRunner(builtin_instance, true)

	address := memory.NewRelocatable(0, 9)

	result, err := builtin.DeduceMemoryCell(address, &mem.Memory)
	expected := memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(14))

	if err != nil {
		t.Errorf("%v", err)
	}

	if *result != *expected {
		t.Errorf("TestDeduceMemoryCellBitwiseForPresetMemoryValidAnd Failed, expected %v, got %v", expected, result)
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

	var ratio uint = 256
	builtin_instance := BitwiseInstanceDef{Ratio: &ratio, TotalNBits: 251}
	builtin := NewBitwiseBuiltinRunner(builtin_instance, true)

	address := memory.NewRelocatable(0, 5)

	result, err := builtin.DeduceMemoryCell(address, &mem.Memory)

	if err != nil {
		t.Errorf("%v", err)
	}

	if result != nil {
		t.Errorf("TestDeduceMemoryCellBitwiseForPresetMemoryValidAnd Failed, expected nil, got %v", result)
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

	var ratio uint = 256
	builtin_instance := BitwiseInstanceDef{Ratio: &ratio, TotalNBits: 251}
	builtin := NewBitwiseBuiltinRunner(builtin_instance, true)

	address := memory.NewRelocatable(0, 5)

	result, err := builtin.DeduceMemoryCell(address, &mem.Memory)

	if err != nil {
		t.Errorf("%v", err)
	}

	if result != nil {
		t.Errorf("TestDeduceMemoryCellBitwiseForPresetMemoryValidAnd Failed, expected nil, got %v", result)
	}

}
