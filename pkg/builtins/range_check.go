package builtins

import (
	"math"

	"github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
	"github.com/lambdaclass/cairo-vm.go/pkg/utils"
	"github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
	"github.com/pkg/errors"
)

const RANGE_CHECK_BUILTIN_NAME = "range_check"
const INNER_RC_BOUND_SHIFT = 16
const INNER_RC_BOUND_MASK = math.MaxUint16
const INNER_RC_BOUND uint64 = 1 << INNER_RC_BOUND_SHIFT
const CELLS_PER_RANGE_CHECK = 1

const RANGE_CHECK_N_PARTS = 8

func RangeCheckError(err error) error {
	return errors.Wrapf(err, "Range check error")
}

func OutsideBoundsError(felt lambdaworks.Felt) error {
	return RangeCheckError(errors.Errorf("Value %d is out of bounds [0, 2^128]", felt))
}

func NotAFeltError(addr memory.Relocatable, val memory.MaybeRelocatable) error {
	rel, _ := val.GetRelocatable()
	return RangeCheckError(errors.Errorf("Value %d found in %d is not a field element", rel, addr))
}

type RangeCheckBuiltinRunner struct {
	base                  memory.Relocatable
	included              bool
	ratio                 uint
	instancesPerComponent uint
	StopPtr               *uint
}

func NewRangeCheckBuiltinRunner(ratio uint) *RangeCheckBuiltinRunner {

	return &RangeCheckBuiltinRunner{ratio: ratio, instancesPerComponent: 1}
}

func DefaultRangeCheckBuiltinRunner() *RangeCheckBuiltinRunner {
	return NewRangeCheckBuiltinRunner(8)
}

func (r *RangeCheckBuiltinRunner) Base() memory.Relocatable {
	return r.base
}

func (r *RangeCheckBuiltinRunner) Name() string {
	return RANGE_CHECK_BUILTIN_NAME
}

func (r *RangeCheckBuiltinRunner) Bound() lambdaworks.Felt {
	bound := lambdaworks.FeltOne().Shl(16 * uint64(RANGE_CHECK_N_PARTS))
	return bound
}

func (r *RangeCheckBuiltinRunner) SetBase(value memory.Relocatable) {
	r.base = value
}

func (r *RangeCheckBuiltinRunner) InitializeSegments(segments *memory.MemorySegmentManager) {
	r.base = segments.AddSegment()
}

func (r *RangeCheckBuiltinRunner) InitialStack() []memory.MaybeRelocatable {
	if r.included {
		return []memory.MaybeRelocatable{*memory.NewMaybeRelocatableRelocatable(r.base)}
	}
	return []memory.MaybeRelocatable{}
}

func (r *RangeCheckBuiltinRunner) DeduceMemoryCell(addr memory.Relocatable, mem *memory.Memory) (*memory.MaybeRelocatable, error) {
	return nil, nil
}

func RangeCheckValidationRule(mem *memory.Memory, address memory.Relocatable) ([]memory.Relocatable, error) {
	res_val, err := mem.Get(address)
	if err != nil {
		return nil, err
	}
	felt, is_felt := res_val.GetFelt()
	if !is_felt {
		return nil, NotAFeltError(address, *res_val)
	}
	if felt.Bits() <= RANGE_CHECK_N_PARTS*INNER_RC_BOUND_SHIFT {
		return []memory.Relocatable{address}, nil
	}
	return nil, OutsideBoundsError(felt)
}

func (r *RangeCheckBuiltinRunner) AddValidationRule(mem *memory.Memory) {
	mem.AddValidationRule(uint(r.base.SegmentIndex), RangeCheckValidationRule)
}

func (r *RangeCheckBuiltinRunner) Include(include bool) {
	r.included = include
}

func (r *RangeCheckBuiltinRunner) Ratio() uint {
	return r.ratio
}

func (r *RangeCheckBuiltinRunner) CellsPerInstance() uint {
	return CELLS_PER_RANGE_CHECK
}

func (r *RangeCheckBuiltinRunner) GetAllocatedMemoryUnits(segments *memory.MemorySegmentManager, currentStep uint) (uint, error) {
	// This condition corresponds to an uninitialized ratio for the builtin, which should only
	// happen when layout is `dynamic`
	if r.Ratio() == 0 {
		// Dynamic layout has the exact number of instances it needs (up to a power of 2).
		used, err := segments.GetSegmentUsedSize(uint(r.base.SegmentIndex))
		if err != nil {
			return 0, err
		}
		instances := used / r.CellsPerInstance()
		components := utils.NextPowOf2(instances / r.instancesPerComponent)
		size := r.CellsPerInstance() * r.instancesPerComponent * components

		return size, nil
	}

	minStep := r.Ratio() * r.instancesPerComponent
	if currentStep < minStep {
		return 0, memory.InsufficientAllocatedCellsErrorMinStepNotReached(minStep, r.Name())
	}
	value, err := utils.SafeDiv(currentStep, r.Ratio())

	if err != nil {
		return 0, errors.Errorf("error calculating builtin memory units: %s", err)
	}

	return r.CellsPerInstance() * value, nil
}

