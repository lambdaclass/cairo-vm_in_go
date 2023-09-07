package cairo_run_test

import (
	"fmt"
	"testing"

	"github.com/lambdaclass/cairo-vm.go/pkg/vm/cairo_run"
)

// Things we are skipping for now:
// - Initializing hint_executor and passing it to `cairo_run`
// - cairo_run_config stuff
// - Asserting expected trace values
// - Asserting memory_holes
func TestFibonacci(t *testing.T) {
	_, err := cairo_run.CairoRun("../../../cairo_programs/fibonacci.json", "plain", false)
	if err != nil {
		t.Errorf("Program execution failed with error: %s", err)
	}
	fmt.Println(err)
}

func TestPoseidonBuiltin(t *testing.T) {
	_, err := cairo_run.CairoRun("../../../cairo_programs/poseidon_builtin.json", "plain", false)
	if err != nil {
		t.Errorf("Program execution failed with error: %s", err)
	}
}

func TestPoseidonHash(t *testing.T) {
	_, err := cairo_run.CairoRun("../../../cairo_programs/poseidon_hash.json", "plain", false)
	if err != nil {
		t.Errorf("Program execution failed with error: %s", err)
	}
}

func TestSimplePrint(t *testing.T) {
	_, err := cairo_run.CairoRun("../../../cairo_programs/simple_print.json", "plain", false)
	if err != nil {
		t.Errorf("Program execution failed with error: %s", err)
	}
}
