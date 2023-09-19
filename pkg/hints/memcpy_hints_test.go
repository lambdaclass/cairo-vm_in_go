package hints_test

import (
	"testing"

	. "github.com/lambdaclass/cairo-vm.go/pkg/hints"
	. "github.com/lambdaclass/cairo-vm.go/pkg/hints/hint_utils"
	. "github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
	. "github.com/lambdaclass/cairo-vm.go/pkg/types"
	. "github.com/lambdaclass/cairo-vm.go/pkg/utils"
	. "github.com/lambdaclass/cairo-vm.go/pkg/vm"
	. "github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
)

func TestAddSegmentHintOk(t *testing.T) {
	vm := NewVirtualMachine()
	initial_segments := vm.Segments.Memory.NumSegments()
	hintProcessor := CairoVmHintProcessor{}
	hintData := any(HintData{
		Code: ADD_SEGMENT,
	})
	err := hintProcessor.ExecuteHint(vm, &hintData, nil, nil)
	if err != nil {
		t.Errorf("ADD_SEGMENT hint test failed with error %s", err)
	}
	if initial_segments+1 != vm.Segments.Memory.NumSegments() {
		t.Errorf("ADD_SEGMENT fail expected: %d segments, got: %d", initial_segments+1, vm.Segments.Memory.NumSegments())
	}
}

func TestExitScopeHintValid(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()
	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{
			"a": {NewMaybeRelocatableFelt(FeltFromUint64(17))},
		},
		vm,
	)
	hintProcessor := CairoVmHintProcessor{}
	hintData := any(HintData{
		Ids:  idsManager,
		Code: VM_EXIT_SCOPE,
	})

	executionScopes := NewExecutionScopes()
	scope := make(map[string]interface{})
	scope["a"] = FeltOne()
	executionScopes.EnterScope(scope)

	err := hintProcessor.ExecuteHint(vm, &hintData, nil, executionScopes)
	if err != nil {
		t.Errorf("TestExitScopeHintValid failed with error %s", err)
	}

}

func TestExitScopeHintInvalid(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()
	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{
			"a": {NewMaybeRelocatableFelt(FeltFromUint64(17))},
		},
		vm,
	)
	hintProcessor := CairoVmHintProcessor{}
	hintData := any(HintData{
		Ids:  idsManager,
		Code: VM_EXIT_SCOPE,
	})

	executionScopes := NewExecutionScopes()
	scope := make(map[string]interface{})
	scope["a"] = FeltOne()

	err := hintProcessor.ExecuteHint(vm, &hintData, nil, executionScopes)
	if err.Error() != ErrCannotExitMainScop.Error() {
		t.Errorf("TestExitScopeHintInvalid should fail with error %s", ErrCannotExitMainScop)
	}

}

func TestMemcpyEnterScopeHintValid(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()
	vm.Segments.AddSegment()
	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{
			"len": {NewMaybeRelocatableFelt(FeltFromUint64(45))},
		},
		vm,
	)
	hintProcessor := CairoVmHintProcessor{}
	hintData := any(HintData{
		Ids:  idsManager,
		Code: MEMCPY_ENTER_SCOPE,
	})

	executionScopes := NewExecutionScopes()
	err := hintProcessor.ExecuteHint(vm, &hintData, nil, executionScopes)
	if err != nil {
		t.Errorf("TestMemcpyEnterScopeHintValid failed with error %s", err)
	}
	res, err := executionScopes.Get("n")
	if err != nil {
		t.Errorf("TestMemcpyEnterScopeHintValid failed with error %s", err)
	}
	if res.(Felt) != FeltFromDecString("45") {
		t.Errorf("TestMemcpyEnterScopeHintValid failed, expected len: %d, got: %d", FeltFromDecString("45"), res.(Felt))
	}
}

