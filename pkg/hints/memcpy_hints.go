package hints

import (
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
