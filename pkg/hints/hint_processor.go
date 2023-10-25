package hints

import (
	"strings"

	. "github.com/lambdaclass/cairo-vm.go/pkg/hints/hint_codes"
	"github.com/lambdaclass/cairo-vm.go/pkg/hints/hint_utils"
	. "github.com/lambdaclass/cairo-vm.go/pkg/hints/hint_utils"
	. "github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
	"github.com/lambdaclass/cairo-vm.go/pkg/parser"
	"github.com/lambdaclass/cairo-vm.go/pkg/types"
	"github.com/lambdaclass/cairo-vm.go/pkg/vm"
	"github.com/pkg/errors"
)

type HintData struct {
	Ids  IdsManager
	Code string
}

type CairoVmHintProcessor struct {
}

func (p *CairoVmHintProcessor) CompileHint(hintParams *parser.HintParams, referenceManager *parser.ReferenceManager) (any, error) {
	references := make(map[string]HintReference, 0)
	for name, n := range hintParams.FlowTrackingData.ReferenceIds {
		if int(n) >= len(referenceManager.References) {
			return nil, errors.New("Reference not found in ReferenceManager")
		}
		split := strings.Split(name, ".")
		name = split[len(split)-1]
		references[name] = ParseHintReference(referenceManager.References[n])
	}
	ids := NewIdsManager(references, hintParams.FlowTrackingData.APTracking, hintParams.AccessibleScopes)
	return HintData{Ids: ids, Code: hintParams.Code}, nil
}

