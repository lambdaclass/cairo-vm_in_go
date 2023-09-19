package hints_test

import (
	"testing"

	. "github.com/lambdaclass/cairo-vm.go/pkg/hints"
	. "github.com/lambdaclass/cairo-vm.go/pkg/hints/hint_utils"
	. "github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"

	. "github.com/lambdaclass/cairo-vm.go/pkg/types"

	. "github.com/lambdaclass/cairo-vm.go/pkg/vm"

	. "github.com/lambdaclass/cairo-vm.go/pkg/utils"

	. "github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
)

func TestMemsetEnterScopeValid(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments = AddNSegments(vm.Segments, 2)
	vm.RunContext.Fp = NewRelocatable(1, 2)
	vm.Segments.Memory.Insert(NewRelocatable(1, 1), NewMaybeRelocatableFeltFromUint64(5))
	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{
			"n": {NewMaybeRelocatableFeltFromUint64(4)},
		},
		vm,
	)
	hintProcessor := CairoVmHintProcessor{}
	hintData := any(HintData{
		Ids:  idsManager,
		Code: MEMSET_ENTER_SCOPE,
	})

	executionScopes := NewExecutionScopes()
	err := hintProcessor.ExecuteHint(vm, &hintData, nil, executionScopes)
	if err != nil {
		t.Errorf("failed with error %s", err)
	}
	val, err := executionScopes.Get("n")
	if err != nil {
		t.Errorf("failed with error %s", err)
	}

	if val.(Felt) != FeltFromUint64(4) {
		t.Errorf("failed, expected value: %s, got: %s", FeltFromUint64(4).ToSignedFeltString(), val.(Felt).ToSignedFeltString())

	}
}

func TestMemsetEnterScopeInvalid(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments = AddNSegments(vm.Segments, 2)
	vm.RunContext.Fp = NewRelocatable(1, 1)
	// insert a relocatable value in the address of ids.len so that it raises an error.
	vm.Segments.Memory.Insert(NewRelocatable(1, 1), NewMaybeRelocatableRelocatableParams(1, 0))
	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{
			"n": {NewMaybeRelocatableRelocatableParams(3, 4)},
		},
		vm,
	)
	hintProcessor := CairoVmHintProcessor{}
	hintData := any(HintData{
		Ids:  idsManager,
		Code: MEMSET_ENTER_SCOPE,
	})

	executionScopes := NewExecutionScopes()
	err := hintProcessor.ExecuteHint(vm, &hintData, nil, executionScopes)
	if err.Error() != ErrIdentifierNotFelt("n").Error() {
		t.Errorf("should fail with: %s", ErrIdentifierNotFelt("n").Error())
	}

}
