package builtins

import (
	"github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
	"github.com/lambdaclass/cairo-vm.go/pkg/utils"

	"github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
	"github.com/pkg/errors"
)

const BITWISE_BUILTIN_NAME = "bitwise"
const BITWISE_CELLS_PER_INSTANCE = 5
const BITWISE_TOTAL_N_BITS = 251
const BITWISE_INPUT_CELLS_PER_INSTANCE = 2

type BitwiseBuiltinRunner struct {
	base                  memory.Relocatable
	ratio                 uint
	instancesPerComponent uint
	included              bool
	TotalNBits            uint
	StopPtr               *uint
}

func BitwiseError(err error) error {
	return errors.Wrapf(err, "Bitwise builtin error\n")
}

func ErrFeltBiggerThanPowerOfTwo(felt lambdaworks.Felt) error {
	return BitwiseError(errors.Errorf("Expected felt %d to be smaller than  2**%d", felt, BITWISE_TOTAL_N_BITS))
}

func NewBitwiseBuiltinRunner(ratio uint) *BitwiseBuiltinRunner {
	return &BitwiseBuiltinRunner{ratio: ratio, instancesPerComponent: 1, TotalNBits: BITWISE_TOTAL_N_BITS}
}

func DefaultBitwiseBuiltinRunner() *BitwiseBuiltinRunner {
	return NewBitwiseBuiltinRunner(256)
}

func (b *BitwiseBuiltinRunner) Base() memory.Relocatable {
	return b.base
}

func (b *BitwiseBuiltinRunner) Name() string {
	return BITWISE_BUILTIN_NAME
}

func (b *BitwiseBuiltinRunner) InitializeSegments(segments *memory.MemorySegmentManager) {
	b.base = segments.AddSegment()
}

func (b *BitwiseBuiltinRunner) InitialStack() []memory.MaybeRelocatable {
	if b.included {
		return []memory.MaybeRelocatable{*memory.NewMaybeRelocatableRelocatable(b.base)}
	} else {
		return []memory.MaybeRelocatable{}
	}
}

func (b *BitwiseBuiltinRunner) DeduceMemoryCell(address memory.Relocatable, mem *memory.Memory) (*memory.MaybeRelocatable, error) {
	index := address.Offset % BITWISE_CELLS_PER_INSTANCE
	if index < BITWISE_INPUT_CELLS_PER_INSTANCE {
		return nil, nil
	}

	x_addr, _ := address.SubUint(index)
	y_addr := x_addr.AddUint(1)

	num_x_felt, err := mem.GetFelt(x_addr)
	if err != nil {
		return nil, nil
	}
	num_y_felt, err := mem.GetFelt(y_addr)
	if err != nil {
		return nil, nil
	}

	if num_x_felt.Bits() > BITWISE_TOTAL_N_BITS {
		return nil, ErrFeltBiggerThanPowerOfTwo(num_x_felt)
	}
	if num_y_felt.Bits() > BITWISE_TOTAL_N_BITS {
		return nil, ErrFeltBiggerThanPowerOfTwo(num_y_felt)
	}

	var res *memory.MaybeRelocatable
	switch index {
	case 2:
		res = memory.NewMaybeRelocatableFelt(num_x_felt.And(num_y_felt))
	case 3:
		res = memory.NewMaybeRelocatableFelt(num_x_felt.Xor(num_y_felt))
	case 4:
		res = memory.NewMaybeRelocatableFelt(num_x_felt.Or(num_y_felt))
	default:
		res = nil
	}
	return res, nil
}

func (b *BitwiseBuiltinRunner) AddValidationRule(*memory.Memory) {}

func (b *BitwiseBuiltinRunner) Include(include bool) {
	b.included = include
}

func (b *BitwiseBuiltinRunner) Ratio() uint {
	return b.ratio
}

func (b *BitwiseBuiltinRunner) CellsPerInstance() uint {
	return BITWISE_CELLS_PER_INSTANCE
}

func (b *BitwiseBuiltinRunner) InputCellsPerInstance() uint {
	return BITWISE_INPUT_CELLS_PER_INSTANCE
}

