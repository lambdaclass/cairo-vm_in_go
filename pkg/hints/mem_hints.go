package hints

import "github.com/lambdaclass/cairo-vm.go/pkg/types"

// Implements hint:
// %{ vm_exit_scope() %}
func vm_exit_scope(executionScopes *types.ExecutionScopes) error {
	err := executionScopes.ExitScope()
	if err != nil {
		return err
	}
	return nil
}
