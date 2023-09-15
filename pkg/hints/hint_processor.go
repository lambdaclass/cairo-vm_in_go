package hints

import (
	"strings"

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
	ids := NewIdsManager(references, hintParams.FlowTrackingData.APTracking)
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
	case IS_POSITIVE:
		return is_positive(data.Ids, vm)
	case ASSERT_NOT_ZERO:
		return assert_not_zero(data.Ids, vm)
	case DEFAULT_DICT_NEW:
		return defaultDictNew(data.Ids, execScopes, vm)
	case DICT_READ:
		return dictRead(data.Ids, execScopes, vm)
	case DICT_WRITE:
		return dictWrite(data.Ids, execScopes, vm)
	case DICT_UPDATE:
		return dictUpdate(data.Ids, execScopes, vm)
	case ASSERT_NOT_EQUAL:
		return assert_not_equal(data.Ids, vm)
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
	case VM_EXIT_SCOPE:
		return vm_exit_scope(execScopes)
	default:
		return errors.Errorf("Unknown Hint: %s", data.Code)
	}
}
