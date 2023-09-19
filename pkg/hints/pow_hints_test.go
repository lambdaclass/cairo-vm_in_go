package hints_test

import (
	. "github.com/lambdaclass/cairo-vm.go/pkg/hints"
	. "github.com/lambdaclass/cairo-vm.go/pkg/hints/hint_codes"
	. "github.com/lambdaclass/cairo-vm.go/pkg/hints/hint_utils"
	. "github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
	. "github.com/lambdaclass/cairo-vm.go/pkg/vm"
	. "github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
	"testing"
)

func TestPowHintOddOk(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()
	vm.Segments.Memory.Insert(NewRelocatable(0, 4), NewMaybeRelocatableFelt(FeltFromUint64(3)))

	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{
			"prev_locs": {NewMaybeRelocatableRelocatable(NewRelocatable(0, 0))},
			"locs":      {nil},
		},
		vm,
	)
	hintProcessor := CairoVmHintProcessor{}
	hintData := any(HintData{
		Ids:  idsManager,
		Code: POW,
	})

	err := hintProcessor.ExecuteHint(vm, &hintData, nil, nil)
	if err != nil {
		t.Errorf("POW hint test failed with error %s", err)
	}

	locs, err := idsManager.GetFelt("locs", vm)
	if err != nil {
		t.Errorf("Failed to get locs.bit with error: %s", err)
	}

	if locs != FeltOne() {
		t.Errorf("locs.bit: %d != 1", locs)
	}
}

func TestPowHintEvenOk(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()
	vm.Segments.Memory.Insert(NewRelocatable(0, 4), NewMaybeRelocatableFelt(FeltFromUint64(2)))

	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{
			"prev_locs": {NewMaybeRelocatableRelocatable(NewRelocatable(0, 0))},
			"locs":      {nil},
		},
		vm,
	)
	hintProcessor := CairoVmHintProcessor{}
	hintData := any(HintData{
		Ids:  idsManager,
		Code: POW,
	})

	err := hintProcessor.ExecuteHint(vm, &hintData, nil, nil)
	if err != nil {
		t.Errorf("POW hint test failed with error %s", err)
	}

	locs, err := idsManager.GetFelt("locs", vm)
	if err != nil {
		t.Errorf("Failed to get locs.bit with error: %s", err)
	}

	if locs != FeltZero() {
		t.Errorf("locs.bit: %d != 0", locs)
	}
}
