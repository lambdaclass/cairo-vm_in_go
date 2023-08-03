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
