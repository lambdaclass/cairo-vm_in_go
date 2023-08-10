package builtins

import (
	"github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
	starknet_crypto "github.com/lambdaclass/cairo-vm.go/pkg/starknet-crypto"
	"github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
)

const POSEIDON_BUILTIN_NAME = "poseidon"
const CELLS_PER_INSTANCE = 6
const INPUT_CELLS_PER_INSTANCE = 3

type PoseidonBuiltinRunner struct {
	base     memory.Relocatable
	included bool
	cache    map[memory.Relocatable]lambdaworks.Felt
}

func NewPoseidonBuiltinRunner(included bool) *PoseidonBuiltinRunner {
	return &PoseidonBuiltinRunner{included: included}
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
	index := address.Offset % CELLS_PER_INSTANCE
	if index < INPUT_CELLS_PER_INSTANCE {
		return nil, nil
	}

	value, ok := p.cache[address]
	if ok {
		return memory.NewMaybeRelocatableFelt(value), nil
	}

	input_start_addr, _ := address.SubUint(index)
	output_start_address := address.AddUint(INPUT_CELLS_PER_INSTANCE)

	// Build the initial poseidon state
	var poseidon_state [3]lambdaworks.Felt

	for i := uint(0); i < INPUT_CELLS_PER_INSTANCE; i++ {
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
