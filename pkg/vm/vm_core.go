package vm

import "github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"

// VirtualMachine represents the Cairo VM.
// Runs Cairo assembly and produces an execution trace.
type VirtualMachine struct {
	runContext  RunContext
	currentStep uint
	segments    memory.MemorySegmentManager
}
