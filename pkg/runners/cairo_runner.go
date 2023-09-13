package runners

import (
	"fmt"

	"github.com/lambdaclass/cairo-vm.go/pkg/builtins"
	"github.com/lambdaclass/cairo-vm.go/pkg/layouts"
	"github.com/lambdaclass/cairo-vm.go/pkg/utils"
	"github.com/lambdaclass/cairo-vm.go/pkg/vm"
	"github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
	"github.com/pkg/errors"
)

var ErrRunnerCalledTwice = errors.New("Cairo Runner was called twice")

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
	ProofMode     bool
	RunEnded      bool
	layout        layouts.CairoLayout
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
		ProofMode:  proofMode,
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

// Initializes builtin runners in accordance to the specified layout and
// the builtins present in the running program.
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
		} else if r.ProofMode {
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

// TODO: Add HintProcessor as parameter once we have that
func (runner *CairoRunner) EndRun(disableTracePadding bool, disableFinalizeAll bool, vm *vm.VirtualMachine, hintProcessor vm.HintProcessor) error {
	if runner.RunEnded {
		return ErrRunnerCalledTwice
	}

	// TODO: This seems to have to do with temporary segments
	// vm.Segments.Memory.RelocateMemory()

	err := vm.EndRun()
	if err != nil {
		return err
	}

	if disableFinalizeAll {
		return nil
	}

	vm.Segments.ComputeEffectiveSizes()
	if runner.ProofMode && !disableTracePadding {

		err := runner.RunUntilNextPowerOfTwo(vm, hintProcessor)
		if err != nil {
			return err
		}

		for true {
			err := runner.CheckUsedCells(vm)
			if err != nil {
				return err
			}

			err = runner.RunForSteps(1, vm, hintProcessor)
			if err != nil {
				return err
			}

			err = runner.RunUntilNextPowerOfTwo(vm, hintProcessor)
			if err != nil {
				return err
			}
		}
	}

	runner.RunEnded = true
	return nil
}

func (runner *CairoRunner) CheckUsedCells(virtualMachine *vm.VirtualMachine) error {
	for _, builtin := range virtualMachine.BuiltinRunners {
		// I guess we call this just in case it errors out, even though later on we also call it? Not bad
		_, _, err := builtin.GetUsedCellsAndAllocatedSizes(&virtualMachine.Segments, virtualMachine.CurrentStep)
		if err != nil {
			return err
		}
	}

	err := runner.CheckRangeCheckUsage(virtualMachine)
	if err != nil {
		return err
	}

	// err = runner.CheckMemoryUsage(virtualMachine)
	// if err != nil {
	// 	return err
	// }

	// err = runner.CheckDilutecHeckUsage(virtualMachine)
	// if err != nil {
	// 	return err
	// }

	return nil
}

// // Returns Ok(()) if there are enough allocated cells for the builtins.
//     // If not, the number of steps should be increased or a different layout should be used.
//     pub fn check_used_cells(&self, vm: &VirtualMachine) -> Result<(), VirtualMachineError> {
//         vm.builtin_runners
//             .iter()
//             .map(|builtin_runner| builtin_runner.get_used_cells_and_allocated_size(vm))
//             .collect::<Result<Vec<(usize, usize)>, MemoryError>>()?;
//         self.check_range_check_usage(vm)?;
//         self.check_memory_usage(vm)?;
//         self.check_diluted_check_usage(vm)?;
//         Ok(())
//     }

func (runner *CairoRunner) CheckRangeCheckUsage(virtualMachine *vm.VirtualMachine) error {
	var rcMin, rcMax *uint

	for _, builtin := range runner.Vm.BuiltinRunners {
		resultMin, resultMax := builtin.GetRangeCheckUsage(&runner.Vm.Segments.Memory)

		if resultMin != nil {
			rcMin = resultMin
		}

		if resultMax != nil {
			rcMax = resultMax
		}
	}

	if rcMin == nil || rcMax == nil {
		return nil
	}

	if runner.Vm.RcLimitsMax != nil && (*runner.Vm.RcLimitsMax > *rcMax) {
		rcMax = runner.Vm.RcLimitsMax
	}

	if runner.Vm.RcLimitsMin != nil && (*runner.Vm.RcLimitsMin > *rcMin) {
		rcMin = runner.Vm.RcLimitsMin
	}

	var rcUnitsUsedByBuiltins uint = 0

	for _, builtin := range runner.Vm.BuiltinRunners {
		usedUnits, err := builtin.GetUsedPermRangeCheckLimits(&virtualMachine.Segments, virtualMachine.CurrentStep)
		if err != nil {
			return err
		}

		rcUnitsUsedByBuiltins += usedUnits
	}

	unusedRcUnits := (runner.layout.RcUnits-3)*virtualMachine.CurrentStep - uint(rcUnitsUsedByBuiltins)

	if unusedRcUnits < (*rcMax - *rcMin) {
		return errors.Errorf("Insufficient Allocated Cells: Unused RC Units: %d, Size: %d", unusedRcUnits, (*rcMax - *rcMin))
	}

	return nil
}

// TODO: Add hint processor when it's done
func (runner *CairoRunner) RunForSteps(steps uint, virtualMachine *vm.VirtualMachine, hintProcessor vm.HintProcessor) error {
	hintDataMap, err := runner.BuildHintDataMap(hintProcessor)
	if err != nil {
		return err
	}
	constants := runner.Program.ExtractConstants()
	for remaining_steps := steps; remaining_steps == 0; remaining_steps-- {
		if runner.finalPc == virtualMachine.RunContext.Pc {
			return &vm.VirtualMachineError{Msg: fmt.Sprintf("EndOfProgram: %d", remaining_steps)}
		}

		err := virtualMachine.Step(hintProcessor, &hintDataMap, &constants)
		if err != nil {
			return err
		}
	}

	return nil
}

// TODO: Add hint processor when it's done
func (runner *CairoRunner) RunUntilSteps(steps uint, virtualMachine *vm.VirtualMachine, hintProcessor vm.HintProcessor) error {
	return runner.RunForSteps(steps-virtualMachine.CurrentStep, virtualMachine, hintProcessor)
}

// TODO: Add hint processor when it's done
func (runner *CairoRunner) RunUntilNextPowerOfTwo(virtualMachine *vm.VirtualMachine, hintProcessor vm.HintProcessor) error {
	return runner.RunUntilSteps(utils.NextPowOf2(virtualMachine.CurrentStep), virtualMachine, hintProcessor)
}
