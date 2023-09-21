package hints_test

import (
	"testing"

	. "github.com/lambdaclass/cairo-vm.go/pkg/hints"
	. "github.com/lambdaclass/cairo-vm.go/pkg/hints/hint_codes"
	. "github.com/lambdaclass/cairo-vm.go/pkg/hints/hint_utils"
	. "github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
	. "github.com/lambdaclass/cairo-vm.go/pkg/types"
	. "github.com/lambdaclass/cairo-vm.go/pkg/vm"
	. "github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
)

func TestFindElementHintOk(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()
	vm.Segments.AddSegment()
	vm.Segments.Memory.Insert(NewRelocatable(1, 0), NewMaybeRelocatableFelt(FeltFromUint64(1)))
	vm.Segments.Memory.Insert(NewRelocatable(1, 1), NewMaybeRelocatableFelt(FeltFromUint64(2)))
	vm.Segments.Memory.Insert(NewRelocatable(1, 2), NewMaybeRelocatableFelt(FeltFromUint64(3)))
	vm.Segments.Memory.Insert(NewRelocatable(1, 3), NewMaybeRelocatableFelt(FeltFromUint64(4)))
	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{
			"array_ptr": {NewMaybeRelocatableRelocatable(NewRelocatable(1, 0))},
			"elm_size":  {NewMaybeRelocatableFelt(FeltFromUint64(2))},
			"n_elms":    {NewMaybeRelocatableFelt(FeltFromUint64(2))},
			"key":       {NewMaybeRelocatableFelt(FeltFromUint64(3))},
			"index":     {nil},
		},
		vm,
	)

	execScopes := NewExecutionScopes()

	hintProcessor := CairoVmHintProcessor{}
	hintData := any(HintData{
		Ids:  idsManager,
		Code: FIND_ELEMENT,
	})
	err := hintProcessor.ExecuteHint(vm, &hintData, nil, execScopes)
	if err != nil {
		t.Errorf("FIND_ELEMENT hint test failed with error: %s", err)
	}
	index, err := idsManager.GetFelt("index", vm)
	if err != nil {
		t.Errorf("%s", err)
	}
	if index.Cmp(FeltFromUint64(1)) != 0 {
		t.Errorf("Index was expected to be 1, got %s", index.ToSignedFeltString())
	}
}

func TestFindElementWithFindElementIndex(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()
	vm.Segments.AddSegment()
	vm.Segments.Memory.Insert(NewRelocatable(1, 0), NewMaybeRelocatableFelt(FeltFromUint64(1)))
	vm.Segments.Memory.Insert(NewRelocatable(1, 1), NewMaybeRelocatableFelt(FeltFromUint64(2)))
	vm.Segments.Memory.Insert(NewRelocatable(1, 2), NewMaybeRelocatableFelt(FeltFromUint64(3)))
	vm.Segments.Memory.Insert(NewRelocatable(1, 3), NewMaybeRelocatableFelt(FeltFromUint64(4)))
	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{
			"array_ptr": {NewMaybeRelocatableRelocatable(NewRelocatable(1, 0))},
			"elm_size":  {NewMaybeRelocatableFelt(FeltFromUint64(2))},
			"n_elms":    {NewMaybeRelocatableFelt(FeltFromUint64(2))},
			"key":       {NewMaybeRelocatableFelt(FeltFromUint64(3))},
			"index":     {nil},
		},
		vm,
	)

	execScopes := NewExecutionScopes()
	scope := make(map[string]interface{})
	scope["find_element_index"] = FeltOne()
	execScopes.EnterScope(scope)

	hintProcessor := CairoVmHintProcessor{}
	hintData := any(HintData{
		Ids:  idsManager,
		Code: FIND_ELEMENT,
	})
	err := hintProcessor.ExecuteHint(vm, &hintData, nil, execScopes)
	if err != nil {
		t.Errorf("FIND_ELEMENT hint test failed with error: %s", err)
	}
	index, err := idsManager.GetFelt("index", vm)
	if err != nil {
		t.Errorf("%s", err)
	}
	if index.Cmp(FeltFromUint64(1)) != 0 {
		t.Errorf("Index was expected to be 1, got %s", index.ToSignedFeltString())
	}
}

