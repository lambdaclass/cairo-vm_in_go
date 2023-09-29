package hints_test

import (
	"math/big"
	"testing"

	. "github.com/lambdaclass/cairo-vm.go/pkg/hints"
	. "github.com/lambdaclass/cairo-vm.go/pkg/hints/hint_codes"
	. "github.com/lambdaclass/cairo-vm.go/pkg/hints/hint_utils"
	. "github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
	"github.com/lambdaclass/cairo-vm.go/pkg/parser"
	. "github.com/lambdaclass/cairo-vm.go/pkg/types"
	. "github.com/lambdaclass/cairo-vm.go/pkg/utils"
	. "github.com/lambdaclass/cairo-vm.go/pkg/vm"
	. "github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
)

func TestUint256AddOk(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()

	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{
			"a": {
				NewMaybeRelocatableFeltFromUint64(2),
				NewMaybeRelocatableFeltFromUint64(3),
			},
			"b": {
				NewMaybeRelocatableFeltFromUint64(4),
				NewMaybeRelocatableFelt(FeltFromDecString("340282366920938463463374607431768211455")),
			},
			"carry_low":  {nil},
			"carry_high": {nil},
		},
		vm,
	)
	hintData := any(HintData{
		Ids:  idsManager,
		Code: UINT256_ADD,
	})
	scopes := NewExecutionScopes()
	hintProcessor := CairoVmHintProcessor{}
	err := hintProcessor.ExecuteHint(vm, &hintData, nil, scopes)
	if err != nil {
		t.Errorf("failed with error: %s", err)
	}

	carry_low, err := idsManager.GetFelt("carry_low", vm)
	if err != nil {
		t.Errorf("failed with error: %s", err)
	}
	if carry_low.Cmp(FeltZero()) != 0 {
		t.Errorf("expected carry_low: 0, got: %s", carry_low.ToSignedFeltString())
	}
	carry_high, err := idsManager.GetFelt("carry_high", vm)
	if err != nil {
		t.Errorf("failed with error: %s", err)
	}
	if carry_high.Cmp(FeltOne()) != 0 {
		t.Errorf("expected carry_high: 0, got: %s", carry_high.ToSignedFeltString())
	}
}

func TestUint256AddLowOnlyOk(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()

	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{
			"a": {
				NewMaybeRelocatableFeltFromUint64(2),
				NewMaybeRelocatableFeltFromUint64(3),
			},
			"b": {
				NewMaybeRelocatableFeltFromUint64(4),
				NewMaybeRelocatableFelt(FeltFromDecString("340282366920938463463374607431768211455")),
			},
			"carry_low": {nil},
		},
		vm,
	)
	hintData := any(HintData{
		Ids:  idsManager,
		Code: UINT256_ADD_LOW,
	})
	scopes := NewExecutionScopes()
	hintProcessor := CairoVmHintProcessor{}
	err := hintProcessor.ExecuteHint(vm, &hintData, nil, scopes)
	if err != nil {
		t.Errorf("failed with error: %s", err)
	}

	carry_low, err := idsManager.GetFelt("carry_low", vm)
	if err != nil {
		t.Errorf("failed with error: %s", err)
	}
	if carry_low.Cmp(FeltZero()) != 0 {
		t.Errorf("expected carry_low: 0, got: %s", carry_low.ToSignedFeltString())
	}
}

func TestUint256AddFailInsert(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()
	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{
			"a": {
				NewMaybeRelocatableFeltFromUint64(2),
				NewMaybeRelocatableFeltFromUint64(3),
			},
			"b": {
				NewMaybeRelocatableFeltFromUint64(4),
				NewMaybeRelocatableFeltFromUint64(2),
			},
			"carry_low": {NewMaybeRelocatableFeltFromUint64(2)},
		},
		vm,
	)
	hintData := any(HintData{
		Ids:  idsManager,
		Code: UINT256_ADD_LOW,
	})
	scopes := NewExecutionScopes()
	hintProcessor := CairoVmHintProcessor{}
	err := hintProcessor.ExecuteHint(vm, &hintData, nil, scopes)
	if err == nil {
		t.Errorf("should fail with error: ErrMemoryWriteOnce")
	}

}

