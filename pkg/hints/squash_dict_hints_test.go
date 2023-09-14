package hints_test

import (
	"reflect"
	"testing"

	. "github.com/lambdaclass/cairo-vm.go/pkg/hints"
	. "github.com/lambdaclass/cairo-vm.go/pkg/hints/hint_utils"
	. "github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
	"github.com/lambdaclass/cairo-vm.go/pkg/types"
	. "github.com/lambdaclass/cairo-vm.go/pkg/vm"
	. "github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
)

func TestSquashDictValidOneKeyDictNoMaxSizeBigKeys(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()
	vm.Segments.AddSegment()
	vm.Segments.AddSegment()
	scopes := types.NewExecutionScopes()
	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{
			"dict_accesses": {NewMaybeRelocatableRelocatable(NewRelocatable(2, 0))},
			"big_keys":      {nil},
			"first_key":     {nil},
			"ptr_diff":      {NewMaybeRelocatableFelt(FeltFromUint64(6))},
			"n_accesses":    {NewMaybeRelocatableFelt(FeltFromUint64(2))},
		},
		vm,
	)
	// Insert dict into memory
	// Dict = {(prime - 1): (1,1), (prime - 1): (1,2)}
	vm.Segments.Memory.Insert(NewRelocatable(2, 0), NewMaybeRelocatableFelt(FeltFromDecString("-1")))
	vm.Segments.Memory.Insert(NewRelocatable(2, 1), NewMaybeRelocatableFelt(FeltOne()))
	vm.Segments.Memory.Insert(NewRelocatable(2, 2), NewMaybeRelocatableFelt(FeltOne()))
	vm.Segments.Memory.Insert(NewRelocatable(2, 3), NewMaybeRelocatableFelt(FeltFromDecString("-1")))
	vm.Segments.Memory.Insert(NewRelocatable(2, 4), NewMaybeRelocatableFelt(FeltOne()))
	vm.Segments.Memory.Insert(NewRelocatable(2, 5), NewMaybeRelocatableFelt(FeltFromUint64(2)))
	hintProcessor := CairoVmHintProcessor{}
	hintData := any(HintData{
		Ids:  idsManager,
		Code: SQUASH_DICT,
	})
	err := hintProcessor.ExecuteHint(vm, &hintData, nil, scopes)
	if err != nil {
		t.Errorf("SQUASH_DICT hint failed with error: %s", err)
	}
	// Check scope
	accessIndicesAny, err := scopes.Get("access_indices")
	accessIndices := accessIndicesAny.(map[MaybeRelocatable][]int)
	// expect access_indices = {prime -1: [0, 1]}
	expectedAccessIndices := map[MaybeRelocatable][]int{
		*NewMaybeRelocatableFelt(FeltFromDecString("-1")): {0, 1},
	}
	if !reflect.DeepEqual(accessIndices, expectedAccessIndices) {
		t.Errorf("SQUASH_DICT wrong access_indices.\n Expected %v, got: %v", expectedAccessIndices, accessIndices)
	}

}
