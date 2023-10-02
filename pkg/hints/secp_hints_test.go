package hints_test

import (
	"math/big"
	"testing"

	. "github.com/lambdaclass/cairo-vm.go/pkg/hints"
	. "github.com/lambdaclass/cairo-vm.go/pkg/hints/hint_codes"
	. "github.com/lambdaclass/cairo-vm.go/pkg/hints/hint_utils"
	. "github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
	. "github.com/lambdaclass/cairo-vm.go/pkg/types"
	. "github.com/lambdaclass/cairo-vm.go/pkg/vm"
	. "github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
)

func TestReduceV1(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()
	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{
			"x": {
				NewMaybeRelocatableFelt(FeltFromDecString("6")),
				NewMaybeRelocatableFelt(FeltFromDecString("6")),
				NewMaybeRelocatableFelt(FeltFromDecString("6")),
			},
		},
		vm,
	)
	scopes := NewExecutionScopes()

	hintProcessor := CairoVmHintProcessor{}
	hintData := any(HintData{
		Ids:  idsManager,
		Code: REDUCE_V1,
	})

	err := hintProcessor.ExecuteHint(vm, &hintData, nil, scopes)
	if err != nil {
		t.Errorf("REDUCE_V1 hint failed with error %s", err)
	}
	// Checked scope variables
	CheckScopeVar[big.Int]("SECP_P", SECP_P(), scopes, t)

	valueUnpacked := Uint384{Limbs: []Felt{FeltFromUint(6), FeltFromUint(6), FeltFromUint(6)}}
	CheckScopeVar[big.Int]("value", valueUnpacked.Pack86(), scopes, t)
}

func TestReduceV2(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()
	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{
			"x": {
				NewMaybeRelocatableFelt(FeltFromDecString("6")),
				NewMaybeRelocatableFelt(FeltFromDecString("6")),
				NewMaybeRelocatableFelt(FeltFromDecString("6")),
			},
		},
		vm,
	)
	scopes := NewExecutionScopes()
	scopes.AssignOrUpdateVariable("SECP_P", SECP_P())

	hintProcessor := CairoVmHintProcessor{}
	hintData := any(HintData{
		Ids:  idsManager,
		Code: REDUCE_V2,
	})

	err := hintProcessor.ExecuteHint(vm, &hintData, nil, scopes)
	if err != nil {
		t.Errorf("REDUCE_V2 hint failed with error %s", err)
	}

	valueUnpacked := Uint384{Limbs: []Felt{FeltFromUint(6), FeltFromUint(6), FeltFromUint(6)}}
	CheckScopeVar[big.Int]("value", valueUnpacked.Pack86(), scopes, t)
}

func TestReduceED(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()
	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{
			"x": {
				NewMaybeRelocatableFelt(FeltFromDecString("6")),
				NewMaybeRelocatableFelt(FeltFromDecString("6")),
				NewMaybeRelocatableFelt(FeltFromDecString("6")),
			},
		},
		vm,
	)
	scopes := NewExecutionScopes()

	hintProcessor := CairoVmHintProcessor{}
	hintData := any(HintData{
		Ids:  idsManager,
		Code: REDUCE_ED25519,
	})

	err := hintProcessor.ExecuteHint(vm, &hintData, nil, scopes)
	if err != nil {
		t.Errorf("REDUCE_ED25519 hint failed with error %s", err)
	}
	// Checked scope variables
	CheckScopeVar[big.Int]("SECP_P", SECP_P_V2(), scopes, t)

	valueUnpacked := Uint384{Limbs: []Felt{FeltFromUint(6), FeltFromUint(6), FeltFromUint(6)}}
	CheckScopeVar[big.Int]("value", valueUnpacked.Pack86(), scopes, t)
}

