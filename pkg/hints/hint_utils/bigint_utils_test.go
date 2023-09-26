package hint_utils_test

import (
	"math/big"
	"testing"

	. "github.com/lambdaclass/cairo-vm.go/pkg/hints"
	. "github.com/lambdaclass/cairo-vm.go/pkg/hints/hint_codes"
	. "github.com/lambdaclass/cairo-vm.go/pkg/hints/hint_utils"
	"github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
	. "github.com/lambdaclass/cairo-vm.go/pkg/types"
	. "github.com/lambdaclass/cairo-vm.go/pkg/vm"
	. "github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
)

func TestNonDetBigInt3Ok(t *testing.T) {
	vm := NewVirtualMachine()

	vm.Segments.AddSegment()
	vm.Segments.AddSegment()
	vm.Segments.AddSegment()

	vm.RunContext.Pc = NewRelocatable(0, 0)
	vm.RunContext.Ap = NewRelocatable(1, 6)
	vm.RunContext.Fp = NewRelocatable(1, 6)

	value, _ := new(big.Int).SetString("7737125245533626718119526477371252455336267181195264773712524553362", 10)
	execScopes := NewExecutionScopes()

	execScopes.AssignOrUpdateVariable("value", *value)

	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{
			"res": {nil},
		},
		vm,
	)

	hintProcessor := CairoVmHintProcessor{}
	hintData := any(HintData{
		Ids:  idsManager,
		Code: NONDET_BIGINT3_V1,
	})
	err := hintProcessor.ExecuteHint(vm, &hintData, nil, execScopes)
	if err != nil {
		t.Errorf("Non Det Big Int 3 hint test failed with error: %s", err)
	} else {
		valueInStruct0, err := idsManager.GetStructFieldFelt("res", 0, vm)
		expected0 := lambdaworks.FeltFromDecString("773712524553362")
		if err != nil {
			t.Errorf("error fetching from ids manager : %s", err)
		}
		if valueInStruct0 != expected0 {
			t.Errorf(" Incorrect field value %s, expected %s", valueInStruct0.ToBigInt().Text(10), expected0.ToBigInt().Text(10))
		}

		valueInStruct1, err := idsManager.GetStructFieldFelt("res", 1, vm)
		expected1 := lambdaworks.FeltFromDecString("57408430697461422066401280")
		if err != nil {
			t.Errorf("error fetching from ids manager : %s", err)
		}
		if valueInStruct1 != expected1 {
			t.Errorf(" Incorrect field value %s, expected %s", valueInStruct1.ToBigInt().Text(10), expected1.ToBigInt().Text(10))
		}

		valueInStruct2, err := idsManager.GetStructFieldFelt("res", 2, vm)
		expected2 := lambdaworks.FeltFromDecString("1292469707114105")
		if err != nil {
			t.Errorf("error fetching from ids manager : %s", err)
		}
		if valueInStruct2 != expected2 {
			t.Errorf(" Incorrect field value %s, expected %s", valueInStruct2.ToBigInt().Text(10), expected2.ToBigInt().Text(10))
		}
	}
}
