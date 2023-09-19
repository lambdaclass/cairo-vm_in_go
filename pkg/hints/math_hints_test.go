package hints_test

import (
	"testing"

	"github.com/lambdaclass/cairo-vm.go/pkg/builtins"
	. "github.com/lambdaclass/cairo-vm.go/pkg/hints"
	. "github.com/lambdaclass/cairo-vm.go/pkg/hints/hint_codes"
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

func TestUnsignedDivRemHintSuccess(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()
	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{
			"div":   {NewMaybeRelocatableFelt(FeltFromDecString("7"))},
			"value": {NewMaybeRelocatableFelt(FeltFromDecString("15"))},
			"r":     {nil},
			"q":     {nil},
		},
		vm,
	)
	rcBuiltin := builtins.DefaultRangeCheckBuiltinRunner()
	vm.BuiltinRunners = []builtins.BuiltinRunner{rcBuiltin}

	hintProcessor := CairoVmHintProcessor{}
	hintData := any(HintData{
		Ids:  idsManager,
		Code: UNSIGNED_DIV_REM,
	})

	err := hintProcessor.ExecuteHint(vm, &hintData, nil, nil)
	if err != nil {
		t.Errorf("UNSIGNED_DIV_REM hint failed with error: %s", err)
	}

	q, err := idsManager.GetFelt("q", vm)
	if err != nil {
		t.Errorf("failed to get `q`: %s", err)
	}

	r, err := idsManager.GetFelt("r", vm)
	if err != nil {
		t.Errorf("failed to get `r`: %s", err)
	}

	if q != FeltFromUint64(2) {
		t.Errorf("Expected q=3, got: %v", q)
	}

	if r != FeltFromUint64(1) {
		t.Errorf("Expected r=1, got: %v", r)
	}
}

func TestUnsignedDivRemHintDivZeroError(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()
	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{
			// This is the condition that should make the hint execution error.
			"div":   {NewMaybeRelocatableFelt(FeltFromDecString("0"))},
			"value": {NewMaybeRelocatableFelt(FeltFromDecString("15"))},
			"r":     {nil},
			"q":     {nil},
		},
		vm,
	)
	rcBuiltin := builtins.DefaultRangeCheckBuiltinRunner()
	vm.BuiltinRunners = []builtins.BuiltinRunner{rcBuiltin}

	hintProcessor := CairoVmHintProcessor{}
	hintData := any(HintData{
		Ids:  idsManager,
		Code: UNSIGNED_DIV_REM,
	})

	err := hintProcessor.ExecuteHint(vm, &hintData, nil, nil)
	if err == nil {
		t.Error("UNSIGNED_DIV_REM should have failed")
	}
}

func TestUnsignedDivRemHintOutOfBoundsError(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()
	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{
			// This is the condition that should make the hint execution error.
			"div":   {NewMaybeRelocatableFelt(FeltFromDecString("10633823966279327296825105735305134081"))},
			"value": {NewMaybeRelocatableFelt(FeltFromDecString("15"))},
			"r":     {nil},
			"q":     {nil},
		},
		vm,
	)
	rcBuiltin := builtins.DefaultRangeCheckBuiltinRunner()
	vm.BuiltinRunners = []builtins.BuiltinRunner{rcBuiltin}

	hintProcessor := CairoVmHintProcessor{}
	hintData := any(HintData{
		Ids:  idsManager,
		Code: UNSIGNED_DIV_REM,
	})

	err := hintProcessor.ExecuteHint(vm, &hintData, nil, nil)
	if err == nil {
		t.Errorf("UNSIGNED_DIV_REM should have failed")
	}
}

