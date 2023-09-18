package hints

import (
	. "github.com/lambdaclass/cairo-vm.go/pkg/hints/hint_utils"
	. "github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
	. "github.com/lambdaclass/cairo-vm.go/pkg/vm"
	. "github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
)

// Implements hint:
// %{ ids.locs.bit = (ids.prev_locs.exp % PRIME) & 1 %}
func pow(ids IdsManager, vm *VirtualMachine) error {
	prev_locs_exp_addr, err := ids.GetAddr("prev_locs", vm)
	prev_locs_exp, _ := vm.Segments.Memory.GetFelt(prev_locs_exp_addr.AddUint(4))

	if err != nil {
		return err
	}

	ids.Insert("locs", NewMaybeRelocatableFelt(prev_locs_exp.And(FeltOne())), vm)
	return nil
}
