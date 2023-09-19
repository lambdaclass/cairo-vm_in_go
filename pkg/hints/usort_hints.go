package hints

import (
	"github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
	"github.com/lambdaclass/cairo-vm.go/pkg/types"
	"github.com/pkg/errors"
)

// Implements hint:
// %{ vm_enter_scope(dict(__usort_max_size = globals().get('__usort_max_size'))) %}
func usort_enter_scope(executionScopes *types.ExecutionScopes) error {
	usort_max_size, err := executionScopes.Get("usort_max_size")

	usort_max_size_felt, cast_ok := usort_max_size.(lambdaworks.Felt)

	if !cast_ok {
		return errors.New("Error casting usort_max_size into a Felt")
	}

	if err != nil {
		return err
	}

	if usort_max_size == nil {
		executionScopes.EnterScope(make(map[string]interface{}))
	}

	scope := make(map[string]interface{})
	scope["usort_max_size"] = usort_max_size_felt
	executionScopes.EnterScope(scope)

	return nil
}
