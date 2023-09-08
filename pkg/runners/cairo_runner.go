package runners

import (
	"github.com/lambdaclass/cairo-vm.go/pkg/builtins"
	"github.com/lambdaclass/cairo-vm.go/pkg/layouts"
	"github.com/lambdaclass/cairo-vm.go/pkg/utils"
	"github.com/lambdaclass/cairo-vm.go/pkg/vm"
	"github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
	"github.com/pkg/errors"
)

type CairoRunner struct {
	Program       vm.Program
	Vm            vm.VirtualMachine
	ProgramBase   memory.Relocatable
	executionBase memory.Relocatable
	initialPc     memory.Relocatable
	initialAp     memory.Relocatable
	initialFp     memory.Relocatable
	finalPc       memory.Relocatable
	mainOffset    uint
	layout        layouts.CairoLayout
	proofMode     bool
}

func NewCairoRunner(program vm.Program, layoutName string, proofMode bool) (*CairoRunner, error) {
	mainIdentifier, ok := (program.Identifiers)["__main__.main"]
	main_offset := uint(0)
	if ok {
		main_offset = uint(mainIdentifier.PC)
	}

	err := utils.CheckBuiltinsSubsequence(program.Builtins)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	layoutBuiltinRunners, err := layouts.GetLayoutBuiltinRunners(layoutName)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	layout := layouts.CairoLayout{Name: layoutName, Builtins: layoutBuiltinRunners}
	runner := CairoRunner{
		Program:    program,
		Vm:         *vm.NewVirtualMachine(),
		mainOffset: main_offset,
		proofMode:  proofMode,
		layout:     layout,
	}
	return &runner, nil
}

// Performs the initialization step, returns the end pointer (pc upon which execution should stop)
func (r *CairoRunner) Initialize() (memory.Relocatable, error) {
	err := r.initializeBuiltins()
	if err != nil {
		return memory.Relocatable{}, errors.New(err.Error())
	}
	r.initializeSegments()
	end, err := r.initializeMainEntrypoint()
	if err == nil {
		err = r.initializeVM()
	}
	return end, err
}

func (r *CairoRunner) initializeBuiltins() error {
	var builtinRunners []builtins.BuiltinRunner
	programBuiltins := map[string]struct{}{}
	for _, builtin := range r.Program.Builtins {
		programBuiltins[builtin] = struct{}{}
	}

	for _, layoutBuiltin := range r.layout.Builtins {
		_, included := programBuiltins[layoutBuiltin.Name()]
		if included {
			delete(programBuiltins, layoutBuiltin.Name())
			layoutBuiltin.Include(true)
			builtinRunners = append(builtinRunners, layoutBuiltin)
		} else if r.proofMode {
			layoutBuiltin.Include(false)
			builtinRunners = append(builtinRunners, layoutBuiltin)
		}
	}

	if len(programBuiltins) != 0 {
		return errors.Errorf("Builtin(s) %v not present in layout %s", programBuiltins, r.layout.Name)
	}

	r.Vm.BuiltinRunners = builtinRunners

	return nil
}

// Creates program, execution and builtin segments
func (r *CairoRunner) initializeSegments() {
	// Program Segment
	r.ProgramBase = r.Vm.Segments.AddSegment()
	// Execution Segment
	r.executionBase = r.Vm.Segments.AddSegment()
	// Builtin Segments
	for i := range r.Vm.BuiltinRunners {
		r.Vm.BuiltinRunners[i].InitializeSegments(&r.Vm.Segments)
	}
}

// Initializes the program segment & initial pc
func (r *CairoRunner) initializeState(entrypoint uint, stack *[]memory.MaybeRelocatable) error {
	r.initialPc = r.ProgramBase
	r.initialPc.Offset += entrypoint
	// Load program data
	_, err := r.Vm.Segments.LoadData(r.ProgramBase, &r.Program.Data)
	if err == nil {
		_, err = r.Vm.Segments.LoadData(r.executionBase, stack)
	}
	// Mark data segment as accessed
	return err
}

// Initializes memory, initial register values & returns the end pointer (final pc) to run from a given pc offset
// (entrypoint)
func (r *CairoRunner) initializeFunctionEntrypoint(entrypoint uint, stack *[]memory.MaybeRelocatable, return_fp memory.Relocatable) (memory.Relocatable, error) {
	end := r.Vm.Segments.AddSegment()
	*stack = append(*stack, *memory.NewMaybeRelocatableRelocatable(return_fp), *memory.NewMaybeRelocatableRelocatable(end))
	r.initialFp = r.executionBase
	r.initialFp.Offset += uint(len(*stack))
	r.initialAp = r.initialFp
	r.finalPc = end
	return end, r.initializeState(entrypoint, stack)
}

// Initializes memory, initial register values & returns the end pointer (final pc) to run from the main entrypoint
func (r *CairoRunner) initializeMainEntrypoint() (memory.Relocatable, error) {
	// When running from main entrypoint, only up to 11 values will be written (9 builtin bases + end + return_fp)
	stack := make([]memory.MaybeRelocatable, 0, 11)
	// Append builtins initial stack to stack
	for i := range r.Vm.BuiltinRunners {
		for _, val := range r.Vm.BuiltinRunners[i].InitialStack() {
			stack = append(stack, val)
		}
	}
	// Handle proof-mode specific behaviour
	return_fp := r.Vm.Segments.AddSegment()
	return r.initializeFunctionEntrypoint(r.mainOffset, &stack, return_fp)
}

// Initializes the vm's run_context, adds builtin validation rules & validates memory
func (r *CairoRunner) initializeVM() error {
	r.Vm.RunContext.Ap = r.initialAp
	r.Vm.RunContext.Fp = r.initialFp
	r.Vm.RunContext.Pc = r.initialPc
	// Add validation rules
	for i := range r.Vm.BuiltinRunners {
		r.Vm.BuiltinRunners[i].AddValidationRule(&r.Vm.Segments.Memory)
	}
	// Apply validation rules to memory
	return r.Vm.Segments.Memory.ValidateExistingMemory()
}

func (r *CairoRunner) BuildHintDataMap(hintProcessor vm.HintProcessor) (map[uint][]any, error) {
	hintDataMap := make(map[uint][]any, 0)
	for pc, hintsParams := range r.Program.Hints {
		hintDatas := make([]any, 0, len(hintsParams))
		for _, hintParam := range hintsParams {
			data, err := hintProcessor.CompileHint(&hintParam, &r.Program.ReferenceManager)
			if err != nil {
				return nil, err
			}
			hintDatas = append(hintDatas, data)
		}
		hintDataMap[pc] = hintDatas
	}

	return hintDataMap, nil
}

func (r *CairoRunner) RunUntilPC(end memory.Relocatable, hintProcessor vm.HintProcessor) error {
	hintDataMap, err := r.BuildHintDataMap(hintProcessor)
	if err != nil {
		return err
	}
	constants := r.Program.ExtractConstants()
	for r.Vm.RunContext.Pc != end {
		err := r.Vm.Step(hintProcessor, &hintDataMap, &constants)
		if err != nil {
			return err
		}
	}
	return nil
}