func TestMemcpyEnterScopeHintInvalid(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()
	vm.Segments.AddSegment()

	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{},
		vm,
	)
	hintProcessor := CairoVmHintProcessor{}
	hintData := any(HintData{
		Ids:  idsManager,
		Code: MEMCPY_ENTER_SCOPE,
	})

	executionScopes := NewExecutionScopes()
	err := hintProcessor.ExecuteHint(vm, &hintData, nil, executionScopes)
	if err.Error() != ErrUnknownIdentifier("len").Error() {
		t.Errorf("TestMemcpyEnterScopeHintInvalid should fail with error %s", ErrUnknownIdentifier("len"))
	}
}

func TestEnterScope(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()
	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{},
		vm,
	)
	hintProcessor := CairoVmHintProcessor{}
	hintData := any(HintData{
		Ids:  idsManager,
		Code: VM_ENTER_SCOPE,
	})

	executionScopes := NewExecutionScopes()
	scope := make(map[string]interface{})
	scope["a"] = FeltOne()

	err := hintProcessor.ExecuteHint(vm, &hintData, nil, executionScopes)
	if err != nil {
		t.Errorf("TestEnterScopeHint failed with error %s", err)
	}
}

func TestMemcpyContinueCopyingValid(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()
	vm.Segments.AddSegment()
	vm.Segments.AddSegment()
	vm.RunContext.Fp = NewRelocatable(2, 0)
	vm.Segments.Memory.Data[NewRelocatable(1, 2)] = *NewMaybeRelocatableFelt(FeltFromUint64(5))

	executionScopes := NewExecutionScopes()
	scope := make(map[string]interface{})
	scope["n"] = FeltOne()
	executionScopes.EnterScope(scope)

	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{
			"continue_copying": nil,
		},
		vm,
	)
	hintProcessor := CairoVmHintProcessor{}
	hintData := any(HintData{
		Ids:  idsManager,
		Code: MEMCPY_CONTINUE_COPYING,
	})

	err := hintProcessor.ExecuteHint(vm, &hintData, nil, executionScopes)
	if err != nil {
		t.Errorf("TestMemsetContinueLoopValidEqual1 failed with error %s", err)
	}
}

func TestMemcpyContinueCopyingVarNotInScope(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()
	vm.RunContext.Fp = NewRelocatable(3, 0)
	vm.Segments.Memory.Insert(NewRelocatable(0, 2), NewMaybeRelocatableFelt(FeltFromUint64(5)))

	executionScopes := NewExecutionScopes()

	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{
			"continue_copying": nil,
		},
		vm,
	)
	hintProcessor := CairoVmHintProcessor{}
	hintData := any(HintData{
		Ids:  idsManager,
		Code: MEMCPY_CONTINUE_COPYING,
	})

	err := hintProcessor.ExecuteHint(vm, &hintData, nil, executionScopes)
	if err.Error() != ErrVariableNotInScope("n").Error() {
		t.Errorf("TestMemcpyContinueCopyingVarNotInScope should fail with error %s", ErrVariableNotInScope("n"))
	}
}

func TestMemcpyContinueCopyingInsertError(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()
	vm.Segments.AddSegment()
	executionScopes := NewExecutionScopes()

	scope := make(map[string]interface{})
	scope["n"] = FeltOne()
	executionScopes.EnterScope(scope)

	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{
			"continue_copying": {NewMaybeRelocatableFelt(FeltFromUint64(5))},
		},
		vm,
	)
	hintProcessor := CairoVmHintProcessor{}
	hintData := any(HintData{
		Ids:  idsManager,
		Code: MEMCPY_CONTINUE_COPYING,
	})

	err := hintProcessor.ExecuteHint(vm, &hintData, nil, executionScopes)
	expected := ErrMemoryWriteOnce(NewRelocatable(0, 0), *NewMaybeRelocatableFeltFromUint64(5), *NewMaybeRelocatableFeltFromUint64(0))
	if err.Error() != expected.Error() {
		t.Errorf("TestMemcpyContinueCopyingInsertError should fail with error %s", expected)
	}
}
