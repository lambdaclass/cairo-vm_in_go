package builtins

import (
	"errors"

	"github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
)

const BITWISE_CELLS_PER_INSTANCE = 5
const BITWISE_TOTAL_N_BITS = 251
const BIWISE_INPUT_CELLS_PER_INSTANCE = 2

type BitwiseBuiltinRunner struct {
	base     memory.Relocatable
	included bool
}

func NewBitwiseBuiltinRunner(included bool) *BitwiseBuiltinRunner {
	return &BitwiseBuiltinRunner{
		included: included,
	}
}

func (b *BitwiseBuiltinRunner) Base() memory.Relocatable {
	return b.base
}

func (b *BitwiseBuiltinRunner) Name() string {
	return "bitwise"
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

func (b *BitwiseBuiltinRunner) DeduceMemoryCell(address memory.Relocatable, segments *memory.Memory) (*memory.MaybeRelocatable, error) {
	index := address.Offset % BITWISE_CELLS_PER_INSTANCE
	if index < BIWISE_INPUT_CELLS_PER_INSTANCE {
		return nil, nil
	}

	x_addr := memory.NewRelocatable(address.SegmentIndex, address.Offset-index)
	y_addr := x_addr.AddUint(1)
	num_x, err := segments.Get(x_addr)
	if err != nil {

		return nil, err
	}

	num_y, err := segments.Get(y_addr)
	if err != nil {

		return nil, err
	}

	num_x_felt, x_is_felt := num_x.GetFelt()
	num_y_felt, y_is_felt := num_y.GetFelt()

	if x_is_felt && y_is_felt {
		if num_x_felt.Bits() > BITWISE_TOTAL_N_BITS {
			return nil, errors.New("Expected Intenger x to be smaller than 2^(total_n_bits)")
		}
		if num_y_felt.Bits() > BITWISE_TOTAL_N_BITS {
			return nil, errors.New("Expected Intenger y to be smaller than 2^(total_n_bits)")
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

	return nil, nil
}

func (b *BitwiseBuiltinRunner) AddValidationRule(*memory.Memory) {}
