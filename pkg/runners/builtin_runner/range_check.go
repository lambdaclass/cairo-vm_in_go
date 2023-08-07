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
	ratio                 *uint32
	base                  int
	stopPtr               *int
	cellsPerInstance      uint32
	nInputCells           uint32
	_bound                *lambdaworks.Felt
	included              bool
	nParts                uint32
	instancesPerComponent uint32
}

func NewRangeCheckBuiltinRunner(ratio *uint32, nParts uint32, included bool) *RangeCheckBuiltinRunner {
	bound := lambdaworks.Felt{}.One().Shl(16 * nParts)
	if nParts != 0 && bound.IsZero() {
		return &RangeCheckBuiltinRunner{
			ratio:                 ratio,
			base:                  0,
			stopPtr:               nil,
			cellsPerInstance:      CELLS_PER_RANGE_CHECK,
			nInputCells:           CELLS_PER_RANGE_CHECK,
			_bound:                nil,
			included:              included,
			nParts:                nParts,
			instancesPerComponent: 1,
		}
	}
	return &RangeCheckBuiltinRunner{
		ratio:                 ratio,
		base:                  0,
		stopPtr:               nil,
		cellsPerInstance:      CELLS_PER_RANGE_CHECK,
		nInputCells:           CELLS_PER_RANGE_CHECK,
		_bound:                &bound,
		included:              included,
		nParts:                nParts,
		instancesPerComponent: 1,
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
