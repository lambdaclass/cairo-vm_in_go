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
	prev_locs_exp_addr, err := ids.GetRelocatable("prev_locs", vm)
	if err != nil {
		return err
	}

	prev_locs_exp, err := vm.Segments.Memory.GetFelt(prev_locs_exp_addr.AddUint(4))
	if err != nil {
		return err
	}

	return ids.Insert("locs", NewMaybeRelocatableFelt(prev_locs_exp.And(FeltOne())), vm)
}
