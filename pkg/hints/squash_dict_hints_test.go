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
	keysAny, err := scopes.Get("keys")
	keys := keysAny.([]MaybeRelocatable)
	expectedKeys := []MaybeRelocatable{}
	if !reflect.DeepEqual(keys, expectedKeys) {
		t.Errorf("SQUASH_DICT wrong keys.\n Expected %v, got: %v", expectedKeys, keys)
	}
	keyAny, err := scopes.Get("key")
	key := keyAny.(MaybeRelocatable)
	expectedKey := *NewMaybeRelocatableFelt(FeltFromDecString("-1"))
	if !reflect.DeepEqual(key, expectedKey) {
		t.Errorf("SQUASH_DICT wrong key.\n Expected %v, got: %v", expectedKey, keys)
	}
}

func TestSquashDictValidOneKeyDictWithMaxSize(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()
	vm.Segments.AddSegment()
	vm.Segments.AddSegment()
	scopes := types.NewExecutionScopes()
	scopes.AssignOrUpdateVariable("__squash_dict_max_size", 12)
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
	vm.Segments.Memory.Insert(NewRelocatable(2, 0), NewMaybeRelocatableFelt(FeltOne()))
	vm.Segments.Memory.Insert(NewRelocatable(2, 1), NewMaybeRelocatableFelt(FeltOne()))
	vm.Segments.Memory.Insert(NewRelocatable(2, 2), NewMaybeRelocatableFelt(FeltOne()))
	vm.Segments.Memory.Insert(NewRelocatable(2, 3), NewMaybeRelocatableFelt(FeltOne()))
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
		*NewMaybeRelocatableFelt(FeltOne()): {0, 1},
	}
	if !reflect.DeepEqual(accessIndices, expectedAccessIndices) {
		t.Errorf("SQUASH_DICT wrong access_indices.\n Expected %v, got: %v", expectedAccessIndices, accessIndices)
	}
	keysAny, err := scopes.Get("keys")
	keys := keysAny.([]MaybeRelocatable)
	expectedKeys := []MaybeRelocatable{}
	if !reflect.DeepEqual(keys, expectedKeys) {
		t.Errorf("SQUASH_DICT wrong keys.\n Expected %v, got: %v", expectedKeys, keys)
	}
	keyAny, err := scopes.Get("key")
	key := keyAny.(MaybeRelocatable)
	expectedKey := *NewMaybeRelocatableFelt(FeltOne())
	if !reflect.DeepEqual(key, expectedKey) {
		t.Errorf("SQUASH_DICT wrong key.\n Expected %v, got: %v", expectedKey, keys)
	}
}

func TestSquashDictValidTwoKeyDictNoMaxSizeBigKeys(t *testing.T) {
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
			"n_accesses":    {NewMaybeRelocatableFelt(FeltFromUint64(4))},
		},
		vm,
	)
	// Insert dict into memory
	//Dict = {1: (1,1), 1: (1,2), 2: (10,10), 2: (10,20)}
	vm.Segments.Memory.Insert(NewRelocatable(2, 0), NewMaybeRelocatableFelt(FeltOne()))
	vm.Segments.Memory.Insert(NewRelocatable(2, 1), NewMaybeRelocatableFelt(FeltOne()))
	vm.Segments.Memory.Insert(NewRelocatable(2, 2), NewMaybeRelocatableFelt(FeltOne()))

	vm.Segments.Memory.Insert(NewRelocatable(2, 3), NewMaybeRelocatableFelt(FeltOne()))
	vm.Segments.Memory.Insert(NewRelocatable(2, 4), NewMaybeRelocatableFelt(FeltOne()))
	vm.Segments.Memory.Insert(NewRelocatable(2, 5), NewMaybeRelocatableFelt(FeltFromUint64(2)))

	vm.Segments.Memory.Insert(NewRelocatable(2, 6), NewMaybeRelocatableFelt(FeltFromUint64(2)))
	vm.Segments.Memory.Insert(NewRelocatable(2, 7), NewMaybeRelocatableFelt(FeltFromUint64(10)))
	vm.Segments.Memory.Insert(NewRelocatable(2, 8), NewMaybeRelocatableFelt(FeltFromUint64(10)))

	vm.Segments.Memory.Insert(NewRelocatable(2, 9), NewMaybeRelocatableFelt(FeltFromUint64(2)))
	vm.Segments.Memory.Insert(NewRelocatable(2, 10), NewMaybeRelocatableFelt(FeltFromUint64(10)))
	vm.Segments.Memory.Insert(NewRelocatable(2, 11), NewMaybeRelocatableFelt(FeltFromUint64(20)))
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
		*NewMaybeRelocatableFelt(FeltOne()):         {0, 1},
		*NewMaybeRelocatableFelt(FeltFromUint64(2)): {2, 3},
	}
	if !reflect.DeepEqual(accessIndices, expectedAccessIndices) {
		t.Errorf("SQUASH_DICT wrong access_indices.\n Expected %v, got: %v", expectedAccessIndices, accessIndices)
	}
	keysAny, err := scopes.Get("keys")
	keys := keysAny.([]MaybeRelocatable)
	expectedKeys := []MaybeRelocatable{
		*NewMaybeRelocatableFelt(FeltFromUint64(2)),
	}
	if !reflect.DeepEqual(keys, expectedKeys) {
		t.Errorf("SQUASH_DICT wrong keys.\n Expected %v, got: %v", expectedKeys, keys)
	}
	keyAny, err := scopes.Get("key")
	key := keyAny.(MaybeRelocatable)
	expectedKey := *NewMaybeRelocatableFelt(FeltOne())
	if !reflect.DeepEqual(key, expectedKey) {
		t.Errorf("SQUASH_DICT wrong key.\n Expected %v, got: %v", expectedKey, keys)
	}
}
