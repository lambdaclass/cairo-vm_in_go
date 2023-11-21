package hints_test

import (
	"testing"

	. "github.com/lambdaclass/cairo-vm.go/pkg/hints"
	. "github.com/lambdaclass/cairo-vm.go/pkg/hints/hint_codes"
	. "github.com/lambdaclass/cairo-vm.go/pkg/hints/hint_utils"
	. "github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
	. "github.com/lambdaclass/cairo-vm.go/pkg/vm"
	. "github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
)

func TestSha256InputFalse(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()
	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{
			"n_bytes":   {NewMaybeRelocatableFelt(FeltFromUint64(2))},
			"full_word": {nil},
		},
		vm,
	)
	hintProcessor := CairoVmHintProcessor{}
	hintData := any(HintData{
		Ids:  idsManager,
		Code: SHA256_INPUT,
	})
	err := hintProcessor.ExecuteHint(vm, &hintData, nil, nil)
	if err != nil {
		t.Errorf("SHA256_INPUT hint test failed with error %s", err)
	}
	// Check ids.full_word
	fullWord, err := idsManager.GetFelt("full_word", vm)
	if err != nil || fullWord.Cmp(FeltZero()) != 0 {
		t.Error("Wrong/No value inserted into ids.full_word")
	}
}

func TestSha256InputTrue(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()
	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{
			"n_bytes":   {NewMaybeRelocatableFelt(FeltFromUint64(8))},
			"full_word": {nil},
		},
		vm,
	)
	hintProcessor := CairoVmHintProcessor{}
	hintData := any(HintData{
		Ids:  idsManager,
		Code: SHA256_INPUT,
	})
	err := hintProcessor.ExecuteHint(vm, &hintData, nil, nil)
	if err != nil {
		t.Errorf("SHA256_INPUT hint test failed with error %s", err)
	}
	// Check ids.full_word
	fullWord, err := idsManager.GetFelt("full_word", vm)
	if err != nil || fullWord.Cmp(FeltOne()) != 0 {
		t.Error("Wrong/No value inserted into ids.full_word")
	}
}
