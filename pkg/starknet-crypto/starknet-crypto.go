package starknet_crypto

/*
#cgo LDFLAGS: pkg/starknet-crypto/lib/libstarknet_crypto.a -ldl
#include "lib/starknet-crypto.h"
#include <stdlib.h>
*/
import "C"
import (
	"errors"
	"fmt"
	"unsafe"

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

func PoseidonPermuteComp(poseidon_state []lambdaworks.Felt) error {
	// Check input args
	if len(poseidon_state) != 3 {
		return errors.New("Poseidon state must have 3 elements")
	}
	// Convert args to c representation
	poseidon_state_c := make([]C.felt_t, 0, 3)
	for i := uint(0); i < 3; i++ {
		poseidon_state_c[i] = toC(poseidon_state[i])
	}
	// Compute hash using starknet-crypto C wrapper
	poseidon_state_ptr := (*C.felt_t)(unsafe.Pointer(&poseidon_state_c))
	C.poseidon_permute(poseidon_state_ptr)
	// convert result to Go representation
	new_poseidon_state := make([]lambdaworks.Felt, 0, 3)
	for i, elem := range poseidon_state_c {
		felt_elem := fromC(elem)
		new_poseidon_state[i] = felt_elem
	}

	fmt.Printf("New state %+v", new_poseidon_state)

	return nil
}

// // Convert args to c representation
// var poseidon_state_c *C.felt_t
// start := unsafe.Pointer(poseidon_state_c)
// for i := uint(0); i < 3; i++ {
// 	c_felt := toC(poseidon_state[i])
// 	*(*C.felt_t)(unsafe.Pointer(uintptr(start) + unsafe.Sizeof(c_felt)*uintptr(i))) = c_felt
// }
// // Compute hash using starknet-crypto C wrapper
// C.poseidon_permute(poseidon_state_c)
// // convert result to Go representation
// new_poseidon_state := make([]lambdaworks.Felt, 0, 3)
// for i := uint(0); i < 3; i++ {

// }
