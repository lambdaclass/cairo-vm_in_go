package runners_test

import (
	"reflect"
	"testing"

	"github.com/lambdaclass/cairo-vm.go/pkg/builtins"
	"github.com/lambdaclass/cairo-vm.go/pkg/hints"
	"github.com/lambdaclass/cairo-vm.go/pkg/hints/hint_utils"
	"github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
	"github.com/lambdaclass/cairo-vm.go/pkg/parser"
	"github.com/lambdaclass/cairo-vm.go/pkg/runners"
	"github.com/lambdaclass/cairo-vm.go/pkg/vm"
	"github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
)

func TestNewCairoRunnerInvalidBuiltin(t *testing.T) {
	// Create a Program with one fake instruction
	program_data := make([]memory.MaybeRelocatable, 1)
	empty_identifiers := make(map[string]vm.Identifier, 0)
	program_data[0] = *memory.NewMaybeRelocatableFelt(lambdaworks.FeltOne())
	program := vm.Program{Data: program_data, Builtins: []string{"fake_builtin"}, Identifiers: empty_identifiers}
	// Create CairoRunner
	_, err := runners.NewCairoRunner(program)
	if err == nil {
		t.Errorf("Expected creating a CairoRunner with fake builtin to fail")
	}
}
func TestInitializeRunnerNoBuiltinsNoProofModeEmptyProgram(t *testing.T) {
	// Create a Program with empty data
	program_data := make([]memory.MaybeRelocatable, 0)
	empty_identifiers := make(map[string]vm.Identifier, 0)
	program := vm.Program{Data: program_data, Identifiers: empty_identifiers}
	// Create CairoRunner
	runner, err := runners.NewCairoRunner(program)
	if err != nil {
		t.Errorf("NewCairoRunner error in test: %s", err)
	}
	// Initialize the runner
	end_ptr, err := runner.Initialize()
	if err != nil {
		t.Errorf("Initialize error in test: %s", err)
	}
	if end_ptr.SegmentIndex != 3 || end_ptr.Offset != 0 {
		t.Errorf("Wrong end ptr value, got %+v", end_ptr)
	}

	// Check CairoRunner values
	if runner.ProgramBase.SegmentIndex != 0 || runner.ProgramBase.Offset != 0 {
		t.Errorf("Wrong ProgramBase value, got %+v", runner.ProgramBase)
	}

	// Check Vm's RunContext values
	if runner.Vm.RunContext.Pc.SegmentIndex != 0 || runner.Vm.RunContext.Pc.Offset != 0 {
		t.Errorf("Wrong Pc value, got %+v", runner.Vm.RunContext.Pc)
	}
	if runner.Vm.RunContext.Ap.SegmentIndex != 1 || runner.Vm.RunContext.Ap.Offset != 2 {
		t.Errorf("Wrong Ap value, got %+v", runner.Vm.RunContext.Ap)
	}
	if runner.Vm.RunContext.Fp.SegmentIndex != 1 || runner.Vm.RunContext.Fp.Offset != 2 {
		t.Errorf("Wrong Fp value, got %+v", runner.Vm.RunContext.Fp)
	}

	// Check memory

	// Program segment
	// 0:0 program_data[0] should be empty
	value, err := runner.Vm.Segments.Memory.Get(memory.Relocatable{SegmentIndex: 0, Offset: 0})
	if err == nil {
		t.Errorf("Expected addr 0:0 to be empty for empty program, got: %+v", value)
	}

	// Execution segment
	// 1:0 return_fp
	value, err = runner.Vm.Segments.Memory.Get(memory.Relocatable{SegmentIndex: 1, Offset: 0})
	if err != nil {
		t.Errorf("Memory Get error in test: %s", err)
	}
	rel, ok := value.GetRelocatable()
	if !ok || rel.SegmentIndex != 2 || rel.Offset != 0 {
		t.Errorf("Wrong value for address 1:0: %d", rel)
	}
	// 1:1 end_ptr
	value, err = runner.Vm.Segments.Memory.Get(memory.Relocatable{SegmentIndex: 1, Offset: 1})
	if err != nil {
		t.Errorf("Memory Get error in test: %s", err)
	}
	rel, ok = value.GetRelocatable()
	if !ok || rel.SegmentIndex != 3 || rel.Offset != 0 {
		t.Errorf("Wrong value for address 1:1: %d", rel)
	}
}

