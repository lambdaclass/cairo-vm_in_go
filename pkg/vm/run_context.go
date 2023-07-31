package vm

import "github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"

// RunContext containts the register states of the
// Cairo VM.
type RunContext struct {
	Pc memory.Relocatable
	Ap memory.Relocatable
	Fp memory.Relocatable
}
