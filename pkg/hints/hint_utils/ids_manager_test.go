package hint_utils_test

import (
	"testing"

	. "github.com/lambdaclass/cairo-vm.go/pkg/hints/hint_utils"
	"github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
	"github.com/lambdaclass/cairo-vm.go/pkg/parser"
	"github.com/lambdaclass/cairo-vm.go/pkg/vm"
	"github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
)

func TestIdsManagerGetAddressSimpleReference(t *testing.T) {
	ids := IdsManager{
		References: map[string]HintReference{
			"val": {
				Offset1: OffsetValue{
					Register:  vm.FP,
					ValueType: Reference,
				},
			},
		},
	}
	vm := vm.NewVirtualMachine()
	addr, err := ids.GetAddr("val", vm)
	if err != nil {
		t.Errorf("Error in test: %s", err)
	}
	if addr != vm.RunContext.Fp {
		t.Errorf("IdsManager.GetAddr returned wrong value")
	}
}

func TestIdsManagerGetAddressWithApTrackingCorrection(t *testing.T) {
	ids := IdsManager{
		References: map[string]HintReference{
			"val": {
				Offset1: OffsetValue{
					Register:  vm.AP,
					ValueType: Reference,
				},
				ApTrackingData: parser.ApTrackingData{Group: 1, Offset: 2},
			},
		},
		HintApTracking: parser.ApTrackingData{Group: 1, Offset: 5},
	}
	// Ap tracking correction ap - (hintOff - idsOff) = (1, 5) - (5 - 2) = (1, 5) - 3 = (1, 2)
	vm := vm.NewVirtualMachine()
	vm.RunContext.Ap = memory.NewRelocatable(1, 5)
	addr, err := ids.GetAddr("val", vm)
	if err != nil {
		t.Errorf("Error in test: %s", err)
	}
	expected := memory.NewRelocatable(1, 2)
	if addr != expected {
		t.Errorf("IdsManager.GetAddr returned wrong value")
	}
}

func TestIdsManagerGetAddressUnknownIdentifier(t *testing.T) {
	ids := IdsManager{
		References: map[string]HintReference{
			"value": {
				Offset1: OffsetValue{
					Register: vm.FP,
				},
			},
		},
	}
	vm := vm.NewVirtualMachine()
	_, err := ids.GetAddr("val", vm)
	if err == nil {
		t.Errorf("IdsManager.GetAddress should have failed")
	}
}

func TestIdsManagerGetAddressComplexReferenceDoubleDeref(t *testing.T) {
	// reference: [ap + 1] + [fp + 2] = (1, 0) + 3 = (1, 3)
	ids := IdsManager{
		References: map[string]HintReference{
			"val": {
				Offset1: OffsetValue{
					Register:    vm.AP,
					Value:       1,
					ValueType:   Reference,
					Dereference: true,
				},
				Offset2: OffsetValue{
					Register:    vm.FP,
					Value:       2,
					ValueType:   Reference,
					Dereference: true,
				},
			},
		},
	}
	vm := vm.NewVirtualMachine()
	vm.Segments.AddSegment()
	// [ap + 1] = (1, 0)
	vm.Segments.Memory.Insert(vm.RunContext.Ap.AddUint(1), memory.NewMaybeRelocatableRelocatable(memory.Relocatable{SegmentIndex: 1, Offset: 0}))
	// [fp + 2] = 3
	vm.Segments.Memory.Insert(vm.RunContext.Fp.AddUint(2), memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(3)))
	addr, err := ids.GetAddr("val", vm)
	if err != nil {
		t.Errorf("Error in test: %s", err)
	}
	expected_addr := memory.Relocatable{SegmentIndex: 1, Offset: 3}
	if addr != expected_addr {
		t.Errorf("IdsManager.GetAddr returned wrong value")
	}
}

