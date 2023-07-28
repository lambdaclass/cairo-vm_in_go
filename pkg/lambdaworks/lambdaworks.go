// Wrapper around the lambdaworks library. A mock, for now.
package lambdaworks

/*
#cgo LDFLAGS: pkg/lambdaworks/lib/liblambdaworks.a -ldl
#include "lib/lambdaworks.h"
*/
import "C"

func Number() int {
	return int(C.number())
}
