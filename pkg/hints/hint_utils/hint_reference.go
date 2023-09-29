package hint_utils

import (
	"fmt"
	"strings"

	. "github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
	"github.com/lambdaclass/cairo-vm.go/pkg/parser"
	"github.com/lambdaclass/cairo-vm.go/pkg/vm"
)

// Contains data used to calculate the address or value of a cairo identifier
// Default value (aka HintReference{}) will consist of a reference to ap without dereference (ap)
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

// Parses a Reference to a HintReference, decoding its Value field
// Returns an emty reference if invalid
// Note: Current implementation is working, but can be quite slow due to Sscanf, should be replaced for a more performant version
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
	// Special subcases: Sometimes if the offset is 0 it gets omitted, so we get [reg] instead of [reg + 0]

	// Reference off omitted: cast(reg, type)
	_, err = fmt.Sscanf(valueString, "cast%c%c, %s", &off1Reg0, &off1Reg1, &valueType)
	if err == nil {
		off1Reg := getRegister(off1Reg0, off1Reg1)
		return HintReference{
			ApTrackingData: reference.ApTrackingData,
			Offset1:        OffsetValue{ValueType: Reference, Register: off1Reg},
			Dereference:    dereference,
			ValueType:      valueType,
		}
	}

	// Reference with deref off omitted: cast([reg], type)
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
	// Reference with deref 2 offsets off1 omitted: cast([reg] + off2, type)
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

	// 2 dereferences off1 omitted: cast([reg] + [reg + off2], type)
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
	// 2 dereferences off2 omitted: cast([reg + off1] + [reg], type)
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
	// 2 dereferences both offs omitted: cast([reg] + [reg], type)
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
	// Reference no dereference 2 offsets - + : cast(reg - off1 + off2, type)
	_, err = fmt.Sscanf(valueString, "cast%c%c - %d + %d, %s", &off1Reg0, &off1Reg1, &off1, &off2, &valueType)
	if err == nil {
		off1Reg := getRegister(off1Reg0, off1Reg1)
		return HintReference{
			ApTrackingData: reference.ApTrackingData,
			Offset1:        OffsetValue{ValueType: Reference, Register: off1Reg, Value: -off1},
			Offset2:        OffsetValue{Value: off2},
			Dereference:    dereference,
			ValueType:      valueType,
		}
	}
	// No matches (aka wrong format)
	return HintReference{ApTrackingData: reference.ApTrackingData}
}

// Returns FP if reg0 is f and reg 1 is p, else returns AP
func getRegister(reg0 byte, reg1 byte) vm.Register {
	reg := vm.AP
	if reg0 == 'f' && reg1 == 'p' {
		reg = vm.FP
	}
	return reg
}
