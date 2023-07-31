package vm

import "github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"

// VirtualMachine represents the Cairo VM.
// Runs Cairo assembly and produces an execution trace.
type VirtualMachine struct {
	runContext  RunContext
	currentStep uint
	Segments    *memory.MemorySegmentManager
}

func NewVirtualMachine() *VirtualMachine {
	runContext := RunContext{Pc: memory.NewRelocatable(0, 0), Ap: memory.NewRelocatable(0, 0), Fp: memory.NewRelocatable(0, 0)}
	segments := memory.NewMemorySegmentManager()

	return &VirtualMachine{runContext: runContext, currentStep: 0, Segments: segments}
}
