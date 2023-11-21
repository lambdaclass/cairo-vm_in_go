package hints_test

import (
	"reflect"
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

func TestBlake2sAddUint256Ok(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()
	data := vm.Segments.AddSegment()
	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{
			"data": {NewMaybeRelocatableRelocatable(data)},
			"high": {NewMaybeRelocatableFelt(FeltFromUint(25))},
			"low":  {NewMaybeRelocatableFelt(FeltFromUint(20))},
		},
		vm,
	)
	hintProcessor := CairoVmHintProcessor{}
	hintData := any(HintData{
		Ids:  idsManager,
		Code: BLAKE2S_ADD_UINT256,
	})
	err := hintProcessor.ExecuteHint(vm, &hintData, nil, nil)
	if err != nil {
		t.Errorf("BLAKE2S_ADD_UINT256 hint test failed with error %s", err)
	}
	// Check the data segment
	dataSegment, err := vm.Segments.GetFeltRange(data, 8)
	expectedDataSegment := []Felt{
		FeltFromUint(20),
		FeltZero(),
		FeltZero(),
		FeltZero(),
		FeltFromUint(25),
		FeltZero(),
		FeltZero(),
		FeltZero(),
	}
	if err != nil || !reflect.DeepEqual(dataSegment, expectedDataSegment) {
		t.Errorf("Wrong/No data loaded.\n Expected %v, got %v", expectedDataSegment, dataSegment)
	}
}

func TestBlake2sAddUint256BigEndOk(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()
	data := vm.Segments.AddSegment()
	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{
			"data": {NewMaybeRelocatableRelocatable(data)},
			"high": {NewMaybeRelocatableFelt(FeltFromUint(25))},
			"low":  {NewMaybeRelocatableFelt(FeltFromUint(20))},
		},
		vm,
	)
	hintProcessor := CairoVmHintProcessor{}
	hintData := any(HintData{
		Ids:  idsManager,
		Code: BLAKE2S_ADD_UINT256_BIGEND,
	})
	err := hintProcessor.ExecuteHint(vm, &hintData, nil, nil)
	if err != nil {
		t.Errorf("BLAKE2S_ADD_UINT256_BIGEND hint test failed with error %s", err)
	}
	// Check the data segment
	dataSegment, err := vm.Segments.GetFeltRange(data, 8)
	expectedDataSegment := []Felt{
		FeltZero(),
		FeltZero(),
		FeltZero(),
		FeltFromUint(25),
		FeltZero(),
		FeltZero(),
		FeltZero(),
		FeltFromUint(20),
	}
	if err != nil || !reflect.DeepEqual(dataSegment, expectedDataSegment) {
		t.Errorf("Wrong/No data loaded.\n Expected %v, got %v", expectedDataSegment, dataSegment)
	}
}

