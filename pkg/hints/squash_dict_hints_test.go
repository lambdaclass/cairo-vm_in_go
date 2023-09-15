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
	scopes.AssignOrUpdateVariable("__squash_dict_max_size", uint64(12))
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

func TestSquashDictInvalidOneKeyDictWithMaxSizeExceeded(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()
	vm.Segments.AddSegment()
	vm.Segments.AddSegment()
	scopes := types.NewExecutionScopes()
	scopes.AssignOrUpdateVariable("__squash_dict_max_size", uint64(1))
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
	if err == nil {
		t.Errorf("SQUASH_DICT hint should have failed")
	}
}

func TestSquashDictInvalidOneKeyDictBadPtrDiff(t *testing.T) {
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
			"ptr_diff":      {NewMaybeRelocatableFelt(FeltFromUint64(7))},
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
	if err == nil {
		t.Errorf("SQUASH_DICT hint should have failed")
	}
}

func TestSquashDictInvalidOneKeyDictNAccessesTooBig(t *testing.T) {
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
			"n_accesses":    {NewMaybeRelocatableFelt(FeltFromDecString("-1"))},
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
	if err == nil {
		t.Errorf("SQUASH_DICT hint should have failed")
	}
}

func TestSquashDictSkipLoopTrue(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()
	vm.Segments.AddSegment()
	vm.Segments.AddSegment()
	scopes := types.NewExecutionScopes()
	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{
			"should_skip_loop": {nil},
		},
		vm,
	)
	hintProcessor := CairoVmHintProcessor{}
	hintData := any(HintData{
		Ids:  idsManager,
		Code: SQUASH_DICT_INNER_SKIP_LOOP,
	})
	err := hintProcessor.ExecuteHint(vm, &hintData, nil, scopes)
	if err != nil {
		t.Errorf("SQUASH_DICT_INNER_SKIP_LOOP hint failed with error: %s", err)
	}
	// Check ids.skip_loop
	skipLoop, err := idsManager.GetFelt("should_skip_loop", vm)
	if err != nil || skipLoop != FeltOne() {
		t.Errorf("SQUASH_DICT_INNER_SKIP_LOOP hint failed. Wrong/No ids.skip_loop")
	}
}

func TestSquashDictSkipLoopFalse(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()
	vm.Segments.AddSegment()
	vm.Segments.AddSegment()
	scopes := types.NewExecutionScopes()
	scopes.AssignOrUpdateVariable("current_access_indices", []MaybeRelocatable{})
	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{
			"should_skip_loop": {nil},
		},
		vm,
	)
	hintProcessor := CairoVmHintProcessor{}
	hintData := any(HintData{
		Ids:  idsManager,
		Code: SQUASH_DICT_INNER_SKIP_LOOP,
	})
	err := hintProcessor.ExecuteHint(vm, &hintData, nil, scopes)
	if err != nil {
		t.Errorf("SQUASH_DICT_INNER_SKIP_LOOP hint failed with error: %s", err)
	}
	// Check ids.skip_loop
	skipLoop, err := idsManager.GetFelt("should_skip_loop", vm)
	if err != nil || skipLoop != FeltZero() {
		t.Errorf("SQUASH_DICT_INNER_SKIP_LOOP hint failed. Wrong/No ids.skip_loop")
	}
}

