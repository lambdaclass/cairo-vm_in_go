package vm_test

import (
	"reflect"
	"testing"

	"github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
	"github.com/lambdaclass/cairo-vm.go/pkg/vm"
)

func TestNewProgram(t *testing.T) {

}

func TestExtractConstantsEmpty(t *testing.T) {
	program := vm.Program{}
	expectedConstants := make(map[string]lambdaworks.Felt)
	if !reflect.DeepEqual(program.ExtractConstants(), expectedConstants) {
		t.Errorf("Wrong Constants, expected %v, got %v", expectedConstants, program.ExtractConstants())
	}

}

func TestExtractConstants(t *testing.T) {
	program := vm.Program{
		Identifiers: map[string]vm.Identifier{
			"start": {
				Value: lambdaworks.FeltFromUint64(4),
				Type:  "label",
			},
			"end": {
				Value: lambdaworks.FeltFromUint64(8),
				Type:  "label",
			},
			"A": {
				Value: lambdaworks.FeltFromUint64(7),
				Type:  "constant",
			},
			"B": {
				Value: lambdaworks.FeltFromUint64(17),
				Type:  "constant",
			},
		},
	}
	expectedConstants := map[string]lambdaworks.Felt{
		"A": lambdaworks.FeltFromUint64(7),
		"B": lambdaworks.FeltFromUint64(17),
	}
	if !reflect.DeepEqual(program.ExtractConstants(), expectedConstants) {
		t.Errorf("Wrong Constants, expected %v, got %v", expectedConstants, program.ExtractConstants())
	}

}
