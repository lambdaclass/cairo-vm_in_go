package hint_utils_test

import (
	"testing"

	. "github.com/lambdaclass/cairo-vm.go/pkg/hints/hint_utils"
	"github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
	"github.com/lambdaclass/cairo-vm.go/pkg/parser"
)

// ParseHintReference tests

func ParseHintReferenceImmediate(t *testing.T) {
	reference := parser.Reference{Value: "cast(17, felt)"}
	expected := HintReference{Offset1: OffsetValue{ValueType: Immediate, Immediate: lambdaworks.FeltFromDecString("17")}}
	if ParseHintReference(reference) != expected {
		t.Errorf("Wrong parsed reference")
	}
}
