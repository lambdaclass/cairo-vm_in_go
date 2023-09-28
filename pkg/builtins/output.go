package builtins

import (
	"github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
)

const OUTPUT_BUILTIN_NAME = "output"

type OutputBuiltinRunner struct {
	base     memory.Relocatable
	included bool
	StopPtr  *uint
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

func (runner *OutputBuiltinRunner) GetRangeCheckUsage(memory *memory.Memory) (*uint, *uint) {
	return nil, nil
}

func (runner *OutputBuiltinRunner) GetUsedPermRangeCheckLimits(segments *memory.MemorySegmentManager, currentStep uint) (uint, error) {
	return 0, nil
}

func (runner *OutputBuiltinRunner) GetUsedDilutedCheckUnits(dilutedSpacing uint, dilutedNBits uint) uint {
	return 0
}

func (runner *OutputBuiltinRunner) GetMemoryAccesses(manager *memory.MemorySegmentManager) ([]memory.Relocatable, error) {
	segmentSize, err := manager.GetSegmentSize(uint(runner.Base().SegmentIndex))
	if err != nil {
		return []memory.Relocatable{}, err
	}

	var ret []memory.Relocatable

	var i uint
	for i = 0; i < segmentSize; i++ {
		ret = append(ret, memory.NewRelocatable(runner.Base().SegmentIndex, i))
	}

	return ret, nil
}

func (r *OutputBuiltinRunner) FinalStack(segments *memory.MemorySegmentManager, pointer memory.Relocatable) (memory.Relocatable, error) {
	if r.included {
		if pointer.Offset == 0 {
			return memory.Relocatable{}, NewErrNoStopPointer(r.Name())
		}

		stopPointerAddr := memory.NewRelocatable(pointer.SegmentIndex, pointer.Offset-1)

		stopPointer, err := segments.Memory.GetRelocatable(stopPointerAddr)
		if err != nil {
			return memory.Relocatable{}, err
		}

		if r.Base().SegmentIndex != stopPointer.SegmentIndex {
			return memory.Relocatable{}, NewErrInvalidStopPointerIndex(r.Name(), stopPointer, r.Base())
		}

		used, err := segments.GetSegmentUsedSize(uint(r.Base().SegmentIndex))
		if err != nil {
			return memory.Relocatable{}, err
		}

		if stopPointer.Offset != used {
			return memory.Relocatable{}, NewErrInvalidStopPointer(r.Name(), used, stopPointer)
		}

		r.StopPtr = &stopPointer.Offset

		return stopPointerAddr, nil
	} else {
		r.StopPtr = new(uint)
		*r.StopPtr = 0
		return pointer, nil
	}
}

func (r *OutputBuiltinRunner) GetUsedInstances(segments *memory.MemorySegmentManager) (uint, error) {
	usedCells, err := segments.GetSegmentUsedSize(uint(r.Base().SegmentIndex))
	if err != nil {
		return 0, nil
	}

	return usedCells, nil
}

func (b *OutputBuiltinRunner) GetMemorySegmentAddresses() (memory.Relocatable, memory.Relocatable, error) {
	if b.StopPtr == nil {
		return memory.Relocatable{}, memory.Relocatable{}, NewErrNoStopPointer(b.Name())
	}
	return b.base, memory.NewRelocatable(b.base.SegmentIndex, *b.StopPtr), nil
}

func (b *OutputBuiltinRunner) RunSecurityChecks(*memory.MemorySegmentManager) error {
	return nil
}
