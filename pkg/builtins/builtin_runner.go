package builtins

import (
	"fmt"
	"sort"

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
	// Cells per builtin instance
	CellsPerInstance() uint
	// Input cells per builtin instance
	InputCellsPerInstance() uint
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
	// Returns the builtin's ratio, is zero if the layout is dynamic
	Ratio() uint
	// Returns the builtin's allocated memory units
	GetAllocatedMemoryUnits(segments *memory.MemorySegmentManager, currentStep uint) (uint, error)
	// Returns the list of memory addresses used by the builtin
	GetMemoryAccesses(*memory.MemorySegmentManager) ([]memory.Relocatable, error)
	GetRangeCheckUsage(*memory.Memory) (*uint, *uint)
	GetUsedPermRangeCheckLimits(segments *memory.MemorySegmentManager, currentStep uint) (uint, error)
	GetUsedDilutedCheckUnits(dilutedSpacing uint, dilutedNBits uint) uint
	GetUsedCellsAndAllocatedSizes(segments *memory.MemorySegmentManager, currentStep uint) (uint, uint, error)
	FinalStack(segments *memory.MemorySegmentManager, pointer memory.Relocatable) (memory.Relocatable, error)
	// Returns the base & stop_ptr
	GetMemorySegmentAddresses() (memory.Relocatable, memory.Relocatable, error)
	// Amount of builtin instances used
	GetUsedInstances(*memory.MemorySegmentManager) (uint, error)
}

func RunSecurityChecksForBuiltin(builtin BuiltinRunner, segments *memory.MemorySegmentManager) error {
	if builtin.Name() == OUTPUT_BUILTIN_NAME {
		return nil
	}

	cellsPerInstance := builtin.CellsPerInstance()
	nInputCells := builtin.InputCellsPerInstance()
	builtinSegmentIndex := builtin.Base().SegmentIndex

	offsets := make([]int, 0)
	// Collect the builtin segment's address' offsets
	for addr := range segments.Memory.Data {
		if addr.SegmentIndex == builtinSegmentIndex {
			offsets = append(offsets, int(addr.Offset))
		}
	}

	if len(offsets) == 0 {
		// No checks to run for empty segment
		return nil
	}
	// Sort offsets for easier comparison
	sort.Ints(offsets)
	// Obtain max offset
	maxOffset := offsets[len(offsets)-1]

	n := (maxOffset / int(cellsPerInstance)) + 1
	//Verify that n is not too large to make sure the expectedOffsets list that is constructed below is not too large.
	if n > len(offsets)/int(nInputCells) {
		return errors.Errorf("Missing memory cells for %s", builtin.Name())
	}

	// Check that the two inputs (x and y) of each instance are set.
	expectedOffsets := make([]int, 0)
	for i := 0; i < n; i++ {
		for j := 0; j < int(nInputCells); j++ {
			expectedOffsets = append(expectedOffsets, int(cellsPerInstance)*i+j)
		}
	}
	// Find the missing offsets (offsets in expectedOffsets but not in offsets)
	missingOffsets := make([]int, 0)
	j := 0
	i := 0
	for i < len(expectedOffsets) && j < len(offsets) {
		if expectedOffsets[i] < offsets[j] {
			missingOffsets = append(missingOffsets, expectedOffsets[i])
		} else {
			j++
		}
		i++
	}
	for i < len(expectedOffsets) {
		missingOffsets = append(missingOffsets, expectedOffsets[i])
		i++
	}
	if len(missingOffsets) != 0 {
		return errors.Errorf("Missing memory cells for builtin: %s: %v", builtin.Name(), missingOffsets)
	}

	// Verify auto deduction rules for the unassigned output cells.
	// Assigned output cells are checked as part of the call to VerifyAutoDeductions().
	for i := uint(0); i < uint(n); i++ {
		for j := uint(nInputCells); j < cellsPerInstance; j++ {
			addr := memory.NewRelocatable(builtinSegmentIndex, cellsPerInstance*i+j)
			_, err := segments.Memory.Get(addr)
			// Output cell not in memory
			if err != nil {
				_, err = builtin.DeduceMemoryCell(addr, &segments.Memory)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}
