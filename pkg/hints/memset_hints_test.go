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

func TestMemsetContinueLoopValidEqual1Hint(t *testing.T) {
	vm := NewVirtualMachine()
	vm.RunContext.Fp = NewRelocatable(1, 0)
	vm.Segments = AddNSegments(vm.Segments, 2)
	executionScopes := NewExecutionScopesWithInitValue("n", FeltFromUint64(1))

	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{
			"continue_loop": {NewMaybeRelocatableFeltFromUint64(0)},
		},
		vm,
	)
	hintProcessor := CairoVmHintProcessor{}
	hintData := any(HintData{
		Ids:  idsManager,
		Code: MEMSET_CONTINUE_LOOP,
	})

	err := hintProcessor.ExecuteHint(vm, &hintData, nil, executionScopes)
	if err != nil {
		t.Errorf("failed with error %s", err)
	}
	val, err := vm.Segments.Memory.GetFelt(NewRelocatable(1, 0))
	if err != nil {
		t.Errorf("failed with error %s", err)
	}
	if val != FeltZero() {
		t.Errorf("failed, expected %d, got: %d", FeltZero(), val)
	}
}

func TestMemsetContinueLoopValidEqual5Hint(t *testing.T) {
	vm := NewVirtualMachine()
	vm.RunContext.Fp = NewRelocatable(1, 0)
	vm.Segments = AddNSegments(vm.Segments, 2)
	executionScopes := NewExecutionScopesWithInitValue("n", FeltFromUint64(5))

	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{
			"continue_loop": {NewMaybeRelocatableFeltFromUint64(1)},
		},
		vm,
	)
	hintProcessor := CairoVmHintProcessor{}
	hintData := any(HintData{
		Ids:  idsManager,
		Code: MEMSET_CONTINUE_LOOP,
	})

	err := hintProcessor.ExecuteHint(vm, &hintData, nil, executionScopes)
	if err != nil {
		t.Errorf("failed with error %s", err)
	}
	val, err := vm.Segments.Memory.GetFelt(NewRelocatable(1, 0))
	if err != nil {
		t.Errorf("failed with error %s", err)
	}
	if val != FeltOne() {
		t.Errorf("failed, expected %d, got: %d", FeltOne(), val)
	}
}

func TestMemsetContinueLoopVarNotInScope(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()
	executionScopes := NewExecutionScopes()
	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{
			"continue_loop": {NewMaybeRelocatableFeltFromUint64(1)},
		},
		vm,
	)
	hintProcessor := CairoVmHintProcessor{}
	hintData := any(HintData{
		Ids:  idsManager,
		Code: MEMSET_CONTINUE_LOOP,
	})

	err := hintProcessor.ExecuteHint(vm, &hintData, nil, executionScopes)
	if err.Error() != ErrVariableNotInScope("n").Error() {
		t.Errorf("should fail with error %s", ErrVariableNotInScope("n"))
	}
}

func TestMemsetContinueLoopInsertError(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments = AddNSegments(vm.Segments, 2)
	executionScopes := NewExecutionScopesWithInitValue("n", FeltOne())
	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{
			"continue_loop": {NewMaybeRelocatableFelt(FeltFromUint64(5))},
		},
		vm,
	)
	hintProcessor := CairoVmHintProcessor{}
	hintData := any(HintData{
		Ids:  idsManager,
		Code: MEMSET_CONTINUE_LOOP,
	})

	err := hintProcessor.ExecuteHint(vm, &hintData, nil, executionScopes)
	expected := ErrMemoryWriteOnce(NewRelocatable(0, 0), *NewMaybeRelocatableFeltFromUint64(5), *NewMaybeRelocatableFeltFromUint64(0))
	if err.Error() != expected.Error() {
		t.Errorf("should fail with error %s", expected)
	}
}
