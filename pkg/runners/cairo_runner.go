package runners

import (
	"github.com/lambdaclass/cairo-vm.go/pkg/vm"
	"github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
)

type CairoRunner struct {
	Program       vm.Program
	Vm            vm.VirtualMachine
	ProgramBase   memory.Relocatable
	ExecutionBase memory.Relocatable
	InitialPc     memory.Relocatable
	InitialAp     memory.Relocatable
	InitialFp     memory.Relocatable
	FinalPc       memory.Relocatable
}

// Creates program, execution and builtin segments
func (r *CairoRunner) InitializeSegments() {
	// Program Segment
	r.ProgramBase = r.Vm.Segments.AddSegment()
	// Execution Segment
	r.ExecutionBase = r.Vm.Segments.AddSegment()
	// Initialize builtin segments
}

// Initializes the program segment & initial pc
func (r *CairoRunner) initializeState(entrypoint uint, stack *[]memory.MaybeRelocatable) error {
	r.InitialPc = r.ProgramBase
	r.InitialPc.Offset += entrypoint
	// Load program data
	_, err := r.Vm.Segments.LoadData(r.ProgramBase, &r.Program.Data)
	if err != nil {
		_, err = r.Vm.Segments.LoadData(r.ExecutionBase, stack)
	}
	// Mark data segment as accessed
	return err
}

// Initializes memory, initial register values & returns the end pointer (final pc) to run from a given pc offset
// (entrypoint)
func (r *CairoRunner) InitializeFunctionEntrypoint(entrypoint uint, stack *[]memory.MaybeRelocatable, return_fp memory.Relocatable) (memory.Relocatable, error) {
	end := r.Vm.Segments.AddSegment()
	*stack = append(*stack, *memory.NewMaybeRelocatableRelocatable(end), *memory.NewMaybeRelocatableRelocatable(return_fp))
	r.InitialFp = r.ExecutionBase
	r.InitialFp.Offset += uint(len(*stack))
	r.InitialAp = r.InitialFp
	r.FinalPc = end
	return end, r.initializeState(entrypoint, stack)
}
