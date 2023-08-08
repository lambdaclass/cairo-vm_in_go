package vm

import (
	"testing"

	"github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
	"github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
)

func TestOpcodeAssertionsResUnconstrained(t *testing.T) {
	instruction := Instruction{
		Off1:     1,
		Off2:     2,
		Off0:     3,
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
		Off1:     1,
		Off2:     2,
		Off0:     3,
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
		Off1:     1,
		Off2:     2,
		Off0:     3,
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
		Off1:     1,
		Off2:     2,
		Off0:     3,
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
		Off1:     1,
		Off2:     2,
		Off0:     3,
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
		Off1:     1,
		Off2:     2,
		Off0:     3,
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

	m1, m2, err := vm.DeduceOp1(instruction, nil, nil)

	if err != nil {
		t.Error(err)
	}

	if m1 != nil {
		t.Error("maybe relocatable of deduced operand is not nil")
	}

	if m2 != nil {
		t.Error("maybe relocatable of deduced operand is not nil")
	}
}

func TestDeduceOp1OpcodeAssertEqResAddWithOptionals(t *testing.T) {
	instruction := Instruction{
		Off1:     1,
		Off2:     2,
		Off0:     3,
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

	expected_dst := memory.NewMaybeRelocatableFelt(lambdaworks.FeltOne())
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
		Off1:     1,
		Off2:     2,
		Off0:     3,
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

	m1, m2, err := vm.DeduceOp1(instruction, nil, nil)

	if err != nil {
		t.Error(err)
	}

	if m1 != nil {
		t.Error("maybe relocatable of deduced operand is not nil")
	}

	if m2 != nil {
		t.Error("maybe relocatable of deduced operand is not nil")
	}
}
func TestDeduceOp1OpcodeAssertEqResMulNonZeroOp0(t *testing.T) {
	instruction := Instruction{
		Off1:     1,
		Off2:     2,
		Off0:     3,
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
		Off1:     1,
		Off2:     2,
		Off0:     3,
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

	m1, m2, err := vm.DeduceOp1(instruction, dst, op0)

	if err != nil {
		t.Error(err)
	}

	if m1 != nil {
		t.Error("maybe relocatable of deduced operand is not nil")
	}

	if m2 != nil {
		t.Error("maybe relocatable of deduced operand is not nil")
	}
}

func TestDeduceOp1OpcodeAssertEqResOp1WithoutDst(t *testing.T) {
	instruction := Instruction{
		Off1:     1,
		Off2:     2,
		Off0:     3,
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

	m1, m2, err := vm.DeduceOp1(instruction, nil, op0)

	if err != nil {
		t.Error(err)
	}

	if m1 != nil {
		t.Error("maybe relocatable of deduced operand is not nil")
	}

	if m2 != nil {
		t.Error("maybe relocatable of deduced operand is not nil")
	}
}
func TestDeduceOp1OpcodeAssertEqResOp1WithDst(t *testing.T) {
	instruction := Instruction{
		Off1:     1,
		Off2:     2,
		Off0:     3,
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

func TestDeduceDstOpcodeAssertEqWithRes(t *testing.T) {
	instruction := Instruction{
		Off1:     1,
		Off2:     2,
		Off0:     3,
		DstReg:   FP,
		Op0Reg:   AP,
		Op1Addr:  Op1SrcAP,
		ResLogic: ResUnconstrained,
		PcUpdate: PcUpdateJump,
		ApUpdate: ApUpdateRegular,
		FpUpdate: FpUpdateRegular,
		Opcode:   AssertEq,
	}

	vm := NewVirtualMachine()

	res := memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(7))
	expected_res := memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(7))

	result_res := vm.DeduceDst(instruction, res)

	if *expected_res != *result_res {
		t.Error("Different Res value")
	}
}

func TestDeduceDstOpcodeAssertEqWithoutRes(t *testing.T) {
	instruction := Instruction{
		Off1:     1,
		Off2:     2,
		Off0:     3,
		DstReg:   FP,
		Op0Reg:   AP,
		Op1Addr:  Op1SrcAP,
		ResLogic: ResUnconstrained,
		PcUpdate: PcUpdateJump,
		ApUpdate: ApUpdateRegular,
		FpUpdate: FpUpdateRegular,
		Opcode:   AssertEq,
	}

	vm := NewVirtualMachine()

	result_res := vm.DeduceDst(instruction, nil)

	if result_res != nil {
		t.Error("Different Res value")
	}
}

func TestDeduceDstOpcodeCall(t *testing.T) {
	instruction := Instruction{
		OffOp0:   1,
		OffOp1:   2,
		OffDst:   3,
		DstReg:   FP,
		Op0Reg:   AP,
		Op1Addr:  Op1SrcAP,
		ResLogic: ResUnconstrained,
		PcUpdate: PcUpdateJump,
		ApUpdate: ApUpdateRegular,
		FpUpdate: FpUpdateRegular,
		Opcode:   Call,
	}

	vm := NewVirtualMachine()
	vm.RunContext.Fp = memory.NewRelocatable(1, 0)

	result_dst := vm.DeduceDst(instruction, nil)
	mr := memory.NewRelocatable(1, 0)
	expected_dst := memory.NewMaybeRelocatableRelocatable(mr)

	if *result_dst != *expected_dst {
		t.Error("Different Dst value")
	}
}

func TestDeduceDstOpcodeRet(t *testing.T) {
	instruction := Instruction{
		OffOp0:   1,
		OffOp1:   2,
		OffDst:   3,
		DstReg:   FP,
		Op0Reg:   AP,
		Op1Addr:  Op1SrcAP,
		ResLogic: ResUnconstrained,
		PcUpdate: PcUpdateJump,
		ApUpdate: ApUpdateRegular,
		FpUpdate: FpUpdateRegular,
		Opcode:   Ret,
	}

	vm := NewVirtualMachine()

	result_dst := vm.DeduceDst(instruction, nil)

	if result_dst != nil {
		t.Error("Different Dst value than nil")
	}
}