func TestIdsManagerGetAddressComplexReferenceOneDeref(t *testing.T) {
	// reference: [ap + 1] + 2 = (1, 0) + 2 = (1, 2)
	ids := IdsManager{
		References: map[string]HintReference{
			"val": {
				Offset1: OffsetValue{
					Register:    vm.AP,
					Value:       1,
					ValueType:   Reference,
					Dereference: true,
				},
				Offset2: OffsetValue{
					Value:     2,
					ValueType: Value,
				},
			},
		},
	}
	vm := vm.NewVirtualMachine()
	vm.Segments.AddSegment()
	// [ap + 1] = (1, 0)
	vm.Segments.Memory.Insert(vm.RunContext.Ap.AddUint(1), memory.NewMaybeRelocatableRelocatable(memory.Relocatable{SegmentIndex: 1, Offset: 0}))
	addr, err := ids.GetAddr("val", vm)
	if err != nil {
		t.Errorf("Error in test: %s", err)
	}
	expected_addr := memory.Relocatable{SegmentIndex: 1, Offset: 2}
	if addr != expected_addr {
		t.Errorf("IdsManager.GetAddr returned wrong value")
	}
}

func TestIdsManagerGetNoDereference(t *testing.T) {
	ids := IdsManager{
		References: map[string]HintReference{
			"val": {
				Offset1: OffsetValue{
					Register:  vm.FP,
					ValueType: Reference,
				},
			},
		},
	}
	vm := vm.NewVirtualMachine()
	val, err := ids.Get("val", vm)
	if err != nil {
		t.Errorf("Error in test: %s", err)
	}
	expected := memory.NewMaybeRelocatableRelocatable(vm.RunContext.Fp)
	if *val != *expected {
		t.Errorf("IdsManager.Get returned wrong value, expected")
	}
}

func TestIdsManagerGetRelocatableNoDereference(t *testing.T) {
	ids := IdsManager{
		References: map[string]HintReference{
			"val": {
				Offset1: OffsetValue{
					Register:  vm.FP,
					ValueType: Reference,
				},
			},
		},
	}
	vm := vm.NewVirtualMachine()
	val, err := ids.GetRelocatable("val", vm)
	if err != nil {
		t.Errorf("Error in test: %s", err)
	}
	expected := vm.RunContext.Fp
	if val != expected {
		t.Errorf("IdsManager.GetRelocatable returned wrong value")
	}
}

func TestIdsManagerGetRelocatableDeref(t *testing.T) {
	ids := IdsManager{
		References: map[string]HintReference{
			"val": {
				Offset1: OffsetValue{
					Register:  vm.FP,
					ValueType: Reference,
				},
				Dereference: true,
			},
		},
	}
	vm := vm.NewVirtualMachine()
	vm.Segments.AddSegment()
	vm.Segments.Memory.Insert(vm.RunContext.Fp, memory.NewMaybeRelocatableRelocatable(memory.NewRelocatable(1, 3)))
	val, err := ids.GetRelocatable("val", vm)
	if err != nil {
		t.Errorf("Error in test: %s", err)
	}
	expected := memory.NewRelocatable(1, 3)
	if val != expected {
		t.Errorf("IdsManager.GetRelocatable returned wrong value")
	}
}

func TestIdsManagerGetFeltImmediate(t *testing.T) {
	ids := IdsManager{
		References: map[string]HintReference{
			"val": {
				Offset1: OffsetValue{
					Immediate: lambdaworks.FeltFromUint64(17),
					ValueType: Immediate,
				},
			},
		},
	}
	vm := vm.NewVirtualMachine()
	val, err := ids.GetFelt("val", vm)
	if err != nil {
		t.Errorf("Error in test: %s", err)
	}
	expected := lambdaworks.FeltFromUint64(17)
	if val != expected {
		t.Errorf("IdsManager.GetFelt returned wrong value")
	}
}