func TestSplit64Ok(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()

	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{
			"a": {
				NewMaybeRelocatableFelt(FeltFromDecString("-3")),
			},
			"low":  {nil},
			"high": {nil},
		},
		vm,
	)

	hintData := any(HintData{
		Ids:  idsManager,
		Code: SPLIT_64,
	})
	scopes := NewExecutionScopes()
	hintProcessor := CairoVmHintProcessor{}
	err := hintProcessor.ExecuteHint(vm, &hintData, nil, scopes)
	if err != nil {
		t.Errorf("failed with error: %s", err)
	}

	low, err := idsManager.GetFelt("low", vm)
	if err != nil {
		t.Errorf("failed with error: %s", err)
	}
	expected_low := FeltFromDecString("-3").And(FeltOne().Shl(64).Sub(FeltOne()))
	if low != expected_low {
		t.Errorf("expected low: %d, got: %d", expected_low, low)
	}

	high, err := idsManager.GetFelt("high", vm)
	if err != nil {
		t.Errorf("failed with error: %s", err)
	}
	expected_high := FeltFromDecString("-3").Shr(64)
	if high != expected_high {
		t.Errorf("expected high: %d, got: %d", expected_high, high)
	}

}

func TestSplit64BigA(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()

	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{
			"a": {
				NewMaybeRelocatableFelt(FeltFromDecString("400066369019890261321163226850167045262")),
			},
			"low":  {nil},
			"high": {nil},
		},
		vm,
	)

	hintData := any(HintData{
		Ids:  idsManager,
		Code: SPLIT_64,
	})
	scopes := NewExecutionScopes()
	hintProcessor := CairoVmHintProcessor{}
	err := hintProcessor.ExecuteHint(vm, &hintData, nil, scopes)
	if err != nil {
		t.Errorf("failed with error: %s", err)
	}

	low, err := idsManager.GetFelt("low", vm)
	if err != nil {
		t.Errorf("failed with error: %s", err)
	}

	expected_low := FeltFromUint64(2279400676465785998)
	if low.Cmp(expected_low) != 0 {
		t.Errorf("expected low: %d, got: %d", expected_low, low)
	}
	high, err := idsManager.GetFelt("high", vm)
	if err != nil {
		t.Errorf("failed with error: %s", err)
	}
	expected_high := FeltFromDecString("21687641321487626429")
	if high.Cmp(expected_high) != 0 {
		t.Errorf("expected high: %d, got: %d", expected_high, high)
	}

}

func TestUint256SqrtOk(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()
	vm.Segments.AddSegment()
	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{
			"n": {
				NewMaybeRelocatableFelt(FeltFromUint64(17)),
				NewMaybeRelocatableFelt(FeltFromUint64(7)),
			},
			"root": {nil, nil},
		},
		vm,
	)
	hintData := any(HintData{
		Ids:  idsManager,
		Code: UINT256_SQRT,
	})
	scopes := NewExecutionScopes()
	hintProcessor := CairoVmHintProcessor{}
	err := hintProcessor.ExecuteHint(vm, &hintData, nil, scopes)
	if err != nil {
		t.Errorf("failed with error: %s", err)
	}

	expected_root, _ := new(big.Int).SetString("48805497317890012913", 10)
	root, err := idsManager.GetFelt("root", vm)
	if err != nil {
		t.Errorf("failed with error: %s", err)
	}
	if root != FeltFromBigInt(expected_root) {
		t.Errorf("failed, expected root: %d, got: %d", FeltFromBigInt(expected_root), root)
	}
}

