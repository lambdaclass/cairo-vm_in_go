package builtinrunner

import "github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
import "errors"

type BitwiseInstanceDef struct {
	Ratio      *uint
	TotalNBits uint
}

type BitwiseBuiltinRunner struct {
	base             memory.Relocatable
	included         bool
	CellsPerInstance uint
	Ratio            *uint
	BitwiseBuiltin   BitwiseInstanceDef
}

func NewBitwiseBuiltinRunner(instance_def BitwiseInstanceDef, included bool) BitwiseBuiltinRunner {
	return BitwiseBuiltinRunner{
		base:             memory.NewRelocatable(0, 0),
		included:         included,
		CellsPerInstance: 5,
		Ratio:            instance_def.Ratio,
		BitwiseBuiltin:   instance_def,
	}
}

func (b *BitwiseBuiltinRunner) Base() memory.Relocatable {
	return b.base
}

func (b *BitwiseBuiltinRunner) Name() string {
	return "bitwise_builtin"
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
	index := address.Offset % b.CellsPerInstance
	if index < 1 {
		return nil, nil
	}

	x_addr := memory.NewRelocatable(address.SegmentIndex, address.Offset-index)
	y_addr, err := (x_addr.AddUint(1))
	if err != nil {
		return nil, err
	}

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
		// TODO: check if this conversion if valid
		if uint64(num_x_felt.Bits()) > uint64(b.BitwiseBuiltin.TotalNBits) {
			return nil, errors.New("Expected Intenger x to be smaller than 2^(total_n_bits)")
		}
		if uint64(num_y_felt.Bits()) > uint64(b.BitwiseBuiltin.TotalNBits) {
			return nil, errors.New("Expected Intenger y to be smaller than 2^(total_n_bits)")
		}

		var res *memory.MaybeRelocatable
		switch index {
		case 2:
			res = memory.NewMaybeRelocatableFelt(num_x_felt.And(num_y_felt))
		case 3:
			res = memory.NewMaybeRelocatableFelt(num_x_felt.Pow(num_y_felt))
		case 4:
			res = memory.NewMaybeRelocatableFelt(num_x_felt.Or(num_y_felt))
		default:
			res = nil
		}
		return res, nil
	}

	return nil, nil
}

func (r *BitwiseBuiltinRunner) AddValidatonRule(segments *memory.Memory) {}
