package runners_test

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/lambdaclass/cairo-vm.go/pkg/builtins"
	"github.com/lambdaclass/cairo-vm.go/pkg/hints"
	"github.com/lambdaclass/cairo-vm.go/pkg/hints/hint_utils"
	"github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
	"github.com/lambdaclass/cairo-vm.go/pkg/parser"
	"github.com/lambdaclass/cairo-vm.go/pkg/runners"
	"github.com/lambdaclass/cairo-vm.go/pkg/vm"
	"github.com/lambdaclass/cairo-vm.go/pkg/vm/cairo_run"
	"github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
)

func TestNewCairoRunnerInvalidBuiltin(t *testing.T) {
	// Create a Program with one fake instruction
	program_data := make([]memory.MaybeRelocatable, 1)
	empty_identifiers := make(map[string]vm.Identifier, 0)
	program_data[0] = *memory.NewMaybeRelocatableFelt(lambdaworks.FeltOne())
	program := vm.Program{Data: program_data, Builtins: []string{"fake_builtin"}, Identifiers: empty_identifiers}
	// Create CairoRunner
	_, err := runners.NewCairoRunner(program, "plain", false)
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
	runner, err := runners.NewCairoRunner(program, "plain", false)
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
	runner, err := runners.NewCairoRunner(program, "plain", false)
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
	program_builtins := []string{builtins.RANGE_CHECK_BUILTIN_NAME}
	program := vm.Program{Data: program_data, Identifiers: empty_identifiers, Builtins: program_builtins}
	// Create CairoRunner
	runner, err := runners.NewCairoRunner(program, "small", false)
	if err != nil {
		t.Errorf("NewCairoRunner error in test: %s", err)
	}
	// Initialize the runner
	_, err = runner.Initialize()
	if err != nil {
		t.Errorf("Initialize error in test: %s", err)
	}

	builtin_runner := runner.Vm.BuiltinRunners[0]
	if builtin_runner.Name() != builtins.RANGE_CHECK_BUILTIN_NAME {
		t.Errorf("Name of runner builtin failed. Expected %s, got %s", builtin_runner.Name(), builtins.RANGE_CHECK_BUILTIN_NAME)
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
	program_builtins := []string{builtins.RANGE_CHECK_BUILTIN_NAME}
	program := vm.Program{Data: program_data, Identifiers: empty_identifiers, Builtins: program_builtins}
	// Create CairoRunner
	runner, err := runners.NewCairoRunner(program, "small", false)
	if err != nil {
		t.Errorf("NewCairoRunner error in test: %s", err)
	}
	// Initialize the runner
	_, err = runner.Initialize()
	if err != nil {
		t.Errorf("Initialize error in test: %s", err)
	}

	builtin_runner := runner.Vm.BuiltinRunners[0]
	if builtin_runner.Name() != builtins.RANGE_CHECK_BUILTIN_NAME {
		t.Errorf("Name of runner builtin failed. Expected %s, got %s", builtin_runner.Name(), builtins.RANGE_CHECK_BUILTIN_NAME)
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

func TestIncludedBuiltinsPlainLayoutNoProofMode(t *testing.T) {
	cairoRunConfig := cairo_run.CairoRunConfig{DisableTracePadding: false, Layout: "small", ProofMode: false}
	// Testing for a program with no builtins
	factorialRunner, err := cairo_run.CairoRun("../../cairo_programs/factorial.json", cairoRunConfig)
	if err != nil {
		t.Errorf("Program execution failed with error: %s", err)
	}
	if len(factorialRunner.Vm.BuiltinRunners) != 0 {
		t.Errorf("The program should not have any builtins included, found %d", len(factorialRunner.Vm.BuiltinRunners))
	}

	// Testing with a program with output builtin
	printRunner, err := cairo_run.CairoRun("../../cairo_programs/simple_print.json", cairoRunConfig)
	if err != nil {
		t.Errorf("Program execution failed with error: %s", err)
	}
	if len(printRunner.Vm.BuiltinRunners) != 1 {
		t.Errorf("Expected only one builtin, found %d", len(factorialRunner.Vm.BuiltinRunners))
	}

	if printRunner.Vm.BuiltinRunners[0].Name() != "output" {
		t.Errorf("Expected output builtin, found: %s", printRunner.Vm.BuiltinRunners[0].Name())
	}
}

// FIXME: This test should changed once the `small` layout is properly implemented. ATM we don't have all
// its builtins implemented.
func TestIncludedBuiltinsSmallLayoutNoProofMode(t *testing.T) {
	cairoRunConfig := cairo_run.CairoRunConfig{DisableTracePadding: false, Layout: "all_cairo", ProofMode: false}
	// Testing for a program with Poseidon builtin
	poseidonRunner, err := cairo_run.CairoRun("../../cairo_programs/poseidon_builtin.json", cairoRunConfig)
	if err != nil {
		t.Errorf("Program execution failed with error: %s", err)
	}
	if len(poseidonRunner.Vm.BuiltinRunners) != 1 {
		t.Errorf("Expected only one builtin found: %d", len(poseidonRunner.Vm.BuiltinRunners))
	}
	if poseidonRunner.Vm.BuiltinRunners[0].Name() != "poseidon" {
		t.Errorf("Expected poseidon buitlin, found %s", poseidonRunner.Vm.BuiltinRunners[0].Name())
	}

	// Testing with a program with bitwise builtin
	bitwiseRunner, err := cairo_run.CairoRun("../../cairo_programs/bitwise_builtin_test.json", cairoRunConfig)
	if err != nil {
		t.Errorf("Program execution failed with error: %s", err)
	}
	if len(bitwiseRunner.Vm.BuiltinRunners) != 1 {
		t.Errorf("Expected only one builtin found: %d", len(bitwiseRunner.Vm.BuiltinRunners))
	}
	if bitwiseRunner.Vm.BuiltinRunners[0].Name() != "bitwise" {
		t.Errorf("Expected poseidon buitlin, found %s", bitwiseRunner.Vm.BuiltinRunners[0].Name())
	}

	// Testing with a program with output, pedersen and range_check builtins
	pedersenRunner, err := cairo_run.CairoRun("../../cairo_programs/pedersen_test.json", cairoRunConfig)
	if err != nil {
		t.Errorf("Program execution failed with error: %s", err)
	}
	if len(pedersenRunner.Vm.BuiltinRunners) != 3 {
		t.Errorf("Expected only one builtin found: %d", len(pedersenRunner.Vm.BuiltinRunners))
	}
	if pedersenRunner.Vm.BuiltinRunners[0].Name() != "output" {
		t.Errorf("Expected output buitlin, found %s", pedersenRunner.Vm.BuiltinRunners[0].Name())
	}
	if pedersenRunner.Vm.BuiltinRunners[1].Name() != "pedersen" {
		t.Errorf("Expected pedersen buitlin, found %s", pedersenRunner.Vm.BuiltinRunners[1].Name())
	}
	if pedersenRunner.Vm.BuiltinRunners[2].Name() != "range_check" {
		t.Errorf("Expected range_check buitlin, found %s", pedersenRunner.Vm.BuiltinRunners[2].Name())
	}
}

func TestBuildHintDataMapEmpty(t *testing.T) {
	program := vm.Program{}
	runner, _ := runners.NewCairoRunner(program, "plain", false)
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

func TestBuildHintDataMapOneHint(t *testing.T) {
	program := vm.Program{
		Hints: map[uint][]parser.HintParams{
			0: {
				{
					Code: "ids.a = ids.b",
					FlowTrackingData: parser.FlowTrackingData{
						APTracking:   parser.ApTrackingData{Group: 1, Offset: 2},
						ReferenceIds: map[string]uint{"a": 0, "b": 1},
					},
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
	runner, _ := runners.NewCairoRunner(program, "plain", false)
	hintProcessor := &hints.CairoVmHintProcessor{}
	expectedHintDataMap := map[uint][]any{
		0: {
			hints.HintData{
				Ids: hint_utils.IdsManager{
					References: map[string]hint_utils.HintReference{
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
					HintApTracking: parser.ApTrackingData{Group: 1, Offset: 2},
				},
				Code: "ids.a = ids.b",
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

func TestWriteOutputFromPresentMemory(t *testing.T) {
	empty_identifiers := make(map[string]vm.Identifier, 0)
	program_builtins := []string{builtins.OUTPUT_BUILTIN_NAME}
	program := vm.Program{Identifiers: empty_identifiers, Builtins: program_builtins}
	// Create CairoRunner
	runner, err := runners.NewCairoRunner(program, "plain", false)
	if err != nil {
		t.Errorf("NewCairoRunner error in test: %s", err)
	}
	// Initialize the runner
	_, err = runner.Initialize()
	if err != nil {
		t.Errorf("Initialize error in test: %s", err)
	}

	builtin := runner.Vm.BuiltinRunners[0]
	if builtin.Name() != builtins.OUTPUT_BUILTIN_NAME {
		t.Errorf("Wrong builtin name: expected: %s, got: %s", builtins.OUTPUT_BUILTIN_NAME, builtin.Name())
	}

	if builtin.Base().SegmentIndex != 2 {
		t.Errorf("Wrong builtin base: expected: %d, got: %d", 2, builtin.Base().SegmentIndex)
	}

	runner.Vm.Segments.Memory.Insert(memory.NewRelocatable(2, 0), memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(1)))
	runner.Vm.Segments.Memory.Insert(memory.NewRelocatable(2, 1), memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(2)))

	var buffer bytes.Buffer
	runner.Vm.WriteOutput(&buffer)

	expected := "1\n2\n"
	result := buffer.String()

	if expected != result {
		t.Errorf("TestWriteOutputFromPresentMemory failed. Expected: %s, got: %s", expected, result)
	}
}

func TestWriteOutputFromProgramGapRelocatableOutput(t *testing.T) {
	program_data := make([]memory.MaybeRelocatable, 4)
	program_data[0] = *memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(4612671187288162301))
	program_data[1] = *memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(5198983563776458752))
	program_data[2] = *memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(2))
	program_data[3] = *memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(2345108766317314046))
	empty_identifiers := make(map[string]vm.Identifier, 0)
	program_builtins := []string{builtins.OUTPUT_BUILTIN_NAME}
	program := vm.Program{Data: program_data, Identifiers: empty_identifiers, Builtins: program_builtins}
	// Create CairoRunner
	runner, err := runners.NewCairoRunner(program, "plain", false)
	if err != nil {
		t.Errorf("NewCairoRunner error in test: %s", err)
	}
	// Initialize the runner
	end, err := runner.Initialize()
	if err != nil {
		t.Errorf("Initialize error in test: %s", err)
	}
	hintProcessor := &hints.CairoVmHintProcessor{}
	err = runner.RunUntilPC(end, hintProcessor)
	if err != nil {
		t.Errorf("RunUntilPC error in test: %s", err)
	}

	var buffer bytes.Buffer
	runner.Vm.WriteOutput(&buffer)

	expected := "<missing>\n{2:0}\n"
	result := buffer.String()

	if expected != result {
		t.Errorf("TestWriteOutputFromProgramGapRelocatableOutput failed. Expected: %s, got: %s", expected, result)
	}
}

func TestWriteOutputFromPresentMemoryNegOutput(t *testing.T) {
	empty_identifiers := make(map[string]vm.Identifier, 0)
	program_builtins := []string{builtins.OUTPUT_BUILTIN_NAME}
	program := vm.Program{Identifiers: empty_identifiers, Builtins: program_builtins}
	// Create CairoRunner
	runner, err := runners.NewCairoRunner(program, "plain", false)
	if err != nil {
		t.Errorf("NewCairoRunner error in test: %s", err)
	}
	// Initialize the runner
	_, err = runner.Initialize()
	if err != nil {
		t.Errorf("Initialize error in test: %s", err)
	}

	builtin := runner.Vm.BuiltinRunners[0]
	if builtin.Name() != builtins.OUTPUT_BUILTIN_NAME {
		t.Errorf("Wrong builtin name: expected: %s, got: %s", builtins.OUTPUT_BUILTIN_NAME, builtin.Name())
	}

	if builtin.Base().SegmentIndex != 2 {
		t.Errorf("Wrong builtin base: expected: %d, got: %d", 2, builtin.Base().SegmentIndex)
	}

	runner.Vm.Segments.Memory.Insert(memory.NewRelocatable(2, 0), memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromDecString("-1")))

	var buffer bytes.Buffer
	runner.Vm.WriteOutput(&buffer)

	expected := "-1\n"
	result := buffer.String()

	if expected != result {
		t.Errorf("TestWriteOutputFromPresentMemoryNegOutput failed. Expected: %s, got: %s", expected, result)
	}
}

// Todo: Uncomment when we can add main entrypoint to program
/*func TestWriteOutputUnorderedBuiltins(t *testing.T) {
	program_data := make([]memory.MaybeRelocatable, 14)
	program_data[0] = *memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(4612671182993129469))
	program_data[1] = *memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(5198983563776458752))
	program_data[2] = *memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(1))
	program_data[3] = *memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(2345108766317314046))
	program_data[4] = *memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(5191102247248822272))
	program_data[5] = *memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(5189976364521848832))
	program_data[6] = *memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(1))
	program_data[7] = *memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(1226245742482522112))
	program_data[8] = *memory.NewMaybeRelocatableRelocatable(memory.NewRelocatable(-7, 10))
	program_data[9] = *memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(5189976364521848832))
	program_data[10] = *memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(17))
	program_data[11] = *memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(1226245742482522112))
	program_data[12] = *memory.NewMaybeRelocatableRelocatable(memory.NewRelocatable(-11, 10))
	program_data[13] = *memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(2345108766317314046))
	empty_identifiers := make(map[string]vm.Identifier, 0)
	program_builtins := []string{builtins.OUTPUT_BUILTIN_NAME, builtins.BITWISE_BUILTIN_NAME}
	program := vm.Program{Data: program_data, Identifiers: empty_identifiers, Builtins: program_builtins}
	// Create CairoRunner
	runner, err := runners.NewCairoRunner(program)
	if err != nil {
		t.Errorf("NewCairoRunner error in test: %s", err)
	}

	// Initialize the runner
	end, err := runner.Initialize()
	if err != nil {
		t.Errorf("Initialize error in test: %s", err)
	}

	runner.Vm.BuiltinRunners[0], runner.Vm.BuiltinRunners[1] = runner.Vm.BuiltinRunners[1], runner.Vm.BuiltinRunners[0]

	err = runner.RunUntilPC(end)
	if err != nil {
		t.Errorf("RunUntilPC error in test: %s", err)
	}
	var buffer bytes.Buffer
	runner.Vm.WriteOutput(&buffer)

	expected := ""
	result := buffer.String()

	if expected != result {
		t.Errorf("TestWriteOutputFromPresentMemoryNegOutput failed. Expected: %s, got: %s", expected, result)
	}
}
*/

func TestCheckRangeCheckUsagePermRangeLimitsNone(t *testing.T) {
	program := vm.Program{Data: nil, Builtins: nil, Identifiers: nil, Hints: nil, ReferenceManager: parser.ReferenceManager{}}

	runner, err := runners.NewCairoRunner(program, "plain", false)
	if err != nil {
		t.Error("Could not initialize Cairo Runner")
	}

	runner.Vm.Trace = make([]vm.TraceEntry, 0)

	err = runner.CheckRangeCheckUsage()
	if err != nil {
		t.Errorf("Check Range Usage Failed With Error %s", err)
	}
}

func TestCheckRangeCheckUsageWithoutBuiltins(t *testing.T) {
	program := vm.Program{Data: nil, Builtins: nil, Identifiers: nil, Hints: nil, ReferenceManager: parser.ReferenceManager{}}

	runner, err := runners.NewCairoRunner(program, "plain", false)
	if err != nil {
		t.Error("Could not initialize Cairo Runner")
	}

	runner.Vm.Trace = make([]vm.TraceEntry, 0)
	runner.Vm.CurrentStep = 1000
	runner.Vm.Segments.Memory.Insert(
		memory.NewRelocatable(0, 0),
		memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromHex("0x80FF80000530")),
	)

	runner.Vm.Trace = make([]vm.TraceEntry, 1)
	runner.Vm.Trace[0] = vm.TraceEntry{Pc: memory.NewRelocatable(0, 0), Ap: memory.NewRelocatable(0, 0), Fp: memory.NewRelocatable(0, 0)}
	err = runner.CheckRangeCheckUsage()
	if err != nil {
		t.Errorf("Check Range Usage Failed With Error %s", err)
	}
}

func TestCheckRangeUsageInsufficientAllocatedCells(t *testing.T) {
	program := vm.Program{Data: nil, Builtins: nil, Identifiers: nil, Hints: nil, ReferenceManager: parser.ReferenceManager{}}

	runner, err := runners.NewCairoRunner(program, "plain", false)
	if err != nil {
		t.Error("Could not initialize Cairo Runner")
	}
	runner.Vm.Trace = make([]vm.TraceEntry, 0)
	builtin := builtins.DefaultRangeCheckBuiltinRunner()
	runner.Vm.BuiltinRunners = append(runner.Vm.BuiltinRunners, builtin)
	runner.Vm.Segments.AddSegment()
	runner.Vm.Segments.Memory.Insert(
		memory.NewRelocatable(0, 0),
		memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromHex("0x80FF80000530")),
	)

	runner.Vm.Trace = make([]vm.TraceEntry, 1)
	runner.Vm.Trace[0] = vm.TraceEntry{Pc: memory.NewRelocatable(0, 0), Ap: memory.NewRelocatable(0, 0), Fp: memory.NewRelocatable(0, 0)}
	runner.Vm.Segments.ComputeEffectiveSizes()
	err = runner.CheckRangeCheckUsage()
	if err == nil {
		t.Error("Check Range Usage Should Have Failed With Insufficient Allocated Cells Error")
	}
}

func TestCheckDilutedCheckUsageWithoutPoolInstance(t *testing.T) {
	program := vm.Program{Data: nil, Builtins: nil, Identifiers: nil, Hints: nil, ReferenceManager: parser.ReferenceManager{}}

	runner, err := runners.NewCairoRunner(program, "plain", false)
	if err != nil {
		t.Error("Could not initialize Cairo Runner")
	}

	runner.Layout.DilutedPoolInstance = nil

	err = runner.CheckDilutedCheckUsage()
	if err != nil {
		t.Errorf("Check Diluted Check Usage Failed With Error %s", err)
	}
}

func TestCheckDilutedCheckUsageWithoutBuiltinRunners(t *testing.T) {
	program := vm.Program{Data: nil, Builtins: nil, Identifiers: nil, Hints: nil, ReferenceManager: parser.ReferenceManager{}}

	runner, err := runners.NewCairoRunner(program, "plain", false)
	if err != nil {
		t.Error("Could not initialize Cairo Runner")
	}

	runner.Vm.CurrentStep = 10000
	runner.Vm.BuiltinRunners = make([]builtins.BuiltinRunner, 0)

	err = runner.CheckDilutedCheckUsage()
	if err != nil {
		t.Errorf("Check Diluted Check Usage Failed With Error %s", err)
	}
}

func TestCheckDilutedCheckUsageInsufficientAllocatedCells(t *testing.T) {
	program := vm.Program{Data: nil, Builtins: nil, Identifiers: nil, Hints: nil, ReferenceManager: parser.ReferenceManager{}}

	runner, err := runners.NewCairoRunner(program, "all_cairo", false)
	if err != nil {
		t.Error("Could not initialize Cairo Runner")
	}

	runner.Vm.CurrentStep = 100
	runner.Vm.BuiltinRunners = make([]builtins.BuiltinRunner, 0)

	err = runner.CheckDilutedCheckUsage()
	if err == nil {
		t.Errorf("Check Diluted Check Usage Should Have failed With Insufficient Allocated Cells Error")
	}
}

func TestCheckDilutedCheckUsage(t *testing.T) {
	program := vm.Program{Data: nil, Builtins: nil, Identifiers: nil, Hints: nil, ReferenceManager: parser.ReferenceManager{}}

	runner, err := runners.NewCairoRunner(program, "plain", false)
	if err != nil {
		t.Error("Could not initialize Cairo Runner")
	}

	runner.Vm.CurrentStep = 8192
	runner.Vm.BuiltinRunners = make([]builtins.BuiltinRunner, 0)
	runner.Vm.BuiltinRunners = append(runner.Vm.BuiltinRunners, builtins.NewBitwiseBuiltinRunner(256))

	err = runner.CheckDilutedCheckUsage()
	if err != nil {
		t.Errorf("Check Diluted Check Usage Failed With Error %s", err)
	}
}

// This test is a huge meme, revisit
func TestCheckUsedCellsDilutedCheckUsageError(t *testing.T) {
	program := vm.Program{Data: nil, Builtins: nil, Identifiers: nil, Hints: nil, ReferenceManager: parser.ReferenceManager{}}

	runner, err := runners.NewCairoRunner(program, "plain", false)
	if err != nil {
		t.Error("Could not initialize Cairo Runner")
	}

	runner.Vm.Segments.SegmentUsedSizes = make(map[uint]uint)
	runner.Vm.Segments.SegmentUsedSizes[0] = 4
	runner.Vm.Trace = []vm.TraceEntry{}

	err = runner.CheckUsedCells()
	if err == nil {
		t.Errorf("Check Used Cells Should Have failed With Insufficient Allocated Cells Error")
	}
}

func TestRunFibonacciGetExecutionResources(t *testing.T) {
	cairoRunConfig := cairo_run.CairoRunConfig{Layout: "all_cairo", ProofMode: false}
	runner, err := cairo_run.CairoRun("../../cairo_programs/fibonacci.json", cairoRunConfig)
	if err != nil {
		t.Errorf("Program execution failed with error: %s", err)
	}
	expectedExecutionResources := runners.ExecutionResources{
		NSteps:                  80,
		BuiltinsInstanceCounter: make(map[string]uint),
	}
	executionResources, _ := runner.GetExecutionResources()
	if !reflect.DeepEqual(executionResources, expectedExecutionResources) {
		t.Errorf("Wong ExecutionResources.\n Expected : %+v, got: %+v", expectedExecutionResources, executionResources)
	}
}

// This test will run the `fib` function in the fibonacci.json program
func TestRunFromEntryPointFibonacci(t *testing.T) {
	compiledProgram, _ := parser.Parse("../../cairo_programs/fibonacci.json")
	programJson := vm.DeserializeProgramJson(compiledProgram)

	entrypoint := programJson.Identifiers["__main__.fib"].PC
	args := []any{
		*memory.NewMaybeRelocatableFelt(lambdaworks.FeltOne()),
		*memory.NewMaybeRelocatableFelt(lambdaworks.FeltOne()),
		*memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint(10)),
	}
	runner, _ := runners.NewCairoRunner(programJson, "all_cairo", false)
	hintProcessor := hints.CairoVmHintProcessor{}

	runner.InitializeBuiltins()
	runner.InitializeSegments()
	err := runner.RunFromEntrypoint(uint(entrypoint), args, &hintProcessor, nil, true, nil)

	if err != nil {
		t.Errorf("Running fib entrypoint failed with error %s", err.Error())
	}

	// Check result
	res, err := runner.Vm.GetReturnValues(1)
	if err != nil {
		t.Errorf("Failed to fetch return values from fib with error %s", err.Error())
	}
	if len(res) != 1 || !reflect.DeepEqual(res[0], *memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint(144))) {
		t.Errorf("Wrong value returned by fib entrypoint.\n Expected [144], got: %+v", res)
	}

}

func TestRunFromEntryPointFibonacciNotEnoughSteps(t *testing.T) {
	compiledProgram, _ := parser.Parse("../../cairo_programs/fibonacci.json")
	programJson := vm.DeserializeProgramJson(compiledProgram)

	entrypoint := programJson.Identifiers["__main__.fib"].PC
	args := []any{
		*memory.NewMaybeRelocatableFelt(lambdaworks.FeltOne()),
		*memory.NewMaybeRelocatableFelt(lambdaworks.FeltOne()),
		*memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint(10)),
	}
	runner, _ := runners.NewCairoRunner(programJson, "all_cairo", false)
	hintProcessor := hints.CairoVmHintProcessor{}

	runner.InitializeBuiltins()
	runner.InitializeSegments()
	runResources := vm.NewRunResources(2)
	err := runner.RunFromEntrypoint(uint(entrypoint), args, &hintProcessor, &runResources, true, nil)

	if err == nil {
		t.Errorf("Running fib entrypoint should have failed")
	}

}
