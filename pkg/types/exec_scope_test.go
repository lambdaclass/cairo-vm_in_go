package types_test

import (
	"reflect"
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

func TestGetRefLocalVariables(t *testing.T) {
	scope := make(map[string]interface{})
	scope["k"] = lambdaworks.FeltOne()

	scopes := types.NewExecutionScopes()
	scopes.EnterScope(scope)

	result, err := scopes.GetRef("k")
	if err != nil {
		t.Errorf("TestGetRefLocalVariables failed with error: %s", err)

	}
	f_one := lambdaworks.FeltOne()
	f_res := (*result).(lambdaworks.Felt)
	expected := &f_one
	if *expected != f_res {
		t.Errorf("TestGetRefLocalVariables failed, expected: %d, got: %d", expected, &f_res)
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
	expected_err := types.ErrVariableNotInScope("a")
	if err.Error() != expected_err.Error() {
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

func TestExitScope(t *testing.T) {
	scope := make(map[string]interface{})
	scope["a"] = lambdaworks.FeltOne()

	scopes := types.NewExecutionScopes()
	scopes.EnterScope(scope)

	locals, err := scopes.GetLocalVariables()
	if err != nil {
		t.Errorf("TestExitScopeTest failed with error: %s", err)
	}

	if len(locals) != 1 {
		t.Errorf("TestExitScopeTest failed, expected length: %d, got: %d", 1, len(locals))
	}

	result, err := scopes.Get("a")
	if err != nil {
		t.Errorf("TestExitScopeTest failed with error: %s", err)

	}

	f_res := result.(lambdaworks.Felt)
	expected := lambdaworks.FeltOne()
	if expected != f_res {
		t.Errorf("TestExitScopeTest failed, expected: %s, got: %s", expected.ToSignedFeltString(), f_res.ToSignedFeltString())
	}

	err = scopes.ExitScope()
	if err != nil {
		t.Errorf("TestExitScopeTest failed with error: %s", err)
	}

	locals, err = scopes.GetLocalVariables()
	if err != nil {
		t.Errorf("TestExitScopeTest failed with error: %s", err)
	}

	if len(locals) != 0 {
		t.Errorf("TestExitScopeTest failed, expected length: %d, got: %d", 0, len(locals))
	}

}

func TestAssignLocalVariable(t *testing.T) {
	scope := make(map[string]interface{})

	scopes := types.NewExecutionScopes()
	scopes.EnterScope(scope)

	scopes.AssignOrUpdateVariable("a", uint64(45))

	locals, err := scopes.GetLocalVariables()
	if err != nil {
		t.Errorf("TestAssignLocalVariable failed with error: %s", err)
	}

	if len(locals) != 1 {
		t.Errorf("TestAssignLocalVariable failed, expected length: %d, got: %d", 1, len(locals))
	}

	result, err := scopes.Get("a")
	if err != nil {
		t.Errorf("TestAssignLocalVariable failed with error: %s", err)

	}

	f_res := result.(uint64)
	expected := uint64(45)
	if expected != f_res {
		t.Errorf("TestAssignLocalVariable failed, expected: uint64(%d), got: %d", expected, f_res)
	}

}

func TestReAssignLocalVariable(t *testing.T) {
	scope := make(map[string]interface{})

	scopes := types.NewExecutionScopes()
	scopes.EnterScope(scope)

	scopes.AssignOrUpdateVariable("a", uint64(45))

	locals, err := scopes.GetLocalVariables()
	if err != nil {
		t.Errorf("TestReAssignLocalVariable failed with error: %s", err)
	}

	if len(locals) != 1 {
		t.Errorf("TestReAssignLocalVariable failed, expected length: %d, got: %d", 1, len(locals))
	}

	result, err := scopes.Get("a")
	if err != nil {
		t.Errorf("TestReAssignLocalVariable failed with error: %s", err)

	}

	res := result.(uint64)
	expected := uint64(45)
	if expected != res {
		t.Errorf("TestReAssignLocalVariable failed, expected: uint64(%d), got: %d", expected, res)
	}

	scopes.AssignOrUpdateVariable("a", lambdaworks.FeltOne())

	locals, err = scopes.GetLocalVariables()
	if err != nil {
		t.Errorf("TestReAssignLocalVariable failed with error: %s", err)
	}

	if len(locals) != 1 {
		t.Errorf("TestReAssignLocalVariable failed, expected length: %d, got: %d", 1, len(locals))
	}

	result, err = scopes.Get("a")
	if err != nil {
		t.Errorf("TestReAssignLocalVariable failed with error: %s", err)

	}

	f_res := result.(lambdaworks.Felt)
	f_expected := lambdaworks.FeltOne()
	if f_expected != f_res {
		t.Errorf("TestReAssignLocalVariable failed, expected: %s, got: %s", f_expected.ToSignedFeltString(), f_res.ToSignedFeltString())
	}

}

func TestDeleteLocalVariable(t *testing.T) {
	scope := make(map[string]interface{})

	scopes := types.NewExecutionScopes()
	scopes.EnterScope(scope)

	scopes.AssignOrUpdateVariable("a", "val")

	locals, err := scopes.GetLocalVariables()
	if err != nil {
		t.Errorf("DeleteLocalVariable failed with error: %s", err)
	}

	if len(locals) != 1 {
		t.Errorf("DeleteLocalVariable failed, expected length: %d, got: %d", 1, len(locals))
	}

	_, err = scopes.Get("a")
	if err != nil {
		t.Errorf("DeleteLocalVariable failed with error: %s", err)

	}

	scopes.DeleteVariable("a")

	// check variable a can't be accessed now
	_, err = scopes.Get("a")
	expected := types.ErrVariableNotInScope("a")
	if err.Error() != expected.Error() {
		t.Errorf("TestDeleteLocalVariable should fail with error: %s", types.ErrVariableNotInScope("a").Error())
	}

}
func TestErrExitMainScope(t *testing.T) {
	scopes := types.NewExecutionScopes()

	err := scopes.ExitScope()
	if err != types.ErrCannotExitMainScop {
		t.Errorf("TestErrExitMainScope should fail with error: %s and fails with: %s", types.ErrCannotExitMainScop, err)
	}
}

func TestGetListU64(t *testing.T) {
	u64_list := []uint64{1, 2, 3, 4, 5}
	scopes := types.NewExecutionScopes()
	scopes.AssignOrUpdateVariable("u64_list", u64_list)

	result, err := scopes.GetList(reflect.TypeOf([]uint64{}), "u64_list")
	if err != nil {
		t.Errorf("TestGetListU64 failed with error: %s", err)
	}

	if reflect.TypeOf(result) != reflect.TypeOf(u64_list) {
		t.Errorf("TestGetListU64 failed with error: Expected list has type %s, got: %s", reflect.TypeOf(u64_list).String(), reflect.TypeOf(result).String())

	}

	if !reflect.DeepEqual(u64_list, result) {
		t.Errorf("TestGetListU64 failed with error: Expected list %v, got: %v", u64_list, result)

	}

}

func TestErrGetList(t *testing.T) {
	u64_list := []uint64{1, 2, 3, 4, 5}
	scopes := types.NewExecutionScopes()
	scopes.AssignOrUpdateVariable("u64_list", u64_list)

	// Correct code
	_, err := scopes.GetList(reflect.TypeOf([]int{}), "u64_list")
	expected := types.ErrVariableNotInScope("u64_list")
	if err.Error() != expected.Error() {
		t.Errorf("TestErrGetList should fail with error: %s\n, got: %s\n", expected, err)
	}

	// Extra code - can be deleated in production
	// If uncommented, comment 'correct code' here and in method GetList
	// _, err := scopes.GetList(reflect.TypeOf([]int{}), "u64_list")
	// expected := types.ErrListTypeNotEqual("u64_list", reflect.TypeOf([]int{}).String(), reflect.TypeOf(u64_list).String())
	// if err.Error() != expected.Error() {
	// 	t.Errorf("TestErrGetList should fail with error: %s\n, got: %s\n", expected, err)
	// }
}
