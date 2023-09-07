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
	Offset1        OffsetValue
	Offset2        OffsetValue
	Dereference    bool
	ApTrackingData parser.ApTrackingData
}

type OffsetValue struct {
	ValueType   offsetValueType
	Immediate   Felt
	Value       int
	Register    vm.Register
	Dereference bool
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
	var immediate uint64
	// Scan various types till a match is found
	// Immediate: cast(number, type) positive
	_, err := fmt.Sscanf(value_string, "cast(%d, %s", &immediate, &value_type)
	if err == nil {
		var felt = FeltFromUint64(immediate)
		return HintReference{
			ApTrackingData: reference.ApTrackingData,
			Offset1:        OffsetValue{Immediate: felt, ValueType: Immediate},
		}
	}
	var off_1_reg_0 byte
	var off_1_reg_1 byte
	var off1 int
	// Reference no deref 1 offset: cast(reg + off, type)
	_, err = fmt.Sscanf(value_string, "cast(%c%c + %d, %s", &off_1_reg_0, &off_1_reg_1, &off1, &value_type)
	if err == nil {
		off1_reg := getRegister(off_1_reg_0, off_1_reg_1)
		return HintReference{
			ApTrackingData: reference.ApTrackingData,
			Offset1:        OffsetValue{ValueType: Reference, Register: off1_reg, Value: off1},
			Dereference:    dereference,
		}
	}
	// Reference no deref 1 offset negative: cast(reg + (-off), type)
	_, err = fmt.Sscanf(value_string, "cast(%c%c + (%d), %s", &off_1_reg_0, &off_1_reg_1, &off1, &value_type)
	if err == nil {
		off1_reg := getRegister(off_1_reg_0, off_1_reg_1)
		return HintReference{
			ApTrackingData: reference.ApTrackingData,
			Offset1:        OffsetValue{ValueType: Reference, Register: off1_reg, Value: off1},
			Dereference:    dereference,
		}
	}
	var off2 int
	// Reference no deref 2 offsets ( - +): cast(reg + (-off1) + off2, type)
	_, err = fmt.Sscanf(value_string, "cast(%c%c + (%d) + %d, %s", &off_1_reg_0, &off_1_reg_1, &off1, &off2, &value_type)
	if err == nil {
		off1_reg := getRegister(off_1_reg_0, off_1_reg_1)
		return HintReference{
			ApTrackingData: reference.ApTrackingData,
			Offset1:        OffsetValue{ValueType: Reference, Register: off1_reg, Value: off1},
			Offset2:        OffsetValue{Value: off2},
			Dereference:    dereference,
		}
	}
	// Reference no deref 2 offsets (+ +): cast(reg + off1 + off2, type)
	_, err = fmt.Sscanf(value_string, "cast(%c%c + %d + %d, %s", &off_1_reg_0, &off_1_reg_1, &off1, &off2, &value_type)
	if err == nil {
		off1_reg := getRegister(off_1_reg_0, off_1_reg_1)
		return HintReference{
			ApTrackingData: reference.ApTrackingData,
			Offset1:        OffsetValue{ValueType: Reference, Register: off1_reg, Value: off1},
			Offset2:        OffsetValue{Value: off2},
			Dereference:    dereference,
		}
	}
	// Reference no deref 2 offsets ( - - ): cast(reg + (-off1) + (-off2), type)
	_, err = fmt.Sscanf(value_string, "cast(%c%c + (%d) + (%d), %s", &off_1_reg_0, &off_1_reg_1, &off1, &off2, &value_type)
	if err == nil {
		off1_reg := getRegister(off_1_reg_0, off_1_reg_1)
		return HintReference{
			ApTrackingData: reference.ApTrackingData,
			Offset1:        OffsetValue{ValueType: Reference, Register: off1_reg, Value: off1},
			Offset2:        OffsetValue{Value: off2},
			Dereference:    dereference,
		}
	}
	// Reference no deref 2 offsets ( + -): cast(reg + off1 + (-off2), type)
	_, err = fmt.Sscanf(value_string, "cast(%c%c + %d + (%d), %s", &off_1_reg_0, &off_1_reg_1, &off1, &off2, &value_type)
	if err == nil {
		off1_reg := getRegister(off_1_reg_0, off_1_reg_1)
		return HintReference{
			ApTrackingData: reference.ApTrackingData,
			Offset1:        OffsetValue{ValueType: Reference, Register: off1_reg, Value: off1},
			Offset2:        OffsetValue{Value: off2},
			Dereference:    dereference,
		}
	}
	// // Reference with deref 1 offset: cast([reg + off1], type)
	// _, err = fmt.Sscanf(value_string, "cast([%c%c + %s], %s", &off_1_reg_0, &off_1_reg_1, &off1, &value_type)
	// if err == nil {
	// 	off1_reg := getRegister(off_1_reg_0, off_1_reg_1)
	// 	num_off1 := offsetValueFromString(off1)
	// 	return HintReference{
	// 		ApTrackingData: reference.ApTrackingData,
	// 		Offset1:        OffsetValue{ValueType: Reference, Register: off1_reg, Value: num_off1, Dereference: true},
	// 		Dereference:    dereference,
	// 	}
	// }
	// // Reference with deref 2 offsets: cast([reg + off1] + off2, type)
	// _, err = fmt.Sscanf(value_string, "cast([%c%c + %s] + %s, %s", &off_1_reg_0, &off_1_reg_1, &off1, off2, &value_type)
	// if err == nil {
	// 	off1_reg := getRegister(off_1_reg_0, off_1_reg_1)
	// 	num_off1 := offsetValueFromString(off1)
	// 	num_off2 := offsetValueFromString(off2)
	// 	return HintReference{
	// 		ApTrackingData: reference.ApTrackingData,
	// 		Offset1:        OffsetValue{ValueType: Reference, Register: off1_reg, Value: num_off1, Dereference: true},
	// 		Offset2:        OffsetValue{Value: num_off2},
	// 		Dereference:    dereference,
	// 	}
	// }
	// // Reference with deref + reference with deref: cast([reg + off1] + [reg + off2], type)
	// var off_2_reg_0 string
	// var off_2_reg_1 string
	// _, err = fmt.Sscanf(value_string, "cast([%c%c + %s] + [%c%c + %s], %s", &off_1_reg_0, &off_1_reg_1, &off1, &off2, &value_type)
	// if err == nil {
	// 	off1_reg := getRegister(off_1_reg_0, off_1_reg_1)
	// 	off2_reg := getRegister(off_2_reg_0, off_2_reg_1)
	// 	num_off1 := offsetValueFromString(off1)
	// 	num_off2 := offsetValueFromString(off2)
	// 	return HintReference{
	// 		ApTrackingData: reference.ApTrackingData,
	// 		Offset1:        OffsetValue{ValueType: Reference, Register: off1_reg, Value: num_off1, Dereference: true},
	// 		Offset2:        OffsetValue{ValueType: Reference, Register: off2_reg, Value: num_off2, Dereference: true},
	// 		Dereference:    dereference,
	// 	}
	// }
	// // Special subcases: Sometimes if the offset is 0 it gets omitted, so we get [reg]

	// // Reference with deref no off: cast([reg], type)
	// _, err = fmt.Sscanf(value_string, "cast([%c%c], %s", &off_1_reg_0, &off_1_reg_1, &value_type)
	// if err == nil {
	// 	off1_reg := getRegister(off_1_reg_0, off_1_reg_1)
	// 	return HintReference{
	// 		ApTrackingData: reference.ApTrackingData,
	// 		Offset1:        OffsetValue{ValueType: Reference, Register: off1_reg, Dereference: true},
	// 		Dereference:    dereference,
	// 	}
	// }
	// // Reference with deref 2 offsets no off1: cast([reg] + off2, type)
	// _, err = fmt.Sscanf(value_string, "cast([%c%c + %s] + %s, %s", &off_1_reg_0, &off_1_reg_1, &off2, &value_type)
	// if err == nil {
	// 	off1_reg := getRegister(off_1_reg_0, off_1_reg_1)
	// 	num_off2 := offsetValueFromString(off2)
	// 	return HintReference{
	// 		ApTrackingData: reference.ApTrackingData,
	// 		Offset1:        OffsetValue{ValueType: Reference, Register: off1_reg, Dereference: true},
	// 		Offset2:        OffsetValue{Value: num_off2},
	// 		Dereference:    dereference,
	// 	}
	// }

	// // 2 dereferences no off1 : cast([reg] + [reg + off2], type)
	// _, err = fmt.Sscanf(value_string, "cast([%c%c] + [%c%c + %s], %s", &off_1_reg_0, &off_1_reg_1, &off2, &value_type)
	// if err == nil {
	// 	off1_reg := getRegister(off_1_reg_0, off_1_reg_1)
	// 	off2_reg := getRegister(off_2_reg_0, off_2_reg_1)
	// 	num_off2 := offsetValueFromString(off2)
	// 	return HintReference{
	// 		ApTrackingData: reference.ApTrackingData,
	// 		Offset1:        OffsetValue{ValueType: Reference, Register: off1_reg, Dereference: true},
	// 		Offset2:        OffsetValue{ValueType: Reference, Register: off2_reg, Value: num_off2, Dereference: true},
	// 		Dereference:    dereference,
	// 	}
	// }
	// // 2 dereferences no off2: cast([reg + off1] + [reg], type)
	// _, err = fmt.Sscanf(value_string, "cast([%c%c + %s] + [%c%c], %s", &off_1_reg_0, &off_1_reg_1, &off1, &value_type)
	// if err == nil {
	// 	off1_reg := getRegister(off_1_reg_0, off_1_reg_1)
	// 	off2_reg := getRegister(off_2_reg_0, off_2_reg_1)
	// 	num_off1 := offsetValueFromString(off1)
	// 	return HintReference{
	// 		ApTrackingData: reference.ApTrackingData,
	// 		Offset1:        OffsetValue{ValueType: Reference, Register: off1_reg, Value: num_off1, Dereference: true},
	// 		Offset2:        OffsetValue{ValueType: Reference, Register: off2_reg, Dereference: true},
	// 		Dereference:    dereference,
	// 	}
	// }
	// // 2 dereferences no offs: cast([reg] + [reg], type)
	// _, err = fmt.Sscanf(value_string, "cast([%c%c] + [%c%c], %s", &off_1_reg_0, &off_1_reg_1, &value_type)
	// if err == nil {
	// 	off1_reg := getRegister(off_1_reg_0, off_1_reg_1)
	// 	off2_reg := getRegister(off_2_reg_0, off_2_reg_1)
	// 	return HintReference{
	// 		ApTrackingData: reference.ApTrackingData,
	// 		Offset1:        OffsetValue{ValueType: Reference, Register: off1_reg, Dereference: true},
	// 		Offset2:        OffsetValue{ValueType: Reference, Register: off2_reg, Dereference: true},
	// 		Dereference:    dereference,
	// 	}
	// }

	return HintReference{ApTrackingData: reference.ApTrackingData}
}

// Parses strings of type num/(-num)
func offsetValueFromString(num string) int {
	value_string, _ := strings.CutPrefix(num, "(")
	value_string, _ = strings.CutSuffix(value_string, ")")
	res, _ := strconv.ParseInt(num, 10, 32)
	return int(res)
}

func getRegister(reg_0 byte, reg_1 byte) vm.Register {
	reg := vm.AP
	if reg_0 == 'f' && reg_1 == 'p' {
		reg = vm.FP
	}
	return reg
}