func (r *RangeCheckBuiltinRunner) GetUsedCellsAndAllocatedSizes(segments *memory.MemorySegmentManager, currentStep uint) (uint, uint, error) {
	used, err := segments.GetSegmentUsedSize(uint(r.base.SegmentIndex))
	if err != nil {
		return 0, 0, err
	}

	size, err := r.GetAllocatedMemoryUnits(segments, currentStep)
	if err != nil {
		return 0, 0, err
	}

	if err != nil {
		return 0, 0, err
	}

	if used > size {
		return 0, 0, memory.InsufficientAllocatedCellsErrorWithBuiltinName(r.Name(), used, size)
	}

	return used, size, nil
}

func (runner *RangeCheckBuiltinRunner) GetRangeCheckUsage(memory *memory.Memory) (*uint, *uint) {
	rangeCheckSegment := memory.GetSegment(runner.base.SegmentIndex)

	if rangeCheckSegment == nil {
		return nil, nil
	}

	var rcMin = new(uint)
	var rcMax = new(uint)

	for _, value := range rangeCheckSegment {
		feltValue, isFelt := value.GetFelt()

		if !isFelt {
			// Is this correct? Feels a little weird, like we should return an error
			return nil, nil
		}

		feltDigits := feltValue.ToLeBytes()
		for i := 0; i < 32; i += 2 {
			var tempValue = (uint16(feltDigits[i+1]) << 8) | uint16((feltDigits[i]))

			if rcMin == nil {
				*rcMin = uint(tempValue)
			}

			if rcMax == nil {
				*rcMax = uint(tempValue)
			}

			if uint(tempValue) < *rcMin {
				*rcMin = uint(tempValue)
			}

			if uint(tempValue) > *rcMax {
				*rcMax = uint(tempValue)
			}
		}
	}

	return rcMin, rcMax
}

func (runner *RangeCheckBuiltinRunner) GetUsedPermRangeCheckLimits(segments *memory.MemorySegmentManager, currentStep uint) (uint, error) {
	usedCells, _, err := runner.GetUsedCellsAndAllocatedSizes(segments, currentStep)

	if err != nil {
		return 0, err
	}

	return usedCells * RANGE_CHECK_N_PARTS, nil
}

func (runner *RangeCheckBuiltinRunner) GetUsedDilutedCheckUnits(dilutedSpacing uint, dilutedNBits uint) uint {
	return 0
}

func (runner *RangeCheckBuiltinRunner) GetMemoryAccesses(manager *memory.MemorySegmentManager) ([]memory.Relocatable, error) {
	segmentSize, err := manager.GetSegmentSize(uint(runner.Base().SegmentIndex))
	if err != nil {
		return []memory.Relocatable{}, err
	}

	var ret []memory.Relocatable

	var i uint
	for i = 0; i < segmentSize; i++ {
		ret = append(ret, memory.NewRelocatable(runner.Base().SegmentIndex, i))
	}

	return ret, nil
}

func (r *RangeCheckBuiltinRunner) FinalStack(segments *memory.MemorySegmentManager, pointer memory.Relocatable) (memory.Relocatable, error) {
	if r.included {
		if pointer.Offset == 0 {
			return memory.Relocatable{}, NewErrNoStopPointer(r.Name())
		}

		stopPointerAddr := memory.NewRelocatable(pointer.SegmentIndex, pointer.Offset-1)

		stopPointer, err := segments.Memory.GetRelocatable(stopPointerAddr)
		if err != nil {
			return memory.Relocatable{}, err
		}

		if r.Base().SegmentIndex != stopPointer.SegmentIndex {
			return memory.Relocatable{}, NewErrInvalidStopPointerIndex(r.Name(), stopPointer, r.Base())
		}

		numInstances, err := r.GetUsedInstances(segments)
		if err != nil {
			return memory.Relocatable{}, err
		}

		used := numInstances * r.CellsPerInstance()

		if stopPointer.Offset != used {
			return memory.Relocatable{}, NewErrInvalidStopPointer(r.Name(), used, stopPointer)
		}

		r.StopPtr = &stopPointer.Offset

		return stopPointerAddr, nil
	} else {
		r.StopPtr = new(uint)
		*r.StopPtr = 0
		return pointer, nil
	}
}

func (r *RangeCheckBuiltinRunner) GetUsedInstances(segments *memory.MemorySegmentManager) (uint, error) {
	usedCells, err := segments.GetSegmentUsedSize(uint(r.Base().SegmentIndex))
	if err != nil {
		return 0, nil
	}

	return usedCells, nil
}
