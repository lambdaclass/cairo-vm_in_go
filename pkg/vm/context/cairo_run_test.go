package context

import (
	"fmt"
	"testing"

	"github.com/lambdaclass/cairo-vm.go/pkg/parser"
)

// Things we are skipping for now:
// - Initializing hint_executor and passing it to `cairo_run`
// - cairo_run_config stuff
// - Asserting expected trace values
// - Asserting memory_holes
func TestFibonacci(t *testing.T) {
	compiledProgram := parser.Parse("../../../cairo_programs/fibonacci.json")
	err := CairoRunProgram(compiledProgram)
	if err != nil {
		t.Errorf("Program execution failed with error: %s", err)
	}
	fmt.Println(err)
}
