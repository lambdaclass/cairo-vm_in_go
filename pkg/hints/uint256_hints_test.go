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
