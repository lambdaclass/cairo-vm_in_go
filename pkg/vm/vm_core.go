package vm

import "github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"

// VirtualMachine represents the Cairo VM.
// Runs Cairo assembly and produces an execution trace.
type VirtualMachine struct {
	runContext  RunContext
	currentStep uint
	segments    memory.MemorySegmentManager
}

func NewVirtualMachine() *VirtualMachine {
	runContext := RunContext{pc: memory.NewRelocatable(0, 0), ap: 0, fp: 0}
	segments := memory.NewMemorySegmentManager()

	return &VirtualMachine{runContext: runContext, currentStep: 0, segments: *segments}
}

func (v *VirtualMachine) Segments() *memory.MemorySegmentManager {
	return &v.segments
}
