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

	result, err := scopes.Get("k")
	if err != nil {
		t.Errorf("TestGetLocalVariables failed with error: %s", err)

	}
	f_res := result.(lambdaworks.Felt)
	expected := lambdaworks.FeltOne()
	if expected != f_res {
		t.Errorf("TestGetLocalVariables failed, expected: %s, got: %s", expected.ToSignedFeltString(), f_res.ToSignedFeltString())
	}
}

func TestEnterNewScope(t *testing.T) {
	scope := make(map[string]interface{})
	scope["a"] = lambdaworks.FeltOne()

	scopes := types.NewExecutionScopes()
	scopes.EnterScope(scope)

	locals, err := scopes.GetLocalVariables()
	if err != nil {
		t.Errorf("TestEnterNewScope failed with error: %s", err)
	}

	if len(locals) != 1 {
		t.Errorf("TestEnterNewScope failed, expected length: %d, got: %d", 1, len(locals))
	}

	result, err := scopes.Get("a")
	if err != nil {
		t.Errorf("TestEnterNewScope failed with error: %s", err)

	}

	f_res := result.(lambdaworks.Felt)
	expected := lambdaworks.FeltOne()
	if expected != f_res {
		t.Errorf("TestEnterNewScope failed, expected: %s, got: %s", expected.ToSignedFeltString(), f_res.ToSignedFeltString())
	}

	snd_scope := make(map[string]interface{})
	snd_scope["b"] = lambdaworks.FeltZero()
	scopes.EnterScope(snd_scope)

	locals, err = scopes.GetLocalVariables()
	if err != nil {
		t.Errorf("TestEnterNewScope failed with error: %s", err)
	}

	if len(locals) != 1 {
		t.Errorf("TestEnterNewScope failed, expected length: %d, got: %d", 1, len(locals))
	}

	// check variable a can't be accessed now
	_, err = scopes.Get("a")
	if err.Error() != types.ErrVariableNotInScope("a").Error() {
		t.Errorf("TestEnterNewScope should fail with error: %s", types.ErrVariableNotInScope("a").Error())

	}

	result, err = scopes.Get("b")
	if err != nil {
		t.Errorf("TestEnterNewScope failed with error: %s", err)
	}

	f_res = result.(lambdaworks.Felt)
	expected = lambdaworks.FeltZero()
	if expected != f_res {
		t.Errorf("TestEnterNewScope failed, expected: %s, got: %s", expected.ToSignedFeltString(), f_res.ToSignedFeltString())
	}

}

func TestExitScopeTest(t *testing.T) {

}
