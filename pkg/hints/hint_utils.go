package hint_utils

import (
	"fmt"
	"strconv"
	"strings"

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
	value       int
	register    vm.Register
	dereference bool
}

type offsetValueType uint

const (
	Value     offsetValueType = 0
	Immediate offsetValueType = 1
	Reference offsetValueType = 2
)

func ParseHintReferenceOld(reference parser.Reference) HintReference {
	var value_string = reference.Value
	// Trim outer brackets if dereference
	// example [cast(reg + offset1, type)] -> cast(reg + offset1, type), dereference = true
	value_string, has_prefix := strings.CutPrefix(value_string, "[")
	value_string, has_suffix := strings.CutSuffix(value_string, "]")
	var dereference = has_prefix && has_suffix
	// The string should always follow with cast(...), so lets trim this part
	value_string, _ = strings.CutPrefix(value_string, "cast(")
	value_string, _ = strings.CutSuffix(value_string, ")")
	// TODO, use the second return to reject invalid strings
	// Now we should consider the possible cases
	//I. Inmediate value (number, type)
	//II. Reference (reference, type)
	// References can be made up of up to 2 offset values
	// The first offset value will always be register + offset1
	// and it can be dereferenced ([reg + off]) or not (ref + off)
	// The second offset can be a register and offset (reg + off), or just an an offset(off).
	// It can be dereferenced or not
	value_string, has_prefix = strings.CutPrefix(value_string, "[")
	if has_prefix {
		// TODO Handle dereference
	} else {
		// No dereference
		var next_components = strings.Split(value_string, " ")
		var register vm.Register
		switch next_components[0] {
		case "ap":
			register = vm.AP
		case "fp":
			register = vm.FP
		// Immediate
		default:
			var felt = FeltFromDecString(next_components[0])
			return HintReference{ap_tracking_data: reference.ApTrackingData, offset1: OffsetValue{immediate: felt, valueType: Immediate}}
		}
		// Here we should have something of the type reg + off/ reg + off + off
		// Beware that the offsets can be positive (num) or negative ((-num))
		if next_components[1] != "+" {
			return HintReference{}
		}
		// Handle first offset
		// check if its negative
		var offset1_value uint
		value_string, has_prefix = strings.CutPrefix(next_components[2], "(-")
		if has_prefix {
			value_string, _ = strings.CutSuffix(value_string, ")")

		}

	}

	return HintReference{ap_tracking_data: reference.ApTrackingData}
}

func ParseHintReference(reference parser.Reference) HintReference {
	var value_string = reference.Value
	// Trim outer brackets if dereference
	// example [cast(reg + offset1, type)] -> cast(reg + offset1, type), dereference = true
	value_string, has_prefix := strings.CutPrefix(value_string, "[")
	value_string, has_suffix := strings.CutSuffix(value_string, "]")
	var dereference = has_prefix && has_suffix
	var value_type string
	var immediate string
	// Scan various types till a match is found
	// Immediate: cast(number, type)
	_, err := fmt.Scanf(value_string, "cast(%s, %s)", &immediate, &value_type)
	if err == nil {
		var felt = FeltFromDecString(immediate)
		return HintReference{ap_tracking_data: reference.ApTrackingData, offset1: OffsetValue{immediate: felt, valueType: Immediate}}
	}
	var off_1_reg_0 string
	var off_1_reg_1 string
	var off1 string
	// Reference no deref 1 offset: cast(reg + off)
	_, err = fmt.Scanf(value_string, "cast(%c%c + %s)", off_1_reg_0, off_1_reg_1, off1)
	if err == nil {
		off1_reg := getRegister(off_1_reg_0, off_1_reg_1)
		num_off1 := offsetValueFromString(off1)
		return HintReference{ap_tracking_data: reference.ApTrackingData, offset1: OffsetValue{valueType: Reference, register: off1_reg, value: num_off1}, dereference: dereference}
	}

	return HintReference{ap_tracking_data: reference.ApTrackingData}
}

func offsetValueFromString(num string) int {
	value_string, has_prefix := strings.CutPrefix(num, "(-")
	if has_prefix {
		value_string, _ = strings.CutSuffix(value_string, ")")
	}
	res, _ := strconv.ParseInt(num, 0, 32)
	return int(res)
}

func getRegister(reg_0 string, reg_1 string) vm.Register {
	reg := vm.AP
	if reg_0 == "f" && reg_1 == "p" {
		reg = vm.FP
	}
	return reg
}
