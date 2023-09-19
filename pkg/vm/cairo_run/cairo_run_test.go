package cairo_run_test

import (
	"bytes"
	"testing"

	"github.com/lambdaclass/cairo-vm.go/pkg/vm/cairo_run"
)

func TestFibonacci(t *testing.T) {
	cairoRunConfig := cairo_run.CairoRunConfig{DisableTracePadding: false, Layout: "all_cairo", ProofMode: false}
	_, err := cairo_run.CairoRun("../../../cairo_programs/fibonacci.json", cairoRunConfig)
	if err != nil {
		t.Errorf("Program execution failed with error: %s", err)
	}
}

func TestFibonacciProofMode(t *testing.T) {
	cairoRunConfig := cairo_run.CairoRunConfig{DisableTracePadding: false, Layout: "all_cairo", ProofMode: true}
	_, err := cairo_run.CairoRun("../../../cairo_programs/proof_programs/fibonacci.json", cairoRunConfig)
	if err != nil {
		t.Errorf("Program execution failed with error: %s", err)
	}
}

func TestFactorial(t *testing.T) {
	cairoRunConfig := cairo_run.CairoRunConfig{DisableTracePadding: false, Layout: "all_cairo", ProofMode: false}
	_, err := cairo_run.CairoRun("../../../cairo_programs/factorial.json", cairoRunConfig)
	if err != nil {
		t.Errorf("Program execution failed with error: %s", err)
	}
}

func TestPoseidonBuiltin(t *testing.T) {
	cairoRunConfig := cairo_run.CairoRunConfig{DisableTracePadding: false, Layout: "all_cairo", ProofMode: false}

	_, err := cairo_run.CairoRun("../../../cairo_programs/poseidon_builtin.json", cairoRunConfig)
	if err != nil {
		t.Errorf("Program execution failed with error: %s", err)
	}
}

func TestPoseidonHash(t *testing.T) {
	cairoRunConfig := cairo_run.CairoRunConfig{DisableTracePadding: false, Layout: "all_cairo", ProofMode: false}

	_, err := cairo_run.CairoRun("../../../cairo_programs/poseidon_hash.json", cairoRunConfig)
	if err != nil {
		t.Errorf("Program execution failed with error: %s", err)
	}
}

func TestSimplePrint(t *testing.T) {
	cairoRunConfig := cairo_run.CairoRunConfig{DisableTracePadding: false, Layout: "all_cairo", ProofMode: false}

	_, err := cairo_run.CairoRun("../../../cairo_programs/simple_print.json", cairoRunConfig)
	if err != nil {
		t.Errorf("Program execution failed with error: %s", err)
	}
}

func TestWriteOutputProgram(t *testing.T) {
	cairoRunConfig := cairo_run.CairoRunConfig{DisableTracePadding: false, Layout: "all_cairo", ProofMode: false}
	runner, err := cairo_run.CairoRun("../../../cairo_programs/bitwise_output.json", cairoRunConfig)
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
	cairoRunConfig := cairo_run.CairoRunConfig{DisableTracePadding: false, ProofMode: false, Layout: "all_cairo"}
	_, err := cairo_run.CairoRun("../../../cairo_programs/pedersen_test.json", cairoRunConfig)
	if err != nil {
		t.Errorf("Program execution failed with error: %s", err)
	}
}

func TestPedersenAndBitwiseBuiltins(t *testing.T) {
	cairoRunConfig := cairo_run.CairoRunConfig{DisableTracePadding: false, ProofMode: false, Layout: "all_cairo"}
	_, err := cairo_run.CairoRun("../../../cairo_programs/pedersen_and_bitwise_builtins.json", cairoRunConfig)
	if err != nil {
		t.Errorf("Program execution failed with error: %s", err)
	}
}

func TestPedersenAndBitwiseBuiltinsWithAllocLocals(t *testing.T) {
	cairoRunConfig := cairo_run.CairoRunConfig{DisableTracePadding: false, ProofMode: false, Layout: "all_cairo"}
	_, err := cairo_run.CairoRun("../../../cairo_programs/pedersen_and_bitwise_builtins_with_alloc_locals.json", cairoRunConfig)
	if err != nil {
		t.Errorf("Program execution failed with error: %s", err)
	}
}

func TestAllocAddSegmentHint(t *testing.T) {
	cairoRunConfig := cairo_run.CairoRunConfig{DisableTracePadding: false, ProofMode: false, Layout: "all_cairo"}
	_, err := cairo_run.CairoRun("../../../cairo_programs/if_reloc_equal.json", cairoRunConfig)
	if err != nil {
		t.Errorf("Program execution failed with error: %s", err)
	}
}

func TestAssertNNHint(t *testing.T) {
	cairoRunConfig := cairo_run.CairoRunConfig{DisableTracePadding: false, ProofMode: false, Layout: "all_cairo"}
	_, err := cairo_run.CairoRun("../../../cairo_programs/assert_nn.json", cairoRunConfig)
	if err != nil {
		t.Errorf("Program execution failed with error: %s", err)
	}
}

