package main

import (
	"math"

	"github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
	"github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
)

const INNER_RC_BOUND_SHIFT = 16
const INNER_RC_BOUND_MASK = math.MaxUint16
const CELLS_PER_RANGE_CHECK = 1

const N_PARTS = 8

type ValidationRule struct {
	// Define the ValidationRule struct as needed.
}

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

func (r *RangeCheckBuiltinRunner) initializeSegments(segments *memory.MemorySegmentManager) {
	r.base = segments.AddSegment().SegmentIndex
}

func (r *RangeCheckBuiltinRunner) initialStack() []memory.MaybeRelocatable {
	if r.included {
		stack := make([]memory.MaybeRelocatable, 1)
		stack[0] = *memory.NewMaybeRelocatableRelocatable(memory.NewRelocatable(r.base, 0))
		return stack
	}
	return []memory.MaybeRelocatable{}
}

func (r *RangeCheckBuiltinRunner) getMemorySegmentAddresses() (int, *int) {
	return r.base, r.stopPtr
}

func (r *RangeCheckBuiltinRunner) getUsedCells(segments *memory.MemorySegmentManager) (uint, error) {
	usedSize, ok := segments.SegmentSizes[uint(r.base)]
	if !ok {
		return 0, memory.MissingSegmentUsize
	}
	return usedSize, nil
}

func (r *RangeCheckBuiltinRunner) getRangeCheckUsage(checkedMemory *memory.Memory) (*struct{uint64, uint64}, error) {
	rangeCheckSegment, err := checkedMemory.Get(memory.NewRelocatable(r.base, 0))
	if err != nil {
		return math.MaxUint64, 0, err
	}
	rcBounds := &struct{ min, max uint64 }{math.MaxInt, math.MinInt}

	// Split value into nParts parts of less than INNER_RC_BOUND size.
	for _, value := range rangeCheckSegment {
		if value == nil || value.value == nil {
			continue
		}
		num := value.value.getInteger()
		if num == nil {
			continue
		}
		if num.bits() <= N_PARTS*INNER_RC_BOUND_SHIFT {
			for i := uint(0); i < N_PARTS; i++ {
				x := (num.Uint64() >> (i * INNER_RC_BOUND_SHIFT)) & INNER_RC_BOUND_MASK
				rcBounds.min = int(math.Min(float64(rcBounds.min), float64(x)))
				rcBounds.max = int(math.Max(float64(rcBounds.max), float64(x)))
			}
		}
	}

	return rcBounds
}

// Implement other methods for RangeCheckBuiltinRunner similar to Rust code.

func main() {
	// Example usage of RangeCheckBuiltinRunner in Go.
	ratio := uint32(10)
	nParts := uint32(8)
	included := true
	builtin := NewRangeCheckBuiltinRunner(&ratio, nParts, included)
	// Initialize segments using builtin.initializeSegments(segments *MemorySegmentManager)
	// Use other methods as needed.
}
