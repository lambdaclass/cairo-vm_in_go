package lambdaworks

/*
#cgo LDFLAGS: pkg/lambdaworks/lib/liblambdaworks.a -ldl
#include "lib/lambdaworks.h"
#include <stdlib.h>
*/
import "C"

import (
	"math/big"
	"strings"
	"unsafe"

	"github.com/pkg/errors"
)

const N_LIMBS_IN_FELT = 4

// Go representation of a single limb (unsigned integer with 64 bits).
type Limb C.limb_t

// Go representation of a 256 bit prime field element (felt).
type Felt struct {
	limbs [N_LIMBS_IN_FELT]Limb
}

func LambdaworksError(err error) error {
	return errors.Wrapf(err, "Lambdaworks Error")
}

func ConversionError(felt Felt, targetType string) error {
	return LambdaworksError(errors.Errorf("Cannot convert felt: %d to %s", felt, targetType))
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
		return 0, ConversionError(felt, "u64")
	}
}

func (felt Felt) ToLeBytes() *[32]byte {
	var result_c [32]C.uint8_t
	var value C.felt_t = felt.toC()
	C.to_le_bytes(&result_c[0], &value[0])

	result := (*[32]byte)(unsafe.Pointer(&result_c))

	return result
}

func (felt Felt) ToBeBytes() *[32]byte {
	var result_c [32]C.uint8_t
	var value C.felt_t = felt.toC()
	C.to_be_bytes(&result_c[0], &value[0])

	result := (*[32]byte)(unsafe.Pointer(&result_c))

	return result
}

func (felt Felt) ToHexString() string {
	// We need to make sure enough space is allocated to fit the longest possible string
	var result_c = C.CString(strings.Repeat(" ", 65))
	defer C.free(unsafe.Pointer(result_c))

	var value C.felt_t = felt.toC()
	C.to_hex_string(result_c, &value[0])
	res := C.GoString(result_c)
	return strings.TrimSpace(res)
}

func FeltFromLeBytes(bytes *[32]byte) Felt {
	var result C.felt_t
	bytes_ptr := (*[32]C.uint8_t)(unsafe.Pointer(bytes))
	C.from_le_bytes(&result[0], &bytes_ptr[0])
	return fromC(result)
}

func FeltFromBeBytes(bytes *[32]byte) Felt {
	var result C.felt_t
	bytes_ptr := (*[32]C.uint8_t)(unsafe.Pointer(bytes))
	C.from_be_bytes(&result[0], &bytes_ptr[0])
	return fromC(result)
}

// Gets a Felt representing 0.
func FeltZero() Felt {
	var result C.felt_t
	C.zero(&result[0])
	return fromC(result)
}

// Gets a Felt representing 1.
func FeltOne() Felt {
	var result C.felt_t
	C.one(&result[0])
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

// Returns the felt
func (f Felt) ToSignedFeltString() string {
	var f_c = f.toC()
	resultPtr := C.to_signed_felt(&f_c[0])
	defer C.free_string(resultPtr)
	goResult := C.GoString(resultPtr)
	return goResult
}

// Returns the number of bits needed to represent the felt
func (a Felt) Bits() C.limb_t {
	if a.IsZero() {
		return 0
	}
	var a_c = a.toC()
	return C.bits(&a_c[0])
}

func (a Felt) And(b Felt) Felt {
	var result C.felt_t
	var a_c C.felt_t = a.toC()
	var b_c C.felt_t = b.toC()
	C.felt_and(&a_c[0], &b_c[0], &result[0])
	return fromC(result)
}

func (a Felt) Xor(b Felt) Felt {
	var result C.felt_t
	var a_c C.felt_t = a.toC()
	var b_c C.felt_t = b.toC()
	C.felt_xor(&a_c[0], &b_c[0], &result[0])
	return fromC(result)
}

func (a Felt) Or(b Felt) Felt {
	var result C.felt_t
	var a_c C.felt_t = a.toC()
	var b_c C.felt_t = b.toC()
	C.felt_or(&a_c[0], &b_c[0], &result[0])
	return fromC(result)
}

func (a Felt) Shr(b uint) Felt {
	var result C.felt_t
	var a_c C.felt_t = a.toC()
	//var b_c C._type_uint = b.toC()
	C.felt_shr(&a_c[0], C.size_t(b), &result[0])
	//C.felt_shr(&a_c[0], &b_c[0], &result[0])
	return fromC(result)
}

func (f Felt) ToBigInt() *big.Int {
	return new(big.Int).SetBytes(f.ToBeBytes()[:32])
}

const CAIRO_PRIME_HEX = "0x800000000000011000000000000000000000000000000000000000000000001"
const SIGNED_FELT_MAX_HEX = "0x400000000000008800000000000000000000000000000000000000000000000"

// Implements `as_int` behaviour
func (f Felt) ToSigned() *big.Int {
	n := f.ToBigInt()
	signedFeltMax, _ := new(big.Int).SetString(SIGNED_FELT_MAX_HEX, 0)
	if n.Cmp(signedFeltMax) == 1 {
		cairoPrime, _ := new(big.Int).SetString(CAIRO_PRIME_HEX, 0)
		return new(big.Int).Neg(new(big.Int).Sub(cairoPrime, n))
	}
	return n
}

/*
Compares x and y and returns:

	-1 if a <  b
	 0 if a == b
	+1 if a >  b
*/
func (a Felt) Cmp(b Felt) int {
	var a_c C.felt_t = a.toC()
	var b_c C.felt_t = b.toC()
	return int(C.cmp(&a_c[0], &b_c[0]))
}
