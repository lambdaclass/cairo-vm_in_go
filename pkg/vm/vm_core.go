package vm

import (
	"errors"
	"fmt"

	"github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
	"github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
)

// VirtualMachine represents the Cairo VM.
// Runs Cairo assembly and produces an execution trace.
type VirtualMachine struct {
	runContext  RunContext
	currentStep uint
	segments    memory.MemorySegmentManager
}

type Operands struct {
	dst memory.MaybeRelocatable
	res memory.MaybeRelocatable
	op0 memory.MaybeRelocatable
	op1 memory.MaybeRelocatable
}

type OperandsAddresses struct {
	dst_addr memory.Relocatable
	op0_addr memory.Relocatable
	op1_addr memory.Relocatable
}

// ------------------------
//  Deduced Operands funcs
// ------------------------

type DeducedOperands struct {
	operands uint8
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
		dst_addr: dst_addr,
		op0_addr: op0_addr,
		op1_addr: op1_addr,
	}

	operands := Operands{
		dst: *dst_op,
		op0: *op0_op,
		op1: *op1_op,
		res: res,
	}
	return operands, accesed_addresses, deduced_operands, nil
}

func (vm VirtualMachine) run_instrucion(instruction Instruction) {
	fmt.Println("hello from instruction")
}
