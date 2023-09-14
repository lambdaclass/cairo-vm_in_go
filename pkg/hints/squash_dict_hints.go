package hints

import (
	"sort"

	"github.com/lambdaclass/cairo-vm.go/pkg/builtins"
	. "github.com/lambdaclass/cairo-vm.go/pkg/hints/hint_utils"
	"github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
	. "github.com/lambdaclass/cairo-vm.go/pkg/types"
	. "github.com/lambdaclass/cairo-vm.go/pkg/vm"
	. "github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
	"github.com/pkg/errors"
	"golang.org/x/exp/maps"
)

// SortMaybeRelocatables implements sort.Interface for []*MaybeRelocatables
type SortMaybeRelocatables []*MaybeRelocatable

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
	address, err := ids.GetAddr("dict_accesses", vm)
	if err != nil {
		return err
	}
	ptrDiff, err := ids.GetFelt("ptr_diff", vm)
	if err != nil {
		return err
	}
	if !ptrDiff.ModFloor(lambdaworks.FeltFromUint64(DICT_ACCESS_SIZE)).IsZero() {
		errors.New("Accesses array size must be divisible by DictAccess.SIZE")
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
	accessIndices := make(map[*MaybeRelocatable][]int)
	for i := 0; i < int(nAccesses); i++ {
		key, err := vm.Segments.Memory.Get(address.AddUint(uint(i) * DICT_ACCESS_SIZE))
		if err != nil {
			return err
		}
		_, hasKey := accessIndices[key]
		if !hasKey {
			accessIndices[key] = make([]int, 0)
		}
		accessIndices[key] = append(accessIndices[key], i)
	}
	//Descending list of keys.
	keys := maps.Keys(accessIndices)
	sort.Sort(sort.Reverse(SortMaybeRelocatables(keys)))
	//Are the keys used bigger than the range_check bound.
	bigKeys := lambdaworks.FeltZero()
	highKeyFelt, isFelt := keys[0].GetFelt()
	if isFelt && highKeyFelt.Bits() >= builtins.RANGE_CHECK_N_PARTS*builtins.INNER_RC_BOUND_SHIFT {
		bigKeys = lambdaworks.FeltOne()
	}
	lowKey := keys[len(keys)-1]
	// Insert new scope variables
	scopes.AssignOrUpdateVariable("access_indices", accessIndices)
	scopes.AssignOrUpdateVariable("keys", keys)
	scopes.AssignOrUpdateVariable("key", lowKey)
	// Insert ids variables
	err = ids.Insert("big_keys", NewMaybeRelocatableFelt(bigKeys), vm)
	if err != nil {
		return err
	}
	return ids.Insert("first_key", lowKey, vm)
}
