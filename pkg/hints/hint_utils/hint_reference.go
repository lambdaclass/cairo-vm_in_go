package hint_utils

import (
	"fmt"
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
	ValueType      string
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
	var valueString = reference.Value
	// Trim outer brackets if dereference
	// example [cast(reg + offset1, type)] -> cast(reg + offset1, type), dereference = true
	valueString, HasPrefix := strings.CutPrefix(valueString, "[")
	valueString, HasSuffix := strings.CutSuffix(valueString, "]")
	var dereference = HasPrefix && HasSuffix
	var valueType string
	var immediate uint64
	// Negative numbers com in the form (-N), we can scan -N, but need a special case for (-N)
	// In order to avoid multipliying our cases, we will remove all "(" & ")" characters from our string
	valueString = strings.ReplaceAll(valueString, "(", "")
	valueString = strings.ReplaceAll(valueString, ")", "")
	// Scan various types till a match is found
	// Immediate: cast(number, type)
	_, err := fmt.Sscanf(valueString, "cast%d, %s", &immediate, &valueType)
	if err == nil {
		var felt = FeltFromUint64(immediate)
		return HintReference{
			ApTrackingData: reference.ApTrackingData,
			Offset1:        OffsetValue{Immediate: felt, ValueType: Immediate},
			ValueType:      valueType,
		}
	}
	var off1Reg0 byte
	var off1Reg1 byte
	var off1 int
	// Reference no deref 1 offset: cast(reg + off, type)
	_, err = fmt.Sscanf(valueString, "cast%c%c + %d, %s", &off1Reg0, &off1Reg1, &off1, &valueType)
	if err == nil {
		off1Reg := getRegister(off1Reg0, off1Reg1)
		return HintReference{
			ApTrackingData: reference.ApTrackingData,
			Offset1:        OffsetValue{ValueType: Reference, Register: off1Reg, Value: off1},
			Dereference:    dereference,
			ValueType:      valueType,
		}
	}
	var off2 int
	// Reference no deref 2 offsets: cast(reg + off1 + off2, type)
	_, err = fmt.Sscanf(valueString, "cast%c%c + %d + %d, %s", &off1Reg0, &off1Reg1, &off1, &off2, &valueType)
	if err == nil {
		off1Reg := getRegister(off1Reg0, off1Reg1)
		return HintReference{
			ApTrackingData: reference.ApTrackingData,
			Offset1:        OffsetValue{ValueType: Reference, Register: off1Reg, Value: off1},
			Offset2:        OffsetValue{Value: off2},
			Dereference:    dereference,
			ValueType:      valueType,
		}
	}
	// Reference with deref 1 offset: cast([reg + off1], type)
	_, err = fmt.Sscanf(valueString, "cast[%c%c + %d], %s", &off1Reg0, &off1Reg1, &off1, &valueType)
	if err == nil {
		off1Reg := getRegister(off1Reg0, off1Reg1)
		return HintReference{
			ApTrackingData: reference.ApTrackingData,
			Offset1:        OffsetValue{ValueType: Reference, Register: off1Reg, Value: off1, Dereference: true},
			Dereference:    dereference,
			ValueType:      valueType,
		}
	}
	// Reference with deref 2 offsets: cast([reg + off1] + off2, type)
	_, err = fmt.Sscanf(valueString, "cast[%c%c + %d] + %d, %s", &off1Reg0, &off1Reg1, &off1, &off2, &valueType)
	if err == nil {
		off1Reg := getRegister(off1Reg0, off1Reg1)
		return HintReference{
			ApTrackingData: reference.ApTrackingData,
			Offset1:        OffsetValue{ValueType: Reference, Register: off1Reg, Value: off1, Dereference: true},
			Offset2:        OffsetValue{Value: off2},
			Dereference:    dereference,
			ValueType:      valueType,
		}
	}
	var off2Reg0 byte
	var off2Reg1 byte
	// Two references with deref: cast([reg + off1] + [reg + off2], type)
	_, err = fmt.Sscanf(valueString, "cast[%c%c + %d] + [%c%c + %d], %s", &off1Reg0, &off1Reg1, &off1, &off2Reg0, &off2Reg1, &off2, &valueType)
	if err == nil {
		off1Reg := getRegister(off1Reg0, off1Reg1)
		off2Reg := getRegister(off2Reg0, off2Reg1)
		return HintReference{
			ApTrackingData: reference.ApTrackingData,
			Offset1:        OffsetValue{ValueType: Reference, Register: off1Reg, Value: off1, Dereference: true},
			Offset2:        OffsetValue{ValueType: Reference, Register: off2Reg, Value: off2, Dereference: true},
			Dereference:    dereference,
			ValueType:      valueType,
		}
	}
	// Special subcases: Sometimes if the offset is 0 it gets omitted, so we get [reg]

	// Reference with deref no off: cast([reg], type)
	_, err = fmt.Sscanf(valueString, "cast[%c%c], %s", &off1Reg0, &off1Reg1, &valueType)
	if err == nil {
		off1Reg := getRegister(off1Reg0, off1Reg1)
		return HintReference{
			ApTrackingData: reference.ApTrackingData,
			Offset1:        OffsetValue{ValueType: Reference, Register: off1Reg, Dereference: true},
			Dereference:    dereference,
			ValueType:      valueType,
		}
	}
	// Reference with deref 2 offsets no off1: cast([reg] + off2, type)
	_, err = fmt.Sscanf(valueString, "cast[%c%c] + %d, %s", &off1Reg0, &off1Reg1, &off2, &valueType)
	if err == nil {
		off1Reg := getRegister(off1Reg0, off1Reg1)
		return HintReference{
			ApTrackingData: reference.ApTrackingData,
			Offset1:        OffsetValue{ValueType: Reference, Register: off1Reg, Dereference: true},
			Offset2:        OffsetValue{Value: off2},
			Dereference:    dereference,
			ValueType:      valueType,
		}
	}

	// 2 dereferences no off1 : cast([reg] + [reg + off2], type)
	_, err = fmt.Sscanf(valueString, "cast[%c%c] + [%c%c + %d], %s", &off1Reg0, &off1Reg1, &off2Reg0, &off2Reg1, &off2, &valueType)
	if err == nil {
		off1Reg := getRegister(off1Reg0, off1Reg1)
		off2Reg := getRegister(off2Reg0, off2Reg1)
		return HintReference{
			ApTrackingData: reference.ApTrackingData,
			Offset1:        OffsetValue{ValueType: Reference, Register: off1Reg, Dereference: true},
			Offset2:        OffsetValue{ValueType: Reference, Register: off2Reg, Value: off2, Dereference: true},
			Dereference:    dereference,
			ValueType:      valueType,
		}
	}
	// 2 dereferences no off2: cast([reg + off1] + [reg], type)
	_, err = fmt.Sscanf(valueString, "cast[%c%c + %d] + [%c%c], %s", &off1Reg0, &off1Reg1, &off1, &off2Reg0, &off2Reg1, &valueType)
	if err == nil {
		off1Reg := getRegister(off1Reg0, off1Reg1)
		off2Reg := getRegister(off2Reg0, off2Reg1)
		return HintReference{
			ApTrackingData: reference.ApTrackingData,
			Offset1:        OffsetValue{ValueType: Reference, Register: off1Reg, Value: off1, Dereference: true},
			Offset2:        OffsetValue{ValueType: Reference, Register: off2Reg, Dereference: true},
			Dereference:    dereference,
			ValueType:      valueType,
		}
	}
	// 2 dereferences no offs: cast([reg] + [reg], type)
	_, err = fmt.Sscanf(valueString, "cast[%c%c] + [%c%c], %s", &off1Reg0, &off1Reg1, &off2Reg0, &off2Reg1, &valueType)
	if err == nil {
		off1Reg := getRegister(off1Reg0, off1Reg1)
		off2Reg := getRegister(off2Reg0, off2Reg1)
		return HintReference{
			ApTrackingData: reference.ApTrackingData,
			Offset1:        OffsetValue{ValueType: Reference, Register: off1Reg, Dereference: true},
			Offset2:        OffsetValue{ValueType: Reference, Register: off2Reg, Dereference: true},
			Dereference:    dereference,
			ValueType:      valueType,
		}
	}

	return HintReference{ApTrackingData: reference.ApTrackingData}
}

func getRegister(reg_0 byte, reg_1 byte) vm.Register {
	reg := vm.AP
	if reg_0 == 'f' && reg_1 == 'p' {
		reg = vm.FP
	}
	return reg
}