func TestUint256SqrtKo(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()

	idsManager := SetupIdsForTest(map[string][]*MaybeRelocatable{
		"n": {
			NewMaybeRelocatableFelt(FeltZero()),
			NewMaybeRelocatableFelt(FeltFromDecString("340282366920938463463374607431768211458")),
		},
		"root": {nil},
	}, vm)

	hintData := any(HintData{Ids: idsManager, Code: UINT256_SQRT})
	hintProcessor := CairoVmHintProcessor{}
	err := hintProcessor.ExecuteHint(vm, &hintData, nil, NewExecutionScopes())
	expectedRoot := FeltFromDecString("340282366920938463463374607431768211456")
	if err.Error() != ErrRootOOR(expectedRoot.ToBigInt()).Error() {
		t.Errorf("failed with error: %s", err)
	}
}

func TestUint256SqrtFeltOk(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()
	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{
			"n": {
				NewMaybeRelocatableFelt(FeltFromUint64(879232)),
				NewMaybeRelocatableFelt(FeltFromUint64(135906)),
			},
			"root": {nil},
		},
		vm,
	)
	hintData := any(HintData{
		Ids:  idsManager,
		Code: UINT256_SQRT_FELT,
	})
	scopes := NewExecutionScopes()
	hintProcessor := CairoVmHintProcessor{}
	err := hintProcessor.ExecuteHint(vm, &hintData, nil, scopes)
	if err != nil {
		t.Errorf("failed with error: %s", err)
	}
	expected_root, _ := new(big.Int).SetString("6800471701195223914689", 10)
	expectedResult := FeltFromBigInt(expected_root)
	root, err := idsManager.GetFelt("root", vm)
	if err != nil {
		t.Errorf("failed with error: %s", err)
	}
	if root != expectedResult {
		t.Errorf("failed, expected root: %d, got: %d", expectedResult, root)
	}
}

func TestUint256SignedNNOkResultOne(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments = AddNSegments(vm.Segments, 5)
	ids := map[string][]*MaybeRelocatable{
		"a": {
			NewMaybeRelocatableFelt(FeltFromUint64(1)),
			NewMaybeRelocatableFelt(FeltFromUint64(1)),
		},
	}
	idsManager := SetupIdsForTest(ids, vm)
	hintData := any(HintData{
		Ids:  idsManager,
		Code: UINT256_SIGNED_NN,
	})
	hintProcessor := CairoVmHintProcessor{}
	err := hintProcessor.ExecuteHint(vm, &hintData, nil, nil)
	if err != nil {
		t.Errorf("failed with error: %s", err)
	}

	result, err := vm.Segments.Memory.GetFelt(vm.RunContext.Ap)
	if err != nil {
		t.Errorf("failed with error: %s", err)
	}

	if result != FeltOne() {
		t.Errorf("failed, expected result: %d, got: %d", FeltOne(), result)
	}
}

func TestUint256SignedNNOkResultZero(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()
	vm.RunContext.Ap = NewRelocatable(0, 5)
	ids := map[string][]*MaybeRelocatable{
		"a": {
			NewMaybeRelocatableFelt(FeltFromUint64(1)),
			NewMaybeRelocatableFelt(FeltFromDecString("-4")),
		},
	}
	idsManager := SetupIdsForTest(ids, vm)
	hintData := any(HintData{
		Ids:  idsManager,
		Code: UINT256_SIGNED_NN,
	})
	hintProcessor := CairoVmHintProcessor{}
	err := hintProcessor.ExecuteHint(vm, &hintData, nil, nil)
	if err != nil {
		t.Errorf("failed with error: %s", err)
	}

	result, err := vm.Segments.Memory.GetFelt(vm.RunContext.Ap)
	if err != nil {
		t.Errorf("failed with error: %s", err)
	}

	if result != FeltZero() {
		t.Errorf("failed, expected result: %d, got: %d", FeltZero(), result)
	}
}

