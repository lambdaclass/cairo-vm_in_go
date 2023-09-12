package hints_test

import (
	"testing"

	. "github.com/lambdaclass/cairo-vm.go/pkg/hints"
	. "github.com/lambdaclass/cairo-vm.go/pkg/hints/hint_utils"
	. "github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
	. "github.com/lambdaclass/cairo-vm.go/pkg/vm"
	. "github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
)

func TestIsNNHintOk(t *testing.T) {
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
		Code: ASSERT_NN,
	})
	err := hintProcessor.ExecuteHint(vm, &hintData, nil)
	if err != nil {
		t.Errorf("ASSERT_NN hint test failed with error %s", err)
	}
}

func TestIsNNHintFail(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()
	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{
			"a": {NewMaybeRelocatableFelt(FeltFromDecString("-1"))},
		},
		vm,
	)
	hintProcessor := CairoVmHintProcessor{}
	hintData := any(HintData{
		Ids:  idsManager,
		Code: ASSERT_NN,
	})
	err := hintProcessor.ExecuteHint(vm, &hintData, nil)
	if err == nil {
		t.Errorf("ASSERT_NN hint should have failed")
	}
}

func TestAssertNotZeroHintOk(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()
	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{
			"value": {NewMaybeRelocatableFelt(FeltFromUint64(17))},
		},
		vm,
	)
	hintProcessor := CairoVmHintProcessor{}
	hintData := any(HintData{
		Ids:  idsManager,
		Code: ASSERT_NOT_ZERO,
	})
	err := hintProcessor.ExecuteHint(vm, &hintData, nil)
	if err != nil {
		t.Errorf("ASSERT_NOT_ZERO hint test failed with error %s", err)
	}
}

func TestAssertNotZeroHintFail(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()
	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{
			"value": {NewMaybeRelocatableFelt(FeltFromUint64(0))},
		},
		vm,
	)
	hintProcessor := CairoVmHintProcessor{}
	hintData := any(HintData{
		Ids:  idsManager,
		Code: ASSERT_NOT_ZERO,
	})
	err := hintProcessor.ExecuteHint(vm, &hintData, nil)
	if err == nil {
		t.Errorf("ASSERT_NOT_ZERO hint should have failed")
	}
}

func TestAssertNotEqualHintNonComparableDiffType(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()
	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{
			"a": {NewMaybeRelocatableFelt(FeltFromUint64(0))},
			"b": {NewMaybeRelocatableRelocatable(NewRelocatable(0, 0))},
		},
		vm,
	)
	hintProcessor := CairoVmHintProcessor{}
	hintData := any(HintData{
		Ids:  idsManager,
		Code: ASSERT_NOT_EQUAL,
	})
	err := hintProcessor.ExecuteHint(vm, &hintData, nil)
	if err == nil {
		t.Errorf("ASSERT_NOT_EQUAL hint should have failed")
	}
}

func TestAssertNotEqualHintNonComparableDiffIndex(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()
	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{
			"a": {NewMaybeRelocatableRelocatable(NewRelocatable(1, 0))},
			"b": {NewMaybeRelocatableRelocatable(NewRelocatable(0, 0))},
		},
		vm,
	)
	hintProcessor := CairoVmHintProcessor{}
	hintData := any(HintData{
		Ids:  idsManager,
		Code: ASSERT_NOT_EQUAL,
	})
	err := hintProcessor.ExecuteHint(vm, &hintData, nil)
	if err == nil {
		t.Errorf("ASSERT_NOT_EQUAL hint should have failed")
	}
}

func TestAssertNotEqualHintEqualRelocatables(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()
	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{
			"a": {NewMaybeRelocatableRelocatable(NewRelocatable(0, 0))},
			"b": {NewMaybeRelocatableRelocatable(NewRelocatable(0, 0))},
		},
		vm,
	)
	hintProcessor := CairoVmHintProcessor{}
	hintData := any(HintData{
		Ids:  idsManager,
		Code: ASSERT_NOT_EQUAL,
	})
	err := hintProcessor.ExecuteHint(vm, &hintData, nil)
	if err == nil {
		t.Errorf("ASSERT_NOT_EQUAL hint should have failed")
	}
}

func TestAssertNotEqualHintEqualFelts(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()
	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{
			"a": {NewMaybeRelocatableFelt(FeltFromUint64(9))},
			"b": {NewMaybeRelocatableFelt(FeltFromUint64(9))},
		},
		vm,
	)
	hintProcessor := CairoVmHintProcessor{}
	hintData := any(HintData{
		Ids:  idsManager,
		Code: ASSERT_NOT_EQUAL,
	})
	err := hintProcessor.ExecuteHint(vm, &hintData, nil)
	if err == nil {
		t.Errorf("ASSERT_NOT_EQUAL hint should have failed")
	}
}

func TestAssertNotEqualHintOkFelts(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()
	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{
			"a": {NewMaybeRelocatableFelt(FeltFromUint64(9))},
			"b": {NewMaybeRelocatableFelt(FeltFromUint64(7))},
		},
		vm,
	)
	hintProcessor := CairoVmHintProcessor{}
	hintData := any(HintData{
		Ids:  idsManager,
		Code: ASSERT_NOT_EQUAL,
	})
	err := hintProcessor.ExecuteHint(vm, &hintData, nil)
	if err != nil {
		t.Errorf("ASSERT_NOT_EQUAL hint failed with error: %s", err)
	}
}

func TestAssertNotEqualHintOkRelocatables(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()
	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{
			"a": {NewMaybeRelocatableRelocatable(NewRelocatable(1, 9))},
			"b": {NewMaybeRelocatableRelocatable(NewRelocatable(1, 7))},
		},
		vm,
	)
	hintProcessor := CairoVmHintProcessor{}
	hintData := any(HintData{
		Ids:  idsManager,
		Code: ASSERT_NOT_EQUAL,
	})
	err := hintProcessor.ExecuteHint(vm, &hintData, nil)
	if err != nil {
		t.Errorf("ASSERT_NOT_EQUAL hint failed with error: %s", err)
	}
}
