package builtins_test

import (
	"testing"

	"github.com/lambdaclass/cairo-vm.go/pkg/builtins"
	"github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
)

func TestBaseSignature(t *testing.T) {
	check_range := builtins.NewSignatureBuiltinRunner()
	if check_range.Base() != memory.NewRelocatable(0, 0) {
		t.Errorf("Wrong base value in %s builtin", check_range.Name())
	}
}

func TestInitializeSegmentsForSignatureBuiltin(t *testing.T) {
	range_check_builtin := builtins.NewSignatureBuiltinRunner()
	segment_manager := memory.NewMemorySegmentManager()
	range_check_builtin.InitializeSegments(&segment_manager)
	if range_check_builtin.Base() != memory.NewRelocatable(0, 0) {
		t.Errorf("Builtin %s base is not 0", range_check_builtin.Name())
	}
}
