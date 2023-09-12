package vm

import (
	"github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
	"github.com/lambdaclass/cairo-vm.go/pkg/parser"
)

type HintProcessor interface {
	// Transforms hint data outputed by the VM into whichever format will be later used by ExecuteHint
	CompileHint(hintParams *parser.HintParams, referenceManager *parser.ReferenceManager) (any, error)
	// Executes the hint which's data is provided by a dynamic structure previously created by CompileHint
	// TODO: add * ExecScopes arg when ready
	ExecuteHint(vm *VirtualMachine, hintData *any, constants *map[string]lambdaworks.Felt) error
}
