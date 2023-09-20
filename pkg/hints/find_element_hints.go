package hints

import (
	. "github.com/lambdaclass/cairo-vm.go/pkg/hints/hint_utils"
	. "github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
	. "github.com/lambdaclass/cairo-vm.go/pkg/vm"
	. "github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
	"github.com/pkg/errors"
)

func find_element(ids IdsManager, vm *VirtualMachine) error {
	key, err := ids.GetFelt("key", vm)
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
	for i := uint(0); i < n_elms_iter; i++ {
		iter_key, err := ids.GetStructFieldFelt("array_start", i, vm)
		if err != nil {
			return err
		}
		if iter_key == key {
			return ids.Insert("index", NewMaybeRelocatableFelt(FeltFromUint(i)), vm)
		}
	}
	return errors.Errorf("Key: %v was not found", key)
}
