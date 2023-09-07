package hints

import (
	"errors"

	. "github.com/lambdaclass/cairo-vm.go/pkg/hints/hint_utils"
	. "github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
	"github.com/lambdaclass/cairo-vm.go/pkg/parser"
	"github.com/lambdaclass/cairo-vm.go/pkg/vm"
)

// HintProcessor Interface Definition

type HintProcessor interface {
	// Transforms hint data outputed by the VM into whichever format will be later used by ExecuteHint
	CompileHint(hintParams *parser.HintParams, referenceManager parser.ReferenceManager) (any, error)
	// Executes the hint which's data is provided by a dynamic structure previously created by CompileHint
	// TODO: add * ExecScopes args when ready
	ExecuteHint(vm *vm.VirtualMachine, hintData *any, constants map[string]Felt) error
}

// CairoVmHintProcessor

type HintData struct {
	Ids        map[string]HintReference
	Code       string
	ApTracking parser.ApTrackingData
}

type CairoVmHintProcessor struct {
}

func (p *CairoVmHintProcessor) CompileHint(hintParams *parser.HintParams, referenceManager parser.ReferenceManager) (any, error) {
	ids := make(map[string]HintReference, 0)
	for name, n := range hintParams.ReferenceIds {
		ids[name] = ParseHintReference(referenceManager.References[n])
	}
	return HintData{Ids: ids, Code: hintParams.Code, ApTracking: hintParams.FlowTrackingData.APTracking}, nil
}

func (p *CairoVmHintProcessor) ExecuteHint(vm *vm.VirtualMachine, hintData *any, constants map[string]Felt) error {
	data, ok := (*hintData).(HintData)
	if !ok {
		return errors.New("Wrong Hint Data")
	}
	switch data.Code {
	default:
		return errors.New("Uknown Hint")
	}
}
