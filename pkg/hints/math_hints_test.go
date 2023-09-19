package hints_test

import (
	"testing"

	. "github.com/lambdaclass/cairo-vm.go/pkg/hints"
	. "github.com/lambdaclass/cairo-vm.go/pkg/hints/hint_utils"
	"github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
	. "github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
	. "github.com/lambdaclass/cairo-vm.go/pkg/types"
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
	err := hintProcessor.ExecuteHint(vm, &hintData, nil, nil)
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
	err := hintProcessor.ExecuteHint(vm, &hintData, nil, nil)
	if err == nil {
		t.Errorf("ASSERT_NN hint should have failed")
	}
}

func TestIsPositiveOkPositive(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()
	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{
			"value":       {NewMaybeRelocatableFelt(FeltFromUint64(17))},
			"is_positive": {nil},
		},
		vm,
	)
	hintProcessor := CairoVmHintProcessor{}
	hintData := any(HintData{
		Ids:  idsManager,
		Code: IS_POSITIVE,
	})
	err := hintProcessor.ExecuteHint(vm, &hintData, nil, nil)
	if err != nil {
		t.Errorf("IS_POSITIVE hint test failed with error %s", err)
	}
	// Check ids.is_positive
	is_positive, err := idsManager.GetFelt("is_positive", vm)
	if err != nil || is_positive != FeltFromUint64(1) {
		t.Errorf("IS_POSITIVE hint test incorrect value for ids.is_positive")
	}
}

func TestIsPositiveOkNegative(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()
	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{
			"value":       {NewMaybeRelocatableFelt(FeltFromDecString("-1"))},
			"is_positive": {nil},
		},
		vm,
	)
	hintProcessor := CairoVmHintProcessor{}
	hintData := any(HintData{
		Ids:  idsManager,
		Code: IS_POSITIVE,
	})
	err := hintProcessor.ExecuteHint(vm, &hintData, nil, nil)
	if err != nil {
		t.Errorf("IS_POSITIVE hint test failed with error %s", err)
	}
	// Check ids.is_positive
	is_positive, err := idsManager.GetFelt("is_positive", vm)
	if err != nil || is_positive != FeltFromUint64(0) {
		t.Errorf("IS_POSITIVE hint test incorrect value for ids.is_positive")
	}
}

func TestIsPositiveOutOfRange(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()
	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{
			"value":       {NewMaybeRelocatableFelt(FeltFromDecString("340282366920938463463374607431768211457"))},
			"is_positive": {nil},
		},
		vm,
	)
	hintProcessor := CairoVmHintProcessor{}
	hintData := any(HintData{
		Ids:  idsManager,
		Code: IS_POSITIVE,
	})
	err := hintProcessor.ExecuteHint(vm, &hintData, nil, nil)
	if err == nil {
		t.Errorf("IS_POSITIVE hint test should have failed")
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
	err := hintProcessor.ExecuteHint(vm, &hintData, nil, nil)
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
	err := hintProcessor.ExecuteHint(vm, &hintData, nil, nil)
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
	err := hintProcessor.ExecuteHint(vm, &hintData, nil, nil)
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
	err := hintProcessor.ExecuteHint(vm, &hintData, nil, nil)
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
	err := hintProcessor.ExecuteHint(vm, &hintData, nil, nil)
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
	err := hintProcessor.ExecuteHint(vm, &hintData, nil, nil)
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
	err := hintProcessor.ExecuteHint(vm, &hintData, nil, nil)
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
	err := hintProcessor.ExecuteHint(vm, &hintData, nil, nil)
	if err != nil {
		t.Errorf("ASSERT_NOT_EQUAL hint failed with error: %s", err)
	}
}

func TestSqrtOk(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()
	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{
			"value": {NewMaybeRelocatableFelt(FeltFromDecString("9"))},
			"root":  {nil},
		},
		vm,
	)
	hintProcessor := CairoVmHintProcessor{}
	hintData := any(HintData{
		Ids:  idsManager,
		Code: SQRT,
	})
	err := hintProcessor.ExecuteHint(vm, &hintData, nil, nil)
	if err != nil {
		t.Errorf("SQRT hint failed with error: %s", err)
	}

	root, err := idsManager.GetFelt("root", vm)
	if err != nil {
		t.Errorf("failed to get root: %s", err)
	}
	if root != FeltFromUint64(3) {
		t.Errorf("Expected sqrt(9) == 3. Got: %v", root)
	}
}

