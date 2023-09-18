package hints_test

import (
	"testing"

	. "github.com/lambdaclass/cairo-vm.go/pkg/hints"
	. "github.com/lambdaclass/cairo-vm.go/pkg/hints/hint_utils"
	"github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
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
	constants := SetupConstantsForTest(map[string]Felt{
		"UPPER_BOUND": lambdaworks.FeltFromUint64(10),
		"SHIFT":       lambdaworks.FeltFromUint64(1),
	},
		&idsManager,
	)

	hintData := any(HintData{
		Ids:  idsManager,
		Code: ASSERT_250_BITS,
	})

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
	constants := SetupConstantsForTest(map[string]Felt{
		"UPPER_BOUND": lambdaworks.FeltFromUint64(10),
		"SHIFT":       lambdaworks.FeltFromUint64(1),
	},
		&idsManager,
	)

	hintData := any(HintData{
		Ids:  idsManager,
		Code: ASSERT_250_BITS,
	})

	err := hintProcessor.ExecuteHint(vm, &hintData, &constants, nil)
	if err == nil {
		t.Errorf("ASSERT_250_BIT hint should have failed with Value outside of 250 bit error")
	}
}

func TestSplitFeltAssertPrimeFailure(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()
	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{
			"value": {NewMaybeRelocatableFelt(FeltFromUint64(1))},
			"high":  {nil},
			"low":   {nil},
		},
		vm,
	)

	hintProcessor := CairoVmHintProcessor{}
	constants := SetupConstantsForTest(map[string]Felt{
		"MAX_HIGH": lambdaworks.FeltFromHex("0xffffffffffffffffffffffffffffffff"),
		"MAX_LOW":  lambdaworks.FeltFromHex("0xffffffffffffffffffffffffffffffff"),
	},
		&idsManager,
	)

	hintData := any(HintData{
		Ids:  idsManager,
		Code: SPLIT_FELT,
	})

	err := hintProcessor.ExecuteHint(vm, &hintData, &constants, nil)
	if err == nil {
		t.Errorf("SPLIT_FELT hint should have failed with assert PRIME - 1 == ids.MAX_HIGH * 2**128 + ids.MAX_LOW error")
	}
}

func TestSplitFeltAssertMaxHighFailedAssertion(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()
	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{
			"value": {NewMaybeRelocatableFelt(FeltFromUint64(1))},
			"high":  {nil},
			"low":   {nil},
		},
		vm,
	)

	hintProcessor := CairoVmHintProcessor{}
	constants := SetupConstantsForTest(map[string]Felt{
		"MAX_HIGH": lambdaworks.FeltFromHex("0xffffffffffffffffffffffffffffffffffff"),
		"MAX_LOW":  lambdaworks.FeltFromHex("0xffffffffffffffffffffffffffffffff"),
	},
		&idsManager,
	)

	hintData := any(HintData{
		Ids:  idsManager,
		Code: SPLIT_FELT,
	})

	err := hintProcessor.ExecuteHint(vm, &hintData, &constants, nil)
	if err == nil {
		t.Errorf("SPLIT_FELT hint should have failed with assert ids.MAX_HIGH < 2**128 and ids.MAX_LOW < 2**128")
	}
}

func TestSplitFeltAssertMaxLowFailedAssertion(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()
	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{
			"value": {NewMaybeRelocatableFelt(FeltFromUint64(1))},
			"high":  {nil},
			"low":   {nil},
		},
		vm,
	)

	hintProcessor := CairoVmHintProcessor{}
	constants := SetupConstantsForTest(map[string]Felt{
		"MAX_HIGH": lambdaworks.FeltFromHex("0xffffffffffffffffffffffffffffffff"),
		"MAX_LOW":  lambdaworks.FeltFromHex("0xffffffffffffffffffffffffffffffffffff"),
	},
		&idsManager,
	)

	hintData := any(HintData{
		Ids:  idsManager,
		Code: SPLIT_FELT,
	})

	err := hintProcessor.ExecuteHint(vm, &hintData, &constants, nil)
	if err == nil {
		t.Errorf("SPLIT_FELT hint should have failed with assert ids.MAX_HIGH < 2**128 and ids.MAX_LOW < 2**128")
	}
}

func TestSplitFeltSuccess(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()

	firstLimb := lambdaworks.FeltFromUint64(1)
	secondLimb := lambdaworks.FeltFromUint64(2)
	thirdLimb := lambdaworks.FeltFromUint64(3)
	fourthLimb := lambdaworks.FeltFromUint64(4)
	value := fourthLimb.Or(thirdLimb.Shl(64).Or(secondLimb.Shl(128).Or(firstLimb.Shl(192))))
	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{
			"value": {NewMaybeRelocatableFelt(value)},
			"high":  {nil},
			"low":   {nil},
		},
		vm,
	)

	hintProcessor := CairoVmHintProcessor{}
	constants := SetupConstantsForTest(map[string]Felt{
		"MAX_HIGH": lambdaworks.FeltFromDecString("10633823966279327296825105735305134080"),
		"MAX_LOW":  lambdaworks.FeltFromUint64(0),
	},
		&idsManager,
	)

	hintData := any(HintData{
		Ids:  idsManager,
		Code: SPLIT_FELT,
	})

	err := hintProcessor.ExecuteHint(vm, &hintData, &constants, nil)
	if err != nil {
		t.Errorf("SPLIT_FELT hint failed with error %s", err)
	}

	high, err := idsManager.GetFelt("high", vm)
	if err != nil {
		t.Errorf("failed to get high: %s", err)
	}

	low, err := idsManager.GetFelt("low", vm)
	if err != nil {
		t.Errorf("failed to get low: %s", err)
	}

	if high != firstLimb.Shl(64).Or(secondLimb) {
		t.Errorf("Expected high == 335438970432432812899076431678123043273. Got: %v", high)
	}

	if low != thirdLimb.Shl(64).Or(fourthLimb) {
		t.Errorf("Expected low == 0. Got: %v", low)
	}
}
