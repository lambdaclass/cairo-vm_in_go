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
	// The string should always follow with cast(...), so lets trim this part
	value_string = string.CutPrefix("cast(")
	value_string = string.CutSuffix(")")
	// TODO, use the second return to reject invalid strings
	// Now we should consider the possible cases
	//I. Inmediate value (number, type)
	//II. Reference (reference, type)
	// References can be made up of up to 2 offset values
	// The first offset value will always be register + offset1
	// and it can be dereferenced ([reg + off]) or not (ref + off)
	// The second offset can be a register and offset (reg + off), or just an an offset(off).
	// It can be dereferenced or not

	return HintReference{ap_tracking_data: reference.ApTrackingData}
}
