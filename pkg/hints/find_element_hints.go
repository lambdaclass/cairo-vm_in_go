package hints

import (
	. "github.com/lambdaclass/cairo-vm.go/pkg/hints/hint_utils"
	. "github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
	. "github.com/lambdaclass/cairo-vm.go/pkg/types"
	. "github.com/lambdaclass/cairo-vm.go/pkg/vm"
	. "github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
	"github.com/pkg/errors"
)

func findElement(ids IdsManager, vm *VirtualMachine, execScopes ExecutionScopes) error {
	arrayPtr, err := ids.GetRelocatable("array_ptr", vm)
	if err != nil {
		return err
	}

	key, err := ids.GetFelt("key", vm)
	if err != nil {
		return err
	}

	elmSizeFelt, err := ids.GetFelt("elm_size", vm)
	if err != nil {
		return err
	}
	elmSize, err := elmSizeFelt.ToUint()
	if err != nil {
		return err
	}

	nElms, err := ids.GetFelt("n_elms", vm)
	if err != nil {
		return err
	}
	nElmsIter, err := nElms.ToUint()
	if err != nil {
		return err
	}

	findElementIndexUncast, err := execScopes.Get("find_element_index")
	if err == nil {
		findElementIndex, ok := findElementIndexUncast.(Felt)
		if !ok {
			return ConversionError(findElementIndex, "felt")
		}
		position, err := arrayPtr.AddFelt(findElementIndex.Mul(elmSizeFelt))
		if err != nil {
			return err
		}

		foundKey, err := vm.Segments.Memory.GetFelt(position)
		if err != nil {
			return err
		}
		if foundKey != key {
			return errors.Errorf(
				"Invalid index found in find_element_index. Index: %s.\nExpected key: %s, found_key %s",
				findElementIndex.ToSignedFeltString(),
				key.ToSignedFeltString(),
				foundKey.ToSignedFeltString(),
			)
		}
		execScopes.DeleteVariable("find_element_index")
		return ids.Insert("index", NewMaybeRelocatableFelt(findElementIndex), vm)
	}

	findElementMaxSizeUncast, err := execScopes.Get("find_element_max_size")
	if err == nil {
		findElementMaxSize, ok := findElementMaxSizeUncast.(Felt)
		if !ok {
			return ConversionError(findElementMaxSize, "felt")
		}
		if nElms.Cmp(findElementMaxSize) == 1 {
			return errors.Errorf(
				"find_element() can only be used with n_elms <= %s.\nGot: n_elms = %s",
				findElementMaxSize.ToSignedFeltString(),
				nElms.ToSignedFeltString(),
			)
		}
	}

	for i := uint(0); i < nElmsIter; i++ {
		iterKey, err := vm.Segments.Memory.GetFelt(arrayPtr.AddUint(i * elmSize))
		if err != nil {
			return err
		}
		if iterKey == key {
			return ids.Insert("index", NewMaybeRelocatableFelt(FeltFromUint(i)), vm)
		}
	}

	return errors.Errorf("Key: %v was not found", key)
}

func searchSortedLower(ids IdsManager, vm *VirtualMachine, execScopes ExecutionScopes) error {
	arrayPtr, err := ids.GetRelocatable("array_ptr", vm)
	if err != nil {
		return err
	}

	key, err := ids.GetFelt("key", vm)
	if err != nil {
		return err
	}

	elmSizeFelt, err := ids.GetFelt("elm_size", vm)
	if err != nil {
		return err
	}
	elmSize, err := elmSizeFelt.ToUint()
	if err != nil {
		return err
	}

	nElms, err := ids.GetFelt("n_elms", vm)
	if err != nil {
		return err
	}
	nElmsIter, err := nElms.ToUint()
	if err != nil {
		return err
	}

	findElementMaxSizeUncast, err := execScopes.Get("find_element_max_size")
	if err == nil {
		findElementMaxSize, ok := findElementMaxSizeUncast.(Felt)
		if !ok {
			return ConversionError(findElementMaxSize, "felt")
		}
		if nElms.Cmp(findElementMaxSize) == 1 {
			return errors.Errorf(
				"find_element() can only be used with n_elms <= %s.\nGot: n_elms = %s",
				findElementMaxSize.ToSignedFeltString(),
				nElms.ToSignedFeltString(),
			)
		}
	}

	for i := uint(0); i < nElmsIter; i++ {
		iterKey, err := vm.Segments.Memory.GetFelt(arrayPtr.AddUint(i * elmSize))
		if err != nil {
			return err
		}
		if iterKey == key || iterKey.Cmp(key) == 1 {
			return ids.Insert("index", NewMaybeRelocatableFelt(FeltFromUint(i)), vm)
		}
	}

	return errors.Errorf("Key: %v was not found", key)
}
