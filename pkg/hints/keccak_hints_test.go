package hints_test

import (
	"testing"

	. "github.com/lambdaclass/cairo-vm.go/pkg/hints"
	. "github.com/lambdaclass/cairo-vm.go/pkg/hints/hint_codes"
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
	scopes.AssignOrUpdateVariable("__keccak_max_size", uint64(500))
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

func TestUnsafeKeccakMaxSizeExceeded(t *testing.T) {
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
	scopes.AssignOrUpdateVariable("__keccak_max_size", uint64(2))
	hintProcessor := CairoVmHintProcessor{}
	hintData := any(HintData{
		Ids:  idsManager,
		Code: UNSAFE_KECCAK,
	})
	err := hintProcessor.ExecuteHint(vm, &hintData, nil, scopes)
	if err == nil {
		t.Errorf("UNSAFE_KECCAK hint test should have failed")
	}
}

func TestUnsafeKeccakInvalidWordSize(t *testing.T) {
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
		*NewMaybeRelocatableFelt(FeltFromDecString("-1")),
		*NewMaybeRelocatableFelt(FeltOne()),
		*NewMaybeRelocatableFelt(FeltOne()),
	}
	vm.Segments.LoadData(data_ptr, &data)
	// Add __keccak_max_size
	scopes := NewExecutionScopes()
	hintProcessor := CairoVmHintProcessor{}
	hintData := any(HintData{
		Ids:  idsManager,
		Code: UNSAFE_KECCAK,
	})
	err := hintProcessor.ExecuteHint(vm, &hintData, nil, scopes)
	if err == nil {
		t.Errorf("UNSAFE_KECCAK hint test should have failed")
	}
}

func TestUnsafeKeccakFinalizeOk(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()
	inputStart := vm.Segments.AddSegment()
	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{
			"keccak_state": {
				NewMaybeRelocatableRelocatable(inputStart),
				NewMaybeRelocatableRelocatable(inputStart.AddUint(2)),
			},
			"high": {nil},
			"low":  {nil},
		},
		vm,
	)
	// Insert keccak input into memory
	input := []MaybeRelocatable{
		*NewMaybeRelocatableFelt(FeltZero()),
		*NewMaybeRelocatableFelt(FeltOne()),
	}
	vm.Segments.LoadData(inputStart, &input)
	hintProcessor := CairoVmHintProcessor{}
	hintData := any(HintData{
		Ids:  idsManager,
		Code: UNSAFE_KECCAK_FINALIZE,
	})
	err := hintProcessor.ExecuteHint(vm, &hintData, nil, nil)
	if err != nil {
		t.Errorf("UNSAFE_KECCAK_FINALIZE hint test failed with error %s", err)
	}
	// Check ids values
	high, err := idsManager.GetFelt("high", vm)
	expectedHigh := FeltFromDecString("235346966651632113557018504892503714354")
	if err != nil || high != expectedHigh {
		t.Errorf("Wrong/No ids.high.\n Expected %s, got %s.", expectedHigh.ToHexString(), high.ToHexString())
	}
	low, err := idsManager.GetFelt("low", vm)
	expectedLow := FeltFromDecString("17219183504112405672555532996650339574")
	if err != nil || low != expectedLow {
		t.Errorf("Wrong/No ids.low\n Expected %s, got %s.", expectedLow.ToHexString(), low.ToHexString())
	}
}

func TestCompareBytesInWordHintEq(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()
	// Advance fp to avoid clashes with values inserted into ap
	vm.RunContext.Fp.Offset += 1
	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{
			"n_bytes": {NewMaybeRelocatableFelt(FeltFromUint64(17))},
		},
		vm,
	)
	constants := SetupConstantsForTest(
		map[string]Felt{
			"BYTES_IN_WORD": FeltFromUint64(17),
		},
		&idsManager)
	hintProcessor := CairoVmHintProcessor{}
	hintData := any(HintData{
		Ids:  idsManager,
		Code: COMPARE_BYTES_IN_WORD_NONDET,
	})
	err := hintProcessor.ExecuteHint(vm, &hintData, &constants, nil)
	if err != nil {
		t.Errorf("COMPARE_BYTES_IN_WORD_NONDET hint test failed with error %s", err)
	}
	// Check the value of memory[ap]
	val, err := vm.Segments.Memory.GetFelt(vm.RunContext.Ap)
	if err != nil || !val.IsZero() {
		t.Error("Wrong/No value inserted into ap")
	}
}

func TestCompareBytesInWordHintGt(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()
	// Advance fp to avoid clashes with values inserted into ap
	vm.RunContext.Fp.Offset += 1
	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{
			"n_bytes": {NewMaybeRelocatableFelt(FeltFromUint64(18))},
		},
		vm,
	)
	constants := SetupConstantsForTest(
		map[string]Felt{
			"BYTES_IN_WORD": FeltFromUint64(17),
		},
		&idsManager)
	hintProcessor := CairoVmHintProcessor{}
	hintData := any(HintData{
		Ids:  idsManager,
		Code: COMPARE_BYTES_IN_WORD_NONDET,
	})
	err := hintProcessor.ExecuteHint(vm, &hintData, &constants, nil)
	if err != nil {
		t.Errorf("COMPARE_BYTES_IN_WORD_NONDET hint test failed with error %s", err)
	}
	// Check the value of memory[ap]
	val, err := vm.Segments.Memory.GetFelt(vm.RunContext.Ap)
	if err != nil || !val.IsZero() {
		t.Error("Wrong/No value inserted into ap")
	}
}

