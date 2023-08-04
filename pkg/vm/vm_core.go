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
	RunContext     RunContext
	CurrentStep    uint
	Segments       memory.MemorySegmentManager
	Trace          []TraceEntry
	RelocatedTrace []RelocatedTraceEntry
}

func NewVirtualMachine() *VirtualMachine {
	segments := memory.NewMemorySegmentManager()
	trace := make([]TraceEntry, 0)
	relocatedTrace := make([]RelocatedTraceEntry, 0)
	return &VirtualMachine{Segments: segments, Trace: trace, RelocatedTrace: relocatedTrace}
}

// Relocates the VM's trace, turning relocatable registers to numbered ones
func (v *VirtualMachine) RelocateTrace(relocationTable *[]uint) error {
	if len(*relocationTable) < 2 {
		return errors.New("no relocation found for execution segment")
	}

	for _, entry := range v.Trace {
		v.RelocatedTrace = append(v.RelocatedTrace, RelocatedTraceEntry{
			Pc: entry.Pc.RelocateAddress(relocationTable),
			Ap: entry.Ap.RelocateAddress(relocationTable),
			Fp: entry.Fp.RelocateAddress(relocationTable),
		})
	}

	return nil
}

func (v *VirtualMachine) GetRelocatedTrace() ([]RelocatedTraceEntry, error) {
	if len(v.RelocatedTrace) > 0 {
		return v.RelocatedTrace, nil
	} else {
		return nil, errors.New("trace not relocated")
	}
}

func (v *VirtualMachine) Relocate() error {
	v.Segments.ComputeEffectiveSizes()
	if len(v.Trace) == 0 {
		return nil
	}

	relocationTable, ok := v.Segments.RelocateSegments()
	// This should be unreachable
	if !ok {
		return errors.New("ComputeEffectiveSizes called but RelocateSegments still returned error")
	}

	_, err := v.Segments.RelocateMemory(&relocationTable)
	if err != nil {
		return err
	}

	v.RelocateTrace(&relocationTable)
	return nil
}

type Operands struct {
	Dst memory.MaybeRelocatable
	Res *memory.MaybeRelocatable
	Op0 memory.MaybeRelocatable
	Op1 memory.MaybeRelocatable
}

// ------------------------
//  Deduced Operands funcs
// ------------------------

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
		new_rel, err := vm.RunContext.Pc.AddUint(instruction.Size())
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

//------------------------
//  virtual machines funcs
// ------------------------

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
		num_op0, m_type := op0.GetFelt()
		num_op1, other_type := op1.GetFelt()
		if m_type && other_type {
			result := memory.NewMaybeRelocatableFelt(num_op0.Mul(num_op1))
			return result, nil
		} else {
			return nil, errors.New("ComputeResRelocatableMul")
		}

	case ResUnconstrained:
		return nil, nil
	}
	return nil, nil
}

func (vm *VirtualMachine) ComputeOperands(instruction Instruction) (Operands, error) {

	dst_addr, err := vm.RunContext.ComputeDstAddr(instruction)
	if err != nil {
		return Operands{}, errors.New("FailedToComputeDstAddr")
	}
	dst_op, dst_err := vm.Segments.Memory.Get(dst_addr)
	if dst_err != nil {
		return Operands{}, err
	}

	op0_addr, err := vm.RunContext.ComputeOp0Addr(instruction)
	if err != nil {
		return Operands{}, errors.New("FailedToComputeOp0Addr")
	}
	op0_op, op_err := vm.Segments.Memory.Get(op0_addr)
	// this should trigger deducde op1
	if op_err != nil {
		return Operands{}, err
	}

	op1_addr, err := vm.RunContext.ComputeOp1Addr(instruction, op0_op)
	if err != nil {
		return Operands{}, errors.New("FailedToComputeOp1Addr")
	}
	// this should trigger deducde op1
	op1_op, op1_err := vm.Segments.Memory.Get(op1_addr)
	if op1_err != nil {
		return Operands{}, err
	}

	res, err := vm.ComputeRes(instruction, *op0_op, *op1_op)

	// uncomment once deduction functions are done

	// var op0 memory.MaybeRelocatable
	// if op0_err != nil {
	// op0 = vm.compute_op0_deductions(op0_addr, res, instruction, &dst_op, &op1_op)
	// } else {
	// op0 = op0_op
	// }

	// var op1 memory.MaybeRelocatable
	// if op1_err != nil {
	// op1 = vm.compute_op1_deductions(op1_addr, res, instruction, &dst_op, &op0)
	// } else {
	// op1 = op1_op
	// }

	// var dst memory.MaybeRelocatable
	// if dst_err != nil {
	// dst = vm.compute_dst_deductions(instruction, &res)
	// } else {
	// dst = dst_op
	// }

	operands := Operands{
		Dst: *dst_op,
		Op0: *op0_op,
		Op1: *op1_op,
		Res: res,
	}
	return operands, nil
}

func (vm VirtualMachine) run_instrucion(instruction Instruction) {
	fmt.Println("hello from instruction")
}

// Updates the value of PC according to the executed instruction
func (vm *VirtualMachine) UpdatePc(instruction *Instruction, operands *Operands) error {
	switch instruction.PcUpdate {
	case PcUpdateRegular:
		vm.RunContext.Pc.Offset += instruction.Size()
	case PcUpdateJump:
		if operands.Res == nil {
			return errors.New("Res.UNCONSTRAINED cannot be used with PcUpdate.JUMP")
		}
		res, ok := operands.Res.GetRelocatable()
		if !ok {
			return errors.New("an integer value as Res cannot be used with PcUpdate.JUMP")
		}
		vm.RunContext.Pc = res
	case PcUpdateJumpRel:
		if operands.Res == nil {
			return errors.New("Res.UNCONSTRAINED cannot be used with PcUpdate.JUMP_REL")
		}
		res, ok := operands.Res.GetFelt()
		if !ok {
			return errors.New("a relocatable value as Res cannot be used with PcUpdate.JUMP_REL")
		}
		new_pc, err := vm.RunContext.Pc.AddFelt(res)
		if err != nil {
			return err
		}
		vm.RunContext.Pc = new_pc
	case PcUpdateJnz:
		if operands.Dst.IsZero() {
			vm.RunContext.Pc.Offset += instruction.Size()
		} else {
			new_pc, err := vm.RunContext.Pc.AddMaybeRelocatable(operands.Op1)
			if err != nil {
				return err
			}
			vm.RunContext.Pc = new_pc
		}

	}
	return nil
}
