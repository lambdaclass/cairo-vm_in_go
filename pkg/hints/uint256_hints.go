package hints

import (
	"math/big"

	. "github.com/lambdaclass/cairo-vm.go/pkg/hints/hint_utils"
	"github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
	. "github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
	. "github.com/lambdaclass/cairo-vm.go/pkg/vm"
	. "github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
	"github.com/pkg/errors"
)

/*
Implements hints:
%{
    sum_low = ids.a.low + ids.b.low
    ids.carry_low = 1 if sum_low >= ids.SHIFT else 0
    sum_high = ids.a.high + ids.b.high + ids.carry_low
    ids.carry_high = 1 if sum_high >= ids.SHIFT else 0
%}
%{
    sum_low = ids.a.low + ids.b.low
    ids.carry_low = 1 if sum_low >= ids.SHIFT else 0
%}
*/

func uint256Add(ids IdsManager, vm *VirtualMachine, lowOnly bool) error {
	shift := FeltOne().Shl(128)
	aLow, err := ids.GetStructFieldFelt("a", 0, vm)
	if err != nil {
		return err
	}

	bLow, err := ids.GetStructFieldFelt("b", 0, vm)
	if err != nil {
		return err
	}

	sumLow := aLow.Add(bLow)
	var carryLow Felt
	switch sumLow.Cmp(shift) {
	case -1:
		carryLow = FeltZero()
	default:
		carryLow = FeltOne()
	}

	if !lowOnly {
		aHigh, err := ids.GetStructFieldFelt("a", 1, vm)
		if err != nil {
			return err
		}
		bHigh, err := ids.GetStructFieldFelt("b", 1, vm)
		if err != nil {
			return err
		}

		sumHigh := aHigh.Add(bHigh.Add(carryLow))
		var carryHigh Felt
		switch sumHigh.Cmp(shift) {
		case -1:
			carryHigh = FeltZero()
		default:
			carryHigh = FeltOne()
		}
		ids.Insert("carry_high", NewMaybeRelocatableFelt(carryHigh), vm)
	}

	return ids.Insert("carry_low", NewMaybeRelocatableFelt(carryLow), vm)

}

/*
Implements hint:

	%{
	    ids.low = ids.a & ((1<<64) - 1)
	    ids.high = ids.a >> 64

%}
*/
func split64(ids IdsManager, vm *VirtualMachine) error {
	a, err := ids.GetFelt("a", vm)
	if err != nil {
		return err
	}
	flag := (FeltOne().Shl(64)).Sub(FeltOne()) // (1 << 64) - 1
	low := a.And(flag)
	high := a.Shr(64) // a >> 64
	err = ids.Insert("low", NewMaybeRelocatableFelt(low), vm)
	if err != nil {
		return err
	}
	err = ids.Insert("high", NewMaybeRelocatableFelt(high), vm)
	if err != nil {
		return err
	}
	return nil

}

/*
Implements hint:

	%{
	    from starkware.python.math_utils import isqrt
	    n = (ids.n.high << 128) + ids.n.low
	    root = isqrt(n)
	    assert 0 <= root < 2 ** 128
	    ids.root.low = root
	    ids.root.high = 0

%}
*/
func uint256Sqrt(ids IdsManager, vm *VirtualMachine, onlyLow bool) error {
	uintN, err := ids.GetUint256("n", vm)
	if err != nil {
		return err
	}

	bHigh := new(big.Int).Lsh(uintN.High.ToBigInt(), 128)
	bLow := uintN.Low.ToBigInt()
	n := new(big.Int).Add(bHigh, bLow)
	root := new(big.Int).Sqrt(n)

	if root.BitLen() > 128 {
		return errors.Errorf("assert 0 <= %d < 2**128", root)
	}

	feltRoot := FeltFromBigInt(root)

	if onlyLow {
		return ids.Insert("root", NewMaybeRelocatableFelt(feltRoot), vm)
	} else {
		return ids.InsertUint256("root", lambdaworks.Uint256{Low: feltRoot, High: FeltZero()}, vm)
	}
}