func TestCompareBytesInWordHintLt(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()
	// Advance fp to avoid clashes with values inserted into ap
	vm.RunContext.Fp.Offset += 1
	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{
			"n_bytes": {NewMaybeRelocatableFelt(FeltFromUint64(16))},
		},
		vm,
	)
	constants := SetupConstantsForTest(
		map[string]Felt{
			"BYTES_IN_WORD": FeltFromUint64(17),
		},
		&idsManager)
	hintProcessor := CairoVmHintProcessor{}
	hintData := any(HintData{
		Ids:  idsManager,
		Code: COMPARE_BYTES_IN_WORD_NONDET,
	})
	err := hintProcessor.ExecuteHint(vm, &hintData, &constants, nil)
	if err != nil {
		t.Errorf("COMPARE_BYTES_IN_WORD_NONDET hint test failed with error %s", err)
	}
	// Check the value of memory[ap]
	val, err := vm.Segments.Memory.GetFelt(vm.RunContext.Ap)
	if err != nil || val != FeltOne() {
		t.Error("Wrong/No value inserted into ap")
	}
}

func TestCompareKeccakFullRateInBytesHintEq(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()
	// Advance fp to avoid clashes with values inserted into ap
	vm.RunContext.Fp.Offset += 1
	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{
			"n_bytes": {NewMaybeRelocatableFelt(FeltFromUint64(17))},
		},
		vm,
	)
	constants := SetupConstantsForTest(
		map[string]Felt{
			"KECCAK_FULL_RATE_IN_BYTES": FeltFromUint64(17),
		},
		&idsManager)
	hintProcessor := CairoVmHintProcessor{}
	hintData := any(HintData{
		Ids:  idsManager,
		Code: COMPARE_KECCAK_FULL_RATE_IN_BYTES_NONDET,
	})
	err := hintProcessor.ExecuteHint(vm, &hintData, &constants, nil)
	if err != nil {
		t.Errorf("COMPARE_KECCAK_FULL_RATE_IN_BYTES_NONDET hint test failed with error %s", err)
	}
	// Check the value of memory[ap]
	val, err := vm.Segments.Memory.GetFelt(vm.RunContext.Ap)
	if err != nil || val != FeltOne() {
		t.Error("Wrong/No value inserted into ap")
	}
}

func TestCompareKeccakFullRateInBytesHintGt(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()
	// Advance fp to avoid clashes with values inserted into ap
	vm.RunContext.Fp.Offset += 1
	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{
			"n_bytes": {NewMaybeRelocatableFelt(FeltFromUint64(18))},
		},
		vm,
	)
	constants := SetupConstantsForTest(
		map[string]Felt{
			"KECCAK_FULL_RATE_IN_BYTES": FeltFromUint64(17),
		},
		&idsManager)
	hintProcessor := CairoVmHintProcessor{}
	hintData := any(HintData{
		Ids:  idsManager,
		Code: COMPARE_KECCAK_FULL_RATE_IN_BYTES_NONDET,
	})
	err := hintProcessor.ExecuteHint(vm, &hintData, &constants, nil)
	if err != nil {
		t.Errorf("COMPARE_KECCAK_FULL_RATE_IN_BYTES_NONDET hint test failed with error %s", err)
	}
	// Check the value of memory[ap]
	val, err := vm.Segments.Memory.GetFelt(vm.RunContext.Ap)
	if err != nil || val != FeltOne() {
		t.Error("Wrong/No value inserted into ap")
	}
}

func TestCompareKeccakFullRateInBytesHintLt(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()
	// Advance fp to avoid clashes with values inserted into ap
	vm.RunContext.Fp.Offset += 1
	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{
			"n_bytes": {NewMaybeRelocatableFelt(FeltFromUint64(16))},
		},
		vm,
	)
	constants := SetupConstantsForTest(
		map[string]Felt{
			"KECCAK_FULL_RATE_IN_BYTES": FeltFromUint64(17),
		},
		&idsManager)
	hintProcessor := CairoVmHintProcessor{}
	hintData := any(HintData{
		Ids:  idsManager,
		Code: COMPARE_KECCAK_FULL_RATE_IN_BYTES_NONDET,
	})
	err := hintProcessor.ExecuteHint(vm, &hintData, &constants, nil)
	if err != nil {
		t.Errorf("COMPARE_KECCAK_FULL_RATE_IN_BYTES_NONDET hint test failed with error %s", err)
	}
	// Check the value of memory[ap]
	val, err := vm.Segments.Memory.GetFelt(vm.RunContext.Ap)
	if err != nil || !val.IsZero() {
		t.Error("Wrong/No value inserted into ap")
	}
}

func TestBlockPermutationOk(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()
	keccak_ptr := vm.Segments.AddSegment()
	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{
			"keccak_ptr": {NewMaybeRelocatableRelocatable(keccak_ptr.AddUint(25))},
		},
		vm,
	)
	data := make([]MaybeRelocatable, 0, 25)
	for i := 0; i < 25; i++ {
		data = append(data, *NewMaybeRelocatableFelt(FeltZero()))
	}
	vm.Segments.LoadData(keccak_ptr, &data)
	hintProcessor := CairoVmHintProcessor{}
	constants := SetupConstantsForTest(
		map[string]Felt{
			"KECCAK_STATE_SIZE_FELTS": FeltFromUint64(25),
		},
		&idsManager)
	hintData := any(HintData{
		Ids:  idsManager,
		Code: BLOCK_PERMUTATION,
	})
	err := hintProcessor.ExecuteHint(vm, &hintData, &constants, nil)
	if err != nil {
		t.Errorf("BLOCK_PERMUTATION hint test failed with error %s", err)
	}
}
