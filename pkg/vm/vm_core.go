package vm

import (
	"errors"
	"fmt"

	"github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
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

// ------------------------
//  Deduced Operands funcs
// ------------------------

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

func (deduced *DeducedOperands) set_dst(value uint8) {
	deduced.operands = deduced.operands | value
}

func (deduced *DeducedOperands) set_op0(value uint8) {
	deduced.operands = deduced.operands | value<<1
}

func (deduced *DeducedOperands) set_op1(value uint8) {
	deduced.operands = deduced.operands | value<<2
}

// ------------------------
//  virtual machines funcs
// ------------------------

func (vm *VirtualMachine) ComputeRes(instruction Instruction, op0 memory.MaybeRelocatable, op1 memory.MaybeRelocatable) (memory.MaybeRelocatable, error) {
	switch instruction.ResLogic {
	case ResOp1:
		return op1, nil

	case ResAdd:
		maybe_rel, err := op0.AddMaybeRelocatable(op1)
		if err != nil {
			return memory.MaybeRelocatable{}, errors.New("adding maybe relocatable")
		}
		return maybe_rel, nil

	case ResMul:
		num_op0, m_type := op0.GetInt()
		num_op1, other_type := op1.GetInt()
		if m_type && other_type {
			result := memory.NewMaybeRelocatableInt(lambdaworks.Add(num_op0.Felt, num_op1.Felt))
			return *result, nil
		} else {
			return memory.MaybeRelocatable{}, errors.New("ComputeResRelocatableMul")
		}

	case ResUnconstrained:
		return memory.MaybeRelocatable{}, nil
	}
	return memory.MaybeRelocatable{}, nil
}

func (vm *VirtualMachine) compute_operands(instruction Instruction) (Operands, OperandsAddresses, DeducedOperands, error) {

	dst_addr, err := vm.runContext.ComputeDstAddr(instruction)
	if err != nil {
		return Operands{}, OperandsAddresses{}, DeducedOperands{}, errors.New("FailtToComputeDstAddr")
	}
	dst_op, _ := vm.segments.Memory.Get(dst_addr)

	op0_addr, err := vm.runContext.ComputeOp0Addr(instruction)
	if err != nil {
		return Operands{}, OperandsAddresses{}, DeducedOperands{}, errors.New("FailtToComputeOp0Addr")
	}
	op0_op, _ := vm.segments.Memory.Get(op0_addr)

	op1_addr, err := vm.runContext.ComputeOp1Addr(instruction, *op0_op)
	if err != nil {
		return Operands{}, OperandsAddresses{}, DeducedOperands{}, errors.New("FailtToComputeOp1Addr")
	}
	op1_op, _ := vm.segments.Memory.Get(op1_addr)

	deduced_operands := DeducedOperands{operands: 0}
	res, err := vm.ComputeRes(instruction, *op0_op, *op1_op)

	accesed_addresses := OperandsAddresses{
		DstAddr: dst_addr,
		Op0Addr: op0_addr,
		Op1Addr: op1_addr,
	}

	operands := Operands{
		Dst: *dst_op,
		Op0: *op0_op,
		Op1: *op1_op,
		Res: &res,
	}
	return operands, accesed_addresses, deduced_operands, nil
}

func (vm VirtualMachine) run_instrucion(instruction Instruction) {
	fmt.Println("hello from instruction")
}
