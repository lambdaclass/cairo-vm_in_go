package hints_test

import (
	"math/big"
	"testing"

	. "github.com/lambdaclass/cairo-vm.go/pkg/hints"
	. "github.com/lambdaclass/cairo-vm.go/pkg/hints/hint_codes"
	. "github.com/lambdaclass/cairo-vm.go/pkg/hints/hint_utils"
	. "github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
	. "github.com/lambdaclass/cairo-vm.go/pkg/types"
	. "github.com/lambdaclass/cairo-vm.go/pkg/vm"
	. "github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
)

func TestDivModNPackedDivMod(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()
	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{
			"a": {
				NewMaybeRelocatableFelt(FeltFromUint64(10)),
				NewMaybeRelocatableFelt(FeltFromUint64(0)),
				NewMaybeRelocatableFelt(FeltFromUint64(0)),
			},
			"b": {
				NewMaybeRelocatableFelt(FeltFromUint64(2)),
				NewMaybeRelocatableFelt(FeltFromUint64(0)),
				NewMaybeRelocatableFelt(FeltFromUint64(0)),
			},
		},
		vm,
	)
	hintProcessor := CairoVmHintProcessor{}
	hintData := any(HintData{
		Ids:  idsManager,
		Code: DIV_MOD_N_PACKED_DIVMOD_V1,
	})
	scopes := NewExecutionScopes()
	err := hintProcessor.ExecuteHint(vm, &hintData, nil, scopes)
	if err != nil {
		t.Errorf("DIV_MOD_N_PACKED_DIVMOD_V1 hint test failed with error %s", err)
	}
	// Check result in scope
	expectedRes := big.NewInt(5)

	res, err := FetchScopeVar[big.Int]("res", scopes)
	if err != nil || res.Cmp(expectedRes) != 0 {
		t.Error("Wrong/No scope value res")
	}

	val, err := FetchScopeVar[big.Int]("val", scopes)
	if err != nil || val.Cmp(expectedRes) != 0 {
		t.Error("Wrong/No scope value val")
	}
}

func TestDivModNPackedDivModExternalN(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()
	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{
			"a": {
				NewMaybeRelocatableFelt(FeltFromUint64(20)),
				NewMaybeRelocatableFelt(FeltFromUint64(0)),
				NewMaybeRelocatableFelt(FeltFromUint64(0)),
			},
			"b": {
				NewMaybeRelocatableFelt(FeltFromUint64(2)),
				NewMaybeRelocatableFelt(FeltFromUint64(0)),
				NewMaybeRelocatableFelt(FeltFromUint64(0)),
			},
		},
		vm,
	)
	hintProcessor := CairoVmHintProcessor{}
	hintData := any(HintData{
		Ids:  idsManager,
		Code: DIV_MOD_N_PACKED_DIVMOD_EXTERNAL_N,
	})
	scopes := NewExecutionScopes()
	scopes.AssignOrUpdateVariable("N", *big.NewInt(7))
	err := hintProcessor.ExecuteHint(vm, &hintData, nil, scopes)
	if err != nil {
		t.Errorf("DIV_MOD_N_PACKED_DIVMOD_EXTERNAL_N hint test failed with error %s", err)
	}
	// Check result in scope
	expectedRes := big.NewInt(3)

	res, err := FetchScopeVar[big.Int]("res", scopes)
	if err != nil || res.Cmp(expectedRes) != 0 {
		t.Error("Wrong/No scope value res")
	}

	val, err := FetchScopeVar[big.Int]("val", scopes)
	if err != nil || val.Cmp(expectedRes) != 0 {
		t.Error("Wrong/No scope value val")
	}
}

func TestDivModSafeDivOk(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()
	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{},
		vm,
	)
	hintProcessor := CairoVmHintProcessor{}
	hintData := any(HintData{
		Ids:  idsManager,
		Code: DIV_MOD_N_SAFE_DIV,
	})
	scopes := NewExecutionScopes()
	scopes.AssignOrUpdateVariable("N", *big.NewInt(5))
	scopes.AssignOrUpdateVariable("a", *big.NewInt(10))
	scopes.AssignOrUpdateVariable("b", *big.NewInt(30))
	scopes.AssignOrUpdateVariable("res", *big.NewInt(2))
	err := hintProcessor.ExecuteHint(vm, &hintData, nil, scopes)
	if err != nil {
		t.Errorf("DIV_MOD_N_SAFE_DIV hint test failed with error %s", err)
	}
	// Check result in scope
	expectedValue := big.NewInt(10) // (2 * 30 - 10) / 5 = 10

	val, err := FetchScopeVar[big.Int]("value", scopes)
	if err != nil || val.Cmp(expectedValue) != 0 {
		t.Error("Wrong/No scope value val")
	}
}

func TestDivModSafeDivPlusOneOk(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()
	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{},
		vm,
	)
	hintProcessor := CairoVmHintProcessor{}
	hintData := any(HintData{
		Ids:  idsManager,
		Code: DIV_MOD_N_SAFE_DIV_PLUS_ONE,
	})
	scopes := NewExecutionScopes()
	scopes.AssignOrUpdateVariable("N", *big.NewInt(5))
	scopes.AssignOrUpdateVariable("a", *big.NewInt(10))
	scopes.AssignOrUpdateVariable("b", *big.NewInt(30))
	scopes.AssignOrUpdateVariable("res", *big.NewInt(2))
	err := hintProcessor.ExecuteHint(vm, &hintData, nil, scopes)
	if err != nil {
		t.Errorf("DIV_MOD_N_SAFE_DIV_PLUS_ONE hint test failed with error %s", err)
	}
	// Check result in scope
	expectedValue := big.NewInt(11) // (2 * 30 - 10) / 5 + 1 = 11

	val, err := FetchScopeVar[big.Int]("value", scopes)
	if err != nil || val.Cmp(expectedValue) != 0 {
		t.Error("Wrong/No scope value val")
	}
}
