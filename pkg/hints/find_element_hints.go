package hints

import (
	. "github.com/lambdaclass/cairo-vm.go/pkg/hints/hint_utils"
	. "github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
	. "github.com/lambdaclass/cairo-vm.go/pkg/types"
	. "github.com/lambdaclass/cairo-vm.go/pkg/vm"
	. "github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
	"github.com/pkg/errors"
)

func find_element(ids IdsManager, vm *VirtualMachine, execScopes ExecutionScopes) error {
	array_ptr, err := ids.GetRelocatable("array_ptr", vm)
	if err != nil {
		return err
	}

	key, err := ids.GetFelt("key", vm)
	if err != nil {
		return err
	}

	elm_size_felt, err := ids.GetFelt("elm_size", vm)
	if err != nil {
		return err
	}
	elm_size, err := elm_size_felt.ToUint()
	if err != nil {
		return err
	}

	n_elms, err := ids.GetFelt("n_elms", vm)
	if err != nil {
		return err
	}
	n_elms_iter, err := n_elms.ToUint()
	if err != nil {
		return err
	}

	find_element_index_uncast, err := execScopes.Get("find_element_index")
	if err == nil {
		find_element_index, ok := find_element_index_uncast.(Felt)
		if !ok {
			return ConversionError(find_element_index, "felt")
		}
		position, err := array_ptr.AddFelt(find_element_index.Mul(elm_size_felt))
		if err != nil {
			return err
		}

		found_key, err := vm.Segments.Memory.GetFelt(position)
		if err != nil {
			return err
		}
		if found_key != key {
			return errors.Errorf(
				"Invalid index found in find_element_index. Index: %s.\nExpected key: %s, found_key %s",
				find_element_index.ToSignedFeltString(),
				key.ToSignedFeltString(),
				found_key.ToSignedFeltString(),
			)
		}
		execScopes.DeleteVariable("find_element_index")
		return ids.Insert("index", NewMaybeRelocatableFelt(find_element_index), vm)
	}

	find_element_max_size_uncast, err := execScopes.Get("find_element_max_size")
	if err == nil {
		find_element_max_size, ok := find_element_max_size_uncast.(Felt)
		if !ok {
			return ConversionError(find_element_max_size, "felt")
		}
		if n_elms.Cmp(find_element_max_size) == 1 {
			return errors.Errorf(
				"find_element() can only be used with n_elms <= %s.\nGot: n_elms = %s",
				find_element_max_size.ToSignedFeltString(),
				n_elms.ToSignedFeltString(),
			)
		}
	}

	for i := uint(0); i < n_elms_iter; i++ {
		iter_key, err := vm.Segments.Memory.GetFelt(array_ptr.AddUint(i * elm_size))
		if err != nil {
			return err
		}
		if iter_key == key {
			return ids.Insert("index", NewMaybeRelocatableFelt(FeltFromUint(i)), vm)
		}
	}

	return errors.Errorf("Key: %v was not found", key)
}

func search_sorted_lower(ids IdsManager, vm *VirtualMachine, execScopes ExecutionScopes) error {
	array_ptr, err := ids.GetRelocatable("array_ptr", vm)
	if err != nil {
		return err
	}

	key, err := ids.GetFelt("key", vm)
	if err != nil {
		return err
	}

	elm_size_felt, err := ids.GetFelt("elm_size", vm)
	if err != nil {
		return err
	}
	elm_size, err := elm_size_felt.ToUint()
	if err != nil {
		return err
	}

	n_elms, err := ids.GetFelt("n_elms", vm)
	if err != nil {
		return err
	}
	n_elms_iter, err := n_elms.ToUint()
	if err != nil {
		return err
	}

	find_element_max_size_uncast, err := execScopes.Get("find_element_max_size")
	if err == nil {
		find_element_max_size, ok := find_element_max_size_uncast.(Felt)
		if !ok {
			return ConversionError(find_element_max_size, "felt")
		}
		if n_elms.Cmp(find_element_max_size) == 1 {
			return errors.Errorf(
				"find_element() can only be used with n_elms <= %s.\nGot: n_elms = %s",
				find_element_max_size.ToSignedFeltString(),
				n_elms.ToSignedFeltString(),
			)
		}
	}

	for i := uint(0); i < n_elms_iter; i++ {
		iter_key, err := vm.Segments.Memory.GetFelt(array_ptr.AddUint(i * elm_size))
		if err != nil {
			return err
		}
		if iter_key == key || iter_key.Cmp(key) == 1 {
			return ids.Insert("index", NewMaybeRelocatableFelt(FeltFromUint(i)), vm)
		}
	}

	return errors.Errorf("Key: %v was not found", key)
}
