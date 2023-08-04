package vm

import (
	"errors"

	"github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
)

// VirtualMachine represents the Cairo VM.
// Runs Cairo assembly and produces an execution trace.
type VirtualMachine struct {
	RunContext     RunContext
	currentStep    uint
	Segments       memory.MemorySegmentManager
	Trace          []TraceEntry
	RelocatedTrace []RelocatedTraceEntry
}

func NewVirtualMachine() *VirtualMachine {
	segments := memory.NewMemorySegmentManager()
	trace := make([]TraceEntry, 0)
	relocatedTrace := make([]RelocatedTraceEntry, 0)
	return &VirtualMachine{Segments: *segments, Trace: trace, RelocatedTrace: relocatedTrace}
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
		res, ok := operands.Res.GetInt()
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