func TestSquashDictInnerFirstIterationOk(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()
	range_check_ptr := vm.Segments.AddSegment()
	scopes := types.NewExecutionScopes()
	scopes.AssignOrUpdateVariable("access_indices", map[MaybeRelocatable][]int{
		*NewMaybeRelocatableFelt(FeltFromUint64(5)): {
			9, 3, 10, 7,
		},
	})
	scopes.AssignOrUpdateVariable("key", *NewMaybeRelocatableFelt(FeltFromUint64(5)))
	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{
			"range_check_ptr": {NewMaybeRelocatableRelocatable(range_check_ptr)},
		},
		vm,
	)
	hintProcessor := CairoVmHintProcessor{}
	hintData := any(HintData{
		Ids:  idsManager,
		Code: SQUASH_DICT_INNER_FIRST_ITERATION,
	})
	err := hintProcessor.ExecuteHint(vm, &hintData, nil, scopes)
	if err != nil {
		t.Errorf("SQUASH_DICT_INNER_FIRST_ITERATION hint failed with error: %s", err)
	}
	// Check scope values
	currentAccessIndicesAny, err := scopes.Get("current_access_indices")
	currentAccessIndices := currentAccessIndicesAny.([]int)
	expectedCurrentAccessIndices := []int{10, 9, 7}
	if !reflect.DeepEqual(currentAccessIndices, expectedCurrentAccessIndices) {
		t.Errorf("Wrong current_access_indices.\n Expected %v, got: %v", expectedCurrentAccessIndices, currentAccessIndices)
	}

	currentAccessIndexAny, err := scopes.Get("current_access_index")
	currentAccessIndex := currentAccessIndexAny.(int)
	expectedCurrentAccessIndex := int(3)
	if !reflect.DeepEqual(currentAccessIndex, expectedCurrentAccessIndex) {
		t.Errorf("Wrong current_access_index.\n Expected %v, got: %v", expectedCurrentAccessIndex, currentAccessIndex)
	}

	// Check memory[ids.range_check_ptr]
	val, err := vm.Segments.Memory.Get(range_check_ptr)
	if err != nil || *val != *NewMaybeRelocatableFelt(FeltFromUint64(3)) {
		t.Errorf("Wrong/No value inserted into memory[ids.range_check_ptr]")
	}
}

func TestSquashDictInnerFirstIterationEmpty(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()
	range_check_ptr := vm.Segments.AddSegment()
	scopes := types.NewExecutionScopes()
	scopes.AssignOrUpdateVariable("access_indices", map[MaybeRelocatable][]int{})
	scopes.AssignOrUpdateVariable("key", *NewMaybeRelocatableFelt(FeltFromUint64(5)))
	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{
			"range_check_ptr": {NewMaybeRelocatableRelocatable(range_check_ptr)},
		},
		vm,
	)
	hintProcessor := CairoVmHintProcessor{}
	hintData := any(HintData{
		Ids:  idsManager,
		Code: SQUASH_DICT_INNER_FIRST_ITERATION,
	})
	err := hintProcessor.ExecuteHint(vm, &hintData, nil, scopes)
	if err == nil {
		t.Error("SQUASH_DICT_INNER_FIRST_ITERATION hint should have failed")
	}
}

func TestSquashDictInnerCheckAccessIndexOk(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()
	scopes := types.NewExecutionScopes()
	scopes.AssignOrUpdateVariable("current_access_indices", []int{
		10, 9, 7, 5,
	},
	)
	scopes.AssignOrUpdateVariable("current_access_index", int(1))
	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{
			"loop_temps": {nil},
		},
		vm,
	)
	hintProcessor := CairoVmHintProcessor{}
	hintData := any(HintData{
		Ids:  idsManager,
		Code: SQUASH_DICT_INNER_CHECK_ACCESS_INDEX,
	})
	err := hintProcessor.ExecuteHint(vm, &hintData, nil, scopes)
	if err != nil {
		t.Errorf("SQUASH_DICT_INNER_CHECK_ACCESS_INDEX hint failed with error: %s", err)
	}
	// Check scope values
	currentAccessIndicesAny, err := scopes.Get("current_access_indices")
	currentAccessIndices := currentAccessIndicesAny.([]int)
	expectedCurrentAccessIndices := []int{10, 9, 7}
	if !reflect.DeepEqual(currentAccessIndices, expectedCurrentAccessIndices) {
		t.Errorf("Wrong current_access_indices.\n Expected %v, got: %v", expectedCurrentAccessIndices, currentAccessIndices)
	}

	currentAccessIndexAny, err := scopes.Get("current_access_index")
	currentAccessIndex := currentAccessIndexAny.(int)
	expectedCurrentAccessIndex := int(5)
	if !reflect.DeepEqual(currentAccessIndex, expectedCurrentAccessIndex) {
		t.Errorf("Wrong current_access_index.\n Expected %v, got: %v", expectedCurrentAccessIndex, currentAccessIndex)
	}

	newAccessIndexAny, err := scopes.Get("new_access_index")
	newAccessIndex := newAccessIndexAny.(int)
	expectedNewAccessIndex := int(5)
	if !reflect.DeepEqual(newAccessIndex, expectedNewAccessIndex) {
		t.Errorf("Wrong new_access_index.\n Expected %v, got: %v", expectedNewAccessIndex, newAccessIndex)
	}

	//Check loop_temps.index_delta_minus_1
	val, err := idsManager.GetFelt("loop_temps", vm)
	if err != nil || val != FeltFromUint64(3) {
		t.Errorf("Wrong/No value inserted into memory[ids.range_check_ptr]")
	}
}

