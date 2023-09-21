package hints_test

import (
	"testing"

	. "github.com/lambdaclass/cairo-vm.go/pkg/hints"
	. "github.com/lambdaclass/cairo-vm.go/pkg/hints/hint_utils"
	"github.com/lambdaclass/cairo-vm.go/pkg/types"
	. "github.com/lambdaclass/cairo-vm.go/pkg/vm"
	. "github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
)

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
