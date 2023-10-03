package hints

import (
	. "github.com/lambdaclass/cairo-vm.go/pkg/hints/hint_utils"
	. "github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
	. "github.com/lambdaclass/cairo-vm.go/pkg/vm"
	. "github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
)

/*
Implements Hint:

	%{
	    ids.low = ids.a & ((1<<128) - 1)
	    ids.high = ids.a >> 128
	%}
*/
func uint384Split128(ids IdsManager, vm *VirtualMachine) error {
	a, err := ids.GetFelt("a", vm)
	if err != nil {
		return err
	}
	low := a.And(FeltFromDecString("340282366920938463463374607431768211455"))
	err = ids.Insert("low", NewMaybeRelocatableFelt(low), vm)
	if err != nil {
		return err
	}
	high := a.Shr(128)
	return ids.Insert("high", NewMaybeRelocatableFelt(high), vm)
}
