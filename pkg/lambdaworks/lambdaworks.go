package lambdaworks

/*
#cgo LDFLAGS: pkg/lambdaworks/lib/liblambdaworks.a -ldl
#include "lib/lambdaworks.h"
*/
import "C"
import "fmt"

// Go representation of a single limb (unsigned integer with 64 bits).
type Limb C.limb_t

// Go representation of a 256 bit prime field element (felt).
type Felt struct {
	limbs [4]Limb
}

func Number() int {
	return int(C.number())
}

// Gets a Felt representing the "value" number, in Montgomery format.
func From(value uint64) Felt {
	var result C.felt_t
	C.from(&result[0], C.uint64_t(value))

	// Convert the result to uint64 (C.uint64_t is mapped to _Ctype_ulonglong).
	goValue := &result[0]

	// Print the result.
	fmt.Println("Result:", uint64(*goValue))
	return Felt{}
}

// // Gets a Felt representing 0.
// func Zero() Felt {
// 	var result Felt
// 	C.zero(&result.limbs[0])
// 	return result
// }

// // Gets a Felt representing 1.
// func One() Felt {
// 	var result Felt
// 	C.one(&result.limbs[0])
// 	return result
// }

// // Writes the result variable with the sum of a and b felts.
// func Add(a, b Felt) Felt {
// 	var result Felt
// 	C.add((*C.felt_t)(&a.limbs[0]), (*C.felt_t)(&b.limbs[0]), (*C.felt_t)(&result.limbs[0]))
// 	return result
// }

// // Writes the result variable with a - b.
// func Sub(a, b Felt) Felt {
// 	var result Felt
// 	C.sub((*C.felt_t)(&a.limbs[0]), (*C.felt_t)(&b.limbs[0]), (*C.felt_t)(&result.limbs[0]))
// 	return result
// }

// // Writes the result variable with a * b.
// func Mul(a, b Felt) Felt {
// 	var result Felt
// 	C.mul((*C.felt_t)(&a.limbs[0]), (*C.felt_t)(&b.limbs[0]), (*C.felt_t)(&result.limbs[0]))
// 	return result
// }

// // Writes the result variable with a / b.
// func Div(a, b Felt) Felt {
// 	var result Felt
// 	C.div((*C.felt_t)(&a.limbs[0]), (*C.felt_t)(&b.limbs[0]), (*C.felt_t)(&result.limbs[0]))
// 	return result
// }