func TestIdsManagerGetFeltDeref(t *testing.T) {
	ids := IdsManager{
		References: map[string]HintReference{
			"val": {
				Offset1: OffsetValue{
					Register:  vm.FP,
					ValueType: Reference,
				},
				Dereference: true,
			},
		},
	}
	vm := vm.NewVirtualMachine()
	vm.Segments.AddSegment()
	vm.Segments.Memory.Insert(vm.RunContext.Fp, memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(17)))
	val, err := ids.GetFelt("val", vm)
	if err != nil {
		t.Errorf("Error in test: %s", err)
	}
	expected := lambdaworks.FeltFromUint64(17)
	if val != expected {
		t.Errorf("IdsManager.GetFelt returned wrong value")
	}
}

func TestIdsManagerGetStructFieldTest(t *testing.T) {
	ids := IdsManager{
		References: map[string]HintReference{
			"cat": {
				Offset1: OffsetValue{
					Register:  vm.FP,
					ValueType: Reference,
				},
				Dereference: true,
			},
		},
	}
	vm := vm.NewVirtualMachine()
	vm.Segments.AddSegment()
	vm.Segments.Memory.Insert(vm.RunContext.Fp, memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(7)))
	vm.Segments.Memory.Insert(vm.RunContext.Fp.AddUint(1), memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(4)))
	lives, err_lives := ids.GetStructFieldFelt("cat", 0, vm)
	paws, err_paws := ids.GetStructFieldFelt("cat", 1, vm)
	if err_lives != nil || err_paws != nil {
		t.Errorf("Error(s) in test: %s, %s", err_lives, err_paws)
	}
	expected_lives := lambdaworks.FeltFromUint64(7)
	expected_paws := lambdaworks.FeltFromUint64(4)
	if lives != expected_lives || paws != expected_paws {
		t.Errorf("IdsManager.GetStructFieldFelt returned wrong values")
	}
}

func TestIdsManagerGetConst(t *testing.T) {
	ids := IdsManager{
		AccessibleScopes: []string{
			"starkware.cairo.common.math",
			"starkware.cairo.common.math.assert_250_bit",
		},
	}
	upperBound := lambdaworks.FeltFromUint64(250)
	constants := map[string]lambdaworks.Felt{
		"starkware.cairo.common.math.assert_250_bit.UPPER_BOUND": upperBound,
	}
	constant, err := ids.GetConst("UPPER_BOUND", &constants)
	if err != nil || constant != upperBound {
		t.Errorf("IdsManager.GetConst returned wrong/no constant")
	}
}

func TestIdsManagerGetConstPrioritizeInnerModule(t *testing.T) {
	ids := IdsManager{
		AccessibleScopes: []string{
			"starkware.cairo.common.math",
			"starkware.cairo.common.math.assert_250_bit",
		},
	}
	upperBound := lambdaworks.FeltFromUint64(250)
	constants := map[string]lambdaworks.Felt{
		"starkware.cairo.common.math.assert_250_bit.UPPER_BOUND": upperBound,
		"starkware.cairo.common.math.UPPER_BOUND":                lambdaworks.FeltZero(),
	}
	constant, err := ids.GetConst("UPPER_BOUND", &constants)
	if err != nil || constant != upperBound {
		t.Errorf("IdsManager.GetConst returned wrong/no constant")
	}
}

func TestIdsManagerGetConstNoMConst(t *testing.T) {
	ids := IdsManager{
		AccessibleScopes: []string{
			"starkware.cairo.common.math",
			"starkware.cairo.common.math.assert_250_bit",
		},
	}
	lowerBound := lambdaworks.FeltFromUint64(250)
	constants := map[string]lambdaworks.Felt{
		"starkware.cairo.common.math.assert_250_bit.LOWER_BOUND": lowerBound,
	}
	_, err := ids.GetConst("UPPER_BOUND", &constants)
	if err == nil {
		t.Errorf("IdsManager.GetConst should have failed")
	}
}
