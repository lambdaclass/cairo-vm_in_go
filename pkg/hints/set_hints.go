package hints

import (
	. "github.com/lambdaclass/cairo-vm.go/pkg/hints/hint_utils"
	. "github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
	. "github.com/lambdaclass/cairo-vm.go/pkg/vm"
	. "github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
	"github.com/pkg/errors"
)

/*
	Implements hint:

	assert ids.elm_size > 0
	assert ids.set_ptr <= ids.set_end_ptr
	elm_list = memory.get_range(ids.elm_ptr, ids.elm_size)
	for i in range(0, ids.set_end_ptr - ids.set_ptr, ids.elm_size):
	    if memory.get_range(ids.set_ptr + i, ids.elm_size) == elm_list:
	        ids.index = i // ids.elm_size
	        ids.is_elm_in_set = 1
	        break
	else:
	    ids.is_elm_in_set = 0
*/
func set_add(ids IdsManager, vm *VirtualMachine) error {
	set_ptr, err := ids.GetRelocatable("set_ptr", vm)
	if err != nil {
		return err
	}
	elm_size_felt, err := ids.GetFelt("elm_size", vm)
	if err != nil {
		return err
	}
	elm_ptr, err := ids.GetRelocatable("elm_ptr", vm)
	if err != nil {
		return err
	}
	set_end_ptr, err := ids.GetRelocatable("set_end_ptr", vm)
	if err != nil {
		return err
	}

	if elm_size_felt.IsZero() {
		return errors.Errorf("assert ids.elm_size > 0")
	}

	elm_size, err := elm_size_felt.ToU64()

	if err != nil {
		return err
	}

	if set_ptr.Offset > set_end_ptr.Offset {
		return errors.Errorf("expected set_ptr: %v <= set_end_ptr: %v", set_ptr, set_end_ptr)
	}

	elem := vm.Segments.Memory.GetRange(elm_ptr, uint(elm_size))

	for i := uint(0); i < set_end_ptr.Offset-set_ptr.Offset; i += uint(elm_size) {
		other_elm := vm.Segments.Memory.GetRange(set_ptr.AddUint(i), uint(elm_size))
		if &elem == &other_elm {
			return ids.Insert("is_elm_in_set", NewMaybeRelocatableFelt(FeltOne()), vm)
		}
	}

	return ids.Insert("is_elm_in_set", NewMaybeRelocatableFelt(FeltZero()), vm)
}
