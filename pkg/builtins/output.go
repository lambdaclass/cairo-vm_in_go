package builtins

import "github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"

const OUTPUT_BUILTIN_NAME = "output"

type OutputBuiltinRunner struct {
	base     memory.Relocatable
	included bool
}

func NewOutputBuiltinRunner() *OutputBuiltinRunner {
	return &OutputBuiltinRunner{}
}

func (r *OutputBuiltinRunner) Base() memory.Relocatable {
	return r.base
}

func (r *OutputBuiltinRunner) Name() string {
	return OUTPUT_BUILTIN_NAME
}

func (r *OutputBuiltinRunner) InitializeSegments(segments *memory.MemorySegmentManager) {
	r.base = segments.AddSegment()
}

func (r *OutputBuiltinRunner) InitialStack() []memory.MaybeRelocatable {
	if r.included {
		return []memory.MaybeRelocatable{*memory.NewMaybeRelocatableRelocatable(r.base)}
	}
	return []memory.MaybeRelocatable{}
}

func (r *OutputBuiltinRunner) DeduceMemoryCell(rel memory.Relocatable, mem *memory.Memory) (*memory.MaybeRelocatable, error) {
	return nil, nil
}

func (r *OutputBuiltinRunner) AddValidationRule(mem *memory.Memory) {}

func (r *OutputBuiltinRunner) Include(include bool) {
	r.included = include
}
