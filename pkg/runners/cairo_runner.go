package runners

import (
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

func (r *CairoRunner) RunUntilPC(end memory.Relocatable) error {
	for r.Vm.RunContext.Pc != end {
		err := r.Vm.Step()
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

	return nil
}

// pub fn end_run(
// 	&mut self,
// 	disable_trace_padding: bool,
// 	disable_finalize_all: bool,
// 	vm: &mut VirtualMachine,
// 	hint_processor: &mut dyn HintProcessor,
// ) -> Result<(), VirtualMachineError> {
// 	if self.run_ended {
// 		return Err(RunnerError::EndRunCalledTwice.into());
// 	}

// 	vm.segments.memory.relocate_memory()?;
// 	vm.end_run(&self.exec_scopes)?;

// 	if disable_finalize_all {
// 		return Ok(());
// 	}

// 	vm.segments.compute_effective_sizes();
// 	if self.proof_mode && !disable_trace_padding {
// 		self.run_until_next_power_of_2(vm, hint_processor)?;
// 		loop {
// 			match self.check_used_cells(vm) {
// 				Ok(_) => break,
// 				Err(e) => match e {
// 					VirtualMachineError::Memory(MemoryError::InsufficientAllocatedCells(_)) => {
// 					}
// 					e => return Err(e),
// 				},
// 			}

// 			self.run_for_steps(1, vm, hint_processor)?;
// 			self.run_until_next_power_of_2(vm, hint_processor)?;
// 		}
// 	}

// 	self.run_ended = true;
// 	Ok(())
// }
