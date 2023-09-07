package hint_utils_test

import (
	"testing"

	. "github.com/lambdaclass/cairo-vm.go/pkg/hints/hint_utils"
	"github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
	"github.com/lambdaclass/cairo-vm.go/pkg/parser"
	"github.com/lambdaclass/cairo-vm.go/pkg/vm"
)

// ParseHintReference tests

func TestParseHintReferenceImmediate(t *testing.T) {
	reference := parser.Reference{Value: "cast(17, felt)"}
	expected := HintReference{
		Offset1:   OffsetValue{ValueType: Immediate, Immediate: lambdaworks.FeltFromDecString("17")},
		ValueType: "felt",
	}

	if ParseHintReference(reference) != expected {
		t.Errorf("Wrong parsed reference, %+v", ParseHintReference(reference))
	}
}

func TestParseHintReferenceImmediateCustomType(t *testing.T) {
	reference := parser.Reference{Value: "cast(17, cat)"}
	expected := HintReference{
		Offset1:   OffsetValue{ValueType: Immediate, Immediate: lambdaworks.FeltFromDecString("17")},
		ValueType: "cat",
	}

	if ParseHintReference(reference) != expected {
		t.Errorf("Wrong parsed reference, %+v", ParseHintReference(reference))
	}
}

func TestParseHintReferenceDereferenceSimpleApBasedPositive(t *testing.T) {
	reference := parser.Reference{Value: "[cast(ap + 1, felt)]"}
	expected := HintReference{
		Offset1:     OffsetValue{ValueType: Reference, Value: 1},
		ValueType:   "felt",
		Dereference: true,
	}
	if ParseHintReference(reference) != expected {
		t.Errorf("Wrong parsed reference, %+v", ParseHintReference(reference))
	}
}
func TestParseHintReferenceSimpleApBasedPositive(t *testing.T) {
	reference := parser.Reference{Value: "cast(ap + 1, felt)"}
	expected := HintReference{
		Offset1:   OffsetValue{ValueType: Reference, Value: 1},
		ValueType: "felt",
	}
	if ParseHintReference(reference) != expected {
		t.Errorf("Wrong parsed reference, %+v", ParseHintReference(reference))
	}
}

func TestParseHintReferenceSimpleApBasedTypePointer(t *testing.T) {
	reference := parser.Reference{Value: "cast(ap + 1, felt*)"}
	expected := HintReference{
		Offset1:   OffsetValue{ValueType: Reference, Value: 1},
		ValueType: "felt*",
	}
	if ParseHintReference(reference) != expected {
		t.Errorf("Wrong parsed reference, %+v", ParseHintReference(reference))
	}
}

func TestParseHintReferenceSimpleFpBasedPositive(t *testing.T) {
	reference := parser.Reference{Value: "cast(fp + 1, felt)"}
	expected := HintReference{
		Offset1:   OffsetValue{ValueType: Reference, Value: 1, Register: vm.FP},
		ValueType: "felt",
	}
	if ParseHintReference(reference) != expected {
		t.Errorf("Wrong parsed reference, %+v", ParseHintReference(reference))
	}
}

func TestParseHintReferenceSimpleApBasedNegative(t *testing.T) {
	reference := parser.Reference{Value: "cast(ap + (-1), felt)"}
	expected := HintReference{
		Offset1:   OffsetValue{ValueType: Reference, Value: -1},
		ValueType: "felt",
	}
	if ParseHintReference(reference) != expected {
		t.Errorf("Wrong parsed reference, %+v", ParseHintReference(reference))
	}
}

func TestParseHintReferenceTwoOffsetsPositives(t *testing.T) {
	reference := parser.Reference{Value: "cast(ap + 1 + 2, felt)"}
	expected := HintReference{
		Offset1:   OffsetValue{ValueType: Reference, Value: 1},
		Offset2:   OffsetValue{ValueType: Value, Value: 2},
		ValueType: "felt",
	}
	if ParseHintReference(reference) != expected {
		t.Errorf("Wrong parsed reference, %+v", ParseHintReference(reference))
	}
}

func TestParseHintReferenceTwoOffsetsNegatives(t *testing.T) {
	reference := parser.Reference{Value: "cast(ap + (-1) + (-2), felt)"}
	expected := HintReference{
		Offset1:   OffsetValue{ValueType: Reference, Value: -1},
		Offset2:   OffsetValue{ValueType: Value, Value: -2},
		ValueType: "felt",
	}
	if ParseHintReference(reference) != expected {
		t.Errorf("Wrong parsed reference, %+v", ParseHintReference(reference))
	}
}

