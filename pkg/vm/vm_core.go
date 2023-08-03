package vm

import "github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"

// VirtualMachine represents the Cairo VM.
// Runs Cairo assembly and produces an execution trace.
type VirtualMachine struct {
	RunContext  RunContext
	currentStep uint
	Segments    memory.MemorySegmentManager
}

func NewVirtualMachine() *VirtualMachine {
	return &VirtualMachine{Segments: *memory.NewMemorySegmentManager()}
}

type Operands struct {
	Dst memory.MaybeRelocatable
	Res memory.MaybeRelocatable
	Op0 memory.MaybeRelocatable
	Op1 memory.MaybeRelocatable
}

func (vm *VirtualMachine) updatePc(instruction *Instruction, operands *Operands) error {
	return nil
}
