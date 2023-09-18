package hints

import (
	"sort"

	"github.com/lambdaclass/cairo-vm.go/pkg/builtins"
	. "github.com/lambdaclass/cairo-vm.go/pkg/hints/hint_utils"
	. "github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
	. "github.com/lambdaclass/cairo-vm.go/pkg/types"
	. "github.com/lambdaclass/cairo-vm.go/pkg/vm"
	. "github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
	"github.com/pkg/errors"
	"golang.org/x/exp/maps"
)

// SortMaybeRelocatables implements sort.Interface for []MaybeRelocatables
type SortMaybeRelocatables []MaybeRelocatable

func (s SortMaybeRelocatables) Len() int      { return len(s) }
func (s SortMaybeRelocatables) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s SortMaybeRelocatables) Less(i, j int) bool {
	isLess := false
	//Integers are considered smaller than all relocatable values.
	a, b := s[i], s[j]
	aFelt, aIsFelt := a.GetFelt()
	bFelt, bIsFelt := b.GetFelt()

	switch true {
	// Both felts
	case aIsFelt && bIsFelt:
		if aFelt.Cmp(bFelt) == -1 {
			isLess = true
		}
	// a Felt, b Relocatable
	case aIsFelt && !bIsFelt:
	// a Relocatable, b Felt
	case !aIsFelt && bIsFelt:
		isLess = true
	// Both Relocatables
	case !aIsFelt && !bIsFelt:
		aRel, _ := a.GetRelocatable()
		bRel, _ := a.GetRelocatable()
		if aRel.SegmentIndex == bRel.SegmentIndex {
			isLess = aRel.Offset < bRel.Offset
		} else {
			isLess = aRel.SegmentIndex < bRel.SegmentIndex
		}
	}

	return isLess
}

func squashDict(ids IdsManager, scopes *ExecutionScopes, vm *VirtualMachine) error {
	address, err := ids.GetRelocatable("dict_accesses", vm)
	if err != nil {
		return err
	}
	ptrDiff, err := ids.GetFelt("ptr_diff", vm)
	if err != nil {
		return err
	}
	if !ptrDiff.ModFloor(FeltFromUint64(DICT_ACCESS_SIZE)).IsZero() {
		return errors.New("Accesses array size must be divisible by DictAccess.SIZE")
	}
	nAccessesFelt, err := ids.GetFelt("n_accesses", vm)
	if err != nil {
		return err
	}
	nAccesses, err := nAccessesFelt.ToU64()
	if err != nil {
		return err
	}
	squashDictMaxSize, err := scopes.Get("__squash_dict_max_size")
	if err == nil {
		maxSize, ok := squashDictMaxSize.(uint64)
		if ok {
			if nAccesses > maxSize {
				return errors.Errorf("squash_dict() can only be used with n_accesses<=%d.\nGot: n_accesses=%d.", maxSize, nAccesses)
			}
		}
	}
	// A map from key to the list of indices accessing it.
	accessIndices := make(map[MaybeRelocatable][]int)
	for i := 0; i < int(nAccesses); i++ {
		key, err := vm.Segments.Memory.Get(address.AddUint(uint(i) * DICT_ACCESS_SIZE))
		if err != nil {
			return err
		}
		_, hasKey := accessIndices[*key]
		if !hasKey {
			accessIndices[*key] = make([]int, 0)
		}
		accessIndices[*key] = append(accessIndices[*key], i)
	}
	//Descending list of keys.
	keys := maps.Keys(accessIndices)
	sort.Sort(sort.Reverse(SortMaybeRelocatables(keys)))
	//Are the keys used bigger than the range_check bound.
	bigKeys := FeltZero()
	highKeyFelt, isFelt := keys[0].GetFelt()
	if isFelt && highKeyFelt.Bits() >= builtins.RANGE_CHECK_N_PARTS*builtins.INNER_RC_BOUND_SHIFT {
		bigKeys = FeltOne()
	}
	if len(keys) == 0 {
		return errors.New("keys is empty")
	}
	lowKey := keys[len(keys)-1]
	// Insert new scope variables
	scopes.AssignOrUpdateVariable("access_indices", accessIndices)
	scopes.AssignOrUpdateVariable("keys", keys[:len(keys)-1])
	scopes.AssignOrUpdateVariable("key", lowKey)
	// Insert ids variables
	err = ids.Insert("big_keys", NewMaybeRelocatableFelt(bigKeys), vm)
	if err != nil {
		return err
	}
	return ids.Insert("first_key", &lowKey, vm)
}

