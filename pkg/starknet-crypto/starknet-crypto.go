package starknet_crypto

/*
#cgo LDFLAGS: pkg/starknet-crypto/lib/libstarknet_crypto.a -ldl
#include "lib/starknet-crypto.h"
#include <stdlib.h>
*/
import "C"
import (
	"errors"

	"github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
)

// Converts a Go Felt to a C felt_t.
func toC(f lambdaworks.Felt) C.felt_t {
	var result C.felt_t
	for i, limb := range f.ToLimbs() {
		result[i] = C.limb_t(limb)
	}
	return result
}

// Converts a C felt_t to a Go Felt.
func fromC(result C.felt_t) lambdaworks.Felt {
	var limbs [4]uint64
	for i, limb := range result {
		limbs[i] = uint64(limb)
	}
	return lambdaworks.FeltFromLimbs(limbs)
}

func PoseidonPermuteComp(poseidon_state *[]lambdaworks.Felt) error {
	// Check input args
	if len(*poseidon_state) != 3 {
		return errors.New("Poseidon state must have 3 elements")
	}
	state := *poseidon_state
	// Convert args to c representation
	first_state_felt := toC(state[0])
	second_state_felt := toC(state[1])
	third_state_felt := toC(state[2])

	// Compute hash using starknet-crypto C wrapper
	C.poseidon_permute(&first_state_felt[0], &second_state_felt[0], &third_state_felt[0])
	// Convert result to Go representation
	new_poseidon_state := make([]lambdaworks.Felt, 3)
	new_poseidon_state[0] = fromC(first_state_felt)
	new_poseidon_state[1] = fromC(second_state_felt)
	new_poseidon_state[2] = fromC(third_state_felt)
	// Update poseidon state
	*poseidon_state = new_poseidon_state

	return nil
}
