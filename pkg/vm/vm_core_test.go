package vm_test

import (
	"testing"

	"github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
	"github.com/lambdaclass/cairo-vm.go/pkg/vm"
	"github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
)

// TODO: This tests checks what happens when Res is not a maybe relocatable
func TestOpcodeAssertionsResUnconstrained(t *testing.T) {
	instruction := vm.Instruction{
		OffOp0:   1,
		OffOp1:   2,
		OffDst:   3,
		DstReg:   vm.FP,
		Op0Reg:   vm.AP,
		Op1Addr:  vm.Op1SrcAP,
		ResLogic: vm.ResAdd,
		PcUpdate: vm.PcUpdateRegular,
		ApUpdate: vm.ApUpdateRegular,
		FpUpdate: vm.FpUpdateAPPlus2,
		Opcode:   vm.AssertEq,
	}

	operands := vm.Operands{
		DST: *memory.NewMaybeRelocatableInt(lambdaworks.From(8)),
		RES: nil,
		OP0: *memory.NewMaybeRelocatableInt(lambdaworks.From(9)),
		OP1: *memory.NewMaybeRelocatableInt(lambdaworks.From(10)),
	}

	vm := vm.NewVirtualMachine()

	err := vm.OpcodeAssertions(instruction, operands)
	if err.Error() != "UnconstrainedResAssertEq" {
		t.Error("Assertion should error out with NonZeroHighBitError")
	}

}
