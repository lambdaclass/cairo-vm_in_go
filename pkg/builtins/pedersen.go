package builtins

import (
	starknet_crypto "github.com/lambdaclass/cairo-vm.go/pkg/starknet_crypto"
	"github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
)

const PEDERSEN_BUILTIN_NAME = "pedersen"
const PEDERSEN_CELLS_PER_INSTANCE = 3

type PedersenBuiltinRunner struct {
	base               memory.Relocatable
	included           bool
	verified_addresses []bool
}

func NewPedersenBuiltinRunner(included bool) *PedersenBuiltinRunner {
	return &PedersenBuiltinRunner{included: included}
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

func (p *PedersenBuiltinRunner) InitialStack() []memory.MaybeRelocatable {
	if p.included {
		return []memory.MaybeRelocatable{*memory.NewMaybeRelocatableRelocatable(p.base)}
	} else {
		return nil
	}
}

func (p *PedersenBuiltinRunner) DeduceMemoryCell(address memory.Relocatable, mem *memory.Memory) (*memory.MaybeRelocatable, error) {
	if address.Offset%PEDERSEN_CELLS_PER_INSTANCE != 2 || p.CheckVerifiedAddresses(address) {
		return nil, nil
	}

	numA, err := mem.GetFelt(memory.Relocatable{SegmentIndex: address.SegmentIndex, Offset: address.Offset - 1})
	if err != nil {
		return nil, nil
	}

	numB, err := mem.GetFelt(memory.Relocatable{SegmentIndex: address.SegmentIndex, Offset: address.Offset - 2})
	if err != nil {
		return nil, nil
	}

	p.ResizeVerifiedAddresses(address)

	x := starknet_crypto.PedersenHash(numA, numB)

	fmt.Printf("starknet_crypto.PedersenHash(numA, numB): = %x \n", x.ToSignedFeltString())
	return memory.NewMaybeRelocatableFelt(x), nil
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
		for i := 0; i <= num+1; i++ {
			p.verified_addresses = append(p.verified_addresses, false)
		}

	}
	p.verified_addresses[address.Offset] = true
}
