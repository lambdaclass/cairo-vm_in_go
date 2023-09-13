package hints_test

import (
	"testing"

	. "github.com/lambdaclass/cairo-vm.go/pkg/hints"
	"github.com/lambdaclass/cairo-vm.go/pkg/hints/dict_manager"
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

func TestDefaultDictNewHasManager(t *testing.T) {
	vm := NewVirtualMachine()
	scopes := types.NewExecutionScopes()
	// Create dictManager & add it to scope
	dictManager := dict_manager.NewDictManager()
	dictManagerRef := &dictManager
	scopes.AssignOrUpdateVariable("__dict_manager", dictManagerRef)
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
	// Check that the manager wan't replaced by a new one
	dictManagerPtr, ok := FetchDictManager(scopes)
	if !ok || dictManagerPtr != dictManagerRef {
		t.Error("DEFAULT_DICT_NEW DictManager replaced")
	}
	// Check that the correct base was inserted into ap
	val, _ := vm.Segments.Memory.Get(vm.RunContext.Ap)
	if val == nil || *val != *NewMaybeRelocatableRelocatable(NewRelocatable(1, 0)) {
		t.Error("DEFAULT_DICT_NEW Wrong/No base inserted into ap")
	}
}

func TestDictReadDefaultValue(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()
	scopes := types.NewExecutionScopes()

	// Create dictManager with a default dictionary & add it to scope
	dictManager := dict_manager.NewDictManager()
	defaultValue := NewMaybeRelocatableFelt(FeltFromUint64(17))
	dict_ptr := dictManager.NewDefaultDictionary(defaultValue, vm)
	scopes.AssignOrUpdateVariable("__dict_manager", &dictManager)

	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{
			"key":      {NewMaybeRelocatableFelt(FeltOne())},
			"dict_ptr": {NewMaybeRelocatableRelocatable(dict_ptr)},
			"value":    {nil},
		},
		vm,
	)
	hintProcessor := CairoVmHintProcessor{}
	hintData := any(HintData{
		Ids:  idsManager,
		Code: DICT_READ,
	})
	err := hintProcessor.ExecuteHint(vm, &hintData, nil, scopes)
	if err != nil {
		t.Errorf("DICT_READ hint test failed with error %s", err)
	}
	// Check ids.value
	val, err := idsManager.GetFelt("value", vm)
	if err != nil || val != FeltFromUint64(17) {
		t.Error("DEFAULT_DICT_NEW Wrong/No ids.value")
	}
}

func TestDictReadOk(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()
	scopes := types.NewExecutionScopes()

	// Create dictManager with a default dictionary & add it to scope
	dictManager := dict_manager.NewDictManager()
	initialDict := map[MaybeRelocatable]MaybeRelocatable{
		*NewMaybeRelocatableFelt(FeltOne()): *NewMaybeRelocatableFelt(FeltFromUint64(7)),
	}
	dict_ptr := dictManager.NewDictionary(&initialDict, vm)
	scopes.AssignOrUpdateVariable("__dict_manager", &dictManager)

	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{
			"key":      {NewMaybeRelocatableFelt(FeltOne())},
			"dict_ptr": {NewMaybeRelocatableRelocatable(dict_ptr)},
			"value":    {nil},
		},
		vm,
	)
	hintProcessor := CairoVmHintProcessor{}
	hintData := any(HintData{
		Ids:  idsManager,
		Code: DICT_READ,
	})
	err := hintProcessor.ExecuteHint(vm, &hintData, nil, scopes)
	if err != nil {
		t.Errorf("DICT_READ hint test failed with error %s", err)
	}
	// Check ids.value
	val, err := idsManager.GetFelt("value", vm)
	if err != nil || val != FeltFromUint64(7) {
		t.Error("DEFAULT_DICT_NEW Wrong/No ids.value")
	}
}

func TestDictReadNoVal(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()
	scopes := types.NewExecutionScopes()

	// Create dictManager with a default dictionary & add it to scope
	dictManager := dict_manager.NewDictManager()
	initialDict := map[MaybeRelocatable]MaybeRelocatable{}
	dict_ptr := dictManager.NewDictionary(&initialDict, vm)
	scopes.AssignOrUpdateVariable("__dict_manager", &dictManager)

	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{
			"key":      {NewMaybeRelocatableFelt(FeltOne())},
			"dict_ptr": {NewMaybeRelocatableRelocatable(dict_ptr)},
			"value":    {nil},
		},
		vm,
	)
	hintProcessor := CairoVmHintProcessor{}
	hintData := any(HintData{
		Ids:  idsManager,
		Code: DICT_READ,
	})
	err := hintProcessor.ExecuteHint(vm, &hintData, nil, scopes)
	if err == nil {
		t.Errorf("DICT_READ hint test should have failed")
	}
}

