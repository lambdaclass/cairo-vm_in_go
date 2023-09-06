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
		return HintReference{
			ap_tracking_data: reference.ApTrackingData,
			offset1:          OffsetValue{immediate: felt, valueType: Immediate},
		}
	}
	var off_1_reg_0 string
	var off_1_reg_1 string
	var off1 string
	// Reference no deref 1 offset: cast(reg + off)
	_, err = fmt.Scanf(value_string, "cast(%c%c + %s)", off_1_reg_0, off_1_reg_1, off1)
	if err == nil {
		off1_reg := getRegister(off_1_reg_0, off_1_reg_1)
		num_off1 := offsetValueFromString(off1)
		return HintReference{
			ap_tracking_data: reference.ApTrackingData,
			offset1:          OffsetValue{valueType: Reference, register: off1_reg, value: num_off1},
			dereference:      dereference,
		}
	}
	var off2 string
	// Reference no deref 2 offsets: cast(reg + off1 + off2)
	_, err = fmt.Scanf(value_string, "cast(%c%c + %s + %s)", off_1_reg_0, off_1_reg_1, off1, off2)
	if err == nil {
		off1_reg := getRegister(off_1_reg_0, off_1_reg_1)
		num_off1 := offsetValueFromString(off1)
		num_off2 := offsetValueFromString(off2)
		return HintReference{
			ap_tracking_data: reference.ApTrackingData,
			offset1:          OffsetValue{valueType: Reference, register: off1_reg, value: num_off1},
			offset2:          OffsetValue{value: num_off2},
			dereference:      dereference,
		}
	}
	// Reference with deref 1 offset: cast([reg + off1])
	_, err = fmt.Scanf(value_string, "cast([%c%c + %s])", off_1_reg_0, off_1_reg_1, off1)
	if err == nil {
		off1_reg := getRegister(off_1_reg_0, off_1_reg_1)
		num_off1 := offsetValueFromString(off1)
		return HintReference{
			ap_tracking_data: reference.ApTrackingData,
			offset1:          OffsetValue{valueType: Reference, register: off1_reg, value: num_off1, dereference: true},
			dereference:      dereference,
		}
	}
	// Reference with deref 2 offsets: cast([reg + off1] + off2)
	_, err = fmt.Scanf(value_string, "cast([%c%c + %s] + %s)", off_1_reg_0, off_1_reg_1, off1, off2)
	if err == nil {
		off1_reg := getRegister(off_1_reg_0, off_1_reg_1)
		num_off1 := offsetValueFromString(off1)
		num_off2 := offsetValueFromString(off2)
		return HintReference{
			ap_tracking_data: reference.ApTrackingData,
			offset1:          OffsetValue{valueType: Reference, register: off1_reg, value: num_off1, dereference: true},
			offset2:          OffsetValue{value: num_off2},
			dereference:      dereference,
		}
	}
	// Reference with deref + reference with deref: cast([reg + off1] + [reg + off2])
	var off_2_reg_0 string
	var off_2_reg_1 string
	_, err = fmt.Scanf(value_string, "cast([%c%c + %s] + [%c%c + %s])", off_1_reg_0, off_1_reg_1, off1, off2)
	if err == nil {
		off1_reg := getRegister(off_1_reg_0, off_1_reg_1)
		off2_reg := getRegister(off_2_reg_0, off_2_reg_1)
		num_off1 := offsetValueFromString(off1)
		num_off2 := offsetValueFromString(off2)
		return HintReference{
			ap_tracking_data: reference.ApTrackingData,
			offset1:          OffsetValue{valueType: Reference, register: off1_reg, value: num_off1, dereference: true},
			offset2:          OffsetValue{valueType: Reference, register: off2_reg, value: num_off2, dereference: true},
			dereference:      dereference,
		}
	}
	// Special subcases: Sometimes if the offset is 0 it gets omitted, so we get [reg]

	// Reference with deref no off: cast([reg])
	_, err = fmt.Scanf(value_string, "cast([%c%c])", off_1_reg_0, off_1_reg_1)
	if err == nil {
		off1_reg := getRegister(off_1_reg_0, off_1_reg_1)
		return HintReference{
			ap_tracking_data: reference.ApTrackingData,
			offset1:          OffsetValue{valueType: Reference, register: off1_reg, dereference: true},
			dereference:      dereference,
		}
	}
	// Reference with deref 2 offsets no off1: cast([reg] + off2)
	_, err = fmt.Scanf(value_string, "cast([%c%c + %s] + %s)", off_1_reg_0, off_1_reg_1, off2)
	if err == nil {
		off1_reg := getRegister(off_1_reg_0, off_1_reg_1)
		num_off2 := offsetValueFromString(off2)
		return HintReference{
			ap_tracking_data: reference.ApTrackingData,
			offset1:          OffsetValue{valueType: Reference, register: off1_reg, dereference: true},
			offset2:          OffsetValue{value: num_off2},
			dereference:      dereference,
		}
	}

	// 2 dereferences no off1 : cast([reg] + [reg + off2])
	_, err = fmt.Scanf(value_string, "cast([%c%c] + [%c%c + %s])", off_1_reg_0, off_1_reg_1, off2)
	if err == nil {
		off1_reg := getRegister(off_1_reg_0, off_1_reg_1)
		off2_reg := getRegister(off_2_reg_0, off_2_reg_1)
		num_off2 := offsetValueFromString(off2)
		return HintReference{
			ap_tracking_data: reference.ApTrackingData,
			offset1:          OffsetValue{valueType: Reference, register: off1_reg, dereference: true},
			offset2:          OffsetValue{valueType: Reference, register: off2_reg, value: num_off2, dereference: true},
			dereference:      dereference,
		}
	}
	// 2 dereferences no off2: cast([reg + off1] + [reg])
	_, err = fmt.Scanf(value_string, "cast([%c%c + %s] + [%c%c])", off_1_reg_0, off_1_reg_1, off1)
	if err == nil {
		off1_reg := getRegister(off_1_reg_0, off_1_reg_1)
		off2_reg := getRegister(off_2_reg_0, off_2_reg_1)
		num_off1 := offsetValueFromString(off1)
		return HintReference{
			ap_tracking_data: reference.ApTrackingData,
			offset1:          OffsetValue{valueType: Reference, register: off1_reg, value: num_off1, dereference: true},
			offset2:          OffsetValue{valueType: Reference, register: off2_reg, dereference: true},
			dereference:      dereference,
		}
	}
	// 2 dereferences no offs: cast([reg] + [reg])
	_, err = fmt.Scanf(value_string, "cast([%c%c] + [%c%c])", off_1_reg_0, off_1_reg_1)
	if err == nil {
		off1_reg := getRegister(off_1_reg_0, off_1_reg_1)
		off2_reg := getRegister(off_2_reg_0, off_2_reg_1)
		return HintReference{
			ap_tracking_data: reference.ApTrackingData,
			offset1:          OffsetValue{valueType: Reference, register: off1_reg, dereference: true},
			offset2:          OffsetValue{valueType: Reference, register: off2_reg, dereference: true},
			dereference:      dereference,
		}
	}

	return HintReference{ap_tracking_data: reference.ApTrackingData}
}

// Parses strings of type num/(-num)
func offsetValueFromString(num string) int {
	value_string, _ := strings.CutPrefix(num, "(")
	value_string, _ = strings.CutSuffix(value_string, ")")
	res, _ := strconv.ParseInt(num, 10, 32)
	return int(res)
}

func getRegister(reg_0 string, reg_1 string) vm.Register {
	reg := vm.AP
	if reg_0 == "f" && reg_1 == "p" {
		reg = vm.FP
	}
	return reg
}
