package builtins_test

import (
	"reflect"
	"testing"

	"github.com/lambdaclass/cairo-vm.go/pkg/builtins"
	"github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
	"github.com/lambdaclass/cairo-vm.go/pkg/vm"
	"github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
)

func TestPoseidonDeduceMemoryCellValid(t *testing.T) {
	poseidon := builtins.NewPoseidonBuiltinRunner(256)
	vmachine := vm.NewVirtualMachine()
	vmachine.BuiltinRunners = append(vmachine.BuiltinRunners, poseidon)

	// Insert input cells into memory
	vmachine.Segments.AddSegment()
	vmachine.Segments.Memory.Insert(
		memory.NewRelocatable(0, 0),
		memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromHex("0x268c44203f1c763bca21beb5aec78b9063cdcdd0fdf6b598bb8e1e8f2b6253f")),
	)
	vmachine.Segments.Memory.Insert(
		memory.NewRelocatable(0, 1),
		memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromHex("0x2b85c9f686f5d3036db55b2ca58a763a3065bc1bc8efbe0e70f3a7171f6cad3")),
	)
	vmachine.Segments.Memory.Insert(
		memory.NewRelocatable(0, 2),
		memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromHex("0x61df3789eef0e1ee0dbe010582a00dd099191e6395dfb976e7be3be2fa9d54b")),
	)

	addr := memory.NewRelocatable(0, 5)
	expected_last_output_cell := memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromHex("0x749d4d0ddf41548e039f183b745a08b80fad54e9ac389021148350bdda70a92"))

	val, err := vmachine.DeduceMemoryCell(addr)
	if !reflect.DeepEqual(val, expected_last_output_cell) || err != nil {
		t.Errorf("Wrong values returned by DeduceMemoryCell")
	}
}

func TestPoseidonDeduceMemoryCellNoInputCells(t *testing.T) {
	poseidon := builtins.NewPoseidonBuiltinRunner(256)
	vmachine := vm.NewVirtualMachine()
	vmachine.BuiltinRunners = append(vmachine.BuiltinRunners, poseidon)
	addr := memory.NewRelocatable(0, 5)
	val, err := vmachine.DeduceMemoryCell(addr)
	if val != nil || err == nil {
		t.Errorf("Wrong values returned by DeduceMemoryCell")
	}
}
func TestPoseidonDeduceMemoryCellInputCell(t *testing.T) {
	poseidon := builtins.NewPoseidonBuiltinRunner(256)
	vmachine := vm.NewVirtualMachine()
	vmachine.BuiltinRunners = append(vmachine.BuiltinRunners, poseidon)
	addr := memory.NewRelocatable(0, 1)
	val, err := vmachine.DeduceMemoryCell(addr)
	if val != nil || err != nil {
		t.Errorf("Wrong values returned by DeduceMemoryCell")
	}
}

func TestPoseidonInitializeSegments(t *testing.T) {
	mem_manager := memory.NewMemorySegmentManager()
	poseidon := builtins.NewPoseidonBuiltinRunner(256)
	poseidon.InitializeSegments(&mem_manager)

	if mem_manager.Memory.NumSegments() != 1 {
		t.Errorf("Wrong number of segments after InitializeSegments")
	}
	if !reflect.DeepEqual(poseidon.Base(), memory.Relocatable{SegmentIndex: 0, Offset: 0}) {
		t.Errorf("Wrong builtin base after InitializeSegments")
	}

}

func TestPoseidonInitialStackIncluded(t *testing.T) {
	poseidon := builtins.NewPoseidonBuiltinRunner(256)
	poseidon.Include(true)
	initial_stack := poseidon.InitialStack()
	expected_stack := []memory.MaybeRelocatable{*memory.NewMaybeRelocatableRelocatable(poseidon.Base())}
	if !reflect.DeepEqual(initial_stack, expected_stack) {
		t.Errorf("Wrong initial stack")
	}
}

func TestPoseidonInitialStackNotIncluded(t *testing.T) {
	poseidon := builtins.NewPoseidonBuiltinRunner(256)
	if len(poseidon.InitialStack()) != 0 {
		t.Errorf("Initial stack should be empty if not included")
	}
}

func TestPoseidonAddValidationRule(t *testing.T) {
	empty_mem := memory.NewMemory()
	mem := memory.NewMemory()
	poseidon := builtins.NewPoseidonBuiltinRunner(256)
	poseidon.AddValidationRule(mem)
	// Check that the memory is equal to a newly created one to check that
	// no validation rules were added
	if !reflect.DeepEqual(mem, empty_mem) {
		t.Errorf("AddValidationRule should do nothing")
	}

}
