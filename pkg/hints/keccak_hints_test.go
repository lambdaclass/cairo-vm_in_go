package hints_test

import (
	"testing"

	. "github.com/lambdaclass/cairo-vm.go/pkg/hints"
	. "github.com/lambdaclass/cairo-vm.go/pkg/hints/hint_utils"
	. "github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
	. "github.com/lambdaclass/cairo-vm.go/pkg/types"
	. "github.com/lambdaclass/cairo-vm.go/pkg/vm"
	. "github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
)

func TestUnsafeKeccakOk(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()
	data_ptr := vm.Segments.AddSegment()
	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{
			"length": {NewMaybeRelocatableFelt(FeltFromUint64(3))},
			"data":   {NewMaybeRelocatableRelocatable(data_ptr)},
			"high":   {nil},
			"low":    {nil},
		},
		vm,
	)
	// Insert data into memory
	data := []MaybeRelocatable{
		*NewMaybeRelocatableFelt(FeltOne()),
		*NewMaybeRelocatableFelt(FeltOne()),
		*NewMaybeRelocatableFelt(FeltOne()),
	}
	vm.Segments.LoadData(data_ptr, &data)
	// Add __keccak_max_size
	scopes := NewExecutionScopes()
	scopes.AssignOrUpdateVariable("__keccak_max_size", uint16(500))
	hintProcessor := CairoVmHintProcessor{}
	hintData := any(HintData{
		Ids:  idsManager,
		Code: UNSAFE_KECCAK,
	})
	err := hintProcessor.ExecuteHint(vm, &hintData, nil, scopes)
	if err != nil {
		t.Errorf("UNSAFE_KECCAK hint test failed with error %s", err)
	}
	// Check ids values
	high, err := idsManager.GetFelt("high", vm)
	expectedHigh := FeltFromDecString("199195598804046335037364682505062700553")
	if err != nil || high != expectedHigh {
		t.Errorf("Wrong/No ids.high.\n Expected %s, got %s.", expectedHigh.ToHexString(), high.ToHexString())
	}
	low, err := idsManager.GetFelt("low", vm)
	expectedLow := FeltFromDecString("259413678945892999811634722593932702747")
	if err != nil || low != expectedLow {
		t.Errorf("Wrong/No ids.low\n Expected %s, got %s.", expectedLow.ToHexString(), low.ToHexString())
	}
}
