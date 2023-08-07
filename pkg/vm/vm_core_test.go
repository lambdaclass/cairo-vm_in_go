package vm

import (
	"testing"

	"github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
	"github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
)

func TestOpcodeAssertionsResUnconstrained(t *testing.T) {
	instruction := Instruction{
		OffOp0:   1,
		OffOp1:   2,
		OffDst:   3,
		DstReg:   FP,
		Op0Reg:   AP,
		Op1Addr:  Op1SrcAP,
		ResLogic: ResAdd,
		PcUpdate: PcUpdateRegular,
		ApUpdate: ApUpdateRegular,
		FpUpdate: FpUpdateAPPlus2,
		Opcode:   AssertEq,
	}

	operands := Operands{
		Dst: *memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(8)),
		Res: nil,
		Op0: *memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(9)),
		Op1: *memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(10)),
	}

	testVm := NewVirtualMachine()

	err := testVm.OpcodeAssertions(instruction, operands)
	if err.Error() != "UnconstrainedResAssertEq" {
		t.Error("Assertion should error out with UnconstrainedResAssertEq")
	}
}

func TestOpcodeAssertionsInstructionFailed(t *testing.T) {
	instruction := Instruction{
		OffOp0:   1,
		OffOp1:   2,
		OffDst:   3,
		DstReg:   FP,
		Op0Reg:   AP,
		Op1Addr:  Op1SrcAP,
		ResLogic: ResAdd,
		PcUpdate: PcUpdateRegular,
		ApUpdate: ApUpdateRegular,
		FpUpdate: FpUpdateAPPlus2,
		Opcode:   AssertEq,
	}

	operands := Operands{
		Dst: *memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(9)),
		Res: memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(8)),
		Op0: *memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(9)),
		Op1: *memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(10)),
	}

	testVm := NewVirtualMachine()
	err := testVm.OpcodeAssertions(instruction, operands)
	if err.Error() != "DiffAssertValues" {
		t.Error("Assertion should error out with DiffAssertValues")
	}

}

func TestOpcodeAssertionsInstructionFailedRelocatables(t *testing.T) {
	instruction := Instruction{
		OffOp0:   1,
		OffOp1:   2,
		OffDst:   3,
		DstReg:   FP,
		Op0Reg:   AP,
		Op1Addr:  Op1SrcAP,
		ResLogic: ResAdd,
		PcUpdate: PcUpdateRegular,
		ApUpdate: ApUpdateRegular,
		FpUpdate: FpUpdateAPPlus2,
		Opcode:   AssertEq,
	}

	operands := Operands{
		Dst: *memory.NewMaybeRelocatableRelocatable(memory.NewRelocatable(1, 1)),
		Res: memory.NewMaybeRelocatableRelocatable(memory.NewRelocatable(1, 2)),
		Op0: *memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(9)),
		Op1: *memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(10)),
	}

	testVm := NewVirtualMachine()
	err := testVm.OpcodeAssertions(instruction, operands)
	if err.Error() != "DiffAssertValues" {
		t.Error("Assertion should error out with DiffAssertValues")
	}
}

func TestOpcodeAssertionsInconsistentOp0(t *testing.T) {
	instruction := Instruction{
		OffOp0:   1,
		OffOp1:   2,
		OffDst:   3,
		DstReg:   FP,
		Op0Reg:   AP,
		Op1Addr:  Op1SrcAP,
		ResLogic: ResAdd,
		PcUpdate: PcUpdateRegular,
		ApUpdate: ApUpdateRegular,
		FpUpdate: FpUpdateAPPlus2,
		Opcode:   Call,
	}

	operands := Operands{
		Dst: *memory.NewMaybeRelocatableRelocatable(memory.NewRelocatable(0, 8)),
		Res: memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(8)),
		Op0: *memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(9)),
		Op1: *memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(10)),
	}

	testVm := NewVirtualMachine()
	testVm.RunContext.Pc = memory.NewRelocatable(0, 4)
	err := testVm.OpcodeAssertions(instruction, operands)
	if err.Error() != "CantWriteReturnPc" {
		t.Error("Assertion should error out with CantWriteReturnPc")
	}
}

func TestOpcodeAssertionsInconsistentDst(t *testing.T) {
	instruction := Instruction{
		OffOp0:   1,
		OffOp1:   2,
		OffDst:   3,
		DstReg:   FP,
		Op0Reg:   AP,
		Op1Addr:  Op1SrcAP,
		ResLogic: ResAdd,
		PcUpdate: PcUpdateRegular,
		ApUpdate: ApUpdateRegular,
		FpUpdate: FpUpdateAPPlus2,
		Opcode:   Call,
	}

	operands := Operands{
		Dst: *memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(8)),
		Res: memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(8)),
		Op0: *memory.NewMaybeRelocatableRelocatable(memory.NewRelocatable(0, 1)),
		Op1: *memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(10)),
	}

	testVm := NewVirtualMachine()
	testVm.RunContext.Fp = memory.NewRelocatable(1, 6)
	err := testVm.OpcodeAssertions(instruction, operands)
	if err.Error() != "CantWriteReturnFp" {
		t.Error("Assertion should error out with CantWriteReturnFp")
	}
}

