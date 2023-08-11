package builtins

import (
	. "github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
	. "github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
)

type KeccakBuiltinRunner struct {
	base     Relocatable
	included bool
	cache    map[Relocatable]Felt
}

func NewKeccakBuiltinRunner(included bool) *KeccakBuiltinRunner {
	return &KeccakBuiltinRunner{included: included, cache: make(map[Relocatable]Felt)}
}

const KECCAK_BUILTIN_NAME = "keccak"

func (k *KeccakBuiltinRunner) Base() Relocatable {
	return k.base
}

func (k *KeccakBuiltinRunner) Name() string {
	return KECCAK_BUILTIN_NAME
}

func (k *KeccakBuiltinRunner) InitializeSegments(segments *MemorySegmentManager) {
	k.base = segments.AddSegment()
}

func (k *KeccakBuiltinRunner) InitialStack() []MaybeRelocatable {
	if k.included {
		return []MaybeRelocatable{*NewMaybeRelocatableRelocatable(k.base)}
	} else {
		return nil
	}
}

func (k *KeccakBuiltinRunner) AddValidationRule(*Memory) {}

func (k *KeccakBuiltinRunner) DeduceMemoryCell(address Relocatable, mem *Memory) (*MaybeRelocatable, error) {
}
