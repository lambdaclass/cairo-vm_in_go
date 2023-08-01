package vm

import (
	"fmt"

	"github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
)

type VirtualMachineError struct {
	Msg string
}

func (e *VirtualMachineError) Error() string {
	return fmt.Sprintf(e.Msg)
}

// VirtualMachine represents the Cairo VM.
// Runs Cairo assembly and produces an execution trace.
type VirtualMachine struct {
	runContext  *RunContext
	currentStep uint
	segments    *memory.MemorySegmentManager
}

func NewVirtualMachine() *VirtualMachine {
	return &VirtualMachine{
		runContext:  NewRunContext(),
		currentStep: 0,
		segments:    memory.NewMemorySegmentManager(),
	}
}

type Operands struct {
	DST memory.MaybeRelocatable
	RES *memory.MaybeRelocatable
	OP0 memory.MaybeRelocatable
	OP1 memory.MaybeRelocatable
}

type OperandsAddresses struct {
	dst_addr memory.Relocatable
	op0_addr memory.Relocatable
	op1_addr memory.Relocatable
}

type DeducedOperands struct {
	operands uint8
}

func (vm *VirtualMachine) OpcodeAssertions(instruction Instruction, operands Operands) error {
	switch instruction.Opcode {
	case AssertEq:
		if operands.RES == nil {
			return &VirtualMachineError{"UnconstrainedResAssertEq"}
		}
		if !operands.RES.IsEqual(&operands.DST) {
			_, isInt := operands.RES.GetInt()
			if isInt {
				return &VirtualMachineError{"IntDiffAssertValues"}
			} else {
				return &VirtualMachineError{"RelocatableDiffAssertValues"}
			}
		}
	case Call:
		new_rel, err := vm.runContext.Pc.AddRelocatable(instruction.size())
		if err != nil {
			return err
		}
		returnPC := memory.NewMaybeRelocatableRelocatable(new_rel)

		if !operands.OP0.IsEqual(returnPC) {
			return &VirtualMachineError{"CantWriteReturnPc"}
		}

		returnFP := vm.runContext.GetFP()
		dstRelocatable, _ := operands.DST.GetRelocatable()
		if !returnFP.IsEqual(&dstRelocatable) {
			return &VirtualMachineError{"CantWriteReturnFp"}
		}
	}

	return nil
}

// func (vm VirtualMachine) compute_operands(instruction Instruction) (Operands, OperandsAddresses, DeducedOperands, VirtualMachineError) {

// 	dst_addr, err := vm.runContext.ComputeDstAddr(instruction)
// 	if err != nil {
// 		return Operands{}, OperandsAddresses{}, DeducedOperands{}, VirtualMachineError{Msg: "FailtToComputeDstAddr"}
// 	}

// 	dst_op, err_dst = vm.segments.Memory.Get(&dst_addr)

// 	op0_addr, err := vm.runContext.ComputeOp0Addr(instruction)
// 	if err != nil {
// 		return Operands{}, OperandsAddresses{}, DeducedOperands{}, VirtualMachineError{Msg: "FailtToComputeOp0Addr"}
// 	}

// 	op0_op, err_op0 := vm.segments.Memory.Get(&op0_addr)

// 	op1_addr, err := vm.runContext.ComputeOp1Addr(instruction)
// 	if err != nil {
// 		return Operands{}, OperandsAddresses{}, DeducedOperands{}, VirtualMachineError{Msg: "FailtToComputeOp1Addr"}
// 	}
// 	op1_op, err_op1 := vm.segments.Memory.Get(&op1_addr)

// 	deduced_operands := DeducedOperands{operands: 0}
// }

func (vm VirtualMachine) run_instrucion(instruction Instruction) {
	fmt.Println("hello from instruction")
}
