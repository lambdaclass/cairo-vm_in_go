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
		Dst: *memory.NewMaybeRelocatableInt(lambdaworks.From(8)),
		Res: nil,
		Op0: *memory.NewMaybeRelocatableInt(lambdaworks.From(9)),
		Op1: *memory.NewMaybeRelocatableInt(lambdaworks.From(10)),
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
		Dst: *memory.NewMaybeRelocatableInt(lambdaworks.From(9)),
		Res: memory.NewMaybeRelocatableInt(lambdaworks.From(8)),
		Op0: *memory.NewMaybeRelocatableInt(lambdaworks.From(9)),
		Op1: *memory.NewMaybeRelocatableInt(lambdaworks.From(10)),
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
		Op0: *memory.NewMaybeRelocatableInt(lambdaworks.From(9)),
		Op1: *memory.NewMaybeRelocatableInt(lambdaworks.From(10)),
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
		Res: memory.NewMaybeRelocatableInt(lambdaworks.From(8)),
		Op0: *memory.NewMaybeRelocatableInt(lambdaworks.From(9)),
		Op1: *memory.NewMaybeRelocatableInt(lambdaworks.From(10)),
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
		Dst: *memory.NewMaybeRelocatableInt(lambdaworks.From(8)),
		Res: memory.NewMaybeRelocatableInt(lambdaworks.From(8)),
		Op0: *memory.NewMaybeRelocatableRelocatable(memory.NewRelocatable(0, 1)),
		Op1: *memory.NewMaybeRelocatableInt(lambdaworks.From(10)),
	}

	testVm := NewVirtualMachine()
	testVm.RunContext.Fp = memory.NewRelocatable(1, 6)
	err := testVm.OpcodeAssertions(instruction, operands)
	if err.Error() != "CantWriteReturnFp" {
		t.Error("Assertion should error out with CantWriteReturnFp")
	}
}
