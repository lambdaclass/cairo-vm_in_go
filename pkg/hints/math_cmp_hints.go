package hints

import (
	. "github.com/lambdaclass/cairo-vm.go/pkg/builtins"
	. "github.com/lambdaclass/cairo-vm.go/pkg/hints/hint_utils"
	. "github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
	. "github.com/lambdaclass/cairo-vm.go/pkg/vm"
	. "github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
)

// memory[ap] = 0 if 0 <= (ids.a % PRIME) < range_check_builtin.bound else 1
func isNN(ids IdsManager, vm *VirtualMachine) error {
	a, err := ids.GetFelt("a", vm)
	if err != nil {
		return err
	}
	if a.Bits() < RANGE_CHECK_N_PARTS*INNER_RC_BOUND_SHIFT {
		return vm.Segments.Memory.Insert(vm.RunContext.Ap, NewMaybeRelocatableFelt(FeltZero()))
	}
	return vm.Segments.Memory.Insert(vm.RunContext.Ap, NewMaybeRelocatableFelt(FeltOne()))
}

// memory[ap] = 0 if 0 <= ((-ids.a - 1) % PRIME) < range_check_builtin.bound else 1
func isNNOutOfRange(ids IdsManager, vm *VirtualMachine) error {
	a, err := ids.GetFelt("a", vm)
	if err != nil {
		return err
	}
	if (FeltZero().Sub(a).Sub(FeltOne())).Bits() < RANGE_CHECK_N_PARTS*INNER_RC_BOUND_SHIFT {
		return vm.Segments.Memory.Insert(vm.RunContext.Ap, NewMaybeRelocatableFelt(FeltZero()))
	}
	return vm.Segments.Memory.Insert(vm.RunContext.Ap, NewMaybeRelocatableFelt(FeltOne()))
}

// memory[ap] = 0 if (ids.a % PRIME) <= (ids.b % PRIME) else 1
func isLeFelt(ids IdsManager, vm *VirtualMachine) error {
	a, err := ids.GetFelt("a", vm)
	if err != nil {
		return err
	}
	b, err := ids.GetFelt("b", vm)
	if err != nil {
		return err
	}
	if a.Cmp(b) != 1 {
		return vm.Segments.Memory.Insert(vm.RunContext.Ap, NewMaybeRelocatableFelt(FeltZero()))
	}
	return vm.Segments.Memory.Insert(vm.RunContext.Ap, NewMaybeRelocatableFelt(FeltOne()))
}