func TestInitializeRunnerNoBuiltinsNoProofModeNonEmptyProgram(t *testing.T) {
	// Create a Program with one fake instruction
	program_data := make([]memory.MaybeRelocatable, 1)
	program_data[0] = *memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(1))
	empty_identifiers := make(map[string]vm.Identifier, 0)
	program := vm.Program{Data: program_data, Identifiers: empty_identifiers}
	// Create CairoRunner
	runner, err := runners.NewCairoRunner(program)
	if err != nil {
		t.Errorf("NewCairoRunner error in test: %s", err)
	}
	// Initialize the runner
	end_ptr, err := runner.Initialize()
	if err != nil {
		t.Errorf("Initialize error in test: %s", err)
	}
	if end_ptr.SegmentIndex != 3 || end_ptr.Offset != 0 {
		t.Errorf("Wrong end ptr value, got %+v", end_ptr)
	}

	// Check CairoRunner values
	if runner.ProgramBase.SegmentIndex != 0 || runner.ProgramBase.Offset != 0 {
		t.Errorf("Wrong ProgramBase value, got %+v", runner.ProgramBase)
	}

	// Check Vm's RunContext values
	if runner.Vm.RunContext.Pc.SegmentIndex != 0 || runner.Vm.RunContext.Pc.Offset != 0 {
		t.Errorf("Wrong Pc value, got %+v", runner.Vm.RunContext.Pc)
	}
	if runner.Vm.RunContext.Ap.SegmentIndex != 1 || runner.Vm.RunContext.Ap.Offset != 2 {
		t.Errorf("Wrong Ap value, got %+v", runner.Vm.RunContext.Ap)
	}
	if runner.Vm.RunContext.Fp.SegmentIndex != 1 || runner.Vm.RunContext.Fp.Offset != 2 {
		t.Errorf("Wrong Fp value, got %+v", runner.Vm.RunContext.Fp)
	}

	// Check memory

	// Program segment
	// 0:0 program_data[0]
	value, err := runner.Vm.Segments.Memory.Get(memory.Relocatable{SegmentIndex: 0, Offset: 0})
	if err != nil {
		t.Errorf("Memory Get error in test: %s", err)
	}
	int, ok := value.GetFelt()
	if !ok || int != lambdaworks.FeltFromUint64(1) {
		t.Errorf("Wrong value for address 0:0: %d", int)
	}

	// Execution segment
	// 1:0 return_fp
	value, err = runner.Vm.Segments.Memory.Get(memory.Relocatable{SegmentIndex: 1, Offset: 0})
	if err != nil {
		t.Errorf("Memory Get error in test: %s", err)
	}
	rel, ok := value.GetRelocatable()
	if !ok || rel.SegmentIndex != 2 || rel.Offset != 0 {
		t.Errorf("Wrong value for address 1:0: %d", rel)
	}
	// 1:1 end_ptr
	value, err = runner.Vm.Segments.Memory.Get(memory.Relocatable{SegmentIndex: 1, Offset: 1})
	if err != nil {
		t.Errorf("Memory Get error in test: %s", err)
	}
	rel, ok = value.GetRelocatable()
	if !ok || rel.SegmentIndex != 3 || rel.Offset != 0 {
		t.Errorf("Wrong value for address 1:1: %d", rel)
	}
}

func TestInitializeRunnerWithRangeCheckValid(t *testing.T) {
	t.Helper()
	// Create a Program with one fake instruction
	program_data := make([]memory.MaybeRelocatable, 1)
	program_data[0] = *memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(1))
	empty_identifiers := make(map[string]vm.Identifier, 0)
	program_builtins := []string{builtins.CHECK_RANGE_BUILTIN_NAME}
	program := vm.Program{Data: program_data, Identifiers: empty_identifiers, Builtins: program_builtins}
	// Create CairoRunner
	runner, err := runners.NewCairoRunner(program)
	if err != nil {
		t.Errorf("NewCairoRunner error in test: %s", err)
	}
	// Initialize the runner
	_, err = runner.Initialize()
	if err != nil {
		t.Errorf("Initialize error in test: %s", err)
	}

	builtin_runner := runner.Vm.BuiltinRunners[0]
	if builtin_runner.Name() != builtins.CHECK_RANGE_BUILTIN_NAME {
		t.Errorf("Name of runner builtin failed. Expected %s, got %s", builtin_runner.Name(), builtins.CHECK_RANGE_BUILTIN_NAME)
	}

	builtin_base := builtin_runner.Base()
	expected_base := memory.NewRelocatable(2, 0)
	if !builtin_base.IsEqual(&expected_base) {
		t.Errorf("Base of runner builtin failed. Expected %d, got %d", expected_base, builtin_base)
	}

	err = runner.Vm.Segments.Memory.Insert(memory.NewRelocatable(2, 0), memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(23)))
	if err != nil {
		t.Errorf("Insertion failed in test with error: %s", err.Error())
	}

	err = runner.Vm.Segments.Memory.Insert(memory.NewRelocatable(2, 1), memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(233)))
	if err != nil {
		t.Errorf("Insert failed in test with error: %s", err.Error())
	}
}

