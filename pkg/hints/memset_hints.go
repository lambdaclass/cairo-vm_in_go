package hints

import (
	. "github.com/lambdaclass/cairo-vm.go/pkg/hints/hint_utils"
	"github.com/lambdaclass/cairo-vm.go/pkg/types"
	. "github.com/lambdaclass/cairo-vm.go/pkg/vm"
	. "github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
)

// Implements hint:
// %{ vm_enter_scope({'n': ids.n}) %}
func memset_enter_scope(ids IdsManager, vm *VirtualMachine, execScopes *types.ExecutionScopes) error {
	n, err := ids.Get("n")
	if err != nil {
		return err
	}
	execScopes.EnterScope(map[string]interface{}{"n": n})
	return nil
}