func TestSignedDivRemHintSuccess(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()
	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{
			"div":      {NewMaybeRelocatableFelt(FeltFromDecString("3"))},
			"value":    {NewMaybeRelocatableFelt(FeltFromDecString("10"))},
			"bound":    {NewMaybeRelocatableFelt(FeltFromDecString("29"))},
			"r":        {nil},
			"biased_q": {nil},
		},
		vm,
	)
	rcBuiltin := builtins.DefaultRangeCheckBuiltinRunner()
	vm.BuiltinRunners = []builtins.BuiltinRunner{rcBuiltin}

	hintProcessor := CairoVmHintProcessor{}
	hintData := any(HintData{
		Ids:  idsManager,
		Code: SIGNED_DIV_REM,
	})

	err := hintProcessor.ExecuteHint(vm, &hintData, nil, nil)
	if err != nil {
		t.Errorf("UNSIGNED_DIV_REM hint failed with error: %s", err)
	}

	biasedQ, err := idsManager.GetFelt("biased_q", vm)
	if err != nil {
		t.Errorf("failed to get `biased_q`: %s", err)
	}

	r, err := idsManager.GetFelt("r", vm)
	if err != nil {
		t.Errorf("failed to get `r`: %s", err)
	}

	if biasedQ != FeltFromUint64(32) {
		t.Errorf("Expected biased_q=32, got: %v", biasedQ)
	}

	if r != FeltFromUint64(1) {
		t.Errorf("Expected r=1, got: %v", r)
	}
}

func TestSignedDivRemHintDivZeroError(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()
	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{
			"div":      {NewMaybeRelocatableFelt(FeltFromDecString("0"))},
			"value":    {NewMaybeRelocatableFelt(FeltFromDecString("10"))},
			"bound":    {NewMaybeRelocatableFelt(FeltFromDecString("29"))},
			"r":        {nil},
			"biased_q": {nil},
		},
		vm,
	)
	rcBuiltin := builtins.DefaultRangeCheckBuiltinRunner()
	vm.BuiltinRunners = []builtins.BuiltinRunner{rcBuiltin}

	hintProcessor := CairoVmHintProcessor{}
	hintData := any(HintData{
		Ids:  idsManager,
		Code: SIGNED_DIV_REM,
	})

	err := hintProcessor.ExecuteHint(vm, &hintData, nil, nil)
	if err == nil {
		t.Errorf("UNSIGNED_DIV_REM hint should have failed")
	}
}

func TestSignedDivRemHintOutOfRcBoundsError(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()
	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{
			"div":      {NewMaybeRelocatableFelt(FeltFromDecString("10633823966279327296825105735305134081"))},
			"value":    {NewMaybeRelocatableFelt(FeltFromDecString("10"))},
			"bound":    {NewMaybeRelocatableFelt(FeltFromDecString("29"))},
			"r":        {nil},
			"biased_q": {nil},
		},
		vm,
	)
	rcBuiltin := builtins.DefaultRangeCheckBuiltinRunner()
	vm.BuiltinRunners = []builtins.BuiltinRunner{rcBuiltin}

	hintProcessor := CairoVmHintProcessor{}
	hintData := any(HintData{
		Ids:  idsManager,
		Code: SIGNED_DIV_REM,
	})

	err := hintProcessor.ExecuteHint(vm, &hintData, nil, nil)
	if err == nil {
		t.Errorf("UNSIGNED_DIV_REM hint should have failed")
	}
}

func TestSignedDivRemHintOutOfBoundsError(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()
	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{
			"div":      {NewMaybeRelocatableFelt(FeltFromDecString("4"))},
			"value":    {NewMaybeRelocatableFelt(FeltFromDecString("16"))},
			"bound":    {NewMaybeRelocatableFelt(FeltFromDecString("2"))},
			"r":        {nil},
			"biased_q": {nil},
		},
		vm,
	)
	rcBuiltin := builtins.DefaultRangeCheckBuiltinRunner()
	vm.BuiltinRunners = []builtins.BuiltinRunner{rcBuiltin}

	hintProcessor := CairoVmHintProcessor{}
	hintData := any(HintData{
		Ids:  idsManager,
		Code: SIGNED_DIV_REM,
	})

	err := hintProcessor.ExecuteHint(vm, &hintData, nil, nil)
	if err == nil {
		t.Errorf("UNSIGNED_DIV_REM hint should have failed")
	}
}
