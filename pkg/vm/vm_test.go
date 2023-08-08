package vm_test

import (
	"bytes"

	"reflect"
	"testing"

	"github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
	"github.com/lambdaclass/cairo-vm.go/pkg/vm"
	"github.com/lambdaclass/cairo-vm.go/pkg/vm/cairo_run"
	"github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
)

func TestDeduceOp0OpcodeRet(t *testing.T) {
	instruction := vm.Instruction{Opcode: vm.Ret}
	vm := vm.NewVirtualMachine()
	op0, res, err := vm.DeduceOp0(&instruction, nil, nil)
	if err != nil {
		t.Errorf("DeduceOp0 failed with error: %s", err)
	}
	if op0 != nil || res != nil {
		t.Errorf("Wrong values returned by DeduceOp0")
	}
}
func TestDeduceOp0OpcodeAssertEqResMulOk(t *testing.T) {
	instruction := vm.Instruction{Opcode: vm.AssertEq, ResLogic: vm.ResMul}
	vm := vm.NewVirtualMachine()
	dst := memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(6))
	op1 := memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(3))
	op0, res, err := vm.DeduceOp0(&instruction, dst, op1)
	if err != nil {
		t.Errorf("DeduceOp0 failed with error: %s", err)
	}
	if !reflect.DeepEqual(op0, memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(2))) || !reflect.DeepEqual(res, dst) {
		t.Errorf("Wrong values returned by DeduceOp0")
	}
}

func TestDeduceOp0OpcodeAssertEqResMulZeroDiv(t *testing.T) {
	instruction := vm.Instruction{Opcode: vm.AssertEq, ResLogic: vm.ResMul}
	vm := vm.NewVirtualMachine()
	dst := memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(6))
	op1 := memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(0))
	op0, res, err := vm.DeduceOp0(&instruction, dst, op1)
	if op0 != nil || res != nil || err != nil {
		t.Errorf("Wrong values returned by DeduceOp0")
	}
}

func TestDeduceOp0OpcodeAssertEqResMulRelValues(t *testing.T) {
	instruction := vm.Instruction{Opcode: vm.AssertEq, ResLogic: vm.ResMul}
	vm := vm.NewVirtualMachine()
	dst := memory.NewMaybeRelocatableRelocatable(memory.Relocatable{})
	op1 := memory.NewMaybeRelocatableRelocatable(memory.Relocatable{})
	op0, res, err := vm.DeduceOp0(&instruction, dst, op1)
	if op0 != nil || res != nil || err != nil {
		t.Errorf("Wrong values returned by DeduceOp0")
	}
}

func TestDeduceOp0OpcodeAssertEqResMulNilValues(t *testing.T) {
	instruction := vm.Instruction{Opcode: vm.AssertEq, ResLogic: vm.ResAdd}
	vm := vm.NewVirtualMachine()
	op0, res, err := vm.DeduceOp0(&instruction, nil, nil)
	if op0 != nil || res != nil || err != nil {
		t.Errorf("Wrong values returned by DeduceOp0")
	}
}

func TestDeduceOp0OpcodeAssertEqResAddOk(t *testing.T) {
	instruction := vm.Instruction{Opcode: vm.AssertEq, ResLogic: vm.ResAdd}
	vm := vm.NewVirtualMachine()
	dst := memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(7))
	op1 := memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(5))
	op0, res, err := vm.DeduceOp0(&instruction, dst, op1)
	if err != nil {
		t.Errorf("DeduceOp0 failed with error: %s", err)
	}
	if !reflect.DeepEqual(op0, memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(2))) || !reflect.DeepEqual(res, dst) {
		t.Errorf("Wrong values returned by DeduceOp0")
	}
}

func TestDeduceOp0OpcodeAssertEqResAddRelValues(t *testing.T) {
	instruction := vm.Instruction{Opcode: vm.AssertEq, ResLogic: vm.ResAdd}
	vm := vm.NewVirtualMachine()
	dst := memory.NewMaybeRelocatableRelocatable(memory.Relocatable{SegmentIndex: 1, Offset: 6})
	op1 := memory.NewMaybeRelocatableRelocatable(memory.Relocatable{SegmentIndex: 1, Offset: 2})
	op0, res, err := vm.DeduceOp0(&instruction, dst, op1)
	if err != nil {
		t.Errorf("DeduceOp0 failed with error: %s", err)
	}
	if !reflect.DeepEqual(op0, memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(4))) || !reflect.DeepEqual(res, dst) {
		t.Errorf("Wrong values returned by DeduceOp0")
	}
}

