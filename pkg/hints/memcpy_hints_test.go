package hints_test

import (
	"testing"

	. "github.com/lambdaclass/cairo-vm.go/pkg/hints"
	. "github.com/lambdaclass/cairo-vm.go/pkg/hints/hint_utils"
	. "github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
	. "github.com/lambdaclass/cairo-vm.go/pkg/types"
	. "github.com/lambdaclass/cairo-vm.go/pkg/vm"
	. "github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
)

func AddSegmentHintOk(t *testing.T) {
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

func TestExitScopeValid(t *testing.T) {
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
		t.Errorf("TestExitScopeValid failed with error %s", err)
	}

}

func TestExitScopeInvalid(t *testing.T) {
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
		t.Errorf("TestExitScopeInvalid should fail with error %s", ErrCannotExitMainScop)
	}

}
