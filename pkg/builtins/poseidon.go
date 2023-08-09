package builtins

import "github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"

const POSEIDON_BUILTIN_NAME = "poseidon"

type PoseidonBuiltinRunner struct {
	base     memory.Relocatable
	included bool
}

func NewPoseidonBuiltinRunner(included bool) *PoseidonBuiltinRunner {
	return &PoseidonBuiltinRunner{included: included}
}

func (p *PoseidonBuiltinRunner) Base() memory.Relocatable {
	return p.base
}

func (p *PoseidonBuiltinRunner) Name() string {
	return POSEIDON_BUILTIN_NAME
}

func (p *PoseidonBuiltinRunner) InitializeSegments(segments *memory.MemorySegmentManager) {
	p.base = segments.AddSegment()
}

func (p *PoseidonBuiltinRunner) InitialStack() []memory.MaybeRelocatable {
	if p.included {
		return []memory.MaybeRelocatable{*memory.NewMaybeRelocatableRelocatable(p.base)}
	} else {
		return nil
	}
}

func (p *PoseidonBuiltinRunner) DeduceMemoryCell(memory.Relocatable, *memory.Memory) (*memory.MaybeRelocatable, error) {
	return nil, nil //TODO
}

func (p *PoseidonBuiltinRunner) AddValidationRule(*memory.Memory) {
}
