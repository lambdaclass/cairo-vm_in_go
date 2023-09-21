package hints_test

import (
	"testing"

	. "github.com/lambdaclass/cairo-vm.go/pkg/hints"
	. "github.com/lambdaclass/cairo-vm.go/pkg/hints/hint_codes"
	. "github.com/lambdaclass/cairo-vm.go/pkg/hints/hint_utils"
	. "github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
	. "github.com/lambdaclass/cairo-vm.go/pkg/vm"
	. "github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
)

func TestSetAddElmInSet(t *testing.T) {
	vm := NewVirtualMachine()
	// Initialize segments
	vm.Segments.AddSegment()
	vm.Segments.AddSegment()
	// element to insert
	vm.Segments.Memory.Insert(NewRelocatable(1, 0), NewMaybeRelocatableFelt(FeltFromUint64(2)))
	// set
	vm.Segments.Memory.Insert(NewRelocatable(1, 1), NewMaybeRelocatableFelt(FeltFromUint64(1)))
	vm.Segments.Memory.Insert(NewRelocatable(1, 2), NewMaybeRelocatableFelt(FeltFromUint64(2)))
	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{
			"elm_ptr":       {NewMaybeRelocatableRelocatable(NewRelocatable(1, 0))},
			"set_ptr":       {NewMaybeRelocatableRelocatable(NewRelocatable(1, 1))},
			"set_end_ptr":   {NewMaybeRelocatableRelocatable(NewRelocatable(1, 4))},
			"elm_size":      {NewMaybeRelocatableFelt(FeltFromUint64(1))},
			"index":         {nil},
			"is_elm_in_set": {nil},
		},
		vm,
	)

	hintProcessor := CairoVmHintProcessor{}
	hintData := any(HintData{
		Ids:  idsManager,
		Code: SET_ADD,
	})
	err := hintProcessor.ExecuteHint(vm, &hintData, nil, nil)
	if err != nil {
		t.Errorf("SET_ADD failed with error: %s", err)
	}

	isElmInSet, err := idsManager.GetFelt("is_elm_in_set", vm)
	if err != nil {
		t.Errorf("SET_ADD couldn't get is_elm_in_set: %s", err)
	}

	if !isElmInSet.IsOne() {
		t.Errorf("Expected is_elm_in_set to be 1, got: %s", isElmInSet.ToSignedFeltString())
	}

	index, err := idsManager.GetFelt("index", vm)
	if err != nil {
		t.Errorf("SET_ADD couldn't get index: %s", err)
	}
	if !index.IsOne() {
		t.Errorf("Expected element to be found at 1, got index: %s", index.ToSignedFeltString())
	}
}

func TestSetAddElmNotInSet(t *testing.T) {
	vm := NewVirtualMachine()
	// Initialize segments
	vm.Segments.AddSegment()
	vm.Segments.AddSegment()
	// element to insert
	vm.Segments.Memory.Insert(NewRelocatable(1, 0), NewMaybeRelocatableFelt(FeltFromUint64(3)))
	// set
	vm.Segments.Memory.Insert(NewRelocatable(1, 1), NewMaybeRelocatableFelt(FeltFromUint64(1)))
	vm.Segments.Memory.Insert(NewRelocatable(1, 2), NewMaybeRelocatableFelt(FeltFromUint64(2)))
	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{
			"elm_ptr":       {NewMaybeRelocatableRelocatable(NewRelocatable(1, 0))},
			"set_ptr":       {NewMaybeRelocatableRelocatable(NewRelocatable(1, 1))},
			"set_end_ptr":   {NewMaybeRelocatableRelocatable(NewRelocatable(1, 3))},
			"elm_size":      {NewMaybeRelocatableFelt(FeltFromUint64(1))},
			"index":         {nil},
			"is_elm_in_set": {nil},
		},
		vm,
	)

	hintProcessor := CairoVmHintProcessor{}
	hintData := any(HintData{
		Ids:  idsManager,
		Code: SET_ADD,
	})
	err := hintProcessor.ExecuteHint(vm, &hintData, nil, nil)
	if err != nil {
		t.Errorf("SET_ADD failed with error: %s", err)
	}

	isElmInSet, err := idsManager.GetFelt("is_elm_in_set", vm)
	if err != nil {
		t.Errorf("SET_ADD couldn't get is_elm_in_set: %s", err)
	}

	if !isElmInSet.IsZero() {
		t.Errorf("Expected is_elm_in_set to be 1, got: %s", isElmInSet.ToSignedFeltString())
	}
}
