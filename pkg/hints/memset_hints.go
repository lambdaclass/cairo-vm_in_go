package hints

import (
	. "github.com/lambdaclass/cairo-vm.go/pkg/hints/hint_utils"
	. "github.com/lambdaclass/cairo-vm.go/pkg/types"
	. "github.com/lambdaclass/cairo-vm.go/pkg/vm"
)

// Implements hint:
// %{ vm_enter_scope({'n': ids.n}) %}
func memset_enter_scope(ids IdsManager, vm *VirtualMachine, execScopes *ExecutionScopes) error {
	n, err := ids.GetFelt("n", vm)
	if err != nil {
		return err
	}
	execScopes.EnterScope(map[string]interface{}{"n": n})
	return nil
}
