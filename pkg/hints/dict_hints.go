package hints

import (
	. "github.com/lambdaclass/cairo-vm.go/pkg/hints/dict_manager"
	. "github.com/lambdaclass/cairo-vm.go/pkg/hints/hint_utils"
	. "github.com/lambdaclass/cairo-vm.go/pkg/types"
	. "github.com/lambdaclass/cairo-vm.go/pkg/vm"
	"github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
)

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
