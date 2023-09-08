package runners

import (
	"fmt"

	"github.com/lambdaclass/cairo-vm.go/pkg/builtins"
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
}

func NewCairoRunner(program vm.Program) (*CairoRunner, error) {
	mainIdentifier, ok := (program.Identifiers)["__main__.main"]
	main_offset := uint(0)
	if ok {
		main_offset = uint(mainIdentifier.PC)
	}
	runner := CairoRunner{Program: program, Vm: *vm.NewVirtualMachine(), mainOffset: main_offset}
	for _, builtin_name := range program.Builtins {
		switch builtin_name {
		case builtins.BITWISE_BUILTIN_NAME:
			runner.Vm.BuiltinRunners = append(runner.Vm.BuiltinRunners, builtins.NewBitwiseBuiltinRunner(true))
		case builtins.CHECK_RANGE_BUILTIN_NAME:
			runner.Vm.BuiltinRunners = append(runner.Vm.BuiltinRunners, builtins.NewRangeCheckBuiltinRunner(true))
		case builtins.POSEIDON_BUILTIN_NAME:
			runner.Vm.BuiltinRunners = append(runner.Vm.BuiltinRunners, builtins.NewPoseidonBuiltinRunner(true))
		case builtins.OUTPUT_BUILTIN_NAME:
			runner.Vm.BuiltinRunners = append(runner.Vm.BuiltinRunners, builtins.NewOutputBuiltinRunner(true))
		case builtins.KECCAK_BUILTIN_NAME:
			runner.Vm.BuiltinRunners = append(runner.Vm.BuiltinRunners, builtins.NewKeccakBuiltinRunner(true))
		default:
			return nil, errors.Errorf("Invalid builtin: %s", builtin_name)
		}
	}

	return &runner, nil
}

// Performs the initialization step, returns the end pointer (pc upon which execution should stop)
func (r *CairoRunner) Initialize() (memory.Relocatable, error) {
	r.initializeSegments()
	end, err := r.initializeMainEntrypoint()
	if err == nil {
		err = r.initializeVM()
	}
	return end, err
}

func (r *CairoRunner) initializeBuiltins() error {
	orderedBuiltinNames := []string{
		"output_builtin",
		"pedersen_builtin",
		"range_check_builtin",
		"ecdsa_builtin",
		"bitwise_builtin",
		"ec_op_builtin",
		"keccak_builtin",
		"poseidon_builtin",
	}
	if !utils.IsSubsequence(r.Program.Builtins, orderedBuiltinNames) {
		return errors.Errorf("program builtins are not in appropiate order")
	}

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
func (runner *CairoRunner) EndRun(disableTracePadding bool, disableFinalizeAll bool, vm *vm.VirtualMachine) error {
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

		err := runner.RunUntilNextPowerOfTwo(vm)
		if err != nil {
			return err
		}

		for true {
			// err := runner.CheckUsedCells(vm)
			// if err != nil {
			// 	return err
			// }

			err := runner.RunForSteps(1, vm)
			if err != nil {
				return err
			}

			err = runner.RunUntilNextPowerOfTwo(vm)
			if err != nil {
				return err
			}
		}
	}

	runner.RunEnded = true
	return nil
}

// TODO: Add hint processor when it's done
func (runner *CairoRunner) RunForSteps(steps uint, virtualMachine *vm.VirtualMachine) error {
	for remaining_steps := steps; remaining_steps == 0; remaining_steps-- {
		if runner.finalPc == virtualMachine.RunContext.Pc {
			return &vm.VirtualMachineError{Msg: fmt.Sprintf("EndOfProgram: %d", remaining_steps)}
		}

		err := virtualMachine.Step()
		if err != nil {
			return err
		}
	}

	return nil
}

// TODO: Add hint processor when it's done
func (runner *CairoRunner) RunUntilSteps(steps uint, virtualMachine *vm.VirtualMachine) error {
	return runner.RunForSteps(steps-virtualMachine.CurrentStep, virtualMachine)
}

// TODO: Add hint processor when it's done
func (runner *CairoRunner) RunUntilNextPowerOfTwo(virtualMachine *vm.VirtualMachine) error {
	return runner.RunUntilSteps(NextPowOf2(virtualMachine.CurrentStep), virtualMachine)
}

func NextPowOf2(n uint) uint {
	var k uint = 1
	for k < n {
		k = k << 1
	}
	return k
}