func TestUint256SignedNNInvalidMemoryInser(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()
	vm.RunContext.Ap = NewRelocatable(0, 5)
	vm.Segments.Memory.Insert(vm.RunContext.Ap, NewMaybeRelocatableFeltFromUint64(10))
	ids := map[string][]*MaybeRelocatable{
		"a": {
			NewMaybeRelocatableFelt(FeltFromUint64(1)),
			NewMaybeRelocatableFelt(FeltFromUint64(1)),
		},
	}
	idsManager := SetupIdsForTest(ids, vm)
	idsManager.HintApTracking = parser.ApTrackingData{Group: 4, Offset: 5}
	hintData := any(HintData{
		Ids:  idsManager,
		Code: UINT256_SIGNED_NN,
	})
	hintProcessor := CairoVmHintProcessor{}
	err := hintProcessor.ExecuteHint(vm, &hintData, nil, nil)

	expectedErr := ErrMemoryWriteOnce(NewRelocatable(0, 5), *NewMaybeRelocatableFeltFromUint64(10), *NewMaybeRelocatableFelt(FeltOne()))
	if err.Error() != expectedErr.Error() {
		t.Errorf("should fail with error: %s", err)
	}
}

func TestUint256UnsignedDivRemOk(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()
	vm.Segments.AddSegment()

	ids := map[string][]*MaybeRelocatable{
		"a": {
			NewMaybeRelocatableFeltFromUint64(89),
			NewMaybeRelocatableFeltFromUint64(72),
		},
		"div": {
			NewMaybeRelocatableFeltFromUint64(3),
			NewMaybeRelocatableFeltFromUint64(7),
		},
		"quotient":  {nil, nil},
		"remainder": {nil, nil},
	}
	idsManager := SetupIdsForTest(ids, vm)
	hintData := any(HintData{
		Ids:  idsManager,
		Code: UINT256_UNSIGNED_DIV_REM,
	})
	hintProcessor := CairoVmHintProcessor{}
	err := hintProcessor.ExecuteHint(vm, &hintData, nil, nil)
	if err != nil {
		t.Errorf("failed with error: %s", err)
	}

	quotient, err := idsManager.GetUint256("quotient", vm)
	if err != nil {
		t.Errorf("failed with error: %s", err)
	}

	expectedQuotient := Uint256{Low: FeltFromUint(10), High: FeltFromUint(0)}
	if quotient != expectedQuotient {
		t.Errorf("expected quotient: %s, got: %s", expectedQuotient.ToString(), quotient.ToString())
	}

	remainder, err := idsManager.GetUint256("remainder", vm)
	if err != nil {
		t.Errorf("failed with error: %s", err)
	}

	expectedRemainder := Uint256{Low: FeltFromUint(59), High: FeltFromUint(2)}
	if remainder != expectedRemainder {
		t.Errorf("expected remainder: %s, got: %s", expectedRemainder.ToString(), remainder.ToString())
	}

}

func TestUint256UnsignedDivRemInvalidMemoryInsert(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()
	vm.Segments.AddSegment()

	ids := map[string][]*MaybeRelocatable{
		"a": {
			NewMaybeRelocatableFeltFromUint64(89),
			NewMaybeRelocatableFeltFromUint64(72),
		},
		"div": {
			NewMaybeRelocatableFeltFromUint64(3),
			NewMaybeRelocatableFeltFromUint64(7),
		},
		"quotient":  {NewMaybeRelocatableFeltFromUint64(2), NewMaybeRelocatableFelt(FeltZero())},
		"remainder": {nil, nil},
	}
	idsManager := SetupIdsForTest(ids, vm)
	hintData := any(HintData{
		Ids:  idsManager,
		Code: UINT256_UNSIGNED_DIV_REM,
	})
	hintProcessor := CairoVmHintProcessor{}
	err := hintProcessor.ExecuteHint(vm, &hintData, nil, nil)
	if err == nil {
		t.Errorf("this test should fail")
	}
}

