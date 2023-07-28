package parser_test

import (
	"reflect"
	"testing"

	"github.com/lambdaclass/cairo-vm.go/pkg/parser"
)

func TestData(t *testing.T) {
	got := parser.Parse("./test_compiled.json")
	expected := []string{"0x480680017fff8000", "0x3e8", "0x480680017fff8000", "0x7d0", "0x48307fff7ffe8000", "0x208b7fff7fff7ffe"}
	if !reflect.DeepEqual(got.Data, expected) {
		t.Errorf("We should have this data %s, got %s", expected, got.Data)
	}
}