func TestDeduceOp0OpcodeAssertEqResAddNilValues(t *testing.T) {
	instruction := vm.Instruction{Opcode: vm.AssertEq, ResLogic: vm.ResAdd}
	vm := vm.NewVirtualMachine()
	op0, res, err := vm.DeduceOp0(&instruction, nil, nil)
	if op0 != nil || res != nil || err != nil {
		t.Errorf("Wrong values returned by DeduceOp0")
	}
}

func TestDeduceOp0OpcodeCall(t *testing.T) {
	instruction := vm.Instruction{Opcode: vm.Call, Op1Addr: vm.Op1SrcAP}
	vm := vm.NewVirtualMachine()
	vm.RunContext.Pc = memory.Relocatable{SegmentIndex: 1, Offset: 7}
	op0, res, err := vm.DeduceOp0(&instruction, nil, nil)
	if err != nil {
		t.Errorf("DeduceOp0 failed with error: %s", err)
	}
	if !reflect.DeepEqual(op0, memory.NewMaybeRelocatableRelocatable(memory.Relocatable{SegmentIndex: 1, Offset: 8})) || res != nil {
		t.Errorf("Wrong values returned by DeduceOp0")
	}
}

func TestUpdateRegistersAllRegularNoImm(t *testing.T) {
	instruction := vm.Instruction{FpUpdate: vm.FpUpdateRegular, ApUpdate: vm.ApUpdateRegular, PcUpdate: vm.PcUpdateRegular, Op1Addr: vm.Op1SrcAP}
	operands := vm.Operands{}
	vm := vm.NewVirtualMachine()
	err := vm.UpdateRegisters(&instruction, &operands)
	if err != nil {
		t.Errorf("UpdateResigters failed with error: %s", err)
	}
	if !reflect.DeepEqual(vm.RunContext.Fp, memory.Relocatable{SegmentIndex: 0, Offset: 0}) {
		t.Errorf("Wrong fp value after registers update")
	}
	if !reflect.DeepEqual(vm.RunContext.Ap, memory.Relocatable{SegmentIndex: 0, Offset: 0}) {
		t.Errorf("Wrong ap value after registers update")
	}
	if !reflect.DeepEqual(vm.RunContext.Pc, memory.Relocatable{SegmentIndex: 0, Offset: 1}) {
		t.Errorf("Wrong pc value after registers update")
	}
}

func TestUpdateRegistersMixedTypes(t *testing.T) {
	instruction := vm.Instruction{FpUpdate: vm.FpUpdateDst, ApUpdate: vm.ApUpdateAdd2, PcUpdate: vm.PcUpdateJumpRel, Op1Addr: vm.Op1SrcAP}
	operands := vm.Operands{Dst: *memory.NewMaybeRelocatableRelocatable(memory.NewRelocatable(1, 11)), Res: memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(8))}
	v := vm.NewVirtualMachine()
	v.RunContext = vm.RunContext{Pc: memory.NewRelocatable(0, 4), Ap: memory.NewRelocatable(1, 5), Fp: memory.NewRelocatable(1, 6)}
	err := v.UpdateRegisters(&instruction, &operands)
	if err != nil {
		t.Errorf("UpdateResigters failed with error: %s", err)
	}
	if !reflect.DeepEqual(v.RunContext.Fp, memory.Relocatable{SegmentIndex: 1, Offset: 11}) {
		t.Errorf("Wrong fp value after registers update")
	}
	if !reflect.DeepEqual(v.RunContext.Ap, memory.Relocatable{SegmentIndex: 1, Offset: 7}) {
		t.Errorf("Wrong ap value after registers update")
	}
	if !reflect.DeepEqual(v.RunContext.Pc, memory.Relocatable{SegmentIndex: 0, Offset: 12}) {
		t.Errorf("Wrong pc value after registers update")
	}
}
func TestUpdateFpRegular(t *testing.T) {
	instruction := vm.Instruction{FpUpdate: vm.FpUpdateRegular}
	operands := vm.Operands{}
	vm := vm.NewVirtualMachine()
	err := vm.UpdateFp(&instruction, &operands)
	if err != nil {
		t.Errorf("UpdateFp failed with error: %s", err)
	}
	if !reflect.DeepEqual(vm.RunContext.Fp, memory.Relocatable{SegmentIndex: 0, Offset: 0}) {
		t.Errorf("Wrong value after fp update")
	}
}

