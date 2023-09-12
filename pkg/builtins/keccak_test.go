package builtins_test

import (
	"reflect"
	"testing"

	"github.com/lambdaclass/cairo-vm.go/pkg/builtins"
	"github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
	"github.com/lambdaclass/cairo-vm.go/pkg/vm"
	"github.com/lambdaclass/cairo-vm.go/pkg/vm/cairo_run"
	"github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
)

func TestKeccakDeduceMemoryCellValid(t *testing.T) {
	keccak := builtins.NewKeccakBuiltinRunner()
	keccak.Include(true)
	vmachine := vm.NewVirtualMachine()
	vmachine.BuiltinRunners = append(vmachine.BuiltinRunners, keccak)

	// Insert input cells into memory
	vmachine.Segments.AddSegment()
	vmachine.Segments.Memory.Insert(
		memory.NewRelocatable(0, 0),
		memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(43)),
	)
	vmachine.Segments.Memory.Insert(
		memory.NewRelocatable(0, 1),
		memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(199)),
	)
	vmachine.Segments.Memory.Insert(
		memory.NewRelocatable(0, 2),
		memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(0)),
	)
	vmachine.Segments.Memory.Insert(
		memory.NewRelocatable(0, 3),
		memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(0)),
	)
	vmachine.Segments.Memory.Insert(
		memory.NewRelocatable(0, 4),
		memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(0)),
	)
	vmachine.Segments.Memory.Insert(
		memory.NewRelocatable(0, 5),
		memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(0)),
	)
	vmachine.Segments.Memory.Insert(
		memory.NewRelocatable(0, 6),
		memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(0)),
	)
	vmachine.Segments.Memory.Insert(
		memory.NewRelocatable(0, 7),
		memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(1)),
	)
	vmachine.Segments.Memory.Insert(
		memory.NewRelocatable(0, 8),
		memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(0)),
	)

	addr := memory.NewRelocatable(0, 9)
	expected_last_output_cell := memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromDecString("1006979841721999878391288827876533441431370448293338267890891"))

	val, err := vmachine.DeduceMemoryCell(addr)
	if !reflect.DeepEqual(val, expected_last_output_cell) || err != nil {
		t.Errorf("Wrong values returned by DeduceMemoryCell")
	}
}

func TestKeccakDeduceMemoryCellNoInputCells(t *testing.T) {
	keccak := builtins.NewKeccakBuiltinRunner()
	keccak.Include(true)
	vmachine := vm.NewVirtualMachine()
	vmachine.BuiltinRunners = append(vmachine.BuiltinRunners, keccak)
	addr := memory.NewRelocatable(0, 10)
	val, err := vmachine.DeduceMemoryCell(addr)
	if val != nil || err == nil {
		t.Errorf("Wrong values returned by DeduceMemoryCell")
	}
}
func TestKeccakDeduceMemoryCellInputCell(t *testing.T) {
	keccak := builtins.NewKeccakBuiltinRunner()
	keccak.Include(true)
	vmachine := vm.NewVirtualMachine()
	vmachine.BuiltinRunners = append(vmachine.BuiltinRunners, keccak)
	addr := memory.NewRelocatable(0, 1)
	val, err := vmachine.DeduceMemoryCell(addr)
	if val != nil || err != nil {
		t.Errorf("Wrong values returned by DeduceMemoryCell")
	}
}

func TestKeccakInitializeSegments(t *testing.T) {
	mem_manager := memory.NewMemorySegmentManager()
	keccak := builtins.NewKeccakBuiltinRunner()
	keccak.Include(true)
	keccak.InitializeSegments(&mem_manager)

	if mem_manager.Memory.NumSegments() != 1 {
		t.Errorf("Wrong number of segments after InitializeSegments")
	}
	if !reflect.DeepEqual(keccak.Base(), memory.Relocatable{SegmentIndex: 0, Offset: 0}) {
		t.Errorf("Wrong builtin base after InitializeSegments")
	}

}

func TestKeccakInitialStackIncluded(t *testing.T) {
	keccak := builtins.NewKeccakBuiltinRunner()
	keccak.Include(true)
	initial_stack := keccak.InitialStack()
	expected_stack := []memory.MaybeRelocatable{*memory.NewMaybeRelocatableRelocatable(keccak.Base())}
	if !reflect.DeepEqual(initial_stack, expected_stack) {
		t.Errorf("Wrong initial stack")
	}
}

func TestKeccakInitialStackNotIncluded(t *testing.T) {
	keccak := builtins.NewKeccakBuiltinRunner()
	keccak.Include(false)
	if len(keccak.InitialStack()) != 0 {
		t.Errorf("Initial stack should be empty if not included")
	}
}

func TestKeccakAddValidationRule(t *testing.T) {
	empty_mem := memory.NewMemory()
	mem := memory.NewMemory()
	keccak := builtins.NewKeccakBuiltinRunner()
	keccak.Include(true)
	keccak.AddValidationRule(mem)
	// Check that the memory is equal to a newly created one to check that
	// no validation rules were added
	if !reflect.DeepEqual(mem, empty_mem) {
		t.Errorf("AddValidationRule should do nothing")
	}

}

func TestIntegrationKeccak(t *testing.T) {
	t.Helper()
	_, err := cairo_run.CairoRun("../../cairo_programs/keccak_builtin.json", "small", false)
	if err != nil {
		t.Errorf("TestIntegrationKeccak failed with error:\n %v", err)
	}
}
