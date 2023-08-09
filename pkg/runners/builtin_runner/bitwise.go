package builtinrunner

import "github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"

type BitwiseBuiltinRunner struct {
	base     memory.Relocatable
	included bool
}

func NewBitwiseBuiltinRunner(included bool) BitwiseBuiltinRunner {
	return BitwiseBuiltinRunner{
		base:     memory.NewRelocatable(0.0),
		included: included,
	}
}

func (r *BitwiseBuiltinRunner) Base() memory.Relocatable {
	return r.base
}

func (r *BitwiseBuiltinRunner) Name() string {
	return "bitwise_builtin"
}

// func (r *BitwiseBuiltinRunner) InitializeSegments(segments *memory.MemorySegmentManager)
