package builtins

import (
	"github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
)

const OUTPUT_BUILTIN_NAME = "output"

type OutputBuiltinRunner struct {
	base     memory.Relocatable
	included bool
}

func NewOutputBuiltinRunner() *OutputBuiltinRunner {
	return &OutputBuiltinRunner{}
}

func (o *OutputBuiltinRunner) Base() memory.Relocatable {
	return o.base
}

func (o *OutputBuiltinRunner) Name() string {
	return OUTPUT_BUILTIN_NAME
}

func (o *OutputBuiltinRunner) InitializeSegments(segments *memory.MemorySegmentManager) {
	o.base = segments.AddSegment()
}

func (o *OutputBuiltinRunner) InitialStack() []memory.MaybeRelocatable {
	if o.included {
		return []memory.MaybeRelocatable{*memory.NewMaybeRelocatableRelocatable(o.base)}
	}
	return []memory.MaybeRelocatable{}
}

func (o *OutputBuiltinRunner) DeduceMemoryCell(rel memory.Relocatable, mem *memory.Memory) (*memory.MaybeRelocatable, error) {
	return nil, nil
}

func (o *OutputBuiltinRunner) AddValidationRule(mem *memory.Memory) {}

func (o *OutputBuiltinRunner) Include(include bool) {
	o.included = include
}

func (o *OutputBuiltinRunner) Ratio() uint {
	return 0
}

// func (r *OutputBuiltinRunner) GetUsedCells(segments *memory.MemorySegmentManager) (uint, error) {
// 	used, err := segments.GetSegmentUsedSize(uint(r.base.SegmentIndex))
// 	if err != nil {
// 		return 0, err
// 	}
// 	return used, nil
// }

func (o *OutputBuiltinRunner) GetAllocatedMemoryUnits(segments *memory.MemorySegmentManager, currentStep uint) (uint, error) {
	return 0, nil
}

func (o *OutputBuiltinRunner) GetUsedCellsAndAllocatedSizes(segments *memory.MemorySegmentManager, currentStep uint) (uint, uint, error) {
	used, err := segments.GetSegmentUsedSize(uint(o.base.SegmentIndex))
	if err != nil {
		return 0, 0, err
	}
	return used, used, nil
}
