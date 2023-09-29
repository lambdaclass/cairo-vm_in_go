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
	// SECP_P
	secpP, err := FetchScopeVar[big.Int]("SECP_P", scopes)
	expectedSecpP := SECP_P()
	if err != nil || secpP.Cmp(&expectedSecpP) != 0 {
		t.Errorf("Wrong/No scope var SECP_P")
	}

	value, err := FetchScopeVar[big.Int]("value", scopes)
	valueUnpacked := Uint384{Limbs: []Felt{FeltFromUint(6), FeltFromUint(6), FeltFromUint(6)}}
	expectedvalue := valueUnpacked.Pack86()
	if err != nil || value.Cmp(&expectedvalue) != 0 {
		t.Errorf("Wrong/No scope var value")
	}
}
