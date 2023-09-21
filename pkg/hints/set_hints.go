package hints

import (
	. "github.com/lambdaclass/cairo-vm.go/pkg/hints/hint_utils"
	. "github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
	. "github.com/lambdaclass/cairo-vm.go/pkg/vm"
	. "github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
	"github.com/pkg/errors"
	"reflect"
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
func setAdd(ids IdsManager, vm *VirtualMachine) error {
	setPtr, err := ids.GetRelocatable("set_ptr", vm)
	if err != nil {
		return err
	}
	elmSizeFelt, err := ids.GetFelt("elm_size", vm)
	if err != nil {
		return err
	}
	elmPtr, err := ids.GetRelocatable("elm_ptr", vm)
	if err != nil {
		return err
	}
	setEndPtr, err := ids.GetRelocatable("set_end_ptr", vm)
	if err != nil {
		return err
	}

	if elmSizeFelt.IsZero() {
		return errors.Errorf("assert ids.elm_size > 0")
	}

	elmSize, err := elmSizeFelt.ToUint()

	if err != nil {
		return err
	}

	if setPtr.Offset > setEndPtr.Offset {
		return errors.Errorf("expected set_ptr: %v <= set_end_ptr: %v", setPtr, setEndPtr)
	}

	elem, err := vm.Segments.Memory.GetRange(elmPtr, elmSize)
	if err != nil {
		return err
	}

	for i := uint(0); i < setEndPtr.Offset-setPtr.Offset-elmSize; i++ {
		otherElm, err := vm.Segments.Memory.GetRange(setPtr.AddUint(i*elmSize), elmSize)
		if err != nil {
			return err
		}
		if reflect.DeepEqual(elem, otherElm) {
			err := ids.Insert("index", NewMaybeRelocatableFelt(FeltFromUint(i)), vm)
			if err != nil {
				return err
			}
			return ids.Insert("is_elm_in_set", NewMaybeRelocatableFelt(FeltOne()), vm)
		}
	}

	return ids.Insert("is_elm_in_set", NewMaybeRelocatableFelt(FeltZero()), vm)
}