func TestFindElementFindElementMaxSizeOk(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()
	vm.Segments.AddSegment()
	vm.Segments.Memory.Insert(NewRelocatable(1, 0), NewMaybeRelocatableFelt(FeltFromUint64(1)))
	vm.Segments.Memory.Insert(NewRelocatable(1, 1), NewMaybeRelocatableFelt(FeltFromUint64(2)))
	vm.Segments.Memory.Insert(NewRelocatable(1, 2), NewMaybeRelocatableFelt(FeltFromUint64(3)))
	vm.Segments.Memory.Insert(NewRelocatable(1, 3), NewMaybeRelocatableFelt(FeltFromUint64(4)))
	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{
			"array_ptr": {NewMaybeRelocatableRelocatable(NewRelocatable(1, 0))},
			"elm_size":  {NewMaybeRelocatableFelt(FeltFromUint64(2))},
			"n_elms":    {NewMaybeRelocatableFelt(FeltFromUint64(2))},
			"key":       {NewMaybeRelocatableFelt(FeltFromUint64(3))},
			"index":     {nil},
		},
		vm,
	)

	execScopes := NewExecutionScopes()
	scope := make(map[string]interface{})
	scope["find_element_max_size"] = FeltFromUint64(2)
	execScopes.EnterScope(scope)

	hintProcessor := CairoVmHintProcessor{}
	hintData := any(HintData{
		Ids:  idsManager,
		Code: FIND_ELEMENT,
	})
	err := hintProcessor.ExecuteHint(vm, &hintData, nil, execScopes)
	if err != nil {
		t.Errorf("FIND_ELEMENT hint test failed with error: %s", err)
	}
	index, err := idsManager.GetFelt("index", vm)
	if err != nil {
		t.Errorf("%s", err)
	}
	if index.Cmp(FeltFromUint64(1)) != 0 {
		t.Errorf("Index was expected to be 1, got %s", index.ToSignedFeltString())
	}
}

func TestFindElementFindElementMaxSizeLessThanNeeded(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()
	vm.Segments.AddSegment()
	vm.Segments.Memory.Insert(NewRelocatable(1, 0), NewMaybeRelocatableFelt(FeltFromUint64(1)))
	vm.Segments.Memory.Insert(NewRelocatable(1, 1), NewMaybeRelocatableFelt(FeltFromUint64(2)))
	vm.Segments.Memory.Insert(NewRelocatable(1, 2), NewMaybeRelocatableFelt(FeltFromUint64(3)))
	vm.Segments.Memory.Insert(NewRelocatable(1, 3), NewMaybeRelocatableFelt(FeltFromUint64(4)))
	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{
			"array_ptr": {NewMaybeRelocatableRelocatable(NewRelocatable(1, 0))},
			"elm_size":  {NewMaybeRelocatableFelt(FeltFromUint64(2))},
			"n_elms":    {NewMaybeRelocatableFelt(FeltFromUint64(2))},
			"key":       {NewMaybeRelocatableFelt(FeltFromUint64(3))},
			"index":     {nil},
		},
		vm,
	)

	execScopes := NewExecutionScopes()
	scope := make(map[string]interface{})
	scope["find_element_max_size"] = FeltOne()
	execScopes.EnterScope(scope)

	hintProcessor := CairoVmHintProcessor{}
	hintData := any(HintData{
		Ids:  idsManager,
		Code: FIND_ELEMENT,
	})
	err := hintProcessor.ExecuteHint(vm, &hintData, nil, execScopes)
	if err == nil {
		t.Errorf("FIND_ELEMENT hint expected to fail with find_element_max_size < n_elms")
	}
}

func TestSearchSortedLowerHintOk(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()
	vm.Segments.AddSegment()
	vm.Segments.Memory.Insert(NewRelocatable(1, 0), NewMaybeRelocatableFelt(FeltFromUint64(1)))
	vm.Segments.Memory.Insert(NewRelocatable(1, 1), NewMaybeRelocatableFelt(FeltFromUint64(2)))
	vm.Segments.Memory.Insert(NewRelocatable(1, 2), NewMaybeRelocatableFelt(FeltFromUint64(3)))
	vm.Segments.Memory.Insert(NewRelocatable(1, 3), NewMaybeRelocatableFelt(FeltFromUint64(4)))
	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{
			"array_ptr": {NewMaybeRelocatableRelocatable(NewRelocatable(1, 0))},
			"elm_size":  {NewMaybeRelocatableFelt(FeltFromUint64(2))},
			"n_elms":    {NewMaybeRelocatableFelt(FeltFromUint64(2))},
			"key":       {NewMaybeRelocatableFelt(FeltFromUint64(0))},
			"index":     {nil},
		},
		vm,
	)

	execScopes := NewExecutionScopes()

	hintProcessor := CairoVmHintProcessor{}
	hintData := any(HintData{
		Ids:  idsManager,
		Code: SEARCH_SORTED_LOWER,
	})
	err := hintProcessor.ExecuteHint(vm, &hintData, nil, execScopes)
	if err != nil {
		t.Errorf("FIND_ELEMENT hint test failed with error: %s", err)
	}
	index, err := idsManager.GetFelt("index", vm)
	if err != nil {
		t.Errorf("%s", err)
	}
	if index.Cmp(FeltZero()) != 0 {
		t.Errorf("Index was expected to be 0, got %s", index.ToSignedFeltString())
	}
}

