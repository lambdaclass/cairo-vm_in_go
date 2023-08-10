package starknet_crypto

/*
#cgo LDFLAGS: pkg/starknet-crypto/lib/libstarknet_crypto.a -ldl
#include "lib/starknet-crypto.h"
#include <stdlib.h>
*/
import "C"
import (
	"github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
)

// Converts a Go Felt to a C felt_t.
func toC(f lambdaworks.Felt) C.felt_t {
	var result C.felt_t
	for i, byte := range f.ToBeBytes() {
		result[i] = C.byte_t(byte)
	}
	return result
}

// Converts a C felt_t to a Go Felt.
func fromC(result C.felt_t) lambdaworks.Felt {
	var bytes [32]uint8
	for i, byte := range result {
		bytes[i] = uint8(byte)
	}
	return lambdaworks.FeltFromBeBytes(&bytes)
}

func PoseidonPermuteComp(poseidon_state *[3]lambdaworks.Felt) {
	state := *poseidon_state
	// Convert args to c representation
	first_state_felt := toC(state[0])
	second_state_felt := toC(state[1])
	third_state_felt := toC(state[2])

	// Compute hash using starknet-crypto C wrapper
	C.poseidon_permute(&first_state_felt[0], &second_state_felt[0], &third_state_felt[0])
	// Convert result to Go representation
	var new_poseidon_state = [3]lambdaworks.Felt{
		fromC(first_state_felt),
		fromC(second_state_felt),
		fromC(third_state_felt),
	}
	// Update poseidon state
	*poseidon_state = new_poseidon_state
}
