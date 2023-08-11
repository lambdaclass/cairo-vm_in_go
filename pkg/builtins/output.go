package builtins

import "github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"

const OUTPUT_BUILTIN_NAME = "output"

type OutputBuiltinRunner struct {
	base     memory.Relocatable
	included bool
}

func (r *OutputBuiltinRunner) Base() memory.Relocatable {
	return r.base
}

func (r *OutputBuiltinRunner) Name() string {
	return OUTPUT_BUILTIN_NAME
}

func (r *OutputBuiltinRunner) InitializeSegments(segments *memory.MemorySegmentManager) {

}

func (r *OutputBuiltinRunner) InitialStack() []memory.MaybeRelocatable {

}

func (r *OutputBuiltinRunner) DeduceMemoryCell(rel memory.Relocatable, mem *memory.Memory) (*memory.MaybeRelocatable, error) {

}

func (r *OutputBuiltinRunner) AddValidationRule(mem *memory.Memory) {

}
