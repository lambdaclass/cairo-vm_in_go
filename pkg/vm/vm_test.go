package vm_test

import (
	"bytes"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/lambdaclass/cairo-vm.go/pkg/vm"
	"github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
)

func TestUpdatePcRegularNoImm(t *testing.T) {
	instruction := vm.Instruction{PcUpdate: vm.PcUpdateRegular, Op1Addr: vm.Op1SrcAP}
	operands := vm.Operands{}
	vm := vm.NewVirtualMachine()
	err := vm.UpdatePc(&instruction, &operands)
	if err != nil {
		t.Errorf("UpdatePc failed with error: %s", err)
	}
	if !reflect.DeepEqual(vm.RunContext.Pc, memory.Relocatable{SegmentIndex: 0, Offset: 1}) {
		t.Errorf("Wrong value after pc update")
	}
}

func TestUpdatePcRegularWithImm(t *testing.T) {
	instruction := vm.Instruction{PcUpdate: vm.PcUpdateRegular, Op1Addr: vm.Op1SrcImm}
	operands := vm.Operands{}
	vm := vm.NewVirtualMachine()
	err := vm.UpdatePc(&instruction, &operands)
	if err != nil {
		t.Errorf("UpdatePc failed with error: %s", err)
	}
	if !reflect.DeepEqual(vm.RunContext.Pc, memory.Relocatable{SegmentIndex: 0, Offset: 2}) {
		t.Errorf("Wrong value after pc update")
	}
}

func TestUpdatePcJumpWithRelRes(t *testing.T) {
	instruction := vm.Instruction{PcUpdate: vm.PcUpdateJump}
	res := memory.Relocatable{SegmentIndex: 0, Offset: 5}
	operands := vm.Operands{Res: memory.NewMaybeRelocatableRelocatable(res)}
	vm := vm.NewVirtualMachine()
	err := vm.UpdatePc(&instruction, &operands)
	if err != nil {
		t.Errorf("UpdatePc failed with error: %s", err)
	}
	if !reflect.DeepEqual(vm.RunContext.Pc, res) {
		t.Errorf("Wrong value after pc update")
	}
}

func TestUpdatePcJumpWithIntRes(t *testing.T) {
	instruction := vm.Instruction{PcUpdate: vm.PcUpdateJump}
	operands := vm.Operands{Res: memory.NewMaybeRelocatableInt(0)}
	vm := vm.NewVirtualMachine()
	err := vm.UpdatePc(&instruction, &operands)
	if err == nil {
		t.Errorf("UpdatePc should have failed")
	}
}

func TestUpdatePcJumpWithoutRes(t *testing.T) {
	instruction := vm.Instruction{PcUpdate: vm.PcUpdateJump}
	operands := vm.Operands{}
	vm := vm.NewVirtualMachine()

	err := vm.UpdatePc(&instruction, &operands)
	if err == nil {
		t.Errorf("UpdatePc should have failed")
	}
}

func TestUpdatePcJumpRelWithIntRes(t *testing.T) {
	instruction := vm.Instruction{PcUpdate: vm.PcUpdateJumpRel}
	operands := vm.Operands{Res: memory.NewMaybeRelocatableInt(5)}
	vm := vm.NewVirtualMachine()
	err := vm.UpdatePc(&instruction, &operands)
	if err != nil {
		t.Errorf("UpdatePc failed with error: %s", err)
	}
	if !reflect.DeepEqual(vm.RunContext.Pc, memory.Relocatable{SegmentIndex: 0, Offset: 5}) {
		t.Errorf("Wrong value after pc update")
	}
}

func TestUpdatePcJumpRelWithRelRes(t *testing.T) {
	instruction := vm.Instruction{PcUpdate: vm.PcUpdateJumpRel}
	res := memory.Relocatable{SegmentIndex: 0, Offset: 5}
	operands := vm.Operands{Res: memory.NewMaybeRelocatableRelocatable(res)}
	vm := vm.NewVirtualMachine()

	err := vm.UpdatePc(&instruction, &operands)
	if err == nil {
		t.Errorf("UpdatePc should have failed")
	}
}

func TestUpdatePcJumpRelNoRes(t *testing.T) {
	instruction := vm.Instruction{PcUpdate: vm.PcUpdateJumpRel}
	operands := vm.Operands{}
	vm := vm.NewVirtualMachine()

	err := vm.UpdatePc(&instruction, &operands)
	if err == nil {
		t.Errorf("UpdatePc should have failed")
	}
}

func TestUpdatePcJnzDstIsZeroNoImm(t *testing.T) {
	instruction := vm.Instruction{PcUpdate: vm.PcUpdateJnz, Op1Addr: vm.Op1SrcAP}
	operands := vm.Operands{Dst: *memory.NewMaybeRelocatableInt(0)}
	vm := vm.NewVirtualMachine()
	err := vm.UpdatePc(&instruction, &operands)
	if err != nil {
		t.Errorf("UpdatePc failed with error: %s", err)
	}
	if !reflect.DeepEqual(vm.RunContext.Pc, memory.Relocatable{SegmentIndex: 0, Offset: 1}) {
		t.Errorf("Wrong value after pc update")
	}
}

