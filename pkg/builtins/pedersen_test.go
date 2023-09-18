package builtins_test

import (
	"reflect"
	"testing"

	"github.com/lambdaclass/cairo-vm.go/pkg/builtins"
	"github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
	"github.com/lambdaclass/cairo-vm.go/pkg/vm"
	"github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
)

func TestPedersenDeduceMemoryCell(t *testing.T) {
	pedersen := builtins.NewPedersenBuiltinRunner(256)
	vmachine := vm.NewVirtualMachine()
	vmachine.BuiltinRunners = append(vmachine.BuiltinRunners, pedersen)
	// Insert input cells into memory
	vmachine.Segments.AddSegment()
	vmachine.Segments.Memory.Insert(
		memory.NewRelocatable(0, 3),
		memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(32)),
	)
	vmachine.Segments.Memory.Insert(
		memory.NewRelocatable(0, 4),
		memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(72)),
	)
	vmachine.Segments.Memory.Insert(
		memory.NewRelocatable(0, 5),
		memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(0)),
	)

	addr := memory.NewRelocatable(0, 5)
	expected_last_output_cell := memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromHex("0x73b3ec210cccbb970f80c6826fb1c40ae9f487617696234ff147451405c339f"))

	// Deduce Memory Cell
	val, err := vmachine.DeduceMemoryCell(addr)
	if !reflect.DeepEqual(val, expected_last_output_cell) || err != nil {
		t.Errorf("Wrong values returned by DeduceMemoryCell: result %s, expected %s", val.ToString(), expected_last_output_cell.ToString())
	}

	// Computing again the deduce memory cell for the same address, must return nil
	second_compute, err := vmachine.DeduceMemoryCell(addr)
	if second_compute != nil || err != nil {
		t.Errorf("Wrong values returned by DeduceMemoryCell: result %s, expected %v", second_compute.ToString(), nil)
	}

}
