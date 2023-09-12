package builtins

import (
	"github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
	starknet_crypto "github.com/lambdaclass/cairo-vm.go/pkg/starknet_crypto"
	"github.com/lambdaclass/cairo-vm.go/pkg/utils"
	"github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
	"github.com/pkg/errors"
)

const POSEIDON_BUILTIN_NAME = "poseidon"
const POSEIDON_CELLS_PER_INSTANCE = 6
const POSEIDON_INPUT_CELLS_PER_INSTANCE = 3

type PoseidonBuiltinRunner struct {
	base                  memory.Relocatable
	included              bool
	cache                 map[memory.Relocatable]lambdaworks.Felt
	ratio                 uint
	instancesPerComponent uint
}

func NewPoseidonBuiltinRunner() *PoseidonBuiltinRunner {
	return &PoseidonBuiltinRunner{cache: make(map[memory.Relocatable]lambdaworks.Felt)}
}

func (p *PoseidonBuiltinRunner) Base() memory.Relocatable {
	return p.base
}

func (p *PoseidonBuiltinRunner) Name() string {
	return POSEIDON_BUILTIN_NAME
}

func (p *PoseidonBuiltinRunner) InitializeSegments(segments *memory.MemorySegmentManager) {
	p.base = segments.AddSegment()
}

func (p *PoseidonBuiltinRunner) InitialStack() []memory.MaybeRelocatable {
	if p.included {
		return []memory.MaybeRelocatable{*memory.NewMaybeRelocatableRelocatable(p.base)}
	} else {
		return nil
	}
}

func (p *PoseidonBuiltinRunner) DeduceMemoryCell(address memory.Relocatable, mem *memory.Memory) (*memory.MaybeRelocatable, error) {
	// Check if its an input cell
	index := address.Offset % POSEIDON_CELLS_PER_INSTANCE
	if index < POSEIDON_INPUT_CELLS_PER_INSTANCE {
		return nil, nil
	}

	value, ok := p.cache[address]
	if ok {
		return memory.NewMaybeRelocatableFelt(value), nil
	}
	// index will always be less or equal to address.Offset so we can ignore the error
	input_start_addr, _ := address.SubUint(index)
	output_start_address := input_start_addr.AddUint(POSEIDON_INPUT_CELLS_PER_INSTANCE)

	// Build the initial poseidon state
	var poseidon_state [POSEIDON_INPUT_CELLS_PER_INSTANCE]lambdaworks.Felt

	for i := uint(0); i < POSEIDON_INPUT_CELLS_PER_INSTANCE; i++ {
		felt, err := mem.GetFelt(input_start_addr.AddUint(i))
		if err != nil {
			return nil, err
		}
		poseidon_state[i] = felt
	}

	// Run the poseidon permutation
	starknet_crypto.PoseidonPermuteComp(&poseidon_state)

	// Insert the new state into the corresponding output cells in the cache
	for i, elem := range poseidon_state {
		p.cache[output_start_address.AddUint(uint(i))] = elem
	}
	return memory.NewMaybeRelocatableFelt(p.cache[address]), nil
}

func (p *PoseidonBuiltinRunner) AddValidationRule(*memory.Memory) {
}

func (p *PoseidonBuiltinRunner) Include(include bool) {
	p.included = include
}

func (p *PoseidonBuiltinRunner) Ratio() uint {
	return p.ratio
}

func (p *PoseidonBuiltinRunner) CellsPerInstance() uint {
	return POSEIDON_CELLS_PER_INSTANCE
}

func (p *PoseidonBuiltinRunner) GetAllocatedMemoryUnits(segments *memory.MemorySegmentManager, currentStep uint) (uint, error) {
	// This condition corresponds to an uninitialized ratio for the builtin, which should only
	// happen when layout is `dynamic`
	if p.Ratio() == 0 {
		// Dynamic layout has the exact number of instances it needs (up to a power of 2).
		used, err := segments.GetSegmentUsedSize(uint(p.base.SegmentIndex))
		if err != nil {
			return 0, err
		}
		instances := used / p.CellsPerInstance()
		components := utils.NextPowOf2(instances / p.instancesPerComponent)
		size := p.CellsPerInstance() * p.instancesPerComponent * components

		return size, nil
	}

	minStep := p.ratio * p.instancesPerComponent
	if currentStep < minStep {
		return 0, errors.Errorf("number of steps must be at least %d for the %s builtin", minStep, p.Name())
	}
	value, err := utils.SafeDiv(currentStep, p.ratio)

	if err != nil {
		return 0, errors.Errorf("error calculating builtin memory units: %s", err)
	}

	return p.CellsPerInstance() * value, nil
}

func (p *PoseidonBuiltinRunner) GetUsedCellsAndAllocatedSizes(segments *memory.MemorySegmentManager, currentStep uint) (uint, uint, error) {
	used, err := segments.GetSegmentUsedSize(uint(p.base.SegmentIndex))
	if err != nil {
		return 0, 0, err
	}

	size, err := p.GetAllocatedMemoryUnits(segments, currentStep)

	if used > size {
		return 0, 0, errors.Errorf("The builtin %s used %d cells but the capacity is %d", p.Name(), used, size)
	}

	return used, size, nil
}
