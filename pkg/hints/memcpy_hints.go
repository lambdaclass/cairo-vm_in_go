package hints

import (
	. "github.com/lambdaclass/cairo-vm.go/pkg/hints/hint_utils"
	"github.com/lambdaclass/cairo-vm.go/pkg/types"
	. "github.com/lambdaclass/cairo-vm.go/pkg/vm"
	. "github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
)

// Implements hint: memory[ap] = segments.add()
func add_segment(vm *VirtualMachine) error {
	new_segment_base := vm.Segments.AddSegment()
	return vm.Segments.Memory.Insert(vm.RunContext.Ap, NewMaybeRelocatableRelocatable(new_segment_base))
}

// Implements hint:
// %{ vm_exit_scope() %}
func vm_exit_scope(executionScopes *types.ExecutionScopes) error {
	return executionScopes.ExitScope()
}

// Implements hint:
// %{ vm_enter_scope({'n': ids.len}) %}
func memcpy_enter_scope(ids IdsManager, vm *VirtualMachine, execScopes *types.ExecutionScopes) error {
	len, err := ids.GetFelt("len", vm)
	if err != nil {
		return err
	}
	scope := map[string]interface{}{"n": len}
	execScopes.EnterScope(scope)
	return nil
}

// Implements hint: vm_enter_scope()
func vm_enter_scope(executionScopes *types.ExecutionScopes) error {
	executionScopes.EnterScope(make(map[string]interface{}))
	return nil
}
