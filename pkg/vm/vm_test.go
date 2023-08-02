package vm_test

import (
	"bytes"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/lambdaclass/cairo-vm.go/pkg/vm"
	"github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
)

// Things we are skipping for now:
// - Initializing hint_executor and passing it to `cairo_run`
// - cairo_run_config stuff
// - Asserting expected trace values
// - Asserting memory_holes
func TestFibonacci(t *testing.T) {
	// compiledProgram := parser.Parse("../../cairo_programs/fibonacci.json")

	// TODO: Uncomment test when we have the bare minimum `CairoRun`
	// err := vm.CairoRun(compiledProgram.Data)
	// if err != nil {
	// 	t.Errorf("Program execution failed with error: %s", err)
	// }
}

func TestRelocateTraceOneEntry(t *testing.T) {
	virtualMachine := vm.NewVirtualMachine()
	buildTestProgramMemory(virtualMachine)

	virtualMachine.Segments.ComputeEffectiveSizes()
	relocationTable, _ := virtualMachine.Segments.RelocateSegments()
	err := virtualMachine.RelocateTrace(&relocationTable)
	if err != nil {
		t.Errorf("Trace relocation error failed with test: %s", err)
	}

	expectedTrace := []vm.TraceEntry{{Pc: 1, Ap: 4, Fp: 4}}
	actualTrace, err := virtualMachine.GetRelocatedTrace()
	if err != nil {
		t.Errorf("Trace relocation error failed with test: %s", err)
	}
	if !reflect.DeepEqual(expectedTrace, *actualTrace) {
		t.Errorf("Relocated trace and expected trace are not the same")
	}
}

func TestWriteBinaryTraceFile(t *testing.T) {
	tracePath, err := filepath.Abs("../../cairo_programs/trace_memory/cairo_trace_struct")
	if err != nil {
		t.Errorf("Trace file writing error failed with test: %s", err)
	}

	expectedTrace, err := ioutil.ReadFile(tracePath)
	if err != nil {
		t.Errorf("Trace file writing error failed with test: %s", err)
	}

	virtualMachine := vm.NewVirtualMachine()
	buildTestProgramMemory(virtualMachine)

	err = virtualMachine.Relocate()
	if err != nil {
		t.Errorf("Trace file writing error failed with test: %s", err)
	}

	relocatedTrace, err := virtualMachine.GetRelocatedTrace()
	if err != nil {
		t.Errorf("Trace file writing error failed with test: %s", err)
	}

	var actualTraceBuffer bytes.Buffer
	vm.WriteEncodedTrace(relocatedTrace, &actualTraceBuffer)

	if !reflect.DeepEqual(expectedTrace, actualTraceBuffer.Bytes()) {
		t.Errorf("Written trace and expected trace are not the same")
	}
}

func buildTestProgramMemory(virtualMachine *vm.VirtualMachine) {
	virtualMachine.Trace = []vm.TraceEntry{{Pc: 0, Ap: 2, Fp: 2}}
	for i := 0; i < 4; i++ {
		virtualMachine.Segments.AddSegment()
	}
	virtualMachine.Segments.Memory.Insert(memory.NewRelocatable(0, 0), memory.NewMaybeRelocatableInt(2345108766317314046))
	virtualMachine.Segments.Memory.Insert(memory.NewRelocatable(1, 0), memory.NewMaybeRelocatableRelocatable(memory.NewRelocatable(2, 0)))
	virtualMachine.Segments.Memory.Insert(memory.NewRelocatable(1, 1), memory.NewMaybeRelocatableRelocatable(memory.NewRelocatable(3, 0)))
}
