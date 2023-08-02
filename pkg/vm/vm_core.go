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
	runContext  RunContext
	currentStep uint
	segments    memory.MemorySegmentManager
}

func NewVirtualMachine() *VirtualMachine {
	return &VirtualMachine{
		runContext:  RunContext{},
		currentStep: 0,
		segments:    memory.NewMemorySegmentManager(),
	}
}

type Operands struct {
	Dst memory.MaybeRelocatable
	Res *memory.MaybeRelocatable
	Op0 memory.MaybeRelocatable
	Op1 memory.MaybeRelocatable
}

type OperandsAddresses struct {
	DstAddr memory.Relocatable
	Op0Addr memory.Relocatable
	Op1Addr memory.Relocatable
}

type DeducedOperands struct {
	operands uint8
}

func (vm *VirtualMachine) OpcodeAssertions(instruction Instruction, operands Operands) error {
	switch instruction.Opcode {
	case AssertEq:
		if operands.Res == nil {
			return &VirtualMachineError{"UnconstrainedResAssertEq"}
		}
		if !operands.Res.IsEqual(&operands.Dst) {
			_, isInt := operands.Res.GetInt()
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

		if !operands.Op0.IsEqual(returnPC) {
			return &VirtualMachineError{"CantWriteReturnPc"}
		}

		returnFP := vm.runContext.Fp
		dstRelocatable, _ := operands.Dst.GetRelocatable()
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
