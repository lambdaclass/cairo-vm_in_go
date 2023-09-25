package hints_test

import (
	"math/big"
	"testing"

	. "github.com/lambdaclass/cairo-vm.go/pkg/hints"
	. "github.com/lambdaclass/cairo-vm.go/pkg/hints/hint_codes"
	. "github.com/lambdaclass/cairo-vm.go/pkg/hints/hint_utils"
	"github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
	. "github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
	. "github.com/lambdaclass/cairo-vm.go/pkg/types"
	. "github.com/lambdaclass/cairo-vm.go/pkg/vm"
	. "github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
)

func TestVerifyZeroWithExternalConst(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()
	vm.Segments.AddSegment()

	vm.RunContext.Pc = NewRelocatable(0, 0)
	vm.RunContext.Ap = NewRelocatable(1, 9)
	vm.RunContext.Fp = NewRelocatable(1, 9)

	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{
			"val": {NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(55)), NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(0)), NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(0))},
			"q":   {nil},
		},
		vm,
	)

	newScepP := big.NewInt(55)
	execScopes := NewExecutionScopes()

	execScopes.AssignOrUpdateVariable("SECP_P", *newScepP)

	hintProcessor := CairoVmHintProcessor{}
	hintData := any(HintData{
		Ids:  idsManager,
		Code: VERIFY_ZERO_EXTERNAL_SECP,
	})
	err := hintProcessor.ExecuteHint(vm, &hintData, nil, execScopes)
	if err != nil {
		t.Errorf("verifyZeroWithExternalConst hint test failed with error: %s", err)
	} else {
		valueInMemory, err := idsManager.GetFelt("q", vm)
		if err != nil {
			t.Errorf("could not fetch value with error: %s", err)
		}
		if valueInMemory != FeltFromUint64(1) {
			t.Errorf("value in memory is not the expected")
		}
	}
}