func TestUint256ExpandedUnsignedDivRemOk(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()
	vm.Segments.AddSegment()

	ids := map[string][]*MaybeRelocatable{
		"a": {
			NewMaybeRelocatableFeltFromUint64(89),
			NewMaybeRelocatableFeltFromUint64(72),
		},

		"div": {
			NewMaybeRelocatableFelt(FeltFromDecString("55340232221128654848")),
			NewMaybeRelocatableFeltFromUint64(3),
			NewMaybeRelocatableFelt(FeltFromDecString("129127208515966861312")),
			NewMaybeRelocatableFeltFromUint64(7),
		},
		"quotient":  {nil, nil},
		"remainder": {nil, nil},
	}
	idsManager := SetupIdsForTest(ids, vm)
	hintData := any(HintData{
		Ids:  idsManager,
		Code: UINT256_EXPANDED_UNSIGNED_DIV_REM,
	})
	hintProcessor := CairoVmHintProcessor{}
	err := hintProcessor.ExecuteHint(vm, &hintData, nil, nil)
	if err != nil {
		t.Errorf("failed with error: %s", err)
	}

	quotient, err := idsManager.GetUint256("quotient", vm)
	if err != nil {
		t.Errorf("failed with error: %s", err)
	}

	expectedQuotient := Uint256{Low: FeltFromUint(10), High: FeltFromUint(0)}
	if quotient != expectedQuotient {
		t.Errorf("expected quotient: %s, got: %s", expectedQuotient.ToString(), quotient.ToString())
	}

	remainder, err := idsManager.GetUint256("remainder", vm)
	if err != nil {
		t.Errorf("failed with error: %s", err)
	}

	expectedRemainder := Uint256{Low: FeltFromUint(59), High: FeltFromUint(2)}
	if remainder != expectedRemainder {
		t.Errorf("expected remainder: %s, got: %s", expectedRemainder.ToString(), remainder.ToString())
	}

}

func TestUint256MulDivOk(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()
	vm.Segments.AddSegment()

	ids := map[string][]*MaybeRelocatable{
		"a": {
			NewMaybeRelocatableFeltFromUint64(89),
			NewMaybeRelocatableFeltFromUint64(72),
		},
		"b": {
			NewMaybeRelocatableFeltFromUint64(3),
			NewMaybeRelocatableFeltFromUint64(7),
		},
		"div": {
			NewMaybeRelocatableFeltFromUint64(107),
			NewMaybeRelocatableFeltFromUint64(114),
		},
		"quotient_low":  {nil, nil},
		"quotient_high": {nil, nil},
		"remainder":     {nil, nil},
	}
	idsManager := SetupIdsForTest(ids, vm)
	hintData := any(HintData{
		Ids:  idsManager,
		Code: UINT256_MUL_DIV_MOD,
	})
	hintProcessor := CairoVmHintProcessor{}
	err := hintProcessor.ExecuteHint(vm, &hintData, nil, nil)
	if err != nil {
		t.Errorf("failed with error: %s", err)
	}

	quotientLow, err := idsManager.GetUint256("quotient_low", vm)
	if err != nil {
		t.Errorf("failed with error: %s", err)
	}
	expectedQuotientLow := Uint256{Low: FeltFromDecString("143276786071974089879315624181797141668"), High: FeltFromUint(4)}
	if !quotientLow.IsEqual(expectedQuotientLow) {
		t.Errorf("expected quotient_low: %s, got: %s", expectedQuotientLow.ToString(), quotientLow.ToString())
	}

	quotientHigh, err := idsManager.GetUint256("quotient_high", vm)
	if err != nil {
		t.Errorf("failed with error: %s", err)
	}
	expectedQuotientHigh := Uint256{Low: FeltFromUint(0), High: FeltFromUint(0)}
	if !quotientHigh.IsEqual(expectedQuotientHigh) {
		t.Errorf("expected quotient_high: %s, got: %s", expectedQuotientHigh.ToString(), quotientHigh.ToString())
	}

	remainder, err := idsManager.GetUint256("remainder", vm)
	if err != nil {
		t.Errorf("failed with error: %s", err)
	}
	expectedRemainder := Uint256{Low: FeltFromDecString("322372768661941702228460154409043568767"), High: FeltFromUint(101)}
	if !remainder.IsEqual(expectedRemainder) {
		t.Errorf("expected remainder: %s, got: %s", expectedRemainder.ToString(), remainder.ToString())
	}
}