func TestInitializeRunnerWithRangeCheckInvalid(t *testing.T) {
	t.Helper()
	// Create a Program with one fake instruction
	program_data := make([]memory.MaybeRelocatable, 1)
	program_data[0] = *memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(1))
	empty_identifiers := make(map[string]vm.Identifier, 0)
	program_builtins := []string{builtins.CHECK_RANGE_BUILTIN_NAME}
	program := vm.Program{Data: program_data, Identifiers: empty_identifiers, Builtins: program_builtins}
	// Create CairoRunner
	runner, err := runners.NewCairoRunner(program)
	if err != nil {
		t.Errorf("NewCairoRunner error in test: %s", err)
	}
	// Initialize the runner
	_, err = runner.Initialize()
	if err != nil {
		t.Errorf("Initialize error in test: %s", err)
	}

	builtin_runner := runner.Vm.BuiltinRunners[0]
	if builtin_runner.Name() != builtins.CHECK_RANGE_BUILTIN_NAME {
		t.Errorf("Name of runner builtin failed. Expected %s, got %s", builtin_runner.Name(), builtins.CHECK_RANGE_BUILTIN_NAME)
	}

	builtin_base := builtin_runner.Base()
	expected_base := memory.NewRelocatable(2, 0)
	if !builtin_base.IsEqual(&expected_base) {
		t.Errorf("Base of runner builtin failed. Expected %d, got %d", expected_base, builtin_base)
	}

	addr := memory.NewRelocatable(2, 0)
	val := memory.NewMaybeRelocatableRelocatable(memory.NewRelocatable(2, 1))
	err = runner.Vm.Segments.Memory.Insert(addr, val)
	expected_error := builtins.NotAFeltError(addr, *val)
	if err.Error() != expected_error.Error() {
		t.Errorf("Test failed: Expected error: %s, Actual error: %s", err.Error(), expected_error.Error())
	}

	addr = memory.NewRelocatable(2, 1)
	val = memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromDecString("-1"))
	err = runner.Vm.Segments.Memory.Insert(addr, val)
	felt, _ := val.GetFelt()
	expected_error = builtins.OutsideBoundsError(felt)
	if err.Error() != expected_error.Error() {
		t.Errorf("Test failed: Expected error: %s, Actual error: %s", err.Error(), expected_error.Error())
	}
}

func TestBuildHintDataMapEmpty(t *testing.T) {
	program := vm.Program{}
	runner, _ := runners.NewCairoRunner(program)
	hintProcessor := &hints.CairoVmHintProcessor{}
	expectedHintDataMap := make(map[uint][]any)
	hintDataMap, err := runner.BuildHintDataMap(hintProcessor)
	if err != nil {
		t.Errorf("Test failed with error: %s", err)
	}
	if !reflect.DeepEqual(hintDataMap, expectedHintDataMap) {
		t.Errorf("Wrong hintDataMap, expected %+v, got %+v", expectedHintDataMap, hintDataMap)
	}
}

func TestBuildHintDataMapHappyPath(t *testing.T) {
	program := vm.Program{
		Hints: map[uint][]parser.HintParams{
			0: {
				{
					Code: "ids.a = ids.b",
					FlowTrackingData: parser.FlowTrackingData{
						APTracking: parser.ApTrackingData{Group: 1, Offset: 2},
					},
					ReferenceIds: map[string]uint{"a": 0, "b": 1},
				},
			},
		},
		ReferenceManager: parser.ReferenceManager{
			References: []parser.Reference{
				{
					Value: "cast(ap + (-2), felt)",
				},
				{
					Value: "cast(ap + (-1), felt)",
				},
			},
		},
	}
	runner, _ := runners.NewCairoRunner(program)
	hintProcessor := &hints.CairoVmHintProcessor{}
	expectedHintDataMap := map[uint][]any{
		0: []any{
			hints.HintData{
				Ids: map[string]hint_utils.HintReference{
					"a": {
						Offset1: hint_utils.OffsetValue{
							ValueType: hint_utils.Reference,
							Value:     -2,
						},
						ValueType: "felt",
					},
					"b": {
						Offset1: hint_utils.OffsetValue{
							ValueType: hint_utils.Reference,
							Value:     -1,
						},
						ValueType: "felt",
					},
				},
				Code:       "ids.a = ids.b",
				ApTracking: parser.ApTrackingData{Group: 1, Offset: 2},
			},
		},
	}
	hintDataMap, err := runner.BuildHintDataMap(hintProcessor)
	if err != nil {
		t.Errorf("Test failed with error: %s", err)
	}
	if !reflect.DeepEqual(hintDataMap, expectedHintDataMap) {
		t.Errorf("Wrong hintDataMap, expected %+v, got %+v", expectedHintDataMap, hintDataMap)
	}
}