func TestUpdatePcJnzDstIsZeroWithImm(t *testing.T) {
	instruction := vm.Instruction{PcUpdate: vm.PcUpdateJnz, Op1Addr: vm.Op1SrcImm}
	operands := vm.Operands{Dst: *memory.NewMaybeRelocatableInt(0)}
	vm := vm.NewVirtualMachine()
	err := vm.UpdatePc(&instruction, &operands)
	if err != nil {
		t.Errorf("UpdatePc failed with error: %s", err)
	}
	if !reflect.DeepEqual(vm.RunContext.Pc, memory.Relocatable{SegmentIndex: 0, Offset: 2}) {
		t.Errorf("Wrong value after pc update")
	}
}

func TestUpdatePcJnzDstNotZeroOp1Int(t *testing.T) {
	instruction := vm.Instruction{PcUpdate: vm.PcUpdateJnz}
	operands := vm.Operands{Dst: *memory.NewMaybeRelocatableInt(1), Op1: *memory.NewMaybeRelocatableInt(3)}
	vm := vm.NewVirtualMachine()
	err := vm.UpdatePc(&instruction, &operands)
	if err != nil {
		t.Errorf("UpdatePc failed with error: %s", err)
	}
	if !reflect.DeepEqual(vm.RunContext.Pc, memory.Relocatable{SegmentIndex: 0, Offset: 3}) {
		t.Errorf("Wrong value after pc update %v", vm.RunContext.Pc)
	}
}

func TestUpdatePcJnzDstNotZeroOp1Rel(t *testing.T) {
	instruction := vm.Instruction{PcUpdate: vm.PcUpdateJnz}
	operands := vm.Operands{Dst: *memory.NewMaybeRelocatableInt(1), Op1: *memory.NewMaybeRelocatableRelocatable(memory.Relocatable{})}
	vm := vm.NewVirtualMachine()
	err := vm.UpdatePc(&instruction, &operands)
	if err == nil {
		t.Errorf("UpdatePc should have failed")
	}
}

// Things we are skipping for now:
// - Initializing hint_executor and passing it to `cairo_run`
// - cairo_run_config stuff
// - Asserting expected trace values
// - Asserting memory_holes
func TestFibonacci(t *testing.T) {
	// compiledProgram := parser.Parse("../../cairo_programs/fibonacci.json")

	// TODO: Uncomment test when we have the bare minimum `CairoRun`
	// err := vm.CairoRun(compiledProgram.Data)
	// if err != nil {
	// 	t.Errorf("Program execution failed with error: %s", err)
	// }
}

func TestRelocateTraceOneEntry(t *testing.T) {
	virtualMachine := vm.NewVirtualMachine()
	buildTestProgramMemory(virtualMachine)

	virtualMachine.Segments.ComputeEffectiveSizes()
	relocationTable, _ := virtualMachine.Segments.RelocateSegments()
	err := virtualMachine.RelocateTrace(&relocationTable)
	if err != nil {
		t.Errorf("Trace relocation error failed with test: %s", err)
	}

	expectedTrace := []vm.RelocatedTraceEntry{{Pc: 1, Ap: 4, Fp: 4}}
	actualTrace, err := virtualMachine.GetRelocatedTrace()
	if err != nil {
		t.Errorf("Trace relocation error failed with test: %s", err)
	}
	if !reflect.DeepEqual(expectedTrace, actualTrace) {
		t.Errorf("Relocated trace and expected trace are not the same")
	}
}

func TestWriteBinaryTraceFile(t *testing.T) {
	tracePath, err := filepath.Abs("../../cairo_programs/struct.trace")
	if err != nil {
		t.Errorf("Trace file writing error failed with test: %s", err)
	}

	expectedTrace, err := os.ReadFile(tracePath)
	if err != nil {
		t.Errorf("Trace file writing error failed with test: %s", err)
	}

	virtualMachine := vm.NewVirtualMachine()
	buildTestProgramMemory(virtualMachine)

	err = virtualMachine.Relocate()
	if err != nil {
		t.Errorf("Trace file writing error failed with test: %s", err)
	}

	relocatedTrace, err := virtualMachine.GetRelocatedTrace()
	if err != nil {
		t.Errorf("Trace file writing error failed with test: %s", err)
	}

	var actualTraceBuffer bytes.Buffer
	vm.WriteEncodedTrace(relocatedTrace, &actualTraceBuffer)

	if !reflect.DeepEqual(expectedTrace, actualTraceBuffer.Bytes()) {
		t.Errorf("Written trace and expected trace are not the same")
	}
}

func TestWriteBinaryMemoryFile(t *testing.T) {
	var relocatedMemory = make(map[uint]uint)
	relocatedMemory[1] = 66
	relocatedMemory[2] = 42
	relocatedMemory[3] = 30

	var actualMemoryBuffer bytes.Buffer
	vm.WriteEncodedMemory(relocatedMemory, &actualMemoryBuffer)
}

func buildTestProgramMemory(virtualMachine *vm.VirtualMachine) {
	virtualMachine.Trace = []vm.TraceEntry{{Pc: memory.NewRelocatable(0, 0), Ap: memory.NewRelocatable(2, 0), Fp: memory.NewRelocatable(2, 0)}}
	for i := 0; i < 4; i++ {
		virtualMachine.Segments.AddSegment()
	}
	virtualMachine.Segments.Memory.Insert(memory.NewRelocatable(0, 0), memory.NewMaybeRelocatableInt(2345108766317314046))
	virtualMachine.Segments.Memory.Insert(memory.NewRelocatable(1, 0), memory.NewMaybeRelocatableRelocatable(memory.NewRelocatable(2, 0)))
	virtualMachine.Segments.Memory.Insert(memory.NewRelocatable(1, 1), memory.NewMaybeRelocatableRelocatable(memory.NewRelocatable(3, 0)))
}
