package hints

import (
	"errors"
	"strings"

	. "github.com/lambdaclass/cairo-vm.go/pkg/hints/hint_utils"
	. "github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
	"github.com/lambdaclass/cairo-vm.go/pkg/parser"
	"github.com/lambdaclass/cairo-vm.go/pkg/types"
	"github.com/lambdaclass/cairo-vm.go/pkg/vm"
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

func (p *CairoVmHintProcessor) ExecuteHint(vm *vm.VirtualMachine, hintData *any, constants *map[string]Felt, executionScopes *types.ExecutionScopes) error {
	data, ok := (*hintData).(HintData)
	if !ok {
		return errors.New("Wrong Hint Data")
	}
	switch data.Code {
	case ASSERT_NN:
		return assert_nn(data.Ids, vm)
	case ASSERT_NOT_ZERO:
		return assert_not_zero(data.Ids, vm)
	case VM_EXIT_SCOPE:
		return vm_exit_scope(executionScopes)
	default:
		return errors.New("Unknown Hint")
	}
}