func TestBlake2sFinaizeOk(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()
	data := vm.Segments.AddSegment()
	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{
			"blake2s_ptr_end": {NewMaybeRelocatableRelocatable(data)},
		},
		vm,
	)
	hintProcessor := CairoVmHintProcessor{}
	hintData := any(HintData{
		Ids:  idsManager,
		Code: BLAKE2S_FINALIZE,
	})
	err := hintProcessor.ExecuteHint(vm, &hintData, nil, nil)
	if err != nil {
		t.Errorf("BLAKE2S_FINALIZE hint test failed with error %s", err)
	}
	// Check the data segment
	dataSegment, err := vm.Segments.GetFeltRange(data, 204)

	expectedDataSegment := []Felt{
		FeltFromUint(1795745351), FeltFromUint(3144134277), FeltFromUint(1013904242), FeltFromUint(2773480762), FeltFromUint(1359893119), FeltFromUint(2600822924), FeltFromUint(528734635),
		FeltFromUint(1541459225), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltFromUint(4294967295), FeltFromUint(813310313),
		FeltFromUint(2491453561), FeltFromUint(3491828193), FeltFromUint(2085238082), FeltFromUint(1219908895), FeltFromUint(514171180), FeltFromUint(4245497115), FeltFromUint(4193177630),

		FeltFromUint(1795745351), FeltFromUint(3144134277), FeltFromUint(1013904242), FeltFromUint(2773480762), FeltFromUint(1359893119), FeltFromUint(2600822924), FeltFromUint(528734635),
		FeltFromUint(1541459225), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltFromUint(4294967295), FeltFromUint(813310313),
		FeltFromUint(2491453561), FeltFromUint(3491828193), FeltFromUint(2085238082), FeltFromUint(1219908895), FeltFromUint(514171180), FeltFromUint(4245497115), FeltFromUint(4193177630),

		FeltFromUint(1795745351), FeltFromUint(3144134277), FeltFromUint(1013904242), FeltFromUint(2773480762), FeltFromUint(1359893119), FeltFromUint(2600822924), FeltFromUint(528734635),
		FeltFromUint(1541459225), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltFromUint(4294967295), FeltFromUint(813310313),
		FeltFromUint(2491453561), FeltFromUint(3491828193), FeltFromUint(2085238082), FeltFromUint(1219908895), FeltFromUint(514171180), FeltFromUint(4245497115), FeltFromUint(4193177630),

		FeltFromUint(1795745351), FeltFromUint(3144134277), FeltFromUint(1013904242), FeltFromUint(2773480762), FeltFromUint(1359893119), FeltFromUint(2600822924), FeltFromUint(528734635),
		FeltFromUint(1541459225), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltFromUint(4294967295), FeltFromUint(813310313),
		FeltFromUint(2491453561), FeltFromUint(3491828193), FeltFromUint(2085238082), FeltFromUint(1219908895), FeltFromUint(514171180), FeltFromUint(4245497115), FeltFromUint(4193177630),

		FeltFromUint(1795745351), FeltFromUint(3144134277), FeltFromUint(1013904242), FeltFromUint(2773480762), FeltFromUint(1359893119), FeltFromUint(2600822924), FeltFromUint(528734635),
		FeltFromUint(1541459225), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltFromUint(4294967295), FeltFromUint(813310313),
		FeltFromUint(2491453561), FeltFromUint(3491828193), FeltFromUint(2085238082), FeltFromUint(1219908895), FeltFromUint(514171180), FeltFromUint(4245497115), FeltFromUint(4193177630),

		FeltFromUint(1795745351), FeltFromUint(3144134277), FeltFromUint(1013904242), FeltFromUint(2773480762), FeltFromUint(1359893119), FeltFromUint(2600822924), FeltFromUint(528734635),
		FeltFromUint(1541459225), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltFromUint(4294967295), FeltFromUint(813310313),
		FeltFromUint(2491453561), FeltFromUint(3491828193), FeltFromUint(2085238082), FeltFromUint(1219908895), FeltFromUint(514171180), FeltFromUint(4245497115), FeltFromUint(4193177630),
	}
	if err != nil || !reflect.DeepEqual(dataSegment, expectedDataSegment) {
		t.Errorf("Wrong/No data loaded.\n Expected: %v.\n Got: %v", expectedDataSegment, dataSegment)
	}
}

