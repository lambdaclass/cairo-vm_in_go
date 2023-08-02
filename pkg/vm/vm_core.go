package vm

import (
	"github.com/lambdaclass/cairo-vm.go/pkg/builtins"
	"github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
)

// VirtualMachine represents the Cairo VM.
// Runs Cairo assembly and produces an execution trace.
type VirtualMachine struct {
	RunContext     RunContext
	currentStep    uint
	Segments       memory.MemorySegmentManager
	BuiltinRunners []builtins.BuiltinRunner
}

func NewVirtualMachine() *VirtualMachine {
	return &VirtualMachine{
		Segments:       *memory.NewMemorySegmentManager(),
		BuiltinRunners: make([]builtins.BuiltinRunner, 0, 9), // There will be at most 9 builtins
	}
}
