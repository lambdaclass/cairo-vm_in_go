package builtins

import (
	"math"

	"github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
	"github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
	"github.com/pkg/errors"
)

const CHECK_RANGE_BUILTIN_NAME = "range_check"
const INNER_RC_BOUND_SHIFT = 16
const INNER_RC_BOUND_MASK = math.MaxUint16
const CELLS_PER_RANGE_CHECK = 1

const N_PARTS = 8

func RangeCheckError(err error) error {
	return errors.Wrapf(err, "Range check error")
}

func OutsideBoundsError(felt lambdaworks.Felt) error {
	upperBound := uint64(math.Pow(2, float64(N_PARTS*INNER_RC_BOUND_SHIFT)))
	return RangeCheckError(errors.Errorf("Value %d is out of bounds [0, %d]", felt, upperBound))
}

func NotAFeltError(addr memory.Relocatable, val memory.MaybeRelocatable) error {
	rel, _ := val.GetRelocatable()
	return RangeCheckError(errors.Errorf("Value %d found in %d is not a felt element", rel, addr))
}

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
		return []memory.MaybeRelocatable{*memory.NewMaybeRelocatableRelocatable(r.base)}
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
		return nil, NotAFeltError(address, *res_val)
	}
	if felt.Bits() <= N_PARTS*INNER_RC_BOUND_SHIFT {
		return []memory.Relocatable{address}, nil
	}
	return nil, OutsideBoundsError(felt)
}

func (r *RangeCheckBuiltinRunner) AddValidationRule(mem *memory.Memory) {
	mem.AddValidationRule(uint(r.base.SegmentIndex), ValidationRule)
}