func TestVerifyZeroV1(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()
	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{
			"val": {
				NewMaybeRelocatableFelt(FeltFromDecString("0")),
				NewMaybeRelocatableFelt(FeltFromDecString("0")),
				NewMaybeRelocatableFelt(FeltFromDecString("0")),
			},
			"q": {nil, nil, nil},
		},
		vm,
	)
	scopes := NewExecutionScopes()

	hintProcessor := CairoVmHintProcessor{}
	hintData := any(HintData{
		Ids:  idsManager,
		Code: VERIFY_ZERO_V1,
	})

	err := hintProcessor.ExecuteHint(vm, &hintData, nil, scopes)
	if err != nil {
		t.Errorf("VERIFY_ZERO_V1 hint failed with error %s", err)
	}
	// Check scope variables
	CheckScopeVar[big.Int]("SECP_P", SECP_P(), scopes, t)
	// Check ids variables
	expectedQ := FeltZero()
	idsQ, err := idsManager.GetFelt("q", vm)
	if err != nil || expectedQ.Cmp(idsQ) != 0 {
		t.Error("Wrong/No ids.q")
	}
}

func TestVerifyZeroV2(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()
	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{
			"val": {
				NewMaybeRelocatableFelt(FeltFromDecString("0")),
				NewMaybeRelocatableFelt(FeltFromDecString("0")),
				NewMaybeRelocatableFelt(FeltFromDecString("0")),
			},
			"q": {nil, nil, nil},
		},
		vm,
	)
	scopes := NewExecutionScopes()

	hintProcessor := CairoVmHintProcessor{}
	hintData := any(HintData{
		Ids:  idsManager,
		Code: VERIFY_ZERO_V2,
	})

	err := hintProcessor.ExecuteHint(vm, &hintData, nil, scopes)
	if err != nil {
		t.Errorf("VERIFY_ZERO_V2 hint failed with error %s", err)
	}
	// Check scope variables
	CheckScopeVar[big.Int]("SECP_P", SECP_P(), scopes, t)
	// Check ids variables
	expectedQ := FeltZero()
	idsQ, err := idsManager.GetFelt("q", vm)
	if err != nil || expectedQ.Cmp(idsQ) != 0 {
		t.Error("Wrong/No ids.q")
	}
}

func TestVerifyZeroV3(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()
	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{
			"val": {
				NewMaybeRelocatableFelt(FeltFromDecString("0")),
				NewMaybeRelocatableFelt(FeltFromDecString("0")),
				NewMaybeRelocatableFelt(FeltFromDecString("0")),
			},
			"q": {nil, nil, nil},
		},
		vm,
	)
	scopes := NewExecutionScopes()

	hintProcessor := CairoVmHintProcessor{}
	hintData := any(HintData{
		Ids:  idsManager,
		Code: VERIFY_ZERO_V3,
	})

	err := hintProcessor.ExecuteHint(vm, &hintData, nil, scopes)
	if err != nil {
		t.Errorf("VERIFY_ZERO_V3 hint failed with error %s", err)
	}
	// Check scope variables
	CheckScopeVar[big.Int]("SECP_P", SECP_P_V2(), scopes, t)
	// Check ids variables
	expectedQ := FeltZero()
	idsQ, err := idsManager.GetFelt("q", vm)
	if err != nil || expectedQ.Cmp(idsQ) != 0 {
		t.Error("Wrong/No ids.q")
	}
}

func TestBigIntTo256(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()
	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{
			"x": {
				NewMaybeRelocatableFelt(FeltFromDecString("1")),
				NewMaybeRelocatableFelt(FeltFromDecString("0")),
				NewMaybeRelocatableFelt(FeltFromDecString("0")),
			},
			"low": {nil},
		},
		vm,
	)
	constants := SetupConstantsForTest(map[string]Felt{
		"BASE": FeltFromUint(1).Shl(86),
	}, &idsManager)

	hintProcessor := CairoVmHintProcessor{}
	hintData := any(HintData{
		Ids:  idsManager,
		Code: BIGINT_TO_UINT256,
	})

	err := hintProcessor.ExecuteHint(vm, &hintData, &constants, nil)
	if err != nil {
		t.Errorf("BIGINT_TO_UINT256 hint failed with error %s", err)
	}
	// Check ids.low
	low, err := idsManager.GetFelt("low", vm)
	if err != nil || low != FeltOne() {
		t.Error("Wrong/No ids.low")
	}
}