func TestUpdateFpDstInt(t *testing.T) {
	instruction := vm.Instruction{FpUpdate: vm.FpUpdateDst}
	operands := vm.Operands{Dst: *memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(9))}
	vm := vm.NewVirtualMachine()
	err := vm.UpdateFp(&instruction, &operands)
	if err != nil {
		t.Errorf("UpdateFp failed with error: %s", err)
	}
	if !reflect.DeepEqual(vm.RunContext.Fp, memory.Relocatable{SegmentIndex: 0, Offset: 9}) {
		t.Errorf("Wrong value after fp update")
	}
}
func TestUpdateFpDstRelocatable(t *testing.T) {
	instruction := vm.Instruction{FpUpdate: vm.FpUpdateDst}
	operands := vm.Operands{Dst: *memory.NewMaybeRelocatableRelocatable(memory.Relocatable{SegmentIndex: 0, Offset: 9})}
	vm := vm.NewVirtualMachine()
	err := vm.UpdateFp(&instruction, &operands)
	if err != nil {
		t.Errorf("UpdateFp failed with error: %s", err)
	}
	if !reflect.DeepEqual(vm.RunContext.Fp, memory.Relocatable{SegmentIndex: 0, Offset: 9}) {
		t.Errorf("Wrong value after fp update")
	}
}

func TestUpdateFpApPlus2(t *testing.T) {
	instruction := vm.Instruction{FpUpdate: vm.FpUpdateAPPlus2}
	operands := vm.Operands{}
	vm := vm.NewVirtualMachine()
	// Change the value of Ap offset
	vm.RunContext.Ap.Offset = 7
	err := vm.UpdateFp(&instruction, &operands)
	if err != nil {
		t.Errorf("UpdateFp failed with error: %s", err)
	}
	if !reflect.DeepEqual(vm.RunContext.Fp, memory.Relocatable{SegmentIndex: 0, Offset: 9}) {
		t.Errorf("Wrong value after fp update")
	}
}

func TestUpdateApRegular(t *testing.T) {
	instruction := vm.Instruction{ApUpdate: vm.ApUpdateRegular}
	operands := vm.Operands{}
	vm := vm.NewVirtualMachine()
	err := vm.UpdateAp(&instruction, &operands)
	if err != nil {
		t.Errorf("UpdateAp failed with error: %s", err)
	}
	if !reflect.DeepEqual(vm.RunContext.Ap, memory.Relocatable{SegmentIndex: 0, Offset: 0}) {
		t.Errorf("Wrong value after ap update")
	}
}

func TestUpdateApAdd2(t *testing.T) {
	instruction := vm.Instruction{ApUpdate: vm.ApUpdateAdd2}
	operands := vm.Operands{}
	vm := vm.NewVirtualMachine()
	err := vm.UpdateAp(&instruction, &operands)
	if err != nil {
		t.Errorf("UpdateAp failed with error: %s", err)
	}
	if !reflect.DeepEqual(vm.RunContext.Ap, memory.Relocatable{SegmentIndex: 0, Offset: 2}) {
		t.Errorf("Wrong value after ap update")
	}
}

func TestUpdateApAdd1(t *testing.T) {
	instruction := vm.Instruction{ApUpdate: vm.ApUpdateAdd1}
	operands := vm.Operands{}
	vm := vm.NewVirtualMachine()
	err := vm.UpdateAp(&instruction, &operands)
	if err != nil {
		t.Errorf("UpdateAp failed with error: %s", err)
	}
	if !reflect.DeepEqual(vm.RunContext.Ap, memory.Relocatable{SegmentIndex: 0, Offset: 1}) {
		t.Errorf("Wrong value after ap update")
	}
}
func TestUpdateApAddWithResInt(t *testing.T) {
	instruction := vm.Instruction{ApUpdate: vm.ApUpdateAdd}
	operands := vm.Operands{Res: memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(5))}
	vm := vm.NewVirtualMachine()
	err := vm.UpdateAp(&instruction, &operands)
	if err != nil {
		t.Errorf("UpdateAp failed with error: %s", err)
	}
	if !reflect.DeepEqual(vm.RunContext.Ap, memory.Relocatable{SegmentIndex: 0, Offset: 5}) {
		t.Errorf("Wrong value after ap update")
	}
}

func TestUpdateApAddWithResRel(t *testing.T) {
	instruction := vm.Instruction{ApUpdate: vm.ApUpdateAdd}
	operands := vm.Operands{Res: memory.NewMaybeRelocatableRelocatable(memory.Relocatable{})}
	vm := vm.NewVirtualMachine()
	err := vm.UpdateAp(&instruction, &operands)
	if err == nil {
		t.Errorf("UpdateAp should have failed")
	}
}

func TestUpdateApAddWithoutRes(t *testing.T) {
	instruction := vm.Instruction{ApUpdate: vm.ApUpdateAdd}
	operands := vm.Operands{}
	vm := vm.NewVirtualMachine()
	err := vm.UpdateAp(&instruction, &operands)
	if err == nil {
		t.Errorf("UpdateAp should have failed")
	}
}

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
	operands := vm.Operands{Res: memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(0))}
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
	operands := vm.Operands{Res: memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(5))}
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
	operands := vm.Operands{Dst: *memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(0))}
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
	operands := vm.Operands{Dst: *memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(0))}
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
	operands := vm.Operands{Dst: *memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(1)), Op1: *memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(3))}
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
	operands := vm.Operands{Dst: *memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(1)), Op1: *memory.NewMaybeRelocatableRelocatable(memory.Relocatable{})}
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

