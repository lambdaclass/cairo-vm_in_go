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
		DST: *memory.NewMaybeRelocatableInt(lambdaworks.From(8)),
		RES: nil,
		OP0: *memory.NewMaybeRelocatableInt(lambdaworks.From(9)),
		OP1: *memory.NewMaybeRelocatableInt(lambdaworks.From(10)),
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
		DST: *memory.NewMaybeRelocatableInt(lambdaworks.From(9)),
		RES: memory.NewMaybeRelocatableInt(lambdaworks.From(8)),
		OP0: *memory.NewMaybeRelocatableInt(lambdaworks.From(9)),
		OP1: *memory.NewMaybeRelocatableInt(lambdaworks.From(10)),
	}

	testVm := NewVirtualMachine()
	err := testVm.OpcodeAssertions(instruction, operands)
	if err.Error() != "IntDiffAssertValues" {
		t.Error("Assertion should error out with IntDiffAssertValues")
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
		DST: *memory.NewMaybeRelocatableRelocatable(memory.NewRelocatable(1, 1)),
		RES: memory.NewMaybeRelocatableRelocatable(memory.NewRelocatable(1, 2)),
		OP0: *memory.NewMaybeRelocatableInt(lambdaworks.From(9)),
		OP1: *memory.NewMaybeRelocatableInt(lambdaworks.From(10)),
	}

	testVm := NewVirtualMachine()
	err := testVm.OpcodeAssertions(instruction, operands)
	if err.Error() != "RelocatableDiffAssertValues" {
		t.Error("Assertion should error out with RelocatableDiffAssertValues")
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
		DST: *memory.NewMaybeRelocatableRelocatable(memory.NewRelocatable(0, 8)),
		RES: memory.NewMaybeRelocatableInt(lambdaworks.From(8)),
		OP0: *memory.NewMaybeRelocatableInt(lambdaworks.From(9)),
		OP1: *memory.NewMaybeRelocatableInt(lambdaworks.From(10)),
	}

	testVm := NewVirtualMachine()
	testVm.runContext.Pc = memory.NewRelocatable(0, 4)
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
		DST: *memory.NewMaybeRelocatableInt(lambdaworks.From(8)),
		RES: memory.NewMaybeRelocatableInt(lambdaworks.From(8)),
		OP0: *memory.NewMaybeRelocatableRelocatable(memory.NewRelocatable(0, 1)),
		OP1: *memory.NewMaybeRelocatableInt(lambdaworks.From(10)),
	}

	testVm := NewVirtualMachine()
	testVm.runContext.Fp = memory.NewRelocatable(1, 6)
	err := testVm.OpcodeAssertions(instruction, operands)
	if err.Error() != "CantWriteReturnFp" {
		t.Error("Assertion should error out with CantWriteReturnFp")
	}
}