func TestAbsValue(t *testing.T) {
	cairoRunConfig := cairo_run.CairoRunConfig{DisableTracePadding: false, ProofMode: false, Layout: "all_cairo"}
	_, err := cairo_run.CairoRun("../../../cairo_programs/abs_value.json", cairoRunConfig)
	if err != nil {
		t.Errorf("Program execution failed with error: %s", err)
	}
}
func TestAssertNotZeroHint(t *testing.T) {
	cairoRunConfig := cairo_run.CairoRunConfig{DisableTracePadding: false, ProofMode: false, Layout: "all_cairo"}
	_, err := cairo_run.CairoRun("../../../cairo_programs/assert_not_zero.json", cairoRunConfig)
	if err != nil {
		t.Errorf("Program execution failed with error: %s", err)
	}
}

func TestBitwiseRecursion(t *testing.T) {
	cairoRunConfig := cairo_run.CairoRunConfig{DisableTracePadding: false, Layout: "all_cairo", ProofMode: false}
	_, err := cairo_run.CairoRun("../../../cairo_programs/bitwise_recursion.json", cairoRunConfig)
	if err != nil {
		t.Errorf("Program execution failed with error: %s", err)
	}
}

func TestBitwiseRecursionProofMode(t *testing.T) {
	cairoRunConfig := cairo_run.CairoRunConfig{DisableTracePadding: false, Layout: "all_cairo", ProofMode: true}
	_, err := cairo_run.CairoRun("../../../cairo_programs/proof_programs/bitwise_recursion.json", cairoRunConfig)
	if err != nil {
		t.Errorf("Program execution failed with error: %s", err)
	}
}

func TestIsQuadResidueoHint(t *testing.T) {
	cairoRunConfig := cairo_run.CairoRunConfig{DisableTracePadding: false, Layout: "small", ProofMode: false}
	_, err := cairo_run.CairoRun("../../../cairo_programs/is_quad_residue.json", cairoRunConfig)
	if err != nil {
		t.Errorf("Program execution failed with error: %s", err)
	}
}

func TestIsQuadResidueoHintProofMode(t *testing.T) {
	cairoRunConfig := cairo_run.CairoRunConfig{DisableTracePadding: false, Layout: "small", ProofMode: true}
	_, err := cairo_run.CairoRun("../../../cairo_programs/proof_programs/is_quad_residue.json", cairoRunConfig)
	if err != nil {
		t.Errorf("Program execution failed with error: %s", err)
	}
}

func TestDict(t *testing.T) {
	cairoRunConfig := cairo_run.CairoRunConfig{DisableTracePadding: false, Layout: "small", ProofMode: false}
	_, err := cairo_run.CairoRun("../../../cairo_programs/dict.json", cairoRunConfig)
	if err != nil {
		t.Errorf("Program execution failed with error: %s", err)
	}
}

func TestDictUpdate(t *testing.T) {
	cairoRunConfig := cairo_run.CairoRunConfig{DisableTracePadding: false, Layout: "small", ProofMode: false}
	_, err := cairo_run.CairoRun("../../../cairo_programs/dict_update.json", cairoRunConfig)
	if err != nil {
		t.Errorf("Program execution failed with error: %s", err)
	}
}

func TestAssertNotEqualHint(t *testing.T) {
	cairoRunConfig := cairo_run.CairoRunConfig{DisableTracePadding: false, Layout: "all_cairo", ProofMode: false}
	_, err := cairo_run.CairoRun("../../../cairo_programs/assert_not_equal.json", cairoRunConfig)
	if err != nil {
		t.Errorf("Program execution failed with error: %s", err)
	}
}

func TestPowHint(t *testing.T) {
	cairoRunConfig := cairo_run.CairoRunConfig{DisableTracePadding: false, Layout: "all_cairo", ProofMode: false}
	_, err := cairo_run.CairoRun("../../../cairo_programs/pow.json", cairoRunConfig)
	if err != nil {
		t.Errorf("Program execution failed with error: %s", err)
	}
}

func TestSqrtHint(t *testing.T) {
	cairoRunConfig := cairo_run.CairoRunConfig{DisableTracePadding: false, Layout: "all_cairo", ProofMode: false}
	_, err := cairo_run.CairoRun("../../../cairo_programs/sqrt.json", cairoRunConfig)
	if err != nil {
		t.Errorf("Program execution failed with error: %s", err)
	}
}

func TestUnsignedDivRemHint(t *testing.T) {
	cairoRunConfig := cairo_run.CairoRunConfig{DisableTracePadding: false, Layout: "all_cairo", ProofMode: false}
	_, err := cairo_run.CairoRun("../../../cairo_programs/sqrt.json", cairoRunConfig)
	if err != nil {
		t.Errorf("Program execution failed with error: %s", err)
	}
}
