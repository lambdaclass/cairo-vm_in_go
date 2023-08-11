package builtins

import (
	"fmt"

	. "github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
	. "github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
)

const KECCAK_CELLS_PER_INSTANCE = 16
const KECCAK_INPUT_CELLS_PER_INSTANCE = 8
const KECCAK_INPUT_BIT_LENTGH = 200
const KECCAK_INPUT_BYTES_LENTGH = 25

type KeccakBuiltinRunner struct {
	base     Relocatable
	included bool
	cache    map[Relocatable]Felt
}

func NewKeccakBuiltinRunner(included bool) *KeccakBuiltinRunner {
	return &KeccakBuiltinRunner{included: included, cache: make(map[Relocatable]Felt)}
}

const KECCAK_BUILTIN_NAME = "keccak"

func (k *KeccakBuiltinRunner) Base() Relocatable {
	return k.base
}

func (k *KeccakBuiltinRunner) Name() string {
	return KECCAK_BUILTIN_NAME
}

func (k *KeccakBuiltinRunner) InitializeSegments(segments *MemorySegmentManager) {
	k.base = segments.AddSegment()
}

func (k *KeccakBuiltinRunner) InitialStack() []MaybeRelocatable {
	if k.included {
		return []MaybeRelocatable{*NewMaybeRelocatableRelocatable(k.base)}
	} else {
		return nil
	}
}

func (k *KeccakBuiltinRunner) AddValidationRule(*Memory) {}

func (k *KeccakBuiltinRunner) DeduceMemoryCell(address Relocatable, mem *Memory) (*MaybeRelocatable, error) {
	index := address.Offset % KECCAK_CELLS_PER_INSTANCE
	if index < KECCAK_INPUT_CELLS_PER_INSTANCE {
		return nil, nil
	}

	value, ok := k.cache[address]
	if ok {
		return NewMaybeRelocatableFelt(value), nil
	}

	input_start_addr, _ := address.SubUint(index)
	output_start_address := input_start_addr.AddUint(KECCAK_INPUT_CELLS_PER_INSTANCE)

	input_message := make([]byte, 0, 25*KECCAK_INPUT_CELLS_PER_INSTANCE)

	for i := uint(0); i < KECCAK_INPUT_CELLS_PER_INSTANCE; i++ {
		felt, err := mem.GetFelt(input_start_addr.AddUint(i))
		if err != nil {
			return nil, err
		}
		// TODO: Check bit length
		le_bytes := felt.ToLeBytes()
		input_message = append(input_message, le_bytes[:25]...)
	}
	// Run keccak
	// Cairo VM here uses the internal k1600 permutation, not the sha3 hash
	// SOLUTION: Couldnt find any cool crates that export keccak_f, so my solution atm will be to link
	// the keccak_f we use in cairo-vm, we would send the 8 felts and receive 8 felts for convenience
	// Checked the keccakf1600 crate, but it uses a state of 50 uint64
	output_message := make([]byte, 0, 200)
	fmt.Printf("INPUT BYTES %v, %v\n", []byte(input_message), len(input_message))
	fmt.Printf("OUTPUT BYTES %v, %v\n", []byte(output_message), len(output_message))
	for i := uint(0); i < KECCAK_INPUT_CELLS_PER_INSTANCE; i++ {
		bytes := (output_message)[25*i : 25*i+25]
		padded_bytes := (*[32]byte)(append(bytes, []byte{0, 0, 0, 0, 0, 0, 0}...))
		felt := FeltFromLeBytes(padded_bytes)
		k.cache[output_start_address.AddUint(i)] = felt
	}
	return NewMaybeRelocatableFelt(k.cache[address]), nil

}
