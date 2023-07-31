package runners

import (
	"github.com/lambdaclass/cairo-vm.go/pkg/vm"
	"github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
)

type CairoRunner struct {
	Vm            vm.VirtualMachine
	ProgramBase   memory.Relocatable
	ExecutionBase memory.Relocatable
}

// Creates program, execution and builtin segments
func (r *CairoRunner) InitializeSegments() {
	// Program Segment
	r.ProgramBase = r.Vm.Segments.AddSegment()
	// Execution Segment
	r.ExecutionBase = r.Vm.Segments.AddSegment()
	// Initialize builtin segments
}