func (p *CairoVmHintProcessor) ExecuteHint(vm *vm.VirtualMachine, hintData *any, constants *map[string]Felt, execScopes *types.ExecutionScopes) error {
	data, ok := (*hintData).(HintData)
	if !ok {
		return errors.New("Wrong Hint Data")
	}
	switch data.Code {
	case ADD_SEGMENT:
		return add_segment(vm)
	case ASSERT_NN:
		return assert_nn(data.Ids, vm)
	case VERIFY_ECDSA_SIGNATURE:
		return verify_ecdsa_signature(data.Ids, vm)
	case IS_POSITIVE:
		return is_positive(data.Ids, vm)
	case ASSERT_NOT_ZERO:
		return assert_not_zero(data.Ids, vm)
	case IS_QUAD_RESIDUE:
		return is_quad_residue(data.Ids, vm)
	case DEFAULT_DICT_NEW:
		return defaultDictNew(data.Ids, execScopes, vm)
	case DICT_READ:
		return dictRead(data.Ids, execScopes, vm)
	case DICT_WRITE:
		return dictWrite(data.Ids, execScopes, vm)
	case DICT_UPDATE:
		return dictUpdate(data.Ids, execScopes, vm)
	case SQUASH_DICT:
		return squashDict(data.Ids, execScopes, vm)
	case SQUASH_DICT_INNER_SKIP_LOOP:
		return squashDictInnerSkipLoop(data.Ids, execScopes, vm)
	case SQUASH_DICT_INNER_FIRST_ITERATION:
		return squashDictInnerFirstIteration(data.Ids, execScopes, vm)
	case SQUASH_DICT_INNER_CHECK_ACCESS_INDEX:
		return squashDictInnerCheckAccessIndex(data.Ids, execScopes, vm)
	case SQUASH_DICT_INNER_CONTINUE_LOOP:
		return squashDictInnerContinueLoop(data.Ids, execScopes, vm)
	case SQUASH_DICT_INNER_ASSERT_LEN_KEYS:
		return squashDictInnerAssertLenKeys(execScopes)
	case SQUASH_DICT_INNER_LEN_ASSERT:
		return squashDictInnerLenAssert(execScopes)
	case SQUASH_DICT_INNER_USED_ACCESSES_ASSERT:
		return squashDictInnerUsedAccessesAssert(data.Ids, execScopes, vm)
	case SQUASH_DICT_INNER_NEXT_KEY:
		return squashDictInnerNextKey(data.Ids, execScopes, vm)
	case DICT_SQUASH_COPY_DICT:
		return dictSquashCopyDict(data.Ids, execScopes, vm)
	case DICT_SQUASH_UPDATE_PTR:
		return dictSquashUpdatePtr(data.Ids, execScopes, vm)
	case DICT_NEW:
		return dictNew(data.Ids, execScopes, vm)
	case VM_EXIT_SCOPE:
		return vm_exit_scope(execScopes)
	case ASSERT_NOT_EQUAL:
		return assert_not_equal(data.Ids, vm)
	case EC_NEGATE:
		return ecNegateImportSecpP(vm, *execScopes, data.Ids)
	case EC_NEGATE_EMBEDDED_SECP:
		return ecNegateEmbeddedSecpP(vm, *execScopes, data.Ids)
	case EC_DOUBLE_ASSIGN_NEW_X_V1, EC_DOUBLE_ASSIGN_NEW_X_V2, EC_DOUBLE_ASSIGN_NEW_X_V3, EC_DOUBLE_ASSIGN_NEW_X_V4:
		return ecDoubleAssignNewX(vm, *execScopes, data.Ids, SECP_P_V2())
	case POW:
		return pow(data.Ids, vm)
	case SQRT:
		return sqrt(data.Ids, vm)
	case MEMCPY_ENTER_SCOPE:
		return memcpy_enter_scope(data.Ids, vm, execScopes)
	case MEMSET_ENTER_SCOPE:
		return memset_enter_scope(data.Ids, vm, execScopes)
	case MEMCPY_CONTINUE_COPYING:
		return memset_step_loop(data.Ids, vm, execScopes, "continue_copying")
	case MEMSET_CONTINUE_LOOP:
		return memset_step_loop(data.Ids, vm, execScopes, "continue_loop")
	case VM_ENTER_SCOPE:
		return vm_enter_scope(execScopes)
	case USORT_ENTER_SCOPE:
		return usortEnterScope(execScopes)
	case USORT_BODY:
		return usortBody(data.Ids, execScopes, vm)
	case USORT_VERIFY:
		return usortVerify(data.Ids, execScopes, vm)
	case USORT_VERIFY_MULTIPLICITY_ASSERT:
		return usortVerifyMultiplicityAssert(execScopes)
	case USORT_VERIFY_MULTIPLICITY_BODY:
		return usortVerifyMultiplicityBody(data.Ids, execScopes, vm)
	case SET_ADD:
		return setAdd(data.Ids, vm)
	case FIND_ELEMENT:
		return findElement(data.Ids, vm, *execScopes)
	case SEARCH_SORTED_LOWER:
		return searchSortedLower(data.Ids, vm, *execScopes)
	case COMPUTE_SLOPE_V1:
		return computeSlopeAndAssingSecpP(vm, *execScopes, data.Ids, "point0", "point1", SECP_P())
	case COMPUTE_SLOPE_V2:
		return computeSlopeAndAssingSecpP(vm, *execScopes, data.Ids, "point0", "point1", SECP_P_V2())
	case COMPUTE_SLOPE_WHITELIST:
		return computeSlopeAndAssingSecpP(vm, *execScopes, data.Ids, "pt0", "pt1", SECP_P())
	case COMPUTE_SLOPE_SECP256R1:
		return computeSlope(vm, *execScopes, data.Ids, "point0", "point1")
	case EC_DOUBLE_SLOPE_V1:
		return computeDoublingSlope(vm, *execScopes, data.Ids, "point", SECP_P(), ALPHA())
	case UNSAFE_KECCAK:
		return unsafeKeccak(data.Ids, vm, *execScopes)
	case UNSAFE_KECCAK_FINALIZE:
		return unsafeKeccakFinalize(data.Ids, vm)
	case COMPARE_BYTES_IN_WORD_NONDET:
		return compareBytesInWordNondet(data.Ids, vm, constants)
	case COMPARE_KECCAK_FULL_RATE_IN_BYTES_NONDET:
		return compareKeccakFullRateInBytesNondet(data.Ids, vm, constants)
	case BLOCK_PERMUTATION:
		return blockPermutation(data.Ids, vm, constants)
	case CAIRO_KECCAK_FINALIZE_V1:
		return cairoKeccakFinalize(data.Ids, vm, constants, 10)
	case CAIRO_KECCAK_FINALIZE_V2:
		return cairoKeccakFinalize(data.Ids, vm, constants, 1000)
	case KECCAK_WRITE_ARGS:
		return keccakWriteArgs(data.Ids, vm)
	case UNSIGNED_DIV_REM:
		return unsignedDivRem(data.Ids, vm)
	case SIGNED_DIV_REM:
		return signedDivRem(data.Ids, vm)
	case ASSERT_LE_FELT:
		return assertLeFelt(data.Ids, vm, execScopes, constants)
	case ASSERT_LE_FELT_EXCLUDED_0:
		return assertLeFeltExcluded0(vm, execScopes)
	case ASSERT_LE_FELT_EXCLUDED_1:
		return assertLeFeltExcluded1(vm, execScopes)
	case ASSERT_LE_FELT_EXCLUDED_2:
		return assertLeFeltExcluded2(vm, execScopes)
	case ASSERT_LT_FELT:
		return assertLtFelt(data.Ids, vm)
	case IS_NN:
		return isNN(data.Ids, vm)
	case IS_NN_OUT_OF_RANGE:
		return isNNOutOfRange(data.Ids, vm)
	case IS_LE_FELT:
		return isLeFelt(data.Ids, vm)
	case ASSERT_250_BITS:
		return Assert250Bit(data.Ids, vm, constants)
	case SPLIT_FELT:
		return SplitFelt(data.Ids, vm, constants)
	case IMPORT_SECP256R1_ALPHA:
		return importSecp256r1Alpha(*execScopes)
	case IMPORT_SECP256R1_N:
		return importSECP256R1N(*execScopes)
	case IMPORT_SECP256R1_P:
		return importSECP256R1P(*execScopes)
	case EC_DOUBLE_SLOPE_EXTERNAL_CONSTS:
		return computeDoublingSlopeExternalConsts(*vm, *execScopes, data.Ids)
	case NONDET_BIGINT3_V1:
		return NondetBigInt3(*vm, *execScopes, data.Ids)
	case SPLIT_INT:
		return splitInt(data.Ids, vm)
	case SPLIT_INT_ASSERT_RANGE:
		return splitIntAssertRange(data.Ids, vm)
	case UINT256_ADD:
		return uint256Add(data.Ids, vm, false)
	case UINT256_ADD_LOW:
		return uint256Add(data.Ids, vm, true)
	case UINT256_SUB:
		return uint256Sub(data.Ids, vm)
	case SPLIT_64:
		return split64(data.Ids, vm)
	case UINT256_SQRT:
		return uint256Sqrt(data.Ids, vm, false)
	case UINT256_SQRT_FELT:
		return uint256Sqrt(data.Ids, vm, true)
	case UINT256_SIGNED_NN:
		return uint256SignedNN(data.Ids, vm)
	case UINT256_UNSIGNED_DIV_REM:
		return uint256UnsignedDivRem(data.Ids, vm)
	case UINT256_EXPANDED_UNSIGNED_DIV_REM:
		return uint256ExpandedUnsignedDivRem(data.Ids, vm)
	case UINT256_MUL_DIV_MOD:
		return uint256MulDivMod(data.Ids, vm)
	case DIV_MOD_N_PACKED_DIVMOD_V1:
		return divModNPackedDivMod(data.Ids, vm, execScopes)
	case DIV_MOD_N_PACKED_DIVMOD_EXTERNAL_N:
		return divModNPackedDivModExternalN(data.Ids, vm, execScopes)
	case XS_SAFE_DIV:
		return divModNSafeDiv(data.Ids, execScopes, "x", "s", false)
	case DIV_MOD_N_SAFE_DIV:
		return divModNSafeDiv(data.Ids, execScopes, "a", "b", false)
	case DIV_MOD_N_SAFE_DIV_PLUS_ONE:
		return divModNSafeDiv(data.Ids, execScopes, "a", "b", true)
	case GET_POINT_FROM_X:
		return getPointFromX(data.Ids, vm, execScopes, constants)
	case VERIFY_ZERO_EXTERNAL_SECP:
		return verifyZeroWithExternalConst(*vm, *execScopes, data.Ids)
	case FAST_EC_ADD_ASSIGN_NEW_X:
		return fastEcAddAssignNewX(data.Ids, vm, execScopes, "point0", "point1", SECP_P())
	case FAST_EC_ADD_ASSIGN_NEW_X_V2:
		return fastEcAddAssignNewX(data.Ids, vm, execScopes, "point0", "point1", SECP_P_V2())
	case FAST_EC_ADD_ASSIGN_NEW_X_V3:
		return fastEcAddAssignNewX(data.Ids, vm, execScopes, "pt0", "pt1", SECP_P())
	case FAST_EC_ADD_ASSIGN_NEW_Y:
		return fastEcAddAssignNewY(execScopes)
	case REDUCE_V1:
		return reduceV1(data.Ids, vm, execScopes)
	case REDUCE_V2:
		return reduceV2(data.Ids, vm, execScopes)
	case REDUCE_ED25519:
		return reduceED25519(data.Ids, vm, execScopes)
	case VERIFY_ZERO_V1, VERIFY_ZERO_V2:
		return verifyZero(data.Ids, vm, execScopes, hint_utils.SECP_P())
	case VERIFY_ZERO_V3:
		return verifyZero(data.Ids, vm, execScopes, hint_utils.SECP_P_V2())
	case BLAKE2S_COMPUTE:
		return blake2sCompute(data.Ids, vm)
	default:
		return errors.Errorf("Unknown Hint: %s", data.Code)
	}
}
