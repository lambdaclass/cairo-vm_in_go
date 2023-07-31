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
	MainOffset    uint
}

func NewCairoRunner(program vm.Program) *CairoRunner {
	// TODO: Fetch main entrypoint offset from program identifiers
	// Placeholder
	main_offset := uint(0)
	return &CairoRunner{Program: program, Vm: *vm.NewVirtualMachine(), MainOffset: main_offset}

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

// Initializes memory, initial register values & returns the end pointer (final pc) to run from the main entrypoint
func (r *CairoRunner) InitializeMainEntrypoint() (memory.Relocatable, error) {
	// When running from main entrypoint, only up to 11 values will be written (9 builtin bases + end + return_fp)
	stack := make([]memory.MaybeRelocatable, 0, 11)
	// Append builtins initial stack to stack
	// Handle proof-mode specific behaviour
	return_fp := r.Vm.Segments.AddSegment()
	return r.InitializeFunctionEntrypoint(r.MainOffset, &stack, return_fp)
}