func (b *BitwiseBuiltinRunner) GetAllocatedMemoryUnits(segments *memory.MemorySegmentManager, currentStep uint) (uint, error) {
	// This condition corresponds to an uninitialized ratio for the builtin, which should only
	// happen when layout is `dynamic`
	if b.Ratio() == 0 {
		// Dynamic layout has the exact number of instances it needs (up to a power of 2).
		used, err := segments.GetSegmentUsedSize(uint(b.base.SegmentIndex))
		if err != nil {
			return 0, err
		}
		instances := used / b.CellsPerInstance()
		components := utils.NextPowOf2(instances / b.instancesPerComponent)
		size := b.CellsPerInstance() * b.instancesPerComponent * components

		return size, nil
	}

	minStep := b.ratio * b.instancesPerComponent
	if currentStep < minStep {
		return 0, memory.InsufficientAllocatedCellsErrorMinStepNotReached(minStep, b.Name())
	}
	value, err := utils.SafeDiv(currentStep, b.ratio)

	if err != nil {
		return 0, errors.Errorf("error calculating builtin memory units: %s", err)
	}

	return b.CellsPerInstance() * value, nil
}

func (b *BitwiseBuiltinRunner) GetUsedCellsAndAllocatedSizes(segments *memory.MemorySegmentManager, currentStep uint) (uint, uint, error) {
	used, err := segments.GetSegmentUsedSize(uint(b.base.SegmentIndex))
	if err != nil {
		return 0, 0, err
	}

	size, err := b.GetAllocatedMemoryUnits(segments, currentStep)
	if err != nil {
		return 0, 0, err
	}

	if used > size {
		return 0, 0, memory.InsufficientAllocatedCellsErrorWithBuiltinName(b.Name(), used, size)
	}

	return used, size, nil
}

func (runner *BitwiseBuiltinRunner) GetRangeCheckUsage(memory *memory.Memory) (*uint, *uint) {
	return nil, nil
}

func (runner *BitwiseBuiltinRunner) GetUsedPermRangeCheckLimits(segments *memory.MemorySegmentManager, currentStep uint) (uint, error) {
	return 0, nil
}

func (runner *BitwiseBuiltinRunner) GetUsedDilutedCheckUnits(dilutedSpacing uint, dilutedNBits uint) uint {
	totalNBits := runner.TotalNBits
	partition := make([]uint, 0)
	var i uint
	for i = 0; i < totalNBits; i += (dilutedSpacing * dilutedNBits) {
		var j uint
		for j = 0; j < dilutedSpacing; j++ {
			if i+j < totalNBits {
				partition = append(partition, i+j)
			}
		}
	}

	partitionLength := uint(len(partition))
	var numTrimmed uint

	for _, element := range partition {
		if (element + dilutedSpacing*(dilutedNBits-1) + 1) > totalNBits {
			numTrimmed++
		}
	}

	return 4*partitionLength + numTrimmed
}

func (runner *BitwiseBuiltinRunner) GetMemoryAccesses(manager *memory.MemorySegmentManager) ([]memory.Relocatable, error) {
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

func (r *BitwiseBuiltinRunner) FinalStack(segments *memory.MemorySegmentManager, pointer memory.Relocatable) (memory.Relocatable, error) {
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

func (r *BitwiseBuiltinRunner) GetUsedInstances(segments *memory.MemorySegmentManager) (uint, error) {
	usedCells, err := segments.GetSegmentUsedSize(uint(r.Base().SegmentIndex))
	if err != nil {
		return 0, nil
	}

	return utils.DivCeil(usedCells, r.CellsPerInstance()), nil
}

func (b *BitwiseBuiltinRunner) GetMemorySegmentAddresses() (memory.Relocatable, memory.Relocatable, error) {
	if b.StopPtr == nil {
		return memory.Relocatable{}, memory.Relocatable{}, NewErrNoStopPointer(b.Name())
	}
	return b.base, memory.NewRelocatable(b.base.SegmentIndex, *b.StopPtr), nil
}
