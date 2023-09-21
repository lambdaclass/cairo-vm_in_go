package hints_test

import (
	"math/big"
	"testing"

	. "github.com/lambdaclass/cairo-vm.go/pkg/hints"
	. "github.com/lambdaclass/cairo-vm.go/pkg/hints/hint_codes"
	. "github.com/lambdaclass/cairo-vm.go/pkg/hints/hint_utils"
	"github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
	. "github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
	. "github.com/lambdaclass/cairo-vm.go/pkg/types"
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
				NewMaybeRelocatableFelt(FeltFromUint64(^uint64(0))),
				NewMaybeRelocatableFelt(FeltFromUint64(^uint64(0))),
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

	l := new(big.Int).SetUint64(^uint64(0))
	h := new(big.Int).Lsh(l, 128)
	expectedRoot := new(big.Int).Sqrt(new(big.Int).Add(l, h))

	expectedResult := lambdaworks.Uint256{Low: FeltFromBigInt(expectedRoot), High: FeltZero()}

	root, err := idsManager.GetUint256("root", vm)
	if err != nil {
		t.Errorf("failed with error: %s", err)
	}
	if root != expectedResult {
		t.Errorf("failed, expected root: %d, got: %d", expectedResult, root)
	}
}