func TestSquashDictInnerCheckAccessIndexEmpty(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()
	scopes := types.NewExecutionScopes()
	scopes.AssignOrUpdateVariable("current_access_indices", []int{})
	scopes.AssignOrUpdateVariable("current_access_index", int(1))
	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{
			"loop_temps": {nil},
		},
		vm,
	)
	hintProcessor := CairoVmHintProcessor{}
	hintData := any(HintData{
		Ids:  idsManager,
		Code: SQUASH_DICT_INNER_CHECK_ACCESS_INDEX,
	})
	err := hintProcessor.ExecuteHint(vm, &hintData, nil, scopes)
	if err == nil {
		t.Errorf("SQUASH_DICT_INNER_CHECK_ACCESS_INDEX hint should have failed")
	}
}

func TestSquashDictContinuepLoopTrue(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()
	vm.Segments.AddSegment()
	vm.Segments.AddSegment()
	scopes := types.NewExecutionScopes()
	scopes.AssignOrUpdateVariable("current_access_indices", []MaybeRelocatable{})
	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{
			"loop_temps": {nil},
		},
		vm,
	)
	hintProcessor := CairoVmHintProcessor{}
	hintData := any(HintData{
		Ids:  idsManager,
		Code: SQUASH_DICT_INNER_CONTINUE_LOOP,
	})
	err := hintProcessor.ExecuteHint(vm, &hintData, nil, scopes)
	if err != nil {
		t.Errorf("SQUASH_DICT_INNER_CONTINUE_LOOP hint failed with error: %s", err)
	}
	// Check ids.loop_temps.should_continue
	skipLoop, err := idsManager.GetStructFieldFelt("loop_temps", 3, vm)
	if err != nil || skipLoop != FeltOne() {
		t.Errorf("SQUASH_DICT_INNER_CONTINUE_LOOP hint failed. Wrong/No ids.loop_temps.should_continue")
	}
}

func TestSquashDictContinueLoopFalse(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()
	vm.Segments.AddSegment()
	vm.Segments.AddSegment()
	scopes := types.NewExecutionScopes()
	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{
			"loop_temps": {nil},
		},
		vm,
	)
	hintProcessor := CairoVmHintProcessor{}
	hintData := any(HintData{
		Ids:  idsManager,
		Code: SQUASH_DICT_INNER_CONTINUE_LOOP,
	})
	err := hintProcessor.ExecuteHint(vm, &hintData, nil, scopes)
	if err != nil {
		t.Errorf("SQUASH_DICT_INNER_CONTINUE_LOOP hint failed with error: %s", err)
	}
	// Check ids.loop_temps.should_continue
	continueLoop, err := idsManager.GetStructFieldFelt("loop_temps", 3, vm)
	if err != nil || continueLoop != FeltZero() {
		t.Errorf("SQUASH_DICT_INNER_CONTINUE_LOOP hint failed. Wrong/No ids.loop_temps.should_continue")
	}
}

func TestSquashDictInnerAssertLenKeysNotEmpty(t *testing.T) {
	vm := NewVirtualMachine()
	scopes := types.NewExecutionScopes()
	scopes.AssignOrUpdateVariable("keys", []MaybeRelocatable{*NewMaybeRelocatableFelt(FeltZero())})
	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{},
		vm,
	)
	hintProcessor := CairoVmHintProcessor{}
	hintData := any(HintData{
		Ids:  idsManager,
		Code: SQUASH_DICT_INNER_ASSERT_LEN_KEYS,
	})
	err := hintProcessor.ExecuteHint(vm, &hintData, nil, scopes)
	if err == nil {
		t.Errorf("SQUASH_DICT_INNER_ASSERT_LEN_KEYS hint should have failed")
	}
}

func TestSquashDictInnerAssertLenKeysEmpty(t *testing.T) {
	vm := NewVirtualMachine()
	scopes := types.NewExecutionScopes()
	scopes.AssignOrUpdateVariable("keys", []MaybeRelocatable{})
	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{},
		vm,
	)
	hintProcessor := CairoVmHintProcessor{}
	hintData := any(HintData{
		Ids:  idsManager,
		Code: SQUASH_DICT_INNER_ASSERT_LEN_KEYS,
	})
	err := hintProcessor.ExecuteHint(vm, &hintData, nil, scopes)
	if err != nil {
		t.Errorf("SQUASH_DICT_INNER_ASSERT_LEN_KEYS hint failed with error: %s", err)
	}
}

