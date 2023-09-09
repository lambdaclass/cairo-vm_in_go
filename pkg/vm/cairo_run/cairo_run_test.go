package cairo_run_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/lambdaclass/cairo-vm.go/pkg/vm/cairo_run"
)

// Things we are skipping for now:
// - Initializing hint_executor and passing it to `cairo_run`
// - cairo_run_config stuff
// - Asserting expected trace values
// - Asserting memory_holes
func TestFibonacci(t *testing.T) {
	_, err := cairo_run.CairoRun("../../../cairo_programs/fibonacci.json")
	if err != nil {
		t.Errorf("Program execution failed with error: %s", err)
	}
	fmt.Println(err)
}

func TestPoseidonBuiltin(t *testing.T) {
	_, err := cairo_run.CairoRun("../../../cairo_programs/poseidon_builtin.json")
	if err != nil {
		t.Errorf("Program execution failed with error: %s", err)
	}
}

func TestPoseidonHash(t *testing.T) {
	_, err := cairo_run.CairoRun("../../../cairo_programs/poseidon_hash.json")
	if err != nil {
		t.Errorf("Program execution failed with error: %s", err)
	}
}

func TestSimplePrint(t *testing.T) {
	_, err := cairo_run.CairoRun("../../../cairo_programs/simple_print.json")
	if err != nil {
		t.Errorf("Program execution failed with error: %s", err)
	}
}

func TestWriteOutputProgram(t *testing.T) {
	runner, err := cairo_run.CairoRun("../../../cairo_programs/bitwise_output.json")
	if err != nil {
		t.Errorf("Program execution failed with error: %s", err)
	}
	var buffer bytes.Buffer
	runner.Vm.WriteOutput(&buffer)

	expected := "0\n"
	result := buffer.String()

	if expected != result {
		t.Errorf("TestWriteOutputProgram failed. Expected: %s, got: %s", expected, result)
	}

}

func TestPedersenBuiltin(t *testing.T) {
	_, err := cairo_run.CairoRun("../../../cairo_programs/pedersen_hash.json")
	if err != nil {
		t.Errorf("Program execution failed with error: %s", err)
	}
}

func TestPedersenTest(t *testing.T) {
	_, err := cairo_run.CairoRun("../../../cairo_programs/pedersen_test.json")
	if err != nil {
		t.Errorf("Program execution failed with error: %s", err)
	}
}
