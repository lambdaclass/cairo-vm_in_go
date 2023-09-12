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
const BIWISE_INPUT_CELLS_PER_INSTANCE = 2

type BitwiseBuiltinRunner struct {
	base                  memory.Relocatable
	ratio                 uint
	instancesPerComponent uint
	included              bool
}

func BitwiseError(err error) error {
	return errors.Wrapf(err, "Bitwise builtin error\n")
}

func ErrFeltBiggerThanPowerOfTwo(felt lambdaworks.Felt) error {
	return BitwiseError(errors.Errorf("Expected felt %d to be smaller than  2**%d", felt, BITWISE_TOTAL_N_BITS))
}

func NewBitwiseBuiltinRunner(ratio uint) *BitwiseBuiltinRunner {
	return &BitwiseBuiltinRunner{ratio: ratio, instancesPerComponent: 1}
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
	if index < BIWISE_INPUT_CELLS_PER_INSTANCE {
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
		return 0, errors.Errorf("number of steps must be at least %d for the %s builtin", minStep, b.Name())
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

	if used > size {
		return 0, 0, errors.Errorf("The builtin %s used %d cells but the capacity is %d", b.Name(), used, size)
	}

	return used, size, nil
}
