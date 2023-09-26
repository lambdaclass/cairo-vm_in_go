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

/*
flag := (1 << 128)

	a := {
		a.low: flag - 5
		a.high = ...
	}

	b := {
		b.low: 4
		b.high = ...
	}

a.low + b.low < flag -> carryLow = 0
*/
func TestUint256AddCarryLow0(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()

	flag := FeltOne().Shl(128)

	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{
			"a": {
				NewMaybeRelocatableFelt(flag.Sub(FeltFromUint64(5))),
				nil,
			},
			"b": {
				NewMaybeRelocatableFelt(FeltFromUint64(4)),
				nil,
			},
			"carry_low": {
				nil,
				nil,
			},
			"carry_high": {
				nil,
				nil,
			},
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

	carry_low, err := idsManager.GetStructFieldFelt("carry_low", 0, vm)
	if err != nil {
		t.Errorf("failed with error: %s", err)
	}
	if carry_low != FeltZero() {
		t.Errorf("expected carry_low: 0, got: %s", carry_low.ToSignedFeltString())
	}
}

/*
flag := (1 << 128)

	a := {
		a.low: flag
		a.high = ...
	}

	b := {
		b.low: 0
		b.high = ...
	}

a.low + b.low >= flag -> carryLow = 1
*/
func TestUint256AddCarryLow1(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()

	flag := FeltOne().Shl(128)

	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{
			"a": {
				NewMaybeRelocatableFelt(flag),
				nil,
			},
			"b": {
				NewMaybeRelocatableFelt(FeltZero()),
				nil,
			},
			"carry_low": {
				nil,
				nil,
			},
			"carry_high": {
				nil,
				nil,
			},
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

	carry_low, err := idsManager.GetStructFieldFelt("carry_low", 0, vm)
	if err != nil {
		t.Errorf("failed with error: %s", err)
	}
	if carry_low != FeltOne() {
		t.Errorf("expected carry_low: 1, got: %s", carry_low.ToSignedFeltString())
	}
}

/*
flag := (1 << 128)

	a := {
		a.low: 0
		a.high = flag / 2
	}

	b := {
		b.low: 0
		b.high = a.high - 1
	}

a.low + b.low < flag -> carryLow = 0
a.high + b.high + carryLow < flag -> carry_high = 0
*/
func TestUint256AddCarryHigh0(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()

	flag := FeltOne().Shl(128)
	aHigh := flag.Div(FeltFromUint64(2))
	bHigh := aHigh.Sub(FeltFromUint64(1))

	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{
			"a": {
				NewMaybeRelocatableFelt(FeltZero()),
				NewMaybeRelocatableFelt(aHigh),
			},
			"b": {
				NewMaybeRelocatableFelt(FeltZero()),
				NewMaybeRelocatableFelt(bHigh),
			},
			"carry_low": {
				nil,
				nil,
			},
			"carry_high": {
				nil,
				nil,
			},
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

	carry_high, err := idsManager.GetStructFieldFelt("carry_high", 0, vm)
	if err != nil {
		t.Errorf("failed with error: %s", err)
	}
	if carry_high != FeltZero() {
		t.Errorf("expected carry_low: 0, got: %s", carry_high.ToSignedFeltString())
	}
}

/*
		flag := (1 << 128)
		a := {
			a.low: flag
			a.high = flag / 2
		}

		b := {
			b.low: 0
			b.high = a.high - 1
		}

		a.low + b.low >= flag -> carryLow = 1
	    a.high + b.high + carryLow > flag -> carry_high = 1
*/
func TestUint256AddCarryHigh1(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()

	flag := FeltOne().Shl(128)
	aHigh := flag.Div(FeltFromUint64(2))
	bHigh := aHigh.Sub(FeltFromUint64(1))

	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{
			"a": {
				NewMaybeRelocatableFelt(flag),
				NewMaybeRelocatableFelt(aHigh),
			},
			"b": {
				NewMaybeRelocatableFelt(FeltZero()),
				NewMaybeRelocatableFelt(bHigh),
			},
			"carry_low": {
				nil,
				nil,
			},
			"carry_high": {
				nil,
				nil,
			},
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
	carry_high, err := idsManager.GetStructFieldFelt("carry_high", 0, vm)
	if err != nil {
		t.Errorf("failed with error: %s", err)
	}
	if carry_high != FeltOne() {
		t.Errorf("expected carry_low: 1, got: %s", carry_high.ToSignedFeltString())
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
			"root": {nil},
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
	idsManager.HintApTracking = parser.ApTrackingData{Group: 4, Offset: 5}
	hintData := any(HintData{
		Ids:  idsManager,
		Code: UINT256_SIGNED_NN,
	})
	hintProcessor := CairoVmHintProcessor{}
	err := hintProcessor.ExecuteHint(vm, &hintData, nil, nil)
	if err != nil {
		t.Errorf("failed with error: %s", err)
	}

	result, err := vm.Segments.Memory.GetFelt(NewRelocatable(4, 5))
	if err != nil {
		t.Errorf("failed with error: %s", err)
	}

	if result != FeltOne() {
		t.Errorf("failed, expected result: %d, got: %d", FeltOne(), result)
	}
}

func TestUint256SignedNNOkResultZero(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments = AddNSegments(vm.Segments, 5)
	ids := map[string][]*MaybeRelocatable{
		"a": {
			NewMaybeRelocatableFelt(FeltFromUint64(1)),
			NewMaybeRelocatableFelt(FeltFromDecString("-4")),
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
	if err != nil {
		t.Errorf("failed with error: %s", err)
	}

	result, err := vm.Segments.Memory.GetFelt(NewRelocatable(4, 5))
	if err != nil {
		t.Errorf("failed with error: %s", err)
	}

	if result != FeltZero() {
		t.Errorf("failed, expected result: %d, got: %d", FeltOne(), result)
	}
}

func TestUint256SignedNNInvalidMemoryInser(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments = AddNSegments(vm.Segments, 5)
	err := vm.Segments.Memory.Insert(NewRelocatable(4, 5), NewMaybeRelocatableFeltFromUint64(10))
	if err != nil {
		t.Errorf("failed with error: %s", err)
	}
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
	err = hintProcessor.ExecuteHint(vm, &hintData, nil, nil)
	expectedErr := ErrMemoryWriteOnce(NewRelocatable(4, 5), *NewMaybeRelocatableFeltFromUint64(10), *NewMaybeRelocatableFelt(FeltOne()))
	if err.Error() != expectedErr.Error() {
		t.Errorf("should fail with error: %s", err)
	}
}

func TestUint256UnsignedDivRemOk(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()
	vm.Segments.AddSegment()

	// add div low
	err := vm.Segments.Memory.Insert(NewRelocatable(1, 6), NewMaybeRelocatableFelt(FeltFromUint64(3)))
	if err != nil {
		t.Errorf("failed with error: %s", err)
	}
	// add div high
	err = vm.Segments.Memory.Insert(NewRelocatable(1, 7), NewMaybeRelocatableFeltFromUint64(7))
	if err != nil {
		t.Errorf("failed with error: %s", err)
	}
	ids := map[string][]*MaybeRelocatable{
		"a": {
			NewMaybeRelocatableFeltFromUint64(89),
			NewMaybeRelocatableFeltFromUint64(72),
		},
		"div":       {NewMaybeRelocatableRelocatable(NewRelocatable(1, 6))},
		"quotient":  {nil, nil},
		"remainder": {nil, nil},
	}
	idsManager := SetupIdsForTest(ids, vm)
	hintData := any(HintData{
		Ids:  idsManager,
		Code: UINT256_UNSIGNED_DIV_REM,
	})
	hintProcessor := CairoVmHintProcessor{}
	err = hintProcessor.ExecuteHint(vm, &hintData, nil, nil)
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

	// add div low
	err := vm.Segments.Memory.Insert(NewRelocatable(1, 6), NewMaybeRelocatableFelt(FeltFromUint64(3)))
	if err != nil {
		t.Errorf("failed with error: %s", err)
	}
	// add div high
	err = vm.Segments.Memory.Insert(NewRelocatable(1, 7), NewMaybeRelocatableFeltFromUint64(7))
	if err != nil {
		t.Errorf("failed with error: %s", err)
	}
	// add hardcoded value on quotient.low
	err = vm.Segments.Memory.Insert(NewRelocatable(0, 3), NewMaybeRelocatableFeltFromUint64(8))
	if err != nil {
		t.Errorf("failed with error: %s", err)
	}
	ids := map[string][]*MaybeRelocatable{
		"a": {
			NewMaybeRelocatableFeltFromUint64(89),
			NewMaybeRelocatableFeltFromUint64(72),
		},
		"div":       {NewMaybeRelocatableRelocatable(NewRelocatable(1, 6))},
		"quotient":  {nil, nil},
		"remainder": {nil, nil},
	}
	idsManager := SetupIdsForTest(ids, vm)
	hintData := any(HintData{
		Ids:  idsManager,
		Code: UINT256_UNSIGNED_DIV_REM,
	})
	hintProcessor := CairoVmHintProcessor{}
	err = hintProcessor.ExecuteHint(vm, &hintData, nil, nil)
	expectedErr := ErrMemoryWriteOnce(NewRelocatable(0, 3), *NewMaybeRelocatableFeltFromUint64(8), *NewMaybeRelocatableFeltFromUint64(10))
	if err.Error() != expectedErr.Error() {
		t.Errorf("failed with error: %s", err)
	}
}
