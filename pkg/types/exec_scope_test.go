package types_test

import (
	"testing"

	"github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
	"github.com/lambdaclass/cairo-vm.go/pkg/types"
)

func TestInitializeExecutionScopes(t *testing.T) {
	scopes := types.NewExecutionScopes()
	if len(scopes.Data()) != 1 {
		t.Errorf("TestInitializeExecutionScopes failed, expected length: %d, got: %d", 1, len((scopes.Data())))
	}
}

func TestGetLocalVariables(t *testing.T) {
	scope := make(map[string]interface{})
	scope["k"] = lambdaworks.FeltOne()

	scopes := types.NewExecutionScopes()
	scopes.EnterScope(scope)

	result, err := scopes.Get("key")
	if err != nil {
		t.Errorf("TestGetLocalVariables failed with error: %s", err)

	}
	f_res := result.(lambdaworks.Felt)
	expected := lambdaworks.FeltOne()
	if expected != f_res {
		t.Errorf("TestGetLocalVariables failed, expected: %s, got: %s", expected.ToSignedFeltString(), f_res.ToSignedFeltString())
	}

}