func TestSquashDictInnerLenAssertKeysNotEmpty(t *testing.T) {
	vm := NewVirtualMachine()
	scopes := types.NewExecutionScopes()
	scopes.AssignOrUpdateVariable("current_access_indices", []int{2, 3})
	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{},
		vm,
	)
	hintProcessor := CairoVmHintProcessor{}
	hintData := any(HintData{
		Ids:  idsManager,
		Code: SQUASH_DICT_INNER_LEN_ASSERT,
	})
	err := hintProcessor.ExecuteHint(vm, &hintData, nil, scopes)
	if err == nil {
		t.Errorf("SQUASH_DICT_INNER_LEN_ASSERT hint should have failed")
	}
}

func TestSquashDictInnerLenAssertEmpty(t *testing.T) {
	vm := NewVirtualMachine()
	scopes := types.NewExecutionScopes()
	scopes.AssignOrUpdateVariable("current_access_indices", []int{})
	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{},
		vm,
	)
	hintProcessor := CairoVmHintProcessor{}
	hintData := any(HintData{
		Ids:  idsManager,
		Code: SQUASH_DICT_INNER_LEN_ASSERT,
	})
	err := hintProcessor.ExecuteHint(vm, &hintData, nil, scopes)
	if err != nil {
		t.Errorf("SQUASH_DICT_INNER_LEN_ASSERT hint failed with error: %s", err)
	}
}

func TestSquashDictInnerUsedAccessesAssertOk(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()
	scopes := types.NewExecutionScopes()
	scopes.AssignOrUpdateVariable("access_indices", map[MaybeRelocatable][]int{
		*NewMaybeRelocatableFelt(FeltZero()): {3},
	})
	scopes.AssignOrUpdateVariable("key", *NewMaybeRelocatableFelt(FeltZero()))
	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{
			"n_used_accesses": {NewMaybeRelocatableFelt(FeltOne())},
		},
		vm,
	)
	hintProcessor := CairoVmHintProcessor{}
	hintData := any(HintData{
		Ids:  idsManager,
		Code: SQUASH_DICT_INNER_USED_ACCESSES_ASSERT,
	})
	err := hintProcessor.ExecuteHint(vm, &hintData, nil, scopes)
	if err != nil {
		t.Errorf("SQUASH_DICT_INNER_USED_ACCESSES_ASSERT hint failed with error: %s", err)
	}
}

func TestSquashDictInnerUsedAccessesAssertBadLen(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()
	scopes := types.NewExecutionScopes()
	scopes.AssignOrUpdateVariable("access_indices", map[MaybeRelocatable][]int{
		*NewMaybeRelocatableFelt(FeltZero()): {3, 5},
	})
	scopes.AssignOrUpdateVariable("key", *NewMaybeRelocatableFelt(FeltZero()))
	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{
			"n_used_accesses": {NewMaybeRelocatableFelt(FeltOne())},
		},
		vm,
	)
	hintProcessor := CairoVmHintProcessor{}
	hintData := any(HintData{
		Ids:  idsManager,
		Code: SQUASH_DICT_INNER_USED_ACCESSES_ASSERT,
	})
	err := hintProcessor.ExecuteHint(vm, &hintData, nil, scopes)
	if err == nil {
		t.Errorf("SQUASH_DICT_INNER_USED_ACCESSES_ASSERT hint should have failed")
	}
}

func TestSquashDictInnerUsedAccessesAssertBadKey(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()
	scopes := types.NewExecutionScopes()
	scopes.AssignOrUpdateVariable("access_indices", map[MaybeRelocatable][]int{
		*NewMaybeRelocatableFelt(FeltZero()): {3},
	})
	scopes.AssignOrUpdateVariable("key", *NewMaybeRelocatableFelt(FeltOne()))
	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{
			"n_used_accesses": {NewMaybeRelocatableFelt(FeltOne())},
		},
		vm,
	)
	hintProcessor := CairoVmHintProcessor{}
	hintData := any(HintData{
		Ids:  idsManager,
		Code: SQUASH_DICT_INNER_USED_ACCESSES_ASSERT,
	})
	err := hintProcessor.ExecuteHint(vm, &hintData, nil, scopes)
	if err == nil {
		t.Errorf("SQUASH_DICT_INNER_USED_ACCESSES_ASSERT hint should have failed")
	}
}
