package hints

import (
	. "github.com/lambdaclass/cairo-vm.go/pkg/hints/hint_utils"
	. "github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
	. "github.com/lambdaclass/cairo-vm.go/pkg/vm"
	. "github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
)

func sha256Input(ids IdsManager, vm *VirtualMachine) error {
	nBytes, err := ids.GetFelt("n_bytes", vm)
	if err != nil {
		return err
	}
	if nBytes.Cmp(FeltFromUint(4)) != -1 {
		return ids.Insert("full_word", NewMaybeRelocatableFelt(FeltOne()), vm)
	}
	return ids.Insert("full_word", NewMaybeRelocatableFelt(FeltZero()), vm)
}
