package hints

import (
	"errors"

	. "github.com/lambdaclass/cairo-vm.go/pkg/hints/hint_utils"
	. "github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
	"github.com/lambdaclass/cairo-vm.go/pkg/parser"
	"github.com/lambdaclass/cairo-vm.go/pkg/vm"
)

type HintData struct {
	Ids        map[string]HintReference
	Code       string
	ApTracking parser.ApTrackingData
}

type CairoVmHintProcessor struct {
}

func (p *CairoVmHintProcessor) CompileHint(hintParams *parser.HintParams, referenceManager *parser.ReferenceManager) (any, error) {
	ids := make(map[string]HintReference, 0)
	for name, n := range hintParams.ReferenceIds {
		if int(n) >= len(referenceManager.References) {
			return nil, errors.New("Reference not found in ReferenceManager")
		}
		ids[name] = ParseHintReference(referenceManager.References[n])
	}
	return HintData{Ids: ids, Code: hintParams.Code, ApTracking: hintParams.FlowTrackingData.APTracking}, nil
}

func (p *CairoVmHintProcessor) ExecuteHint(vm *vm.VirtualMachine, hintData *any, constants *map[string]Felt) error {
	data, ok := (*hintData).(HintData)
	if !ok {
		return errors.New("Wrong Hint Data")
	}
	switch data.Code {
	default:
		return errors.New("Unknown Hint")
	}
}
