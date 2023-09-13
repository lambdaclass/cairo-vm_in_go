package hints_test

import (
	"testing"

	. "github.com/lambdaclass/cairo-vm.go/pkg/hints"
	. "github.com/lambdaclass/cairo-vm.go/pkg/hints/hint_utils"
	. "github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
	"github.com/lambdaclass/cairo-vm.go/pkg/types"
	. "github.com/lambdaclass/cairo-vm.go/pkg/vm"
	. "github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
)

func TestDefaultDictNewCreateManager(t *testing.T) {
	vm := NewVirtualMachine()
	scopes := types.NewExecutionScopes()
	vm.Segments.AddSegment()
	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{
			"default_value": {NewMaybeRelocatableFelt(FeltFromUint64(17))},
		},
		vm,
	)
	hintProcessor := CairoVmHintProcessor{}
	hintData := any(HintData{
		Ids:  idsManager,
		Code: DEFAULT_DICT_NEW,
	})
	// Advance AP so that values don't clash with FP-based ids
	vm.RunContext.Ap = NewRelocatable(0, 5)
	err := hintProcessor.ExecuteHint(vm, &hintData, nil, scopes)
	if err != nil {
		t.Errorf("DEFAULT_DICT_NEW hint test failed with error %s", err)
	}
	// Check that a manager was created in the scope
	_, ok := FetchDictManager(scopes)
	if !ok {
		t.Error("DEFAULT_DICT_NEW No DictManager created")
	}
	// Check that the correct base was inserted into ap
	val, _ := vm.Segments.Memory.Get(vm.RunContext.Ap)
	if val == nil || *val != *NewMaybeRelocatableRelocatable(NewRelocatable(1, 0)) {
		t.Error("DEFAULT_DICT_NEW Wrong/No base inserted into ap")
	}
}
