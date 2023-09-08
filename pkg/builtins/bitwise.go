package builtins

import (
	"github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
	"github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
	"github.com/pkg/errors"
)

const BITWISE_BUILTIN_NAME = "bitwise_builtin"
const BITWISE_CELLS_PER_INSTANCE = 5
const BITWISE_TOTAL_N_BITS = 251
const BIWISE_INPUT_CELLS_PER_INSTANCE = 2

type BitwiseBuiltinRunner struct {
	base     memory.Relocatable
	included bool
}

func BitwiseError(err error) error {
	return errors.Wrapf(err, "Bitwise builtin error\n")
}

func ErrFeltBiggerThanPowerOfTwo(felt lambdaworks.Felt) error {
	return BitwiseError(errors.Errorf("Expected felt %d to be smaller than  2**%d", felt, BITWISE_TOTAL_N_BITS))
}

func NewBitwiseBuiltinRunner() *BitwiseBuiltinRunner {
	return &BitwiseBuiltinRunner{}
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

func (r *BitwiseBuiltinRunner) Include(include bool) {
	r.included = include
}
