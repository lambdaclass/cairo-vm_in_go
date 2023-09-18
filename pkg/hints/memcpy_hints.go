package hints

import (
	. "github.com/lambdaclass/cairo-vm.go/pkg/hints/hint_utils"
	"github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
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

/*
	Implements hint:

	%{
	    n -= 1
	    ids.`i_name` = 1 if n > 0 else 0

%}
*/
func memset_step_loop(ids IdsManager, vm *VirtualMachine, execScoes *types.ExecutionScopes, i_name string) error {
	// get `n` variable from vm scope
	n, err := execScoes.GetRef("n")
	if err != nil {
		return err
	}
	// this variable will hold the value of `n - 1`
	*n = (*n).(lambdaworks.Felt).Sub(lambdaworks.FeltOne())
	// if `new_n` is positive, insert 1 in the address of `continue_loop`
	// else, insert 0
	var flag *MaybeRelocatable
	if (*n).(lambdaworks.Felt).IsPositive() {
		flag = NewMaybeRelocatableFelt(lambdaworks.FeltOne())
	} else {
		flag = NewMaybeRelocatableFelt(lambdaworks.FeltZero())
	}
	addr, err := ids.GetAddr(i_name, vm)
	if err != nil {
		return err
	}
	err = vm.Segments.Memory.Insert(addr, flag)
	if err != nil {
		return err
	}
	return nil
}

// Implements hint: vm_enter_scope()
func vm_enter_scope(executionScopes *types.ExecutionScopes) error {
	executionScopes.EnterScope(make(map[string]interface{}))
	return nil
}
