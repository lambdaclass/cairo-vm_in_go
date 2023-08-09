package lambdaworks

/*
#cgo LDFLAGS: pkg/lambdaworks/lib/liblambdaworks.a -ldl
#include "lib/lambdaworks.h"
#include <stdlib.h>
*/
import "C"

import (
	"errors"
	"unsafe"
)

const N_LIMBS_IN_FELT = 4

// Go representation of a single limb (unsigned integer with 64 bits).
type Limb C.limb_t

// Go representation of a 256 bit prime field element (felt).
type Felt struct {
	limbs [N_LIMBS_IN_FELT]Limb
}

// Converts a Go Felt to a C felt_t.
func (f Felt) toC() C.felt_t {
	var result C.felt_t
	for i, limb := range f.limbs {
		result[i] = C.limb_t(limb)
	}
	return result
}

// Converts a C felt_t to a Go Felt.
func fromC(result C.felt_t) Felt {
	var limbs [N_LIMBS_IN_FELT]Limb
	for i, limb := range result {
		limbs[i] = Limb(limb)
	}
	return Felt{limbs: limbs}
}

// Gets a Felt representing the "value" number, in Montgomery format.
func FeltFromUint64(value uint64) Felt {
	var result C.felt_t
	C.from(&result[0], C.uint64_t(value))
	return fromC(result)
}

func FeltFromHex(value string) Felt {
	cs := C.CString(value)
	defer C.free(unsafe.Pointer(cs))

	var result C.felt_t
	C.from_hex(&result[0], cs)
	return fromC(result)
}

func FeltFromDecString(value string) Felt {
	cs := C.CString(value)
	defer C.free(unsafe.Pointer(cs))

	var result C.felt_t
	C.from_dec_str(&result[0], cs)
	return fromC(result)
}

// turns a felt to usize
func (felt Felt) ToU64() (uint64, error) {
	if felt.limbs[0] == 0 && felt.limbs[1] == 0 && felt.limbs[2] == 0 {
		return uint64(felt.limbs[3]), nil
	} else {
		return 0, errors.New("Cannot convert felt to u64")
	}
}

// Gets a Felt representing 1.
func FeltOne() Felt {
	var result C.felt_t
	C.one(&result[0])
	return fromC(result)
}

// Gets a Felt representing 0.
func FeltZero() Felt {
	var result C.felt_t
	C.zero(&result[0])
	return fromC(result)
}

func (f Felt) IsZero() bool {
	return f == FeltZero()
}

// Writes the result variable with the sum of a and b felts.
func (a Felt) Add(b Felt) Felt {
	var result C.felt_t
	var a_c C.felt_t = a.toC()
	var b_c C.felt_t = b.toC()
	C.add(&a_c[0], &b_c[0], &result[0])
	return fromC(result)
}

// Writes the result variable with a - b.
func (a Felt) Sub(b Felt) Felt {
	var result C.felt_t
	var a_c C.felt_t = a.toC()
	var b_c C.felt_t = b.toC()
	C.sub(&a_c[0], &b_c[0], &result[0])
	return fromC(result)
}

// Writes the result variable with a * b.
func (a Felt) Mul(b Felt) Felt {
	var result C.felt_t
	var a_c C.felt_t = a.toC()
	var b_c C.felt_t = b.toC()
	C.mul(&a_c[0], &b_c[0], &result[0])
	return fromC(result)
}

// Writes the result variable with a / b.
func (a Felt) Div(b Felt) Felt {
	var result C.felt_t
	var a_c C.felt_t = a.toC()
	var b_c C.felt_t = b.toC()
	C.lw_div(&a_c[0], &b_c[0], &result[0])
	return fromC(result)
}

// Returns the number of bits needed to represent the felt
func (a Felt) Bits() C.limb_t {
	if a.IsZero() {
		return 0
	}
	var a_c = a.toC()
	return C.bits(&a_c[0])

}
