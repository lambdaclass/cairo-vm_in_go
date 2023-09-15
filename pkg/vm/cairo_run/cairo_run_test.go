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
	_, err := cairo_run.CairoRun("../../../cairo_programs/fibonacci.json", "small", false)
	if err != nil {
		t.Errorf("Program execution failed with error: %s", err)
	}
	fmt.Println(err)
}

func TestPoseidonBuiltin(t *testing.T) {
	_, err := cairo_run.CairoRun("../../../cairo_programs/poseidon_builtin.json", "small", false)
	if err != nil {
		t.Errorf("Program execution failed with error: %s", err)
	}
}

func TestPoseidonHash(t *testing.T) {
	_, err := cairo_run.CairoRun("../../../cairo_programs/poseidon_hash.json", "small", false)
	if err != nil {
		t.Errorf("Program execution failed with error: %s", err)
	}
}

func TestSimplePrint(t *testing.T) {
	_, err := cairo_run.CairoRun("../../../cairo_programs/simple_print.json", "small", false)
	if err != nil {
		t.Errorf("Program execution failed with error: %s", err)
	}
}

func TestWriteOutputProgram(t *testing.T) {
	runner, err := cairo_run.CairoRun("../../../cairo_programs/bitwise_output.json", "small", false)
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
	_, err := cairo_run.CairoRun("../../../cairo_programs/pedersen_test.json", "small", false)
	if err != nil {
		t.Errorf("Program execution failed with error: %s", err)
	}
}

func TestPedersenAndBitwiseBuiltins(t *testing.T) {
	_, err := cairo_run.CairoRun("../../../cairo_programs/pedersen_and_bitwise_builtins.json", "small", false)
	if err != nil {
		t.Errorf("Program execution failed with error: %s", err)
	}
}

func TestPedersenAndBitwiseBuiltinsWithAllocLocals(t *testing.T) {
	_, err := cairo_run.CairoRun("../../../cairo_programs/pedersen_and_bitwise_builtins_with_alloc_locals.json", "small", false)
	if err != nil {
		t.Errorf("Program execution failed with error: %s", err)
	}
}

func TestAllocAddSegmentHint(t *testing.T) {
	_, err := cairo_run.CairoRun("../../../cairo_programs/if_reloc_equal.json", "small", false)
	if err != nil {
		t.Errorf("Program execution failed with error: %s", err)
	}
}

func TestAssertNNHint(t *testing.T) {
	_, err := cairo_run.CairoRun("../../../cairo_programs/assert_nn.json", "small", false)
	if err != nil {
		t.Errorf("Program execution failed with error: %s", err)
	}
}

func TestAbsValue(t *testing.T) {
	_, err := cairo_run.CairoRun("../../../cairo_programs/abs_value.json", "small", false)
	if err != nil {
		t.Errorf("Program execution failed with error: %s", err)
	}
}
func TestAssertNotZeroHint(t *testing.T) {
	_, err := cairo_run.CairoRun("../../../cairo_programs/assert_not_zero.json", "small", false)
	if err != nil {
		t.Errorf("Program execution failed with error: %s", err)
	}
}

func TestIsQuadResidueoHint(t *testing.T) {
	_, err := cairo_run.CairoRun("../../../cairo_programs/is_quad_residue.json", "small", false)
	if err != nil {
		t.Errorf("Program execution failed with error: %s", err)
	}
}

func TestDict(t *testing.T) {
	_, err := cairo_run.CairoRun("../../../cairo_programs/dict.json", "small", false)
	if err != nil {
		t.Errorf("Program execution failed with error: %s", err)
	}
}

func TestAssertNotEqualHint(t *testing.T) {
	_, err := cairo_run.CairoRun("../../../cairo_programs/assert_not_equal.json", "small", false)
	if err != nil {
		t.Errorf("Program execution failed with error: %s", err)
	}
}