func TestBlake2sFinaizeV3Ok(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()
	data := vm.Segments.AddSegment()
	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{
			"blake2s_ptr_end": {NewMaybeRelocatableRelocatable(data)},
		},
		vm,
	)
	hintProcessor := CairoVmHintProcessor{}
	hintData := any(HintData{
		Ids:  idsManager,
		Code: BLAKE2S_FINALIZE_V3,
	})
	err := hintProcessor.ExecuteHint(vm, &hintData, nil, nil)
	if err != nil {
		t.Errorf("BLAKE2S_FINALIZE_V3 hint test failed with error %s", err)
	}
	// Check the data segment
	dataSegment, err := vm.Segments.GetFeltRange(data, 204)

	expectedDataSegment := []Felt{
		FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(),
		FeltFromUint(1795745351), FeltFromUint(3144134277), FeltFromUint(1013904242), FeltFromUint(2773480762), FeltFromUint(1359893119), FeltFromUint(2600822924), FeltFromUint(528734635), FeltFromUint(1541459225), FeltZero(), FeltFromUint(4294967295), FeltFromUint(813310313),
		FeltFromUint(2491453561), FeltFromUint(3491828193), FeltFromUint(2085238082), FeltFromUint(1219908895), FeltFromUint(514171180), FeltFromUint(4245497115), FeltFromUint(4193177630),

		FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(),
		FeltFromUint(1795745351), FeltFromUint(3144134277), FeltFromUint(1013904242), FeltFromUint(2773480762), FeltFromUint(1359893119), FeltFromUint(2600822924), FeltFromUint(528734635), FeltFromUint(1541459225), FeltZero(), FeltFromUint(4294967295), FeltFromUint(813310313),
		FeltFromUint(2491453561), FeltFromUint(3491828193), FeltFromUint(2085238082), FeltFromUint(1219908895), FeltFromUint(514171180), FeltFromUint(4245497115), FeltFromUint(4193177630),

		FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(),
		FeltFromUint(1795745351), FeltFromUint(3144134277), FeltFromUint(1013904242), FeltFromUint(2773480762), FeltFromUint(1359893119), FeltFromUint(2600822924), FeltFromUint(528734635), FeltFromUint(1541459225), FeltZero(), FeltFromUint(4294967295), FeltFromUint(813310313),
		FeltFromUint(2491453561), FeltFromUint(3491828193), FeltFromUint(2085238082), FeltFromUint(1219908895), FeltFromUint(514171180), FeltFromUint(4245497115), FeltFromUint(4193177630),

		FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(),
		FeltFromUint(1795745351), FeltFromUint(3144134277), FeltFromUint(1013904242), FeltFromUint(2773480762), FeltFromUint(1359893119), FeltFromUint(2600822924), FeltFromUint(528734635), FeltFromUint(1541459225), FeltZero(), FeltFromUint(4294967295), FeltFromUint(813310313),
		FeltFromUint(2491453561), FeltFromUint(3491828193), FeltFromUint(2085238082), FeltFromUint(1219908895), FeltFromUint(514171180), FeltFromUint(4245497115), FeltFromUint(4193177630),

		FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(),
		FeltFromUint(1795745351), FeltFromUint(3144134277), FeltFromUint(1013904242), FeltFromUint(2773480762), FeltFromUint(1359893119), FeltFromUint(2600822924), FeltFromUint(528734635), FeltFromUint(1541459225), FeltZero(), FeltFromUint(4294967295), FeltFromUint(813310313),
		FeltFromUint(2491453561), FeltFromUint(3491828193), FeltFromUint(2085238082), FeltFromUint(1219908895), FeltFromUint(514171180), FeltFromUint(4245497115), FeltFromUint(4193177630),

		FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(), FeltZero(),
		FeltFromUint(1795745351), FeltFromUint(3144134277), FeltFromUint(1013904242), FeltFromUint(2773480762), FeltFromUint(1359893119), FeltFromUint(2600822924), FeltFromUint(528734635), FeltFromUint(1541459225), FeltZero(), FeltFromUint(4294967295), FeltFromUint(813310313),
		FeltFromUint(2491453561), FeltFromUint(3491828193), FeltFromUint(2085238082), FeltFromUint(1219908895), FeltFromUint(514171180), FeltFromUint(4245497115), FeltFromUint(4193177630),
	}
	if err != nil || !reflect.DeepEqual(dataSegment, expectedDataSegment) {
		t.Errorf("Wrong/No data loaded.\n Expected: %v.\n Got: %v", expectedDataSegment, dataSegment)
	}
}

func TestExampleBlake2sCompressEmptyInput(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()
	output := vm.Segments.AddSegment()
	blake2sStart := vm.Segments.AddSegment()
	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{
			"output":        {NewMaybeRelocatableRelocatable(output)},
			"blake2s_start": {NewMaybeRelocatableRelocatable(blake2sStart)},
			"n_bytes":       {NewMaybeRelocatableFelt(FeltOne())},
		},
		vm,
	)
	hintProcessor := CairoVmHintProcessor{}
	hintData := any(HintData{
		Ids:  idsManager,
		Code: EXAMPLE_BLAKE2S_COMPRESS,
	})
	err := hintProcessor.ExecuteHint(vm, &hintData, nil, nil)
	if err == nil {
		t.Errorf("EXAMPLE_BLAKE2S_COMPRESS hint test should have failed")
	}
}