func squashDictInnerSkipLoop(ids IdsManager, scopes *ExecutionScopes, vm *VirtualMachine) error {
	currentAccessIndicesAny, err := scopes.Get("current_access_indices")
	if err != nil {
		return err
	}
	currentAccessIndices, ok := currentAccessIndicesAny.([]int)
	if !ok {
		return errors.New("current_access_indices not in scope")
	}
	// Hint Logic
	if len(currentAccessIndices) != 0 {
		return ids.Insert("should_skip_loop", NewMaybeRelocatableFelt(FeltZero()), vm)
	}
	return ids.Insert("should_skip_loop", NewMaybeRelocatableFelt(FeltOne()), vm)
}

func squashDictInnerFirstIteration(ids IdsManager, scopes *ExecutionScopes, vm *VirtualMachine) error {
	// Fetch scope variables
	accessIndicesAny, err := scopes.Get("access_indices")
	if err != nil {
		return err
	}
	accessIndices, ok := accessIndicesAny.(map[MaybeRelocatable][]int)
	if !ok {
		return errors.New("access_indices not in scope")
	}

	keyAny, err := scopes.Get("key")
	if err != nil {
		return err
	}
	key, ok := keyAny.(MaybeRelocatable)
	if !ok {
		return errors.New("key not in scope")
	}
	// Fetch ids variables
	rangeCheckPtr, err := ids.GetRelocatable("range_check_ptr", vm)
	if err != nil {
		return err
	}
	// Hint Logic
	currentAccessIndices := accessIndices[key]
	sort.Sort(sort.Reverse(sort.IntSlice(currentAccessIndices)))
	if len(currentAccessIndices) == 0 {
		return errors.New("current_access_indices is empty")
	}
	currentAccessIndex := currentAccessIndices[len(currentAccessIndices)-1]
	currentAccessIndices = currentAccessIndices[:len(currentAccessIndices)-1]
	// Add variables to scope
	scopes.AssignOrUpdateVariable("current_access_indices", currentAccessIndices)
	scopes.AssignOrUpdateVariable("current_access_index", currentAccessIndex)
	//Insert current_accesss_index into range_check_ptr
	return vm.Segments.Memory.Insert(rangeCheckPtr, NewMaybeRelocatableFelt(FeltFromUint64(uint64(currentAccessIndex))))
}

func squashDictInnerCheckAccessIndex(ids IdsManager, scopes *ExecutionScopes, vm *VirtualMachine) error {
	// Fetch scope variables
	currentAccessIndicesAny, err := scopes.Get("current_access_indices")
	if err != nil {
		return err
	}
	currentAccessIndices, ok := currentAccessIndicesAny.([]int)
	if !ok {
		return errors.New("current_access_indices not in scope")
	}
	currentAccessIndexAny, err := scopes.Get("current_access_index")
	if err != nil {
		return err
	}
	currentAccessIndex, ok := currentAccessIndexAny.(int)
	if !ok {
		return errors.New("current_access_index not in scope")
	}
	// Hint Logic
	if len(currentAccessIndices) == 0 {
		return errors.New("current_access_indices is empty")
	}
	newAccessIndex := currentAccessIndices[len(currentAccessIndices)-1]
	currentAccessIndices = currentAccessIndices[:len(currentAccessIndices)-1]
	deltaMinusOne := newAccessIndex - currentAccessIndex - 1
	// Update scope variables
	scopes.AssignOrUpdateVariable("current_access_indices", currentAccessIndices)
	scopes.AssignOrUpdateVariable("current_access_index", newAccessIndex)
	scopes.AssignOrUpdateVariable("new_access_index", newAccessIndex)
	// Update ids variables
	return ids.Insert("loop_temps", NewMaybeRelocatableFelt(FeltFromUint64(uint64(deltaMinusOne))), vm)
}