func TestAssertLeFeltOk(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()
	vm.Segments.AddSegment()
	scopes := NewExecutionScopes()
	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{
			"a":               {NewMaybeRelocatableFelt(FeltOne())},
			"b":               {NewMaybeRelocatableFelt(FeltFromUint64(2))},
			"range_check_ptr": {NewMaybeRelocatableRelocatable(NewRelocatable(1, 0))},
		},
		vm,
	)
	constants := SetupConstantsForTest(map[string]Felt{
		"PRIME_OVER_3_HIGH": FeltFromHex("4000000000000088000000000000001"),
		"PRIME_OVER_2_HIGH": FeltFromHex("2AAAAAAAAAAAAB05555555555555556"),
	},
		&idsManager,
	)
	hintProcessor := CairoVmHintProcessor{}
	hintData := any(HintData{
		Ids:  idsManager,
		Code: ASSERT_LE_FELT,
	})
	err := hintProcessor.ExecuteHint(vm, &hintData, &constants, scopes)
	if err != nil {
		t.Errorf("ASSERT_LE_FELT hint failed with error: %s", err)
	}
}

func TestAssertLeFeltExcluded0Zero(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()
	scopes := NewExecutionScopes()
	scopes.AssignOrUpdateVariable("excluded", int(0))
	hintProcessor := CairoVmHintProcessor{}
	hintData := any(HintData{
		Code: ASSERT_LE_FELT_EXCLUDED_0,
	})
	err := hintProcessor.ExecuteHint(vm, &hintData, nil, scopes)
	if err != nil {
		t.Errorf("ASSERT_LE_FELT_EXCLUDED_0 hint test failed with error %s", err)
	}
	// Check the value of memory[ap]
	val, err := vm.Segments.Memory.GetFelt(vm.RunContext.Ap)
	if err != nil || !val.IsZero() {
		t.Error("Wrong/No value inserted into ap")
	}
}

func TestAssertLeFeltExcluded0One(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()
	scopes := NewExecutionScopes()
	scopes.AssignOrUpdateVariable("excluded", int(1))
	hintProcessor := CairoVmHintProcessor{}
	hintData := any(HintData{
		Code: ASSERT_LE_FELT_EXCLUDED_0,
	})
	err := hintProcessor.ExecuteHint(vm, &hintData, nil, scopes)
	if err != nil {
		t.Errorf("ASSERT_LE_FELT_EXCLUDED_0 hint test failed with error %s", err)
	}
	// Check the value of memory[ap]
	val, err := vm.Segments.Memory.GetFelt(vm.RunContext.Ap)
	if err != nil || val != FeltOne() {
		t.Error("Wrong/No value inserted into ap")
	}
}

func TestAssertLeFeltExcluded1Zero(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()
	scopes := NewExecutionScopes()
	scopes.AssignOrUpdateVariable("excluded", int(1))
	hintProcessor := CairoVmHintProcessor{}
	hintData := any(HintData{
		Code: ASSERT_LE_FELT_EXCLUDED_1,
	})
	err := hintProcessor.ExecuteHint(vm, &hintData, nil, scopes)
	if err != nil {
		t.Errorf("ASSERT_LE_FELT_EXCLUDED_1 hint test failed with error %s", err)
	}
	// Check the value of memory[ap]
	val, err := vm.Segments.Memory.GetFelt(vm.RunContext.Ap)
	if err != nil || !val.IsZero() {
		t.Error("Wrong/No value inserted into ap")
	}
}

func TestAssertLeFeltExcluded1One(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()
	scopes := NewExecutionScopes()
	scopes.AssignOrUpdateVariable("excluded", int(0))
	hintProcessor := CairoVmHintProcessor{}
	hintData := any(HintData{
		Code: ASSERT_LE_FELT_EXCLUDED_1,
	})
	err := hintProcessor.ExecuteHint(vm, &hintData, nil, scopes)
	if err != nil {
		t.Errorf("ASSERT_LE_FELT_EXCLUDED_1 hint test failed with error %s", err)
	}
	// Check the value of memory[ap]
	val, err := vm.Segments.Memory.GetFelt(vm.RunContext.Ap)
	if err != nil || val != FeltOne() {
		t.Error("Wrong/No value inserted into ap")
	}
}

