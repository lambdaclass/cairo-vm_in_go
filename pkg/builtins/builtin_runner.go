package builtins

import (
	"fmt"

	"github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
	"github.com/pkg/errors"
)

var ErrNoStopPointer = errors.New("No Stop Pointer")
var ErrInvalidStopPointerIndex = errors.New("Invalid Stop Pointer Index")
var ErrInvalidStopPointer = errors.New("Invalid Stop Pointer")

func NewErrNoStopPointer(builtinName string) error {
	return fmt.Errorf("%w builtin: %s", ErrNoStopPointer, builtinName)
}

func NewErrInvalidStopPointerIndex(builtinName string, stopPtr memory.Relocatable, base memory.Relocatable) error {
	return fmt.Errorf("%w builtin: %s stopPtr: (%d, %d) base: (%d, %d)", ErrInvalidStopPointerIndex, builtinName, stopPtr.SegmentIndex, stopPtr.Offset, base.SegmentIndex, base.Offset)
}

func NewErrInvalidStopPointer(builtinName string, used uint, stopPtr memory.Relocatable) error {
	return fmt.Errorf("%w builtin: %s used: (%d, %d) stopPtr: (%d, %d)", ErrInvalidStopPointer, builtinName, stopPtr.SegmentIndex, used, stopPtr.SegmentIndex, stopPtr.Offset)
}

type BuiltinRunner interface {
	// Returns the first address of the builtin's memory segment
	Base() memory.Relocatable
	// Returns the name of the builtin
	Name() string
	// Creates a memory segment for the builtin and initializes its base
	InitializeSegments(*memory.MemorySegmentManager)
	// Returns the builtin's initial stack
	InitialStack() []memory.MaybeRelocatable
	// Attempts to deduce the value of a memory cell given by its address. Can return either a nil pointer and an error, if an error arises during the deduction,
	// a valid pointer and nil if the deduction was successful, or a nil pointer and nil if there is no deduction for the memory cell
	DeduceMemoryCell(memory.Relocatable, *memory.Memory) (*memory.MaybeRelocatable, error)
	// Adds a validation rule to the memory
	// Validation rules are applied when a value is inserted into the builtin's segment
	AddValidationRule(*memory.Memory)
	// Sets the inclusion of the Builtin Runner in the Cairo Runner
	Include(bool)
	// TODO: Later additions -> Some of them could depend on a Default Implementation
	// // Most of them depend on Layouts being implemented
	// // Use cases:
	// // I. PROOF_MODE
	// Returns the builtin's ratio, is zero if the layout is dynamic
	Ratio() uint
	// Returns the builtin's allocated memory units
	GetAllocatedMemoryUnits(segments *memory.MemorySegmentManager, currentStep uint) (uint, error)
	// // Returns the list of memory addresses used by the builtin
	GetMemoryAccesses(*memory.MemorySegmentManager) ([]memory.Relocatable, error)
	GetRangeCheckUsage(*memory.Memory) (*uint, *uint)
	GetUsedPermRangeCheckLimits(segments *memory.MemorySegmentManager, currentStep uint) (uint, error)
	GetUsedDilutedCheckUnits(dilutedSpacing uint, dilutedNBits uint) uint
	GetUsedCellsAndAllocatedSizes(segments *memory.MemorySegmentManager, currentStep uint) (uint, uint, error)
	FinalStack(segments *memory.MemorySegmentManager, pointer memory.Relocatable) (memory.Relocatable, error)
	// // II. SECURITY (secure-run flag cairo-run || verify-secure flag run_from_entrypoint)
	// RunSecurityChecks(*vm.VirtualMachine) error // verify_secure_runner logic
	// // Returns the base & stop_ptr, stop_ptr can be nil
	GetMemorySegmentAddresses() (memory.Relocatable, memory.Relocatable, error)
	// // III. STARKNET-SPECIFIC
	GetUsedInstances(*memory.MemorySegmentManager) (uint, error)
	// // IV. GENERAL CASE (but not critical)
	// FinalStack(*memory.MemorySegmentManager, memory.Relocatable) (memory.Relocatable, error) // read_return_values
}
