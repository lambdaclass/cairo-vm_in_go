package hints_test

import (
	"testing"

	. "github.com/lambdaclass/cairo-vm.go/pkg/hints"
	. "github.com/lambdaclass/cairo-vm.go/pkg/hints/hint_utils"
	. "github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
	. "github.com/lambdaclass/cairo-vm.go/pkg/vm"
	. "github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
)

func TestPowHintOk(t *testing.T) {
	vm := NewVirtualMachine()

	vm.Segments.AddSegment()
	vm.Segments.AddSegment()
	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{
			"prev_locs": {NewMaybeRelocatableRelocatable(NewRelocatable(0, 0))},
			"exp":       {NewMaybeRelocatableFelt(FeltFromUint64(3))},
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

	cast_locs, err := locs.ToU64()
	if err != nil {
		t.Errorf("Couldn't cast locs.big to u64: %s", err)
	}

	if cast_locs != 1 {
		t.Errorf("locs.bit: %d != 1", cast_locs)
	}
}
