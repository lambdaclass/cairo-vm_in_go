package builtinrunner

import (
	"errors"
	"math"

	"github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
	"github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
)

const INNER_RC_BOUND_SHIFT = 16
const INNER_RC_BOUND_MASK = math.MaxUint16
const CELLS_PER_RANGE_CHECK = 1

const N_PARTS = 8

type RangeCheckBuiltinRunner struct {
	base     int
	included bool
}

func NewRangeCheckBuiltinRunner(ratio *uint32, nParts uint32, included bool) *RangeCheckBuiltinRunner {
	f_one := lambdaworks.FeltOne()
	bound := f_one.Shl(16 * nParts)
	if nParts != 0 && bound.IsZero() {
		return &RangeCheckBuiltinRunner{
			base:     0,
			included: included,
		}
	}
	return &RangeCheckBuiltinRunner{
		base:     0,
		included: included,
	}
}

func (r *RangeCheckBuiltinRunner) Base() int {
	return r.base
}

func (r *RangeCheckBuiltinRunner) Name() string {
	return "RangeCheck"
}

func (r *RangeCheckBuiltinRunner) InitializeSegments(segments *memory.MemorySegmentManager) {
	r.base = segments.AddSegment().SegmentIndex
}

func (r *RangeCheckBuiltinRunner) InitialStack() []memory.MaybeRelocatable {
	if r.included {
		stack := []memory.MaybeRelocatable{*memory.NewMaybeRelocatableRelocatable(memory.NewRelocatable(r.base, 0))}
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
		return nil, errors.New("RangeCheckFoundNonInt")
	}
	felt, is_felt := res_val.GetFelt()
	if !is_felt {
		return nil, errors.New("NotFeltElement")
	}
	if felt.Bits() <= N_PARTS*INNER_RC_BOUND_SHIFT {
		return []memory.Relocatable{address}, nil
	}
	return nil, errors.New("RangeCheckNumOutOfBounds")
}

func (r *RangeCheckBuiltinRunner) AddValidationRule(mem *memory.Memory) {
	mem.AddValidationRule(uint(r.base), ValidationRule)
}
