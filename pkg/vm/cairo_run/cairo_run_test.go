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

func TestPoseidonBuiltinProofMode(t *testing.T) {
	testProgramProof("poseidon_builtin", t)
}

func TestPoseidonHash(t *testing.T) {
	testProgram("poseidon_hash", t)
}

func TestPoseidonHashProofMode(t *testing.T) {
	testProgramProof("poseidon_hash", t)
}

func TestSimplePrint(t *testing.T) {
	testProgram("simple_print", t)
}

func TestSimplePrintProofMode(t *testing.T) {
	testProgramProof("simple_print", t)
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

func TestPedersenBuiltinProofMode(t *testing.T) {
	testProgramProof("pedersen_test", t)
}

func TestPedersenAndBitwiseBuiltins(t *testing.T) {
	testProgram("pedersen_and_bitwise_builtins", t)
}

func TestPedersenAndBitwiseBuiltinsProofMode(t *testing.T) {
	testProgramProof("pedersen_and_bitwise_builtins", t)
}

func TestPedersenAndBitwiseBuiltinsWithAllocLocals(t *testing.T) {
	testProgram("pedersen_and_bitwise_builtins_with_alloc_locals", t)
}

func TestPedersenAndBitwiseBuiltinsWithAllocLocalsProofMode(t *testing.T) {
	testProgramProof("pedersen_and_bitwise_builtins_with_alloc_locals", t)
}

func TestAllocAddSegmentHint(t *testing.T) {
	testProgram("if_reloc_equal", t)
}

func TestAllocAddSegmentHintProofMode(t *testing.T) {
	testProgramProof("if_reloc_equal", t)
}

func TestAssertNNHint(t *testing.T) {
	testProgram("assert_nn", t)
}

func TestAssertNNHintProofMode(t *testing.T) {
	testProgramProof("assert_nn", t)
}

func TestAbsValue(t *testing.T) {
	testProgram("abs_value", t)
}

func TestAbsValueProofMode(t *testing.T) {
	testProgramProof("abs_value", t)
}

func TestCommonSignature(t *testing.T) {
	testProgram("common_signature", t)
}

func TestCommonSignatureProofMode(t *testing.T) {
	testProgramProof("common_signature", t)
}

func TestAssertNotZeroHint(t *testing.T) {
	testProgram("assert_not_zero", t)
}

func TestAssertNotZeroHintProofMode(t *testing.T) {
	testProgramProof("assert_not_zero", t)
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

func TestDictProofMode(t *testing.T) {
	testProgramProof("dict", t)
}

func TestDictUpdate(t *testing.T) {
	testProgram("dict_update", t)
}

func TestDictUpdateProofMode(t *testing.T) {
	testProgramProof("dict_update", t)
}

func TestAssertNotEqualHint(t *testing.T) {
	testProgram("assert_not_equal", t)
}

func TestAssertNotEqualHintProofMode(t *testing.T) {
	testProgramProof("assert_not_equal", t)
}

func TestPowHint(t *testing.T) {
	testProgram("pow", t)
}

func TestPowHintProofMode(t *testing.T) {
	testProgramProof("pow", t)
}

func TestSqrtHint(t *testing.T) {
	testProgram("sqrt", t)
}

func TestSqrtHintProofMode(t *testing.T) {
	testProgramProof("sqrt", t)
}

func TestUnsafeKeccak(t *testing.T) {
	testProgram("unsafe_keccak", t)
}

func TestUnsafeKeccakProofMode(t *testing.T) {
	testProgramProof("unsafe_keccak", t)
}

func TestUnsafeKeccakFinalize(t *testing.T) {
	testProgram("unsafe_keccak_finalize", t)
}

func TestUnsafeKeccakFinalizeProofMode(t *testing.T) {
	testProgramProof("unsafe_keccak_finalize", t)
}

func TestUnsignedDivRemHint(t *testing.T) {
	testProgram("unsigned_div_rem", t)
}

func TestUnsignedDivRemHintProofMode(t *testing.T) {
	testProgramProof("unsigned_div_rem", t)
}

func TestSignedDivRemHint(t *testing.T) {
	testProgram("signed_div_rem", t)
}

func TestSignedDivRemHintProofMode(t *testing.T) {
	testProgramProof("signed_div_rem", t)
}

func TestSetAddHint(t *testing.T) {
	testProgram("set_add", t)
}

func TestMemcpyHint(t *testing.T) {
	testProgram("memcpy_test", t)
}

func TestMemcpyHintProofMode(t *testing.T) {
	testProgramProof("memcpy_test", t)
}

func TestAssertLeFelt(t *testing.T) {
	testProgram("assert_le_felt", t)
}

func TestAssertLeFeltProofMode(t *testing.T) {
	testProgramProof("assert_le_felt", t)
}

func TestAssertLtFelt(t *testing.T) {
	testProgram("assert_lt_felt", t)
}

func TestAssertLtFeltProofMode(t *testing.T) {
	testProgramProof("assert_lt_felt", t)
}

func TestMemsetHint(t *testing.T) {
	testProgram("memset", t)
}

func TestMemsetHintProofMode(t *testing.T) {
	testProgramProof("memset", t)
}

func TestMathCmp(t *testing.T) {
	testProgram("math_cmp", t)
}

func TestMathCmpProofMode(t *testing.T) {
	testProgramProof("math_cmp", t)
}

func TestSquashDict(t *testing.T) {
	testProgram("squash_dict", t)
}

func TestSquashDictProofMode(t *testing.T) {
	testProgramProof("squash_dict", t)
}

func TestFindElementHint(t *testing.T) {
	testProgram("find_element", t)
}

func TestSearchSortedLowerHint(t *testing.T) {
	testProgram("search_sorted_lower", t)
}

func TestAssert250BitHint(t *testing.T) {
	testProgram("assert_250_bit_element_array", t)
}

func TestAssert250BitHintProofMode(t *testing.T) {
	testProgramProof("assert_250_bit_element_array", t)
}

func TestDictSquash(t *testing.T) {
	testProgram("dict_squash", t)
}

func TestDictSquashProofMode(t *testing.T) {
	testProgramProof("dict_squash", t)
}

func TestSplitFeltHint(t *testing.T) {
	testProgram("split_felt", t)
}

func TestUsort(t *testing.T) {
	testProgram("usort", t)
}

func TestUsortProofMode(t *testing.T) {
	testProgramProof("usort", t)
}

func TestSplitFeltHintProofMode(t *testing.T) {
	testProgramProof("split_felt", t)
}

func TestSplitIntHint(t *testing.T) {
	testProgram("split_int", t)
}

func TestSplitIntHintProofMode(t *testing.T) {
	testProgramProof("split_int", t)
}

func TestEcDoubleAssign(t *testing.T) {
	testProgram("ec_double_assign", t)
}

func TestIntegrationEcDoubleSlope(t *testing.T) {
	testProgram("ec_double_slope", t)
}

func TestKeccakIntegrationTests(t *testing.T) {
	testProgram("keccak_integration_tests", t)
}

func TestCairoKeccak(t *testing.T) {
	testProgram("cairo_keccak", t)
}

func TestKeccakAddUint256(t *testing.T) {
	testProgram("keccak_add_uint256", t)
}
