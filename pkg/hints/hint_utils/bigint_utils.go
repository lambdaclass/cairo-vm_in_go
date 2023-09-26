package hint_utils

import (
	"math/big"

	"github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
	. "github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
	. "github.com/lambdaclass/cairo-vm.go/pkg/types"
	. "github.com/lambdaclass/cairo-vm.go/pkg/vm"
	"github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
	. "github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
	"github.com/pkg/errors"
)

/*
Implements hint:
%{
    from starkware.cairo.common.cairo_secp.secp_utils import split

    segments.write_arg(ids.res.address_, split(value))
%}
*/

func NondetBigInt3(virtual_machine VirtualMachine, exec_scopes ExecutionScopes, ids_data IdsManager) error {
	res_relloc, err := ids_data.GetAddr("res", &virtual_machine)
	if err != nil {
		return err
	}

	value_uncast, err := exec_scopes.Get("value")
	if err != nil {
		return err
	}
	value := value_uncast.(big.Int)

	bigint3_split, err := Bigint3Split(value)
	if err != nil {
		return err
	}
	arg := make([]memory.MaybeRelocatable, 0)

	for i := 0; i < 3; i++ {
		m := memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromBigInt(&bigint3_split[i]))
		arg = append(arg, *m)
	}

	virtual_machine.Segments.LoadData(res_relloc, &arg)
	return nil
}

// In cairo, various structs are used to represent big integers, all of them have numbered fields of Felt type (d0, d1,...) and share the same behaviours
// This file contains an implementation of each behaviour at the limbs level, and the wrappers for each specific type

// Generic methods for all types
func limbsFromVarName(nLimbs int, name string, ids IdsManager, vm *VirtualMachine) ([]Felt, error) {
	baseAddr, err := ids.GetAddr(name, vm)
	if err != nil {
		return nil, err
	}
	return limbsFromBaseAddress(nLimbs, name, baseAddr, vm)
}

func limbsFromBaseAddress(nLimbs int, name string, addr Relocatable, vm *VirtualMachine) ([]Felt, error) {
	//fmt.Println("addr in libms base addr: ", addr)
	limbs := make([]Felt, 0)
	for i := 0; i < nLimbs; i++ {
		felt, err := vm.Segments.Memory.GetFelt(addr.AddUint(uint(i)))
		//fmt.Println("value in memory: ", felt.ToBigInt().Text(10), addr.AddUint(uint(i)))
		if err == nil {
			limbs = append(limbs, felt)
		} else {
			//fmt.Println("error name: ", name)
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

func limbsInsertFromVarName(limbs []Felt, name string, ids IdsManager, vm *VirtualMachine) error {
	baseAddr, err := ids.GetAddr(name, vm)
	if err != nil {
		return err
	}
	for i := 0; i < len(limbs); i++ {
		err = vm.Segments.Memory.Insert(baseAddr.AddUint(uint(i)), NewMaybeRelocatableFelt(limbs[i]))
		if err != nil {
			return err
		}
	}
	return nil
}

func splitIntoLimbs(num *big.Int, numLimbs int) []Felt {
	limbs := make([]Felt, 0, numLimbs)
	bitmask := new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 128), big.NewInt(1))
	for i := 0; i < numLimbs; i++ {
		limbs[i] = FeltFromBigInt(new(big.Int).Lsh(new(big.Int).And(num, bitmask), 128))
	}
	return limbs
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
