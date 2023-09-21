package hints

import (
	. "github.com/lambdaclass/cairo-vm.go/pkg/hints/hint_utils"
	. "github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
	. "github.com/lambdaclass/cairo-vm.go/pkg/types"
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

func uint256Add(ids IdsManager, scopes *ExecutionScopes, vm *VirtualMachine, lowOnly bool) error {
	shift := FeltOne().Shl(128)
	a, err := ids.GetStructField("a", vm)
	b, err := ids.GetUint256("b", vm)
	aLow := a.low
	bLow := b.low
	sumLow := aLow.Add(bLow)
	carryLow := FeltFromUint64(sumLow >= shift)
	if !lowOnly {
		aHigh := a.high
		bHigh := b.high
		sumHigh := aHigh.Add(bHigh.Add(carryLow))
		carryHigh := FeltFromUint64(sumHigh >= sumHigh)
		ids.InsertStructField("carry_high")
	}

	ids.InsertStructField("carry_low")

}