func TestParseHintReferenceTwoOffsetsPosNeg(t *testing.T) {
	reference := parser.Reference{Value: "cast(ap + 1 + (-2), felt)"}
	expected := HintReference{
		Offset1:   OffsetValue{ValueType: Reference, Value: 1},
		Offset2:   OffsetValue{ValueType: Value, Value: -2},
		ValueType: "felt",
	}
	if ParseHintReference(reference) != expected {
		t.Errorf("Wrong parsed reference, %+v", ParseHintReference(reference))
	}
}

func TestParseHintReferenceTwoOffsetsNegPos(t *testing.T) {
	reference := parser.Reference{Value: "cast(ap + (-1) + 2, felt)"}
	expected := HintReference{
		Offset1:   OffsetValue{ValueType: Reference, Value: -1},
		Offset2:   OffsetValue{ValueType: Value, Value: 2},
		ValueType: "felt",
	}
	if ParseHintReference(reference) != expected {
		t.Errorf("Wrong parsed reference, %+v", ParseHintReference(reference))
	}
}

func TestParseHintReferenceDerefOneOffset(t *testing.T) {
	reference := parser.Reference{Value: "cast([ap + 1], felt)"}
	expected := HintReference{
		Offset1:   OffsetValue{ValueType: Reference, Value: 1, Dereference: true},
		ValueType: "felt",
	}
	if ParseHintReference(reference) != expected {
		t.Errorf("Wrong parsed reference, %+v", ParseHintReference(reference))
	}
}

func TestParseHintReferenceDerefOneOffsetNegative(t *testing.T) {
	reference := parser.Reference{Value: "cast([ap + (-1)], felt)"}
	expected := HintReference{
		Offset1:   OffsetValue{ValueType: Reference, Value: -1, Dereference: true},
		ValueType: "felt",
	}
	if ParseHintReference(reference) != expected {
		t.Errorf("Wrong parsed reference, %+v", ParseHintReference(reference))
	}
}

func TestParseHintReferenceDerefTwoOffsetsNegatives(t *testing.T) {
	reference := parser.Reference{Value: "cast([ap + (-1)] + (-2), felt)"}
	expected := HintReference{
		Offset1:   OffsetValue{ValueType: Reference, Value: -1, Dereference: true},
		Offset2:   OffsetValue{ValueType: Value, Value: -2},
		ValueType: "felt",
	}
	if ParseHintReference(reference) != expected {
		t.Errorf("Wrong parsed reference, %+v", ParseHintReference(reference))
	}
}

func TestParseHintReferenceDerefTwoOffsetsPositives(t *testing.T) {
	reference := parser.Reference{Value: "cast([ap + 1] + 2, felt)"}
	expected := HintReference{
		Offset1:   OffsetValue{ValueType: Reference, Value: 1, Dereference: true},
		Offset2:   OffsetValue{ValueType: Value, Value: 2},
		ValueType: "felt",
	}
	if ParseHintReference(reference) != expected {
		t.Errorf("Wrong parsed reference, %+v", ParseHintReference(reference))
	}
}

func TestParseHintReferenceDerefTwoOffsetsNegPos(t *testing.T) {
	reference := parser.Reference{Value: "cast([ap + (-1)] + 2, felt)"}
	expected := HintReference{
		Offset1:   OffsetValue{ValueType: Reference, Value: -1, Dereference: true},
		Offset2:   OffsetValue{ValueType: Value, Value: 2},
		ValueType: "felt",
	}
	if ParseHintReference(reference) != expected {
		t.Errorf("Wrong parsed reference, %+v", ParseHintReference(reference))
	}
}

func TestParseHintReferenceDerefTwoOffsetsPosNeg(t *testing.T) {
	reference := parser.Reference{Value: "cast([ap + 1] + (-2), felt)"}
	expected := HintReference{
		Offset1:   OffsetValue{ValueType: Reference, Value: 1, Dereference: true},
		Offset2:   OffsetValue{ValueType: Value, Value: -2},
		ValueType: "felt",
	}
	if ParseHintReference(reference) != expected {
		t.Errorf("Wrong parsed reference, %+v", ParseHintReference(reference))
	}
}

func TestParseHintTwoDereferencesPositives(t *testing.T) {
	reference := parser.Reference{Value: "cast([ap + 1] + [ap + 2], felt)"}
	expected := HintReference{
		Offset1:   OffsetValue{ValueType: Reference, Value: 1, Dereference: true},
		Offset2:   OffsetValue{ValueType: Reference, Value: 2, Dereference: true},
		ValueType: "felt",
	}
	if ParseHintReference(reference) != expected {
		t.Errorf("Wrong parsed reference, %+v", ParseHintReference(reference))
	}
}

