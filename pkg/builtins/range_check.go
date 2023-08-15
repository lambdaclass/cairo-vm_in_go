package builtins

import (
	"errors"
	"math"

	"github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
)

var ErrRangeOutOfBounds = errors.New("RangeCheckNumOutOfBounds")

const CHECK_RANGE_BUILTIN_NAME = "range_check"
const INNER_RC_BOUND_SHIFT = 16
const INNER_RC_BOUND_MASK = math.MaxUint16
const CELLS_PER_RANGE_CHECK = 1

const N_PARTS = 8

type RangeCheckBuiltinRunner struct {
	base     memory.Relocatable
	included bool
}

func NewRangeCheckBuiltinRunner(included bool) *RangeCheckBuiltinRunner {
	return &RangeCheckBuiltinRunner{
		included: included,
	}
}

func (r *RangeCheckBuiltinRunner) Base() memory.Relocatable {
	return r.base
}

func (r *RangeCheckBuiltinRunner) Name() string {
	return CHECK_RANGE_BUILTIN_NAME
}

func (r *RangeCheckBuiltinRunner) InitializeSegments(segments *memory.MemorySegmentManager) {
	r.base = segments.AddSegment()
}

func (r *RangeCheckBuiltinRunner) InitialStack() []memory.MaybeRelocatable {
	if r.included {
		stack := []memory.MaybeRelocatable{*memory.NewMaybeRelocatableRelocatable(r.base)}
		return stack
	}
	return []memory.MaybeRelocatable{}
}

func (r *RangeCheckBuiltinRunner) DeduceMemoryCell(addr memory.Relocatable, mem *memory.Memory) (*memory.MaybeRelocatable, error) {
	return nil, nil
}

func ValidationRule(mem *memory.Memory, address memory.Relocatable) ([]memory.Relocatable, error) {
	res_val, err := mem.Get(address)
	if err != nil {
		return nil, err
	}
	felt, is_felt := res_val.GetFelt()
	if !is_felt {
		return nil, errors.New("NotFeltElement")
	}
	if felt.Bits() <= N_PARTS*INNER_RC_BOUND_SHIFT {
		return []memory.Relocatable{address}, nil
	}
	return nil, ErrRangeOutOfBounds
}

func (r *RangeCheckBuiltinRunner) AddValidationRule(mem *memory.Memory) {
	mem.AddValidationRule(uint(r.base.SegmentIndex), ValidationRule)
}
