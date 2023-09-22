package cairo_run_test

import (
	"bytes"
	"testing"

	"github.com/lambdaclass/cairo-vm.go/pkg/vm/cairo_run"
)

func testProgram(programName string, t *testing.T) {
	cairoRunConfig := cairo_run.CairoRunConfig{DisableTracePadding: false, Layout: "all_cairo", ProofMode: false}
	_, err := cairo_run.CairoRun("../../../cairo_programs/"+programName+".json", cairoRunConfig)
	if err != nil {
		t.Errorf("Program execution failed with error: %s", err)
	}
}

func testProgramProof(programName string, t *testing.T) {
	cairoRunConfig := cairo_run.CairoRunConfig{DisableTracePadding: false, Layout: "all_cairo", ProofMode: true}
	_, err := cairo_run.CairoRun("../../../cairo_programs/proof_programs/"+programName+".json", cairoRunConfig)
	if err != nil {
		t.Errorf("Program execution failed with error: %s", err)
	}
}

func TestFibonacci(t *testing.T) {
	testProgram("fibonacci", t)
}

func TestFibonacciProofMode(t *testing.T) {
	testProgramProof("fibonacci", t)
}

func TestFactorial(t *testing.T) {
	testProgram("factorial", t)
}

func TestPoseidonBuiltin(t *testing.T) {
	testProgram("poseidon_builtin", t)
}

func TestPoseidonHash(t *testing.T) {
	testProgram("poseidon_hash", t)
}

func TestSimplePrint(t *testing.T) {
	testProgram("simple_print", t)
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
	testProgram("pedersen_test", t)
}

func TestPedersenAndBitwiseBuiltins(t *testing.T) {
	testProgram("pedersen_and_bitwise_builtins", t)
}

func TestPedersenAndBitwiseBuiltinsWithAllocLocals(t *testing.T) {
	testProgram("pedersen_and_bitwise_builtins_with_alloc_locals", t)
}

func TestAllocAddSegmentHint(t *testing.T) {
	testProgram("if_reloc_equal", t)
}

func TestAssertNNHint(t *testing.T) {
	testProgram("assert_nn", t)
}

func TestAbsValue(t *testing.T) {
	testProgram("abs_value", t)
}
func TestCommonSignature(t *testing.T) {
	testProgram("common_signature", t)
}
func TestAssertNotZeroHint(t *testing.T) {
	testProgram("assert_not_zero", t)
}

func TestBitwiseRecursion(t *testing.T) {
	testProgram("bitwise_recursion", t)
}

func TestBitwiseRecursionProofMode(t *testing.T) {
	testProgramProof("bitwise_recursion", t)
}

func TestIsQuadResidueoHint(t *testing.T) {
	testProgram("is_quad_residue", t)
}

func TestIsQuadResidueoHintProofMode(t *testing.T) {
	testProgramProof("is_quad_residue", t)
}

func TestDict(t *testing.T) {
	testProgram("dict", t)
}

func TestDictUpdate(t *testing.T) {
	testProgram("dict_update", t)
}

func TestAssertNotEqualHint(t *testing.T) {
	testProgram("assert_not_equal", t)
}

func TestPowHint(t *testing.T) {
	testProgram("pow", t)
}

func TestSqrtHint(t *testing.T) {
	testProgram("sqrt", t)
}

func TestUnsafeKeccak(t *testing.T) {
	testProgram("unsafe_keccak", t)
}

func TestUnsafeKeccakFinalize(t *testing.T) {
	testProgram("unsafe_keccak_finalize", t)
}

func TestUnsignedDivRemHint(t *testing.T) {
	cairoRunConfig := cairo_run.CairoRunConfig{DisableTracePadding: false, Layout: "all_cairo", ProofMode: false}
	_, err := cairo_run.CairoRun("../../../cairo_programs/unsigned_div_rem.json", cairoRunConfig)
	if err != nil {
		t.Errorf("Program execution failed with error: %s", err)
	}
}

func TestMemcpyHint(t *testing.T) {
	testProgram("memcpy_test", t)
}

func TestAssertLeFelt(t *testing.T) {
	testProgram("assert_le_felt", t)
}

func TestAssertLtFelt(t *testing.T) {
	testProgram("assert_lt_felt", t)
}

func TestMemsetHint(t *testing.T) {
	testProgram("memset", t)
}

func TestMathCmp(t *testing.T) {
	testProgram("math_cmp", t)
}
func TestSquashDict(t *testing.T) {
	testProgram("squash_dict", t)
}

func TestSignedDivRemHint(t *testing.T) {
	cairoRunConfig := cairo_run.CairoRunConfig{DisableTracePadding: false, Layout: "all_cairo", ProofMode: false}
	_, err := cairo_run.CairoRun("../../../cairo_programs/signed_div_rem.json", cairoRunConfig)
	if err != nil {
		t.Errorf("Program execution failed with error: %s", err)
	}
}

func TestAssert250BitHint(t *testing.T) {
	testProgram("assert_250_bit_element_array", t)
}

func TestDictSquash(t *testing.T) {
	testProgram("dict_squash", t)
}

func TestSplitFeltHint(t *testing.T) {
	testProgram("split_felt", t)
}

func TestSplitIntHint(t *testing.T) {
	cairoRunConfig := cairo_run.CairoRunConfig{DisableTracePadding: false, Layout: "all_cairo", ProofMode: false}
	_, err := cairo_run.CairoRun("../../../cairo_programs/split_int.json", cairoRunConfig)
	if err != nil {
		t.Errorf("Program execution failed with error: %s", err)
	}
}

func TestFastEcAddAssignNewXHint(t *testing.T) {
	testProgram("fast_ec_add_v2", t)
}
