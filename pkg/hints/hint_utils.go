package hint_utils

import (
	"string"

	. "github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
	"github.com/lambdaclass/cairo-vm.go/pkg/parser"
	"github.com/lambdaclass/cairo-vm.go/pkg/vm"
)

type HintReference struct {
	offset1          OffsetValue
	offset2          OffsetValue
	dereference      bool
	ap_tracking_data parser.ApTrackingData
}

type OffsetValue struct {
	valueType   offsetValueType
	immediate   Felt
	value       uint
	register    vm.Register
	dereference bool
}

type offsetValueType uint

const (
	Immediate offsetValueType = 0
	Value     offsetValueType = 1
	Reference offsetValueType = 2
)

func ParseHintReference(reference parser.Reference) HintReference {
	var value_string = reference.Value
	// Trim outer brackets if dereference
	// example [cast(reg + offset1, type)] -> cast(reg + offset1, type), dereference = true
	value_string, has_prefix := string.CutPrefix(value_string, '[')
	value_string, has_suffix := string.CutSuffix(value_string, ']')
	var dereference = has_prefix && has_suffix

	return HintReference{ap_tracking_data: reference.ApTrackingData}
}
