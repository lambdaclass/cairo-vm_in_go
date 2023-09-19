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

func TestIsNNHintZero(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()
	// Advance fp to avoid clashes with values inserted into ap
	vm.RunContext.Fp.Offset += 1
	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{
			"a": {NewMaybeRelocatableFelt(FeltFromUint64(17))},
		},
		vm,
	)
	hintProcessor := CairoVmHintProcessor{}
	hintData := any(HintData{
		Ids:  idsManager,
		Code: IS_NN,
	})
	err := hintProcessor.ExecuteHint(vm, &hintData, nil, nil)
	if err != nil {
		t.Errorf("IS_NN hint test failed with error %s", err)
	}
	// Check the value of memory[ap]
	val, err := vm.Segments.Memory.GetFelt(vm.RunContext.Ap)
	if err != nil || !val.IsZero() {
		t.Error("Wrong/No value inserted into ap")
	}
}

func TestIsNNHintOne(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()
	// Advance fp to avoid clashes with values inserted into ap
	vm.RunContext.Fp.Offset += 1
	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{
			"a": {NewMaybeRelocatableFelt(FeltFromDecString("-1"))},
		},
		vm,
	)
	hintProcessor := CairoVmHintProcessor{}
	hintData := any(HintData{
		Ids:  idsManager,
		Code: IS_NN,
	})
	err := hintProcessor.ExecuteHint(vm, &hintData, nil, nil)
	if err != nil {
		t.Errorf("IS_NN hint test failed with error %s", err)
	}
	// Check the value of memory[ap]
	val, err := vm.Segments.Memory.GetFelt(vm.RunContext.Ap)
	if err != nil || val != FeltOne() {
		t.Error("Wrong/No value inserted into ap")
	}
}

func TestIsNNOutOfRangeHintZero(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()
	// Advance fp to avoid clashes with values inserted into ap
	vm.RunContext.Fp.Offset += 1
	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{
			"a": {NewMaybeRelocatableFelt(FeltFromDecString("-1"))},
		},
		vm,
	)
	hintProcessor := CairoVmHintProcessor{}
	hintData := any(HintData{
		Ids:  idsManager,
		Code: IS_NN_OUT_OF_RANGE,
	})
	err := hintProcessor.ExecuteHint(vm, &hintData, nil, nil)
	if err != nil {
		t.Errorf("IS_NN_OUT_OF_RANGE hint test failed with error %s", err)
	}
	// Check the value of memory[ap]
	val, err := vm.Segments.Memory.GetFelt(vm.RunContext.Ap)
	if err != nil || !val.IsZero() {
		t.Error("Wrong/No value inserted into ap")
	}
}

func TestIsNNOutOfRangeHintOne(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()
	// Advance fp to avoid clashes with values inserted into ap
	vm.RunContext.Fp.Offset += 1
	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{
			"a": {NewMaybeRelocatableFelt(FeltFromUint64(17))},
		},
		vm,
	)
	hintProcessor := CairoVmHintProcessor{}
	hintData := any(HintData{
		Ids:  idsManager,
		Code: IS_NN_OUT_OF_RANGE,
	})
	err := hintProcessor.ExecuteHint(vm, &hintData, nil, nil)
	if err != nil {
		t.Errorf("IS_NN_OUT_OF_RANGE hint test failed with error %s", err)
	}
	// Check the value of memory[ap]
	val, err := vm.Segments.Memory.GetFelt(vm.RunContext.Ap)
	if err != nil || val != FeltOne() {
		t.Error("Wrong/No value inserted into ap")
	}
}

func TestIsLeFeltEq(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()
	// Advance fp to avoid clashes with values inserted into ap
	vm.RunContext.Fp.Offset += 1
	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{
			"a": {NewMaybeRelocatableFelt(FeltFromUint64(17))},
			"b": {NewMaybeRelocatableFelt(FeltFromUint64(17))},
		},
		vm,
	)
	hintProcessor := CairoVmHintProcessor{}
	hintData := any(HintData{
		Ids:  idsManager,
		Code: IS_LE_FELT,
	})
	err := hintProcessor.ExecuteHint(vm, &hintData, nil, nil)
	if err != nil {
		t.Errorf("IS_LE_FELT hint test failed with error %s", err)
	}
	// Check the value of memory[ap]
	val, err := vm.Segments.Memory.GetFelt(vm.RunContext.Ap)
	if err != nil || !val.IsZero() {
		t.Error("Wrong/No value inserted into ap")
	}
}

func TestIsLeFeltLt(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()
	// Advance fp to avoid clashes with values inserted into ap
	vm.RunContext.Fp.Offset += 1
	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{
			"a": {NewMaybeRelocatableFelt(FeltFromUint64(16))},
			"b": {NewMaybeRelocatableFelt(FeltFromUint64(17))},
		},
		vm,
	)
	hintProcessor := CairoVmHintProcessor{}
	hintData := any(HintData{
		Ids:  idsManager,
		Code: IS_LE_FELT,
	})
	err := hintProcessor.ExecuteHint(vm, &hintData, nil, nil)
	if err != nil {
		t.Errorf("IS_LE_FELT hint test failed with error %s", err)
	}
	// Check the value of memory[ap]
	val, err := vm.Segments.Memory.GetFelt(vm.RunContext.Ap)
	if err != nil || !val.IsZero() {
		t.Error("Wrong/No value inserted into ap")
	}
}

func TestIsLeFeltGt(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()
	// Advance fp to avoid clashes with values inserted into ap
	vm.RunContext.Fp.Offset += 1
	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{
			"a": {NewMaybeRelocatableFelt(FeltFromUint64(18))},
			"b": {NewMaybeRelocatableFelt(FeltFromUint64(17))},
		},
		vm,
	)
	hintProcessor := CairoVmHintProcessor{}
	hintData := any(HintData{
		Ids:  idsManager,
		Code: IS_LE_FELT,
	})
	err := hintProcessor.ExecuteHint(vm, &hintData, nil, nil)
	if err != nil {
		t.Errorf("IS_LE_FELT hint test failed with error %s", err)
	}
	// Check the value of memory[ap]
	val, err := vm.Segments.Memory.GetFelt(vm.RunContext.Ap)
	if err != nil || val != FeltOne() {
		t.Error("Wrong/No value inserted into ap")
	}
}