func TestDeduceOp1OpcodeCall(t *testing.T) {
	instruction := Instruction{
		OffOp0:   1,
		OffOp1:   2,
		OffDst:   3,
		DstReg:   FP,
		Op0Reg:   AP,
		Op1Addr:  Op1SrcAP,
		ResLogic: ResAdd,
		PcUpdate: PcUpdateJump,
		ApUpdate: ApUpdateRegular,
		FpUpdate: FpUpdateRegular,
		Opcode:   Call,
	}

	vm := NewVirtualMachine()

	_, _, err := vm.DeduceOp1(instruction, nil, nil)

	if err != nil {
		t.Error(err)
	}
}

func TestDeduceOp1OpcodeAssertEqResAddWithOptionals(t *testing.T) {
	instruction := Instruction{
		OffOp0:   1,
		OffOp1:   2,
		OffDst:   3,
		DstReg:   FP,
		Op0Reg:   AP,
		Op1Addr:  Op1SrcAP,
		ResLogic: ResAdd,
		PcUpdate: PcUpdateJump,
		ApUpdate: ApUpdateRegular,
		FpUpdate: FpUpdateRegular,
		Opcode:   AssertEq,
	}

	vm := NewVirtualMachine()

	dst := memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(3))
	op0 := memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(2))

	expected_dst := memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(1))
	expected_op0 := memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(3))

	m1, m2, err := vm.DeduceOp1(instruction, dst, op0)

	if err != nil {
		t.Error(err)
	}

	if *m1 != *expected_dst {
		t.Error("Different dst value")
	}
	if *m2 != *expected_op0 {
		t.Error("Different op0 value")
	}
}

func TestDeduceOp1OpcodeAssertEqResAddWithoutOptionals(t *testing.T) {
	instruction := Instruction{
		OffOp0:   1,
		OffOp1:   2,
		OffDst:   3,
		DstReg:   FP,
		Op0Reg:   AP,
		Op1Addr:  Op1SrcAP,
		ResLogic: ResAdd,
		PcUpdate: PcUpdateJump,
		ApUpdate: ApUpdateRegular,
		FpUpdate: FpUpdateRegular,
		Opcode:   AssertEq,
	}

	vm := NewVirtualMachine()

	_, _, err := vm.DeduceOp1(instruction, nil, nil)

	if err != nil {
		t.Error(err)
	}
}
func TestDeduceOp1OpcodeAssertEqResMulNonZeroOp0(t *testing.T) {
	instruction := Instruction{
		OffOp0:   1,
		OffOp1:   2,
		OffDst:   3,
		DstReg:   FP,
		Op0Reg:   AP,
		Op1Addr:  Op1SrcAP,
		ResLogic: ResMul,
		PcUpdate: PcUpdateJump,
		ApUpdate: ApUpdateRegular,
		FpUpdate: FpUpdateRegular,
		Opcode:   AssertEq,
	}

	vm := NewVirtualMachine()

	dst := memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(4))
	op0 := memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(2))

	expected_dst := memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(2))
	expected_op0 := memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(4))

	m1, m2, err := vm.DeduceOp1(instruction, dst, op0)

	if err != nil {
		t.Error(err)
	}

	if *m1 != *expected_dst {
		t.Error("Different dst value")
	}
	if *m2 != *expected_op0 {
		t.Error("Different op0 value")
	}
}

func TestDeduceOp1OpcodeAssertEqResMulZeroOp0(t *testing.T) {
	instruction := Instruction{
		OffOp0:   1,
		OffOp1:   2,
		OffDst:   3,
		DstReg:   FP,
		Op0Reg:   AP,
		Op1Addr:  Op1SrcAP,
		ResLogic: ResMul,
		PcUpdate: PcUpdateJump,
		ApUpdate: ApUpdateRegular,
		FpUpdate: FpUpdateRegular,
		Opcode:   AssertEq,
	}

	vm := NewVirtualMachine()

	dst := memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(4))
	op0 := memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(0))

	_, _, err := vm.DeduceOp1(instruction, dst, op0)

	if err != nil {
		t.Error(err)
	}
}

func TestDeduceOp1OpcodeAssertEqResOp1WithoutDst(t *testing.T) {
	instruction := Instruction{
		OffOp0:   1,
		OffOp1:   2,
		OffDst:   3,
		DstReg:   FP,
		Op0Reg:   AP,
		Op1Addr:  Op1SrcAP,
		ResLogic: ResOp1,
		PcUpdate: PcUpdateJump,
		ApUpdate: ApUpdateRegular,
		FpUpdate: FpUpdateRegular,
		Opcode:   AssertEq,
	}

	vm := NewVirtualMachine()

	op0 := memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(0))

	_, _, err := vm.DeduceOp1(instruction, nil, op0)

	if err != nil {
		t.Error(err)
	}
}
func TestDeduceOp1OpcodeAssertEqResOp1WithDst(t *testing.T) {
	instruction := Instruction{
		OffOp0:   1,
		OffOp1:   2,
		OffDst:   3,
		DstReg:   FP,
		Op0Reg:   AP,
		Op1Addr:  Op1SrcAP,
		ResLogic: ResOp1,
		PcUpdate: PcUpdateJump,
		ApUpdate: ApUpdateRegular,
		FpUpdate: FpUpdateRegular,
		Opcode:   AssertEq,
	}

	vm := NewVirtualMachine()

	dst := memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(7))

	expected_dst := memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(7))
	expected_op0 := memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(7))

	m1, m2, err := vm.DeduceOp1(instruction, dst, nil)

	if err != nil {
		t.Error(err)
	}

	if *m1 != *expected_dst {
		t.Error("Different dst value")
	}
	if *m2 != *expected_op0 {
		t.Error("Different op0 value")
	}
}