func TestAssertLeFeltExcluded2Ok(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()
	scopes := NewExecutionScopes()
	scopes.AssignOrUpdateVariable("excluded", int(2))
	hintProcessor := CairoVmHintProcessor{}
	hintData := any(HintData{
		Code: ASSERT_LE_FELT_EXCLUDED_2,
	})
	err := hintProcessor.ExecuteHint(vm, &hintData, nil, scopes)
	if err != nil {
		t.Errorf("ASSERT_LE_FELT_EXCLUDED_2 hint test failed with error %s", err)
	}
}

func TestAssertLeFeltExcluded2Err(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()
	scopes := NewExecutionScopes()
	scopes.AssignOrUpdateVariable("excluded", int(0))
	hintProcessor := CairoVmHintProcessor{}
	hintData := any(HintData{
		Code: ASSERT_LE_FELT_EXCLUDED_2,
	})
	err := hintProcessor.ExecuteHint(vm, &hintData, nil, scopes)
	if err == nil {
		t.Errorf("ASSERT_LE_FELT_EXCLUDED_2 hint test should have failed")
	}
}

func TestAssertLeFeltHintOk(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()
	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{
			"a": {NewMaybeRelocatableFelt(FeltFromUint64(17))},
			"b": {NewMaybeRelocatableFelt(FeltFromUint64(18))},
		},
		vm,
	)
	hintProcessor := CairoVmHintProcessor{}
	hintData := any(HintData{
		Ids:  idsManager,
		Code: ASSERT_LT_FELT,
	})
	err := hintProcessor.ExecuteHint(vm, &hintData, nil, nil)
	if err != nil {
		t.Errorf("ASSERT_LT_FELT hint test failed with error %s", err)
	}
}

func TestAssertLeFeltHintErr(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()
	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{
			"a": {NewMaybeRelocatableFelt(FeltFromUint64(17))},
			"b": {NewMaybeRelocatableFelt(FeltFromUint64(16))},
		},
		vm,
	)
	hintProcessor := CairoVmHintProcessor{}
	hintData := any(HintData{
		Ids:  idsManager,
		Code: ASSERT_LT_FELT,
	})
	err := hintProcessor.ExecuteHint(vm, &hintData, nil, nil)
	if err == nil {
		t.Errorf("ASSERT_LT_FELT hint test should have failed")
	}
}
func TestAssert250BitHintSuccess(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()
	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{
			"value": {NewMaybeRelocatableFelt(FeltFromUint64(3))},
			"high":  {nil},
			"low":   {nil},
		},
		vm,
	)

	hintProcessor := CairoVmHintProcessor{}
	hintData := any(HintData{
		Ids:  idsManager,
		Code: ASSERT_250_BITS,
	})

	constants := make(map[string]Felt)
	constants["UPPER_BOUND"] = lambdaworks.FeltFromUint64(10)
	constants["SHIFT"] = lambdaworks.FeltFromUint64(1)

	err := hintProcessor.ExecuteHint(vm, &hintData, &constants, nil)
	if err != nil {
		t.Errorf("ASSERT_250_BIT hint failed with error %s", err)
	}

	high, err := idsManager.GetFelt("high", vm)
	if err != nil {
		t.Errorf("failed to get high: %s", err)
	}

	low, err := idsManager.GetFelt("low", vm)
	if err != nil {
		t.Errorf("failed to get low: %s", err)
	}

	if high != FeltFromUint64(3) {
		t.Errorf("Expected high == 3. Got: %v", high)
	}

	if low != FeltFromUint64(0) {
		t.Errorf("Expected low == 0. Got: %v", low)
	}
}

func TestAssert250BitHintFail(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()
	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{
			"value": {NewMaybeRelocatableFelt(FeltFromUint64(20))},
			"high":  {nil},
			"low":   {nil},
		},
		vm,
	)

	hintProcessor := CairoVmHintProcessor{}
	hintData := any(HintData{
		Ids:  idsManager,
		Code: ASSERT_250_BITS,
	})

	constants := make(map[string]Felt)
	constants["UPPER_BOUND"] = lambdaworks.FeltFromUint64(10)
	constants["SHIFT"] = lambdaworks.FeltFromUint64(1)

	err := hintProcessor.ExecuteHint(vm, &hintData, &constants, nil)
	if err == nil {
		t.Errorf("ASSERT_250_BIT hint should have failed with Value outside of 250 bit error")
	}
}
