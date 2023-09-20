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
				Type:  "const",
			},
			"B": {
				Value: lambdaworks.FeltFromUint64(17),
				Type:  "const",
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

func TestExtractConstantsWithAliasedConstants(t *testing.T) {
	program := vm.Program{
		Identifiers: map[string]vm.Identifier{
			"path.A": {
				Value: lambdaworks.FeltFromUint64(7),
				Type:  "const",
			},
			"path.b": {
				Value: lambdaworks.FeltFromUint64(17),
				Type:  "label",
			},
			"other_path.A": {
				Destination: "path.A",
				Type:        "alias",
			},
			"other_path.b": {
				Destination: "path.b",
				Type:        "alias",
			},
		},
	}
	expectedConstants := map[string]lambdaworks.Felt{
		"path.A":       lambdaworks.FeltFromUint64(7),
		"other_path.A": lambdaworks.FeltFromUint64(7),
	}
	if !reflect.DeepEqual(program.ExtractConstants(), expectedConstants) {
		t.Errorf("Wrong Constants, expected %v, got %v", expectedConstants, program.ExtractConstants())
	}
}
