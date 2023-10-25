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

func TestBlake2sComputeOutputOffsetZero(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()
	output := vm.Segments.AddSegment()
	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{
			"output": {NewMaybeRelocatableRelocatable(output)},
		},
		vm,
	)
	hintProcessor := CairoVmHintProcessor{}
	hintData := any(HintData{
		Ids:  idsManager,
		Code: BLAKE2S_COMPUTE,
	})
	err := hintProcessor.ExecuteHint(vm, &hintData, nil, nil)
	if err == nil {
		t.Errorf("BLAKE2S_COMPUTE hint test should have failed")
	}
}

func TestBlake2sComputeOutputSegmentEmpty(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()
	output := vm.Segments.AddSegment()
	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{
			"output": {NewMaybeRelocatableRelocatable(output.AddUint(26))},
		},
		vm,
	)
	hintProcessor := CairoVmHintProcessor{}
	hintData := any(HintData{
		Ids:  idsManager,
		Code: BLAKE2S_COMPUTE,
	})
	err := hintProcessor.ExecuteHint(vm, &hintData, nil, nil)
	if err == nil {
		t.Errorf("BLAKE2S_COMPUTE hint test should have failed")
	}
}

func TestBlake2sComputeBigOutput(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()
	output := vm.Segments.AddSegment()
	data := []MaybeRelocatable{
		*NewMaybeRelocatableFelt(FeltFromDecString("7842562439562793675803603603688959")),
		*NewMaybeRelocatableFelt(FeltFromDecString("7842562439562793675803603603688959")),
		*NewMaybeRelocatableFelt(FeltFromDecString("7842562439562793675803603603688959")),
		*NewMaybeRelocatableFelt(FeltFromDecString("7842562439562793675803603603688959")),
		*NewMaybeRelocatableFelt(FeltFromDecString("7842562439562793675803603603688959")),
		*NewMaybeRelocatableFelt(FeltFromDecString("7842562439562793675803603603688959")),
		*NewMaybeRelocatableFelt(FeltFromDecString("7842562439562793675803603603688959")),
		*NewMaybeRelocatableFelt(FeltFromDecString("7842562439562793675803603603688959")),
		*NewMaybeRelocatableFelt(FeltFromDecString("7842562439562793675803603603688959")),
		*NewMaybeRelocatableFelt(FeltFromDecString("7842562439562793675803603603688959")),
		*NewMaybeRelocatableFelt(FeltFromDecString("7842562439562793675803603603688959")),
		*NewMaybeRelocatableFelt(FeltFromDecString("7842562439562793675803603603688959")),
		*NewMaybeRelocatableFelt(FeltFromDecString("7842562439562793675803603603688959")),
		*NewMaybeRelocatableFelt(FeltFromDecString("7842562439562793675803603603688959")),
		*NewMaybeRelocatableFelt(FeltFromDecString("7842562439562793675803603603688959")),
		*NewMaybeRelocatableFelt(FeltFromDecString("7842562439562793675803603603688959")),
		*NewMaybeRelocatableFelt(FeltFromDecString("7842562439562793675803603603688959")),
		*NewMaybeRelocatableFelt(FeltFromDecString("7842562439562793675803603603688959")),
		*NewMaybeRelocatableFelt(FeltFromDecString("7842562439562793675803603603688959")),
		*NewMaybeRelocatableFelt(FeltFromDecString("7842562439562793675803603603688959")),
		*NewMaybeRelocatableFelt(FeltFromDecString("7842562439562793675803603603688959")),
		*NewMaybeRelocatableFelt(FeltFromDecString("7842562439562793675803603603688959")),
		*NewMaybeRelocatableFelt(FeltFromDecString("7842562439562793675803603603688959")),
		*NewMaybeRelocatableFelt(FeltFromDecString("7842562439562793675803603603688959")),
		*NewMaybeRelocatableFelt(FeltFromDecString("7842562439562793675803603603688959")),
		*NewMaybeRelocatableFelt(FeltFromDecString("7842562439562793675803603603688959")),
	}
	output, _ = vm.Segments.LoadData(output, &data)
	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{
			"output": {NewMaybeRelocatableRelocatable(output)},
		},
		vm,
	)
	hintProcessor := CairoVmHintProcessor{}
	hintData := any(HintData{
		Ids:  idsManager,
		Code: BLAKE2S_COMPUTE,
	})
	err := hintProcessor.ExecuteHint(vm, &hintData, nil, nil)
	if err == nil {
		t.Errorf("BLAKE2S_COMPUTE hint test should have failed")
	}
}

func TestBlake2sComputeOk(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()
	output := vm.Segments.AddSegment()
	data := []MaybeRelocatable{
		*NewMaybeRelocatableFelt(FeltFromDecString("17")),
		*NewMaybeRelocatableFelt(FeltFromDecString("17")),
		*NewMaybeRelocatableFelt(FeltFromDecString("17")),
		*NewMaybeRelocatableFelt(FeltFromDecString("17")),
		*NewMaybeRelocatableFelt(FeltFromDecString("17")),
		*NewMaybeRelocatableFelt(FeltFromDecString("17")),
		*NewMaybeRelocatableFelt(FeltFromDecString("17")),
		*NewMaybeRelocatableFelt(FeltFromDecString("17")),
		*NewMaybeRelocatableFelt(FeltFromDecString("17")),
		*NewMaybeRelocatableFelt(FeltFromDecString("17")),
		*NewMaybeRelocatableFelt(FeltFromDecString("17")),
		*NewMaybeRelocatableFelt(FeltFromDecString("17")),
		*NewMaybeRelocatableFelt(FeltFromDecString("17")),
		*NewMaybeRelocatableFelt(FeltFromDecString("17")),
		*NewMaybeRelocatableFelt(FeltFromDecString("17")),
		*NewMaybeRelocatableFelt(FeltFromDecString("17")),
		*NewMaybeRelocatableFelt(FeltFromDecString("17")),
		*NewMaybeRelocatableFelt(FeltFromDecString("17")),
		*NewMaybeRelocatableFelt(FeltFromDecString("17")),
		*NewMaybeRelocatableFelt(FeltFromDecString("17")),
		*NewMaybeRelocatableFelt(FeltFromDecString("17")),
		*NewMaybeRelocatableFelt(FeltFromDecString("17")),
		*NewMaybeRelocatableFelt(FeltFromDecString("17")),
		*NewMaybeRelocatableFelt(FeltFromDecString("17")),
		*NewMaybeRelocatableFelt(FeltFromDecString("17")),
		*NewMaybeRelocatableFelt(FeltFromDecString("17")),
	}
	output, _ = vm.Segments.LoadData(output, &data)
	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{
			"output": {NewMaybeRelocatableRelocatable(output)},
		},
		vm,
	)
	hintProcessor := CairoVmHintProcessor{}
	hintData := any(HintData{
		Ids:  idsManager,
		Code: BLAKE2S_COMPUTE,
	})
	err := hintProcessor.ExecuteHint(vm, &hintData, nil, nil)
	if err != nil {
		t.Errorf("BLAKE2S_COMPUTE hint test failed with error %s", err)
	}
}