func TestDictWriteOk(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()
	scopes := types.NewExecutionScopes()

	// Create dictManager with a default dictionary & add it to scope
	dictManager := dict_manager.NewDictManager()
	initialDict := map[MaybeRelocatable]MaybeRelocatable{
		*NewMaybeRelocatableFelt(FeltOne()): *NewMaybeRelocatableFelt(FeltFromUint64(7)),
	}
	dict_ptr := dictManager.NewDictionary(&initialDict, vm)
	scopes.AssignOrUpdateVariable("__dict_manager", &dictManager)

	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{
			"key":       {NewMaybeRelocatableFelt(FeltOne())},
			"dict_ptr":  {NewMaybeRelocatableRelocatable(dict_ptr)},
			"new_value": {NewMaybeRelocatableFelt(FeltFromUint64(17))},
		},
		vm,
	)
	hintProcessor := CairoVmHintProcessor{}
	hintData := any(HintData{
		Ids:  idsManager,
		Code: DICT_WRITE,
	})
	err := hintProcessor.ExecuteHint(vm, &hintData, nil, scopes)
	if err != nil {
		t.Errorf("DICT_WRITE hint test failed with error %s", err)
	}
	// Check ids.prev_value
	val, err := vm.Segments.Memory.GetFelt(dict_ptr.AddUint(1))
	if err != nil || val != FeltFromUint64(7) {
		t.Error("DICT_WRITE Wrong/No ids.value")
	}
}

func TestDictWriteNoPrevValue(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()
	scopes := types.NewExecutionScopes()

	// Create dictManager with a default dictionary & add it to scope
	dictManager := dict_manager.NewDictManager()
	initialDict := map[MaybeRelocatable]MaybeRelocatable{}
	dict_ptr := dictManager.NewDictionary(&initialDict, vm)
	scopes.AssignOrUpdateVariable("__dict_manager", &dictManager)

	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{
			"key":       {NewMaybeRelocatableFelt(FeltOne())},
			"dict_ptr":  {NewMaybeRelocatableRelocatable(dict_ptr)},
			"new_value": {NewMaybeRelocatableFelt(FeltFromUint64(17))},
		},
		vm,
	)
	hintProcessor := CairoVmHintProcessor{}
	hintData := any(HintData{
		Ids:  idsManager,
		Code: DICT_WRITE,
	})
	err := hintProcessor.ExecuteHint(vm, &hintData, nil, scopes)
	if err == nil {
		t.Error("DICT_WRITE hint test should have failed")
	}
}

func TestDictWriteNewWriteDefault(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()
	scopes := types.NewExecutionScopes()

	// Create dictManager with a default dictionary & add it to scope
	dictManager := dict_manager.NewDictManager()
	defaultValue := FeltFromUint64(17)
	dict_ptr := dictManager.NewDefaultDictionary(NewMaybeRelocatableFelt(defaultValue), vm)
	scopes.AssignOrUpdateVariable("__dict_manager", &dictManager)

	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{
			"key":       {NewMaybeRelocatableFelt(FeltOne())},
			"dict_ptr":  {NewMaybeRelocatableRelocatable(dict_ptr)},
			"new_value": {NewMaybeRelocatableFelt(FeltFromUint64(17))},
		},
		vm,
	)
	hintProcessor := CairoVmHintProcessor{}
	hintData := any(HintData{
		Ids:  idsManager,
		Code: DICT_WRITE,
	})
	err := hintProcessor.ExecuteHint(vm, &hintData, nil, scopes)
	if err != nil {
		t.Errorf("DICT_WRITE hint test failed with error %s", err)
	}
	// Check ids.prev_value
	val, err := vm.Segments.Memory.GetFelt(dict_ptr.AddUint(1))
	if err != nil || val != defaultValue {
		t.Error("DICT_WRITE Wrong/No ids.value")
	}
}
