package hints_test

import (
	"reflect"
	"sort"
	"testing"

	. "github.com/lambdaclass/cairo-vm.go/pkg/hints"
	. "github.com/lambdaclass/cairo-vm.go/pkg/hints/hint_utils"
	"github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
	"github.com/lambdaclass/cairo-vm.go/pkg/types"
	. "github.com/lambdaclass/cairo-vm.go/pkg/vm"
	. "github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
)

func TestReverseFeltArray(t *testing.T) {
	array := []lambdaworks.Felt{lambdaworks.FeltFromUint(6), lambdaworks.FeltFromUint(0), lambdaworks.FeltFromUint(100), lambdaworks.FeltFromUint(1), lambdaworks.FeltFromUint(50)}

	sort.Sort(SortFelt(array))

	sortedarray := []lambdaworks.Felt{lambdaworks.FeltFromUint(0), lambdaworks.FeltFromUint(1), lambdaworks.FeltFromUint(6), lambdaworks.FeltFromUint(50), lambdaworks.FeltFromUint(100)}

	if !reflect.DeepEqual(array, sortedarray) {
		t.Errorf("Error sorting felt array")
	}

}

func TestUsortWithMaxSize(t *testing.T) {
	vm := NewVirtualMachine()
	scopes := types.NewExecutionScopes()
	scopes.AssignOrUpdateVariable("usort_max_size", uint64(1))
	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{},
		vm,
	)
	hintProcessor := CairoVmHintProcessor{}
	hintData := any(HintData{
		Ids:  idsManager,
		Code: USORT_ENTER_SCOPE,
	})
	err := hintProcessor.ExecuteHint(vm, &hintData, nil, scopes)
	if err != nil {
		t.Errorf("USORT_ENTER_SCOPE hint execution failed")
	}

	usort_max_size_interface, err := scopes.Get("usort_max_size")

	if err != nil {
		t.Errorf("Error assigning usort_max_size")
	}

	usort_max_size := usort_max_size_interface.(uint64)

	if usort_max_size != uint64(1) {
		t.Errorf("Error assigning usort_max_size")
	}

}
func TestUsortOutOfRange(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()
	vm.Segments.AddSegment()
	vm.Segments.AddSegment()
	scopes := types.NewExecutionScopes()
	scopes.AssignOrUpdateVariable("usort_max_size", uint64(1))
	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{
			"input":     {NewMaybeRelocatableRelocatable(NewRelocatable(2, 1))},
			"input_len": {NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(5))},
		},
		vm,
	)
	hintProcessor := CairoVmHintProcessor{}
	hintData := any(HintData{
		Ids:  idsManager,
		Code: USORT_BODY,
	})
	err := hintProcessor.ExecuteHint(vm, &hintData, nil, scopes)
	if err == nil {
		t.Errorf("USORT_BODY hint should have failed")
	}

}

func TestUsortVerify(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()
	vm.Segments.AddSegment()
	vm.Segments.AddSegment()
	scopes := types.NewExecutionScopes()
	positions_dict := make(map[lambdaworks.Felt][]uint64)
	positions_dict[lambdaworks.FeltFromUint64(0)] = []uint64{2}
	positions_dict[lambdaworks.FeltFromUint64(1)] = []uint64{1}
	positions_dict[lambdaworks.FeltFromUint64(2)] = []uint64{0}

	scopes.AssignOrUpdateVariable("positions_dict", positions_dict)

	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{
			"value": {NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(0))},
		},
		vm,
	)
	hintProcessor := CairoVmHintProcessor{}
	hintData := any(HintData{
		Ids:  idsManager,
		Code: USORT_VERIFY,
	})
	err := hintProcessor.ExecuteHint(vm, &hintData, nil, scopes)
	if err != nil {
		t.Errorf("USORT_VERIFY failed")
	}

	positions_interface, err := scopes.Get("positions")

	if err != nil {
		t.Errorf("Error assigning positions_interface")
	}

	positions := positions_interface.([]uint64)

	if !reflect.DeepEqual(positions, []uint64{2}) {
		t.Errorf("Error assigning positions")
	}

	last_pos_interface, err := scopes.Get("last_pos")

	if err != nil {
		t.Errorf("Error assigning last_pos")
	}

	last_pos := last_pos_interface.(uint64)

	if last_pos != uint64(0) {
		t.Errorf("Error assigning last_pos")
	}

}

// const USORT_VERIFY_MULTIPLICITY_ASSERT = "assert len(positions) == 0"

func TestUsortVerifyMultiplicityAssert(t *testing.T) {
	vm := NewVirtualMachine()
	scopes := types.NewExecutionScopes()

	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{},
		vm,
	)
	hintProcessor := CairoVmHintProcessor{}
	hintData := any(HintData{
		Ids:  idsManager,
		Code: USORT_VERIFY_MULTIPLICITY_ASSERT,
	})
	err := hintProcessor.ExecuteHint(vm, &hintData, nil, scopes)
	if err == nil {
		t.Errorf("USORT_VERIFY_MULTIPLICITY_ASSERT should have failed")
	}

	positions := []uint64{0}

	scopes.AssignOrUpdateVariable("positions", positions)

	err = hintProcessor.ExecuteHint(vm, &hintData, nil, scopes)
	if err == nil {
		t.Errorf("USORT_VERIFY_MULTIPLICITY_ASSERT should have failed")
	}

	positions = []uint64{}

	scopes.AssignOrUpdateVariable("positions", positions)

	err = hintProcessor.ExecuteHint(vm, &hintData, nil, scopes)
	if err != nil {
		t.Errorf("USORT_VERIFY_MULTIPLICITY_ASSERT  failed")
	}

}
