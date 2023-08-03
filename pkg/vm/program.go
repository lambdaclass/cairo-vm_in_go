package vm

import (
	"github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
	"github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
)

type ApTracking struct {
	Group  uint64
	Offset uint64
}

type FlowTrackingData struct {
	ApTracking   ApTracking
	ReferenceIDs map[string]uint64 `json:"reference_ids"`
}

type HintParams struct {
	Code             string
	AccessibleScopes []string `json:"accesible_scopes"`
	FlowTrackingData FlowTrackingData
}

type SharedProgramData struct {
	Data                   []memory.MaybeRelocatable
	Hints                  []HintParams
	HintsRange            []HintRange
	Main                   *int
	Start, end             *int
	errorMessageAttributes []Attribute
	instructionLocations   map[int]InstructionLocation
	identifiers            map[string]Identifier
	referenceManager       []HintReference
}

type Program struct {
	sharedProgramData *SharedProgramData
	constants         map[string]lambdaworks.Felt
	builtins          []BuiltinName
}

type BuiltinName struct {
	// Your BuiltinName fields here, if any
	// Example:
	// anotherField float64
}
