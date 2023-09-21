package hints

import (
	. "github.com/lambdaclass/cairo-vm.go/pkg/hints/hint_utils"
	. "github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
	. "github.com/lambdaclass/cairo-vm.go/pkg/vm"
	. "github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
)

type Uint256 struct {
	low  Felt
	high Felt
}

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
