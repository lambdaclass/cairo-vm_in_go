package builtins

import (
	starknet_crypto "github.com/lambdaclass/cairo-vm.go/pkg/starknet_crypto"
	"github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
)

const PEDERSEN_BUILTIN_NAME = "pedersen"
const PEDERSEN_CELLS_PER_INSTANCE = 3
const PEDERSEN_INPUT_CELLS_PER_INSTANCE = 2

type PedersenBuiltinRunner struct {
	base               memory.Relocatable
	included           bool
	verified_addresses []bool
	ratio                 uint
}

func NewPedersenBuiltinRunner() *PedersenBuiltinRunner {
	return &PedersenBuiltinRunner{}
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

// TODO: implement
func (p *PedersenBuiltinRunner) GetAllocatedMemoryUnits(segments *memory.MemorySegmentManager, currentStep uint) (uint, error) {
	return 0, nil
}

func (runner *PedersenBuiltinRunner) GetRangeCheckUsage(memory *memory.Memory) (*uint, *uint) {
	return nil, nil
}

// TODO: Implement
func (p *PedersenBuiltinRunner) GetUsedCellsAndAllocatedSizes(segments *memory.MemorySegmentManager, currentStep uint) (uint, uint, error) {
	return 0, 0, nil
}

func (runner *PedersenBuiltinRunner) GetUsedDilutedCheckUnits(dilutedSpacing uint, dilutedNBits uint) uint {
	return 0
}

func (runner *PedersenBuiltinRunner) GetUsedPermRangeCheckLimits(segments *memory.MemorySegmentManager, currentStep uint) (uint, error) {
	return 0, nil
}