func TestSearchSortedLowerFindElementMaxSizeOk(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()
	vm.Segments.AddSegment()
	vm.Segments.Memory.Insert(NewRelocatable(1, 0), NewMaybeRelocatableFelt(FeltFromUint64(1)))
	vm.Segments.Memory.Insert(NewRelocatable(1, 1), NewMaybeRelocatableFelt(FeltFromUint64(2)))
	vm.Segments.Memory.Insert(NewRelocatable(1, 2), NewMaybeRelocatableFelt(FeltFromUint64(3)))
	vm.Segments.Memory.Insert(NewRelocatable(1, 3), NewMaybeRelocatableFelt(FeltFromUint64(4)))
	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{
			"array_ptr": {NewMaybeRelocatableRelocatable(NewRelocatable(1, 0))},
			"elm_size":  {NewMaybeRelocatableFelt(FeltFromUint64(2))},
			"n_elms":    {NewMaybeRelocatableFelt(FeltFromUint64(2))},
			"key":       {NewMaybeRelocatableFelt(FeltFromUint64(3))},
			"index":     {nil},
		},
		vm,
	)

	execScopes := NewExecutionScopes()
	scope := make(map[string]interface{})
	scope["find_element_max_size"] = FeltFromUint64(2)
	execScopes.EnterScope(scope)

	hintProcessor := CairoVmHintProcessor{}
	hintData := any(HintData{
		Ids:  idsManager,
		Code: SEARCH_SORTED_LOWER,
	})
	err := hintProcessor.ExecuteHint(vm, &hintData, nil, execScopes)
	if err != nil {
		t.Errorf("FIND_ELEMENT hint test failed with error: %s", err)
	}
	index, err := idsManager.GetFelt("index", vm)
	if err != nil {
		t.Errorf("%s", err)
	}
	if index.Cmp(FeltFromUint64(1)) != 0 {
		t.Errorf("Index was expected to be 1, got %s", index.ToSignedFeltString())
	}
}

func TestSearchSortedLowerFindElementMaxSizeLessThanNeeded(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()
	vm.Segments.AddSegment()
	vm.Segments.Memory.Insert(NewRelocatable(1, 0), NewMaybeRelocatableFelt(FeltFromUint64(1)))
	vm.Segments.Memory.Insert(NewRelocatable(1, 1), NewMaybeRelocatableFelt(FeltFromUint64(2)))
	vm.Segments.Memory.Insert(NewRelocatable(1, 2), NewMaybeRelocatableFelt(FeltFromUint64(3)))
	vm.Segments.Memory.Insert(NewRelocatable(1, 3), NewMaybeRelocatableFelt(FeltFromUint64(4)))
	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{
			"array_ptr": {NewMaybeRelocatableRelocatable(NewRelocatable(1, 0))},
			"elm_size":  {NewMaybeRelocatableFelt(FeltFromUint64(2))},
			"n_elms":    {NewMaybeRelocatableFelt(FeltFromUint64(2))},
			"key":       {NewMaybeRelocatableFelt(FeltFromUint64(3))},
			"index":     {nil},
		},
		vm,
	)

	execScopes := NewExecutionScopes()
	scope := make(map[string]interface{})
	scope["find_element_max_size"] = FeltOne()
	execScopes.EnterScope(scope)

	hintProcessor := CairoVmHintProcessor{}
	hintData := any(HintData{
		Ids:  idsManager,
		Code: SEARCH_SORTED_LOWER,
	})
	err := hintProcessor.ExecuteHint(vm, &hintData, nil, execScopes)
	if err == nil {
		t.Errorf("FIND_ELEMENT hint expected to fail with find_element_max_size < n_elms")
	}
}