func squashDictInnerContinueLoop(ids IdsManager, scopes *ExecutionScopes, vm *VirtualMachine) error {
	currentAccessIndicesAny, err := scopes.Get("current_access_indices")
	if err != nil {
		return err
	}
	currentAccessIndices, ok := currentAccessIndicesAny.([]int)
	if !ok {
		return errors.New("current_access_indices not in scope")
	}
	// Hint Logic
	if len(currentAccessIndices) == 0 {
		return ids.InsertStructField("loop_temps", 3, NewMaybeRelocatableFelt(FeltZero()), vm)
	}
	return ids.InsertStructField("loop_temps", 3, NewMaybeRelocatableFelt(FeltOne()), vm)
}

func squashDictInnerAssertLenKeys(scopes *ExecutionScopes) error {
	// Fetch scope variables
	keysAny, err := scopes.Get("keys")
	if err != nil {
		return err
	}
	keys, ok := keysAny.([]MaybeRelocatable)
	if !ok {
		return errors.New("keys not in scope")
	}
	// Hint logic
	if len(keys) != 0 {
		return errors.New("Assertion failed: len(keys) == 0")
	}
	return nil
}

func squashDictInnerLenAssert(scopes *ExecutionScopes) error {
	// Fetch scope variables
	currentAccessIndicesAny, err := scopes.Get("current_access_indices")
	if err != nil {
		return err
	}
	currentAccessIndices, ok := currentAccessIndicesAny.([]int)
	if !ok {
		return errors.New("current_access_indices not in scope")
	}
	// Hint logic
	if len(currentAccessIndices) != 0 {
		return errors.New("Assertion failed: len(current_access_indices) == 0")
	}
	return nil
}

func squashDictInnerUsedAccessesAssert(ids IdsManager, scopes *ExecutionScopes, vm *VirtualMachine) error {
	// Fetch scope variables
	accessIndicesAny, err := scopes.Get("access_indices")
	if err != nil {
		return err
	}
	accessIndices, ok := accessIndicesAny.(map[MaybeRelocatable][]int)
	if !ok {
		return errors.New("access_indices not in scope")
	}
	keyAny, err := scopes.Get("key")
	if err != nil {
		return err
	}
	key, ok := keyAny.(MaybeRelocatable)
	if !ok {
		return errors.New("key not in scope")
	}
	// Fetch ids variable
	nUsedAccesses, err := ids.GetFelt("n_used_accesses", vm)
	if err != nil {
		return err
	}
	// Hint logic
	if FeltFromUint64(uint64(len(accessIndices[key]))) != nUsedAccesses {
		return errors.New("Assertion failed: ids.n_used_accesses == len(access_indices[key])")
	}
	return nil
}

func squashDictInnerNextKey(ids IdsManager, scopes *ExecutionScopes, vm *VirtualMachine) error {
	// Fetch scope variables
	keysAny, err := scopes.Get("keys")
	if err != nil {
		return err
	}
	keys, ok := keysAny.([]MaybeRelocatable)
	if !ok {
		return errors.New("keys not in scope")
	}
	// Hint logic
	if len(keys) <= 0 {
		return errors.New("Assertion failed: len(keys) > 0.\nNo keys left but remaining_accesses > 0.)")
	}
	key := keys[len(keys)-1]
	keys = keys[:len(keys)-1]
	ids.Insert("next_key", &key, vm)
	// Update scope variables
	scopes.AssignOrUpdateVariable("keys", keys)
	scopes.AssignOrUpdateVariable("key", key)
	return nil
}
