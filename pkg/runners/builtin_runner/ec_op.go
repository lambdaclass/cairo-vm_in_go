package builtinrunner

import "github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"

type EcOpBuiltinRunner struct {
	included bool
	base     memory.Relocatable
}

func NewEcOpBuiltinRunner(included bool) *EcOpBuiltinRunner {
	return &EcOpBuiltinRunner{
		included: included,
	}
}

func (ec *EcOpBuiltinRunner) Base() memory.Relocatable {
	return ec.base
}

func (ec *EcOpBuiltinRunner) Name() string {
	return "bitwise"
}

func (ec *EcOpBuiltinRunner) InitializeSegments(segments *memory.MemorySegmentManager) {
	ec.base = segments.AddSegment()
}

func (ec *EcOpBuiltinRunner) InitialStack() []memory.MaybeRelocatable {
	if ec.included {
		return []memory.MaybeRelocatable{*memory.NewMaybeRelocatableRelocatable(ec.base)}
	} else {
		return []memory.MaybeRelocatable{}
	}
}

func (ec *EcOpBuiltinRunner) DeduceMemoryCell(address memory.Relocatable, segments *memory.Memory) (*memory.MaybeRelocatable, error) {
	return nil, nil
}