func VmNew(run_context vm.RunContext, current_step uint, segments_manager memory.MemorySegmentManager) vm.VirtualMachine {
	return vm.VirtualMachine{
		RunContext:  run_context,
		CurrentStep: current_step,
		Segments:    segments_manager,
	}
}

func TestComputeOperandsAddAp(t *testing.T) {
	instruction := vm.Instruction{
		Off0:     0,
		Off1:     1,
		Off2:     2,
		DstReg:   vm.AP,
		Op0Reg:   vm.FP,
		Op1Addr:  vm.Op1SrcAP,
		ResLogic: vm.ResAdd,
		PcUpdate: vm.PcUpdateRegular,
		ApUpdate: vm.ApUpdateRegular,
		FpUpdate: vm.FpUpdateRegular,
		Opcode:   vm.NOp,
	}

	memory_manager := memory.NewMemorySegmentManager()
	run_context := vm.RunContext{
		Ap: memory.NewRelocatable(1, 0),
		Fp: memory.NewRelocatable(1, 0),
		Pc: memory.NewRelocatable(0, 0),
	}
	vmachine := VmNew(run_context, 0, memory_manager)

	for i := 0; i < 2; i++ {
		vmachine.Segments.AddSegment()
	}

	dst_addr := memory.NewRelocatable(1, 0)
	dst_addr_value := memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(5))
	op0_addr := memory.NewRelocatable(1, 1)
	op0_addr_value := memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(2))
	op1_addr := memory.NewRelocatable(1, 2)
	op1_addr_value := memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(3))

	vmachine.Segments.Memory.Insert(dst_addr, dst_addr_value)
	vmachine.Segments.Memory.Insert(op0_addr, op0_addr_value)
	vmachine.Segments.Memory.Insert(op1_addr, op1_addr_value)

	expected_operands := vm.Operands{
		Dst: *dst_addr_value,
		Res: dst_addr_value,
		Op0: *op0_addr_value,
		Op1: *op1_addr_value,
	}

	operands, _ := vmachine.ComputeOperands(instruction)

	if operands.Dst != expected_operands.Dst {
		t.Errorf("Different Dst register")
	}
	if operands.Op0 != expected_operands.Op0 {
		t.Errorf("Different op0 register")
	}
	if operands.Op1 != expected_operands.Op1 {
		t.Errorf("Different op1 register")
	}
	if *operands.Res != *expected_operands.Res {
		t.Errorf("Different res register")
	}
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

	expectedTrace := []vm.RelocatedTraceEntry{{Pc: lambdaworks.FeltFromUint64(1), Ap: lambdaworks.FeltFromUint64(4), Fp: lambdaworks.FeltFromUint64(4)}}
	actualTrace, err := virtualMachine.GetRelocatedTrace()
	if err != nil {
		t.Errorf("Trace relocation error failed with test: %s", err)
	}
	if !reflect.DeepEqual(expectedTrace, actualTrace) {
		t.Errorf("Relocated trace and expected trace are not the same")
	}
}

func TestWriteBinaryMemoryFile(t *testing.T) {
	var relocatedMemory = make(map[uint]uint)
	relocatedMemory[1] = 66
	relocatedMemory[2] = 42
	relocatedMemory[3] = 30

	var actualMemoryBuffer bytes.Buffer
	cairo_run.WriteEncodedMemory(relocatedMemory, &actualMemoryBuffer)
}

func buildTestProgramMemory(virtualMachine *vm.VirtualMachine) {
	virtualMachine.Trace = []vm.TraceEntry{{Pc: memory.NewRelocatable(0, 0), Ap: memory.NewRelocatable(2, 0), Fp: memory.NewRelocatable(2, 0)}}
	for i := 0; i < 4; i++ {
		virtualMachine.Segments.AddSegment()
	}
	virtualMachine.Segments.Memory.Insert(memory.NewRelocatable(0, 0), memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(2345108766317314046)))
	virtualMachine.Segments.Memory.Insert(memory.NewRelocatable(1, 0), memory.NewMaybeRelocatableRelocatable(memory.NewRelocatable(2, 0)))
	virtualMachine.Segments.Memory.Insert(memory.NewRelocatable(1, 1), memory.NewMaybeRelocatableRelocatable(memory.NewRelocatable(3, 0)))
}
