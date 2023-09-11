package hints_test

import (
	"reflect"
	"testing"

	. "github.com/lambdaclass/cairo-vm.go/pkg/hints"
	. "github.com/lambdaclass/cairo-vm.go/pkg/hints/hint_utils"
	"github.com/lambdaclass/cairo-vm.go/pkg/parser"
	"github.com/lambdaclass/cairo-vm.go/pkg/vm"
)

func TestCompileHintEmpty(t *testing.T) {
	hintProcessor := &CairoVmHintProcessor{}
	hintParams := &parser.HintParams{}
	referenceManager := &parser.ReferenceManager{}
	expectedData := HintData{}
	data, err := hintProcessor.CompileHint(hintParams, referenceManager)
	if err != nil {
		t.Errorf("Error in test: %s", err)
	}
	if reflect.DeepEqual(data.(HintData), expectedData) {
		t.Errorf("Wrong hint data, %+v", data)
	}
}

func TestCompileHintHappyPath(t *testing.T) {
	hintProcessor := &CairoVmHintProcessor{}
	hintParams := &parser.HintParams{
		Code: "ids.a = ids.b",
		FlowTrackingData: parser.FlowTrackingData{
			APTracking: parser.ApTrackingData{Group: 1, Offset: 2},
		},
		ReferenceIds: map[string]uint{"__main.__.a": 0, "__main__.b": 1},
	}
	referenceManager := &parser.ReferenceManager{
		References: []parser.Reference{
			{
				Value: "cast(ap + (-2), felt)",
			},
			{
				Value: "cast(ap + (-1), felt)",
			},
		},
	}
	expectedData := HintData{
		Ids: IdsManager{References: map[string]HintReference{
			"a": {
				Offset1: OffsetValue{
					ValueType: Reference,
					Value:     -2,
				},
				ValueType: "felt",
			},
			"b": {
				Offset1: OffsetValue{
					ValueType: Reference,
					Value:     -1,
				},
				ValueType: "felt",
			},
		},
			HintApTracking: parser.ApTrackingData{Group: 1, Offset: 2},
		},
		Code: "ids.a = ids.b",
	}
	data, err := hintProcessor.CompileHint(hintParams, referenceManager)
	if err != nil {
		t.Errorf("Error in test: %s", err)
	}
	if !reflect.DeepEqual(data.(HintData), expectedData) {
		t.Errorf("Wrong hint data, %+v", data)
	}
}

func TestCompileHintMissingReference(t *testing.T) {
	hintProcessor := &CairoVmHintProcessor{}
	hintParams := &parser.HintParams{
		ReferenceIds: map[string]uint{"a": 0, "b": 1},
	}
	referenceManager := &parser.ReferenceManager{}
	_, err := hintProcessor.CompileHint(hintParams, referenceManager)
	if err == nil {
		t.Errorf("Should have failed")
	}
}

func TestExecuteHintWrongHintData(t *testing.T) {
	hintProcessor := &CairoVmHintProcessor{}
	hintData := any("Mistake")
	vm := vm.NewVirtualMachine()
	err := hintProcessor.ExecuteHint(vm, &hintData, nil)
	if err == nil {
		t.Errorf("Should have failed")
	}
}

func TestExecuteHintUnknownHint(t *testing.T) {
	hintProcessor := &CairoVmHintProcessor{}
	hintData := any(HintData{Code: "print(Hello World)"})
	vm := vm.NewVirtualMachine()
	err := hintProcessor.ExecuteHint(vm, &hintData, nil)
	if err == nil {
		t.Errorf("Should have failed")
	}
}