func TestParseHintTwoDereferencesPositivesDiffRegister(t *testing.T) {
	reference := parser.Reference{Value: "cast([ap + 1] + [fp + 2], felt)"}
	expected := HintReference{
		Offset1:   OffsetValue{ValueType: Reference, Value: 1, Dereference: true},
		Offset2:   OffsetValue{ValueType: Reference, Value: 2, Dereference: true, Register: vm.FP},
		ValueType: "felt",
	}
	if ParseHintReference(reference) != expected {
		t.Errorf("Wrong parsed reference, %+v", ParseHintReference(reference))
	}
}

func TestParseHintTwoDereferencesPosNeg(t *testing.T) {
	reference := parser.Reference{Value: "cast([ap + 1] + [ap + (-2)], felt)"}
	expected := HintReference{
		Offset1:   OffsetValue{ValueType: Reference, Value: 1, Dereference: true},
		Offset2:   OffsetValue{ValueType: Reference, Value: -2, Dereference: true},
		ValueType: "felt",
	}
	if ParseHintReference(reference) != expected {
		t.Errorf("Wrong parsed reference, %+v", ParseHintReference(reference))
	}
}

func TestParseHintReferenceDerefOneOffsetOmitted(t *testing.T) {
	reference := parser.Reference{Value: "cast([ap], felt)"}
	expected := HintReference{
		Offset1:   OffsetValue{ValueType: Reference, Value: 0, Dereference: true},
		ValueType: "felt",
	}
	if ParseHintReference(reference) != expected {
		t.Errorf("Wrong parsed reference, %+v", ParseHintReference(reference))
	}
}

func TestParseHintReferenceDerefTwoOffsetOmitted(t *testing.T) {
	reference := parser.Reference{Value: "cast([ap] + 1, felt)"}
	expected := HintReference{
		Offset1:   OffsetValue{ValueType: Reference, Value: 0, Dereference: true},
		Offset2:   OffsetValue{Value: 1},
		ValueType: "felt",
	}
	if ParseHintReference(reference) != expected {
		t.Errorf("Wrong parsed reference, %+v", ParseHintReference(reference))
	}
}

func TestParseHintReferenceDoubleDerefFirstOffOmitted(t *testing.T) {
	reference := parser.Reference{Value: "cast([ap] + [ap + 1], felt)"}
	expected := HintReference{
		Offset1:   OffsetValue{ValueType: Reference, Value: 0, Dereference: true},
		Offset2:   OffsetValue{ValueType: Reference, Value: 1, Dereference: true},
		ValueType: "felt",
	}
	if ParseHintReference(reference) != expected {
		t.Errorf("Wrong parsed reference, %+v", ParseHintReference(reference))
	}
}

func TestParseHintReferenceDoubleDerefSecondOffOmitted(t *testing.T) {
	reference := parser.Reference{Value: "cast([ap + 1] + [ap], felt)"}
	expected := HintReference{
		Offset1:   OffsetValue{ValueType: Reference, Value: 1, Dereference: true},
		Offset2:   OffsetValue{ValueType: Reference, Value: 0, Dereference: true},
		ValueType: "felt",
	}
	if ParseHintReference(reference) != expected {
		t.Errorf("Wrong parsed reference, %+v", ParseHintReference(reference))
	}
}

func TestParseHintReferenceDoubleDerefBothOffOmitted(t *testing.T) {
	reference := parser.Reference{Value: "cast([ap] + [ap], felt)"}
	expected := HintReference{
		Offset1:   OffsetValue{ValueType: Reference, Value: 0, Dereference: true},
		Offset2:   OffsetValue{ValueType: Reference, Value: 0, Dereference: true},
		ValueType: "felt",
	}
	if ParseHintReference(reference) != expected {
		t.Errorf("Wrong parsed reference, %+v", ParseHintReference(reference))
	}
}

func TestParseHintDereferenceReferenceDoubleDerefBothOffOmitted(t *testing.T) {
	reference := parser.Reference{Value: "[cast([ap] + [ap], felt)]"}
	expected := HintReference{
		Offset1:     OffsetValue{ValueType: Reference, Value: 0, Dereference: true},
		Offset2:     OffsetValue{ValueType: Reference, Value: 0, Dereference: true},
		ValueType:   "felt",
		Dereference: true,
	}
	if ParseHintReference(reference) != expected {
		t.Errorf("Wrong parsed reference, %+v", ParseHintReference(reference))
	}
}
