package vm

import (
	"errors"
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
	RunContext  RunContext
	CurrentStep uint
	Segments    memory.MemorySegmentManager
}

func NewVirtualMachine() VirtualMachine {
	return VirtualMachine{Segments: memory.NewMemorySegmentManager()}
}

type Operands struct {
	Dst memory.MaybeRelocatable
	Res *memory.MaybeRelocatable
	Op0 memory.MaybeRelocatable
	Op1 memory.MaybeRelocatable
}

type OperandsAddresses struct {
	Dst_addr memory.Relocatable
	Op0_addr memory.Relocatable
	Op1_addr memory.Relocatable
}

// ------------------------
//  Deduced Operands funcs
// ------------------------

type DeducedOperands struct {
	Operands uint8
}

func (vm *VirtualMachine) OpcodeAssertions(instruction Instruction, operands Operands) error {
	switch instruction.Opcode {
	case AssertEq:
		if operands.Res == nil {
			return &VirtualMachineError{"UnconstrainedResAssertEq"}
		}
		if !operands.Res.IsEqual(&operands.Dst) {
			return &VirtualMachineError{"DiffAssertValues"}
		}
	case Call:
		new_rel, err := vm.RunContext.Pc.AddUint(instruction.size())
		if err != nil {
			return err
		}
		returnPC := memory.NewMaybeRelocatableRelocatable(new_rel)

		if !operands.Op0.IsEqual(returnPC) {
			return &VirtualMachineError{"CantWriteReturnPc"}
		}

		returnFP := vm.RunContext.Fp
		dstRelocatable, _ := operands.Dst.GetRelocatable()
		if !returnFP.IsEqual(&dstRelocatable) {
			return &VirtualMachineError{"CantWriteReturnFp"}
		}
	}

	return nil
}

func (deduced *DeducedOperands) set_dst(value uint8) {
	deduced.Operands = deduced.Operands | value
}

func (deduced *DeducedOperands) set_op0(value uint8) {
	deduced.Operands = deduced.Operands | value<<1
}

func (deduced *DeducedOperands) set_op1(value uint8) {
	deduced.Operands = deduced.Operands | value<<2
}

//------------------------
//  virtual machines funcs
// ------------------------

func VmNew(run_context RunContext, current_step uint, segments_manager memory.MemorySegmentManager) VirtualMachine {
	return VirtualMachine{
		RunContext:  run_context,
		CurrentStep: current_step,
		Segments:    segments_manager,
	}
}

func (vm *VirtualMachine) ComputeRes(instruction Instruction, op0 memory.MaybeRelocatable, op1 memory.MaybeRelocatable) (*memory.MaybeRelocatable, error) {
	switch instruction.ResLogic {
	case ResOp1:
		return &op1, nil

	case ResAdd:
		maybe_rel, err := op0.AddMaybeRelocatable(op1)
		if err != nil {
			return nil, err
		}
		return &maybe_rel, nil

	case ResMul:
		num_op0, m_type := op0.GetInt()
		num_op1, other_type := op1.GetInt()
		if m_type && other_type {
			result := memory.NewMaybeRelocatableInt(num_op0.Felt.Mul(num_op1.Felt))
			return result, nil
		} else {
			return nil, errors.New("ComputeResRelocatableMul")
		}

	case ResUnconstrained:
		return nil, nil
	}
	return nil, nil
}

func (vm *VirtualMachine) ComputeOperands(instruction Instruction) (Operands, OperandsAddresses, DeducedOperands, error) {

	dst_addr, err := vm.RunContext.ComputeDstAddr(instruction)
	if err != nil {
		return Operands{}, OperandsAddresses{}, DeducedOperands{}, errors.New("FailedToComputeDstAddr")
	}
	dst_op, _ := vm.Segments.Memory.Get(dst_addr)

	op0_addr, err := vm.RunContext.ComputeOp0Addr(instruction)
	if err != nil {
		return Operands{}, OperandsAddresses{}, DeducedOperands{}, errors.New("FailedToComputeOp0Addr")
	}
	op0_op, _ := vm.Segments.Memory.Get(op0_addr)

	op1_addr, err := vm.RunContext.ComputeOp1Addr(instruction, *op0_op)
	if err != nil {
		return Operands{}, OperandsAddresses{}, DeducedOperands{}, errors.New("FailedToComputeOp1Addr")
	}
	op1_op, _ := vm.Segments.Memory.Get(op1_addr)

	deduced_operands := DeducedOperands{Operands: 0}
	res, err := vm.ComputeRes(instruction, *op0_op, *op1_op)

	accesed_addresses := OperandsAddresses{
		Dst_addr: dst_addr,
		Op0_addr: op0_addr,
		Op1_addr: op1_addr,
	}

	operands := Operands{
		Dst: *dst_op,
		Op0: *op0_op,
		Op1: *op1_op,
		Res: res,
	}
	return operands, accesed_addresses, deduced_operands, nil
}

func (vm VirtualMachine) run_instrucion(instruction Instruction) {
	fmt.Println("hello from instruction")
}
