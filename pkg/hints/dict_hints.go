package hints

import (
	"errors"

	. "github.com/lambdaclass/cairo-vm.go/pkg/hints/dict_manager"
	. "github.com/lambdaclass/cairo-vm.go/pkg/hints/hint_utils"
	. "github.com/lambdaclass/cairo-vm.go/pkg/types"
	. "github.com/lambdaclass/cairo-vm.go/pkg/vm"
	"github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
)

const DICT_ACCESS_SIZE = 3

func FetchDictManager(scopes *ExecutionScopes) (*DictManager, bool) {
	dictManager, err := scopes.Get("__dict_manager")
	if err != nil {
		return nil, false
	}
	val, ok := dictManager.(*DictManager)
	return val, ok
}

func defaultDictNew(ids IdsManager, scopes *ExecutionScopes, vm *VirtualMachine) error {
	defaultValue, err := ids.Get("default_value", vm)
	if err != nil {
		return err
	}
	dictManager, ok := FetchDictManager(scopes)
	if !ok {
		newDictManager := NewDictManager()
		dictManager = &newDictManager
		scopes.AssignOrUpdateVariable("__dict_manager", dictManager)
	}
	base := dictManager.NewDefaultDictionary(defaultValue, vm)
	return vm.Segments.Memory.Insert(vm.RunContext.Ap, memory.NewMaybeRelocatableRelocatable(base))
}

func dictRead(ids IdsManager, scopes *ExecutionScopes, vm *VirtualMachine) error {
	// Extract Variables
	dictManager, ok := FetchDictManager(scopes)
	if !ok {
		return errors.New("Variable __dict_manager not present in current execution scope")
	}
	dict_ptr, err := ids.GetRelocatable("dict_ptr", vm)
	if err != nil {
		return err
	}
	key, err := ids.Get("key", vm)
	if err != nil {
		return err
	}
	// Hint Logic
	tracker, err := dictManager.GetTracker(dict_ptr)
	if err != nil {
		return err
	}
	tracker.CurrentPtr.Offset += DICT_ACCESS_SIZE
	val, err := tracker.GetValue(key)
	if err != nil {
		return err
	}
	return ids.Insert("value", val, vm)
}

func dictWrite(ids IdsManager, scopes *ExecutionScopes, vm *VirtualMachine) error {
	// Extract Variables
	dictManager, ok := FetchDictManager(scopes)
	if !ok {
		return errors.New("Variable __dict_manager not present in current execution scope")
	}
	dict_ptr, err := ids.GetRelocatable("dict_ptr", vm)
	if err != nil {
		return err
	}
	key, err := ids.Get("key", vm)
	if err != nil {
		return err
	}
	new_value, err := ids.Get("new_value", vm)
	if err != nil {
		return err
	}
	/* dict_ptr has type *DictAccess
	struct DictAccess {
		key: felt,
		prev_value: felt,
		new_value: felt,
	}
	so ids.dict_ptr.prev_value = [dict_ptr + 1]
	*/
	prev_val_addr := dict_ptr.AddUint(1)
	// Hint Logic
	tracker, err := dictManager.GetTracker(dict_ptr)
	if err != nil {
		return err
	}
	tracker.CurrentPtr.Offset += DICT_ACCESS_SIZE
	prev_val, err := tracker.GetValue(key)
	if err != nil {
		return err
	}
	tracker.InsertValue(key, new_value)
	return vm.Segments.Memory.Insert(prev_val_addr, prev_val)
}
