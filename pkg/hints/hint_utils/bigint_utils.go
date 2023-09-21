package hint_utils

import (
	"math/big"

	. "github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
	. "github.com/lambdaclass/cairo-vm.go/pkg/vm"
	. "github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
	"github.com/pkg/errors"
)

// Generic methods for all types
func limbsFromVarName(nLimbs int, name string, ids IdsManager, vm *VirtualMachine) ([]Felt, error) {
	baseAddr, err := ids.GetAddr(name, vm)
	if err != nil {
		return nil, err
	}
	return limbsFromBaseAddress(nLimbs, name, baseAddr, vm)
}

func limbsFromBaseAddress(nLimbs int, name string, addr Relocatable, vm *VirtualMachine) ([]Felt, error) {
	limbs := make([]Felt, 0)
	for i := 0; i < nLimbs; i++ {
		felt, err := vm.Segments.Memory.GetFelt(addr.AddUint(uint(i)))
		if err == nil {
			limbs = append(limbs, felt)
		} else {
			return nil, errors.Errorf("Identifier %s has no member d%d", name, i)
		}
	}
	return limbs, nil
}

func limbsPack86(limbs []Felt) big.Int {
	sum := big.NewInt(0)
	for i := 0; i < 3; i++ {
		felt := limbs[i]
		shifed := new(big.Int).Lsh(felt.ToSigned(), uint(i*86))
		sum.Add(sum, shifed)
	}
	return *sum
}

func limbsPack(limbs []Felt) big.Int {
	sum := big.NewInt(0)
	for i := 0; i < len(limbs); i++ {
		felt := limbs[i]
		shifed := new(big.Int).Lsh(felt.ToSigned(), uint(i*128))
		sum.Add(sum, shifed)
	}
	return *sum
}

// Concrete type definitions

// BigInt3

type BigInt3 struct {
	Limbs []Felt
}

func (b *BigInt3) Pack86() big.Int {
	return limbsPack86(b.Limbs)
}

func BigInt3FromBaseAddr(addr Relocatable, name string, vm *VirtualMachine) (BigInt3, error) {
	limbs, err := limbsFromBaseAddress(3, name, addr, vm)
	return BigInt3{Limbs: limbs}, err
}
