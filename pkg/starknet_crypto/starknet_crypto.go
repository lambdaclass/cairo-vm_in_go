package starknet_crypto

/*
#cgo LDFLAGS: pkg/starknet_crypto/lib/libstarknet_crypto.a -ldl
#include "lib/starknet_crypto.h"
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

	// Compute hash using starknet_crypto C wrapper
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

func PedersenHash(f1 lambdaworks.Felt, f2 lambdaworks.Felt) lambdaworks.Felt {
	felt_1 := toC(f1)
	felt_2 := toC(f2)
	var result C.felt_t

	C.pedersen_hash(&felt_1[0], &felt_2[0], &result[0])

	hash := fromC(result)

	return hash
}

func VerifySignature(public_key lambdaworks.Felt, message lambdaworks.Felt, r lambdaworks.Felt, s lambdaworks.Felt) bool {
	public_key_for_c := toC(public_key)
	message_for_c := toC(message)
	r_for_c := toC(r)
	s_for_c := toC(s)

	c_verify_status := C.verify_signature(&public_key_for_c[0], &message_for_c[0], &r_for_c[0], &s_for_c[0])

	return bool(c_verify_status)
}
