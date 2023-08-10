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
	fmt.Println("Input state %+v", poseidon_state)
	// Check input args
	if len(poseidon_state) != 3 {
		return errors.New("Poseidon state must have 3 elements")
	}
	// Convert args to c representation
	first_state_felt := toC(poseidon_state[0])
	second_state_felt := toC(poseidon_state[1])
	third_state_felt := toC(poseidon_state[2])

	// Compute hash using starknet-crypto C wrapper
	C.poseidon_permute(first_state_felt, second_state_felt, third_state_felt)
	fmt.Println("C run success")
	// convert result to Go representation
	new_poseidon_state := make([]lambdaworks.Felt, 0, 3)
	new_poseidon_state[0] = fromC(first_state_felt)
	new_poseidon_state[1] = fromC(second_state_felt)
	new_poseidon_state[2] = fromC(third_state_felt)

	fmt.Println("New state %+v", new_poseidon_state)

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

// // Convert args to c representation
// poseidon_state_c := make([]C.felt_t, 3, 3)

// for i := uint(0); i < 3; i++ {
// 	poseidon_state_c[i] = toC(poseidon_state[i])
// }
// fmt.Println("poseidon_state_c %+v", poseidon_state_c)
// // Compute hash using starknet-crypto C wrapper
// poseidon_state_ptr := (*C.felt_t)(unsafe.Pointer(&poseidon_state_c))
// fmt.Println("Unsafe ptr cast")
// C.poseidon_permute(poseidon_state_ptr)
// fmt.Println("C run success")
// // convert result to Go representation