func TestIsZeroTrue(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()
	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{
			"x": {
				NewMaybeRelocatableFelt(FeltZero()),
			},
		},
		vm,
	)
	// Advance ap to avoid clashes
	vm.RunContext.Ap.Offset += 1

	hintProcessor := CairoVmHintProcessor{}
	hintData := any(HintData{
		Ids:  idsManager,
		Code: IS_ZERO_NONDET,
	})

	err := hintProcessor.ExecuteHint(vm, &hintData, nil, nil)
	if err != nil {
		t.Errorf("IS_ZERO_NONDET hint failed with error %s", err)
	}
	// Check memory[ap]
	val, err := vm.Segments.Memory.GetFelt(vm.RunContext.Ap)
	if err != nil || val != FeltOne() {
		t.Error("Wrong/No value inserted into ap")
	}
}

func TestIsZeroFalse(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()
	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{
			"x": {
				NewMaybeRelocatableFelt(FeltOne()),
			},
		},
		vm,
	)
	// Advance ap to avoid clashes
	vm.RunContext.Ap.Offset += 1

	hintProcessor := CairoVmHintProcessor{}
	hintData := any(HintData{
		Ids:  idsManager,
		Code: IS_ZERO_NONDET,
	})

	err := hintProcessor.ExecuteHint(vm, &hintData, nil, nil)
	if err != nil {
		t.Errorf("IS_ZERO_NONDET hint failed with error %s", err)
	}
	// Check memory[ap]
	val, err := vm.Segments.Memory.GetFelt(vm.RunContext.Ap)
	if err != nil || val != FeltZero() {
		t.Error("Wrong/No value inserted into ap")
	}
}

func TestIsZeroPack(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()
	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{
			"x": {
				NewMaybeRelocatableFelt(FeltFromUint(6)),
				NewMaybeRelocatableFelt(FeltFromUint(6)),
				NewMaybeRelocatableFelt(FeltFromUint(6)),
			},
		},
		vm,
	)

	hintProcessor := CairoVmHintProcessor{}
	hintData := any(HintData{
		Ids:  idsManager,
		Code: IS_ZERO_PACK_V1,
	})
	scopes := NewExecutionScopes()

	err := hintProcessor.ExecuteHint(vm, &hintData, nil, scopes)
	if err != nil {
		t.Errorf("IS_ZERO_PACK_V1 hint failed with error %s", err)
	}
	// Check scope vars
	CheckScopeVar[big.Int]("SECP_P", SECP_P(), scopes, t)

	xUnpacked := Uint384{Limbs: []Felt{FeltFromUint(6), FeltFromUint(6), FeltFromUint(6)}}
	CheckScopeVar[big.Int]("x", xUnpacked.Pack86(), scopes, t)
}

func TestIsZeroAssignScopeVars(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()
	idsManager := SetupIdsForTest(map[string][]*MaybeRelocatable{}, vm)

	hintProcessor := CairoVmHintProcessor{}
	hintData := any(HintData{
		Ids:  idsManager,
		Code: IS_ZERO_ASSIGN_SCOPE_VARS,
	})
	scopes := NewExecutionScopes()
	x, _ := new(big.Int).SetString("52621538839140286024584685587354966255185961783273479086367", 10)
	scopes.AssignOrUpdateVariable("x", *x)

	err := hintProcessor.ExecuteHint(vm, &hintData, nil, scopes)
	if err != nil {
		t.Errorf("IS_ZERO_ASSIGN_SCOPE_VARS hint failed with error %s", err)
	}
	// Check scope vars
	CheckScopeVar[big.Int]("SECP_P", SECP_P(), scopes, t)

	value, _ := new(big.Int).SetString("19429627790501903254364315669614485084365347064625983303617500144471999752609", 10)
	CheckScopeVar[big.Int]("value", *value, scopes, t)
	CheckScopeVar[big.Int]("x_inv", *value, scopes, t)
}
