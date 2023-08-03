package vm

import "github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"

type TraceEntry struct {
	Pc memory.Relocatable
	Ap memory.Relocatable
	Fp memory.Relocatable
}

type RelocatedTraceEntry struct {
	Pc uint
	Ap uint
	Fp uint
}
