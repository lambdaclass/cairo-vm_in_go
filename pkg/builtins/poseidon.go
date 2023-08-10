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

func (p *PoseidonBuiltinRunner) DeduceMemoryCell(address memory.Relocatable, memory *memory.Memory) (*memory.MaybeRelocatable, error) {
	// Check if its an input cell
	index := address.Offset % CELLS_PER_INSTANCE
	if index < INPUT_CELLS_PER_INSTANCE {
		return nil, nil
	}

	// TODO: fetch from cache
	input_start_addr, _ := address.SubUint(index)
	first_output_address := address.AddUint(INPUT_CELLS_PER_INSTANCE)

	// Build the initial poseidon state
	var poseidon_state [3]lambdaworks.Felt

	for i := uint(0); i < INPUT_CELLS_PER_INSTANCE; i++ {
		felt, err := memory.GetFelt(first_output_address.AddUint(i))
		if err != nil {
			return nil, err
		}
		poseidon_state[i] = felt
	}

	// Run the poseidon permutation

	starknet_crypto.PoseidonPermuteComp(&poseidon_state)

	// Insert the new state into the output cells
	// TODO insert into cache
	//TODO fetch result from output cache

	return nil, nil //TODO
}

func (p *PoseidonBuiltinRunner) AddValidationRule(*memory.Memory) {
}
