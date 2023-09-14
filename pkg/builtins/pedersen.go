package builtins

import (
	starknet_crypto "github.com/lambdaclass/cairo-vm.go/pkg/starknet_crypto"
	"github.com/lambdaclass/cairo-vm.go/pkg/utils"
	"github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
	"github.com/pkg/errors"
)

const PEDERSEN_BUILTIN_NAME = "pedersen"
const PEDERSEN_CELLS_PER_INSTANCE = 3
const PEDERSEN_INPUT_CELLS_PER_INSTANCE = 2

type PedersenBuiltinRunner struct {
	base                  memory.Relocatable
	included              bool
	verified_addresses    []bool
	ratio                 uint
	instancesPerComponent uint
}

func NewPedersenBuiltinRunner() *PedersenBuiltinRunner {
	return &PedersenBuiltinRunner{instancesPerComponent: 1}
}

func DefaultPedersenBuiltinRunner() *PedersenBuiltinRunner {
	return &PedersenBuiltinRunner{
		ratio:                 8,
		instancesPerComponent: 1,
	}
}

func (r *PedersenBuiltinRunner) Include(include bool) {
	r.included = include
}

func (p *PedersenBuiltinRunner) Base() memory.Relocatable {
	return p.base
}

func (p *PedersenBuiltinRunner) Name() string {
	return PEDERSEN_BUILTIN_NAME
}

func (p *PedersenBuiltinRunner) InitializeSegments(segments *memory.MemorySegmentManager) {
	p.base = segments.AddSegment()
}

func (p *PedersenBuiltinRunner) Ratio() uint {
	return p.ratio
}

func (p *PedersenBuiltinRunner) InitialStack() []memory.MaybeRelocatable {
	if p.included {
		return []memory.MaybeRelocatable{*memory.NewMaybeRelocatableRelocatable(p.base)}
	} else {
		return nil
	}
}

func (p *PedersenBuiltinRunner) DeduceMemoryCell(address memory.Relocatable, mem *memory.Memory) (*memory.MaybeRelocatable, error) {
	if address.Offset%PEDERSEN_CELLS_PER_INSTANCE != PEDERSEN_INPUT_CELLS_PER_INSTANCE || p.CheckVerifiedAddresses(address) {
		return nil, nil
	}

	feltA, err := mem.GetFelt(memory.Relocatable{SegmentIndex: address.SegmentIndex, Offset: address.Offset - 1})
	if err != nil {
		return nil, nil
	}

	feltB, err := mem.GetFelt(memory.Relocatable{SegmentIndex: address.SegmentIndex, Offset: address.Offset - 2})
	if err != nil {
		return nil, nil
	}

	p.ResizeVerifiedAddresses(address)

	hash := starknet_crypto.PedersenHash(feltB, feltA)

	return memory.NewMaybeRelocatableFelt(hash), nil
}

func (p *PedersenBuiltinRunner) AddValidationRule(*memory.Memory) {
}

func (p *PedersenBuiltinRunner) CheckVerifiedAddresses(address memory.Relocatable) bool {
	if len(p.verified_addresses) < int(address.Offset) {
		return false
	}

	return p.verified_addresses[address.Offset]
}

func (p *PedersenBuiltinRunner) ResizeVerifiedAddresses(address memory.Relocatable) {
	num := int(address.Offset) - len(p.verified_addresses)
	if num > 0 {
		for i := 0; i <= num; i++ {
			p.verified_addresses = append(p.verified_addresses, false)
		}

	}
	p.verified_addresses[address.Offset] = true
}

func (p *PedersenBuiltinRunner) CellsPerInstance() uint {
	return PEDERSEN_CELLS_PER_INSTANCE
}

func (p *PedersenBuiltinRunner) GetAllocatedMemoryUnits(segments *memory.MemorySegmentManager, currentStep uint) (uint, error) {
	// This condition corresponds to an uninitialized ratio for the builtin, which should only
	// happen when layout is `dynamic`
	if p.Ratio() == 0 {
		// Dynamic layout has the exact number of instances it needs (up to a power of 2).
		used, err := segments.GetSegmentUsedSize(uint(p.base.SegmentIndex))
		if err != nil {
			return 0, err
		}
		instances := used / p.CellsPerInstance()
		components := utils.NextPowOf2(instances / p.instancesPerComponent)
		size := p.CellsPerInstance() * p.instancesPerComponent * components

		return size, nil
	}

	minStep := p.ratio * p.instancesPerComponent
	if currentStep < minStep {
		return 0, errors.Errorf("number of steps must be at least %d for the %s builtin", minStep, p.Name())
	}
	value, err := utils.SafeDiv(currentStep, p.ratio)

	if err != nil {
		return 0, errors.Errorf("error calculating builtin memory units: %s", err)
	}

	return p.CellsPerInstance() * value, nil
}

func (runner *PedersenBuiltinRunner) GetRangeCheckUsage(memory *memory.Memory) (*uint, *uint) {
	return nil, nil
}

func (p *PedersenBuiltinRunner) GetUsedCellsAndAllocatedSizes(segments *memory.MemorySegmentManager, currentStep uint) (uint, uint, error) {
	used, err := segments.GetSegmentUsedSize(uint(p.base.SegmentIndex))
	if err != nil {
		return 0, 0, err
	}

	size, err := p.GetAllocatedMemoryUnits(segments, currentStep)

	if used > size {
		return 0, 0, errors.Errorf("The builtin %s used %d cells but the capacity is %d", p.Name(), used, size)
	}

	return used, size, nil
}

func (runner *PedersenBuiltinRunner) GetUsedDilutedCheckUnits(dilutedSpacing uint, dilutedNBits uint) uint {
	return 0
}

func (runner *PedersenBuiltinRunner) GetUsedPermRangeCheckLimits(segments *memory.MemorySegmentManager, currentStep uint) (uint, error) {
	return 0, nil
}

func (runner *PedersenBuiltinRunner) GetMemoryAccesses(manager *memory.MemorySegmentManager) ([]memory.Relocatable, error) {
	segmentSize := manager.SegmentSizes[uint(runner.Base().SegmentIndex)]

	var ret []memory.Relocatable

	var i uint
	for i = 0; i < segmentSize; i++ {
		ret = append(ret, memory.NewRelocatable(runner.Base().SegmentIndex, i))
	}

	return ret, nil
}
