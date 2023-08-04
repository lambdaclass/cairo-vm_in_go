package parser_test

import (
	"reflect"
	"testing"

	"github.com/lambdaclass/cairo-vm.go/pkg/parser"
)

func TestData(t *testing.T) {
	got := parser.Parse("../../cairo_programs/fibonacci.json")
	expected := []string{"0x480680017fff8000",
		"0x1",
		"0x480680017fff8000",
		"0x1",
		"0x480680017fff8000",
		"0xa",
		"0x1104800180018000",
		"0x5",
		"0x400680017fff7fff",
		"0x90",
		"0x208b7fff7fff7ffe",
		"0x20780017fff7ffd",
		"0x5",
		"0x480a7ffc7fff8000",
		"0x480a7ffc7fff8000",
		"0x208b7fff7fff7ffe",
		"0x482a7ffc7ffb8000",
		"0x480a7ffc7fff8000",
		"0x48127ffe7fff8000",
		"0x482680017ffd8000",
		"0x800000000000011000000000000000000000000000000000000000000000000",
		"0x1104800180018000",
		"0x800000000000010fffffffffffffffffffffffffffffffffffffffffffffff7",
		"0x208b7fff7fff7ffe"}
	if !reflect.DeepEqual(got.Data, expected) {
		t.Errorf("We should have this data %s, got %s", expected, got.Data)
	}
}
