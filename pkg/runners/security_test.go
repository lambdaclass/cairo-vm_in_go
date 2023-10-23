package runners_test

import (
	"testing"

	"github.com/lambdaclass/cairo-vm.go/pkg/builtins"
	. "github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
	. "github.com/lambdaclass/cairo-vm.go/pkg/runners"
	. "github.com/lambdaclass/cairo-vm.go/pkg/vm"
	. "github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
)

func TestVerifySecureRunnerEmptyMemory(t *testing.T) {
	runner, _ := NewCairoRunner(Program{}, "all_cairo", false)
	runner.Initialize()
	err := VerifySecureRunner(runner, true, nil)
	if err != nil {
		t.Errorf("VerifySecureRunner failed with error: %s", err.Error())
	}
}

func TestVerifySecureRunnerOutOfBoundsAccessProgram(t *testing.T) {
	runner, _ := NewCairoRunner(Program{}, "all_cairo", false)
	runner.Initialize()
	// Insert an element into the program segment to trigger out of bounds access to program segment error
	runner.Vm.Segments.Memory.Insert(runner.ProgramBase, NewMaybeRelocatableFelt(FeltOne()))
	err := VerifySecureRunner(runner, true, nil)
	if err == nil {
		t.Errorf("VerifySecureRunner should have failed")
	}
}

func TestVerifySecureRunnerOutOfBoundsAccessBuiltin(t *testing.T) {
	runner, _ := NewCairoRunner(Program{Builtins: []string{builtins.OUTPUT_BUILTIN_NAME}}, "all_cairo", false)
	runner.Initialize()
	stopPtr := uint(0)
	runner.Vm.BuiltinRunners[0].(*builtins.OutputBuiltinRunner).StopPtr = &stopPtr
	// Insert an element into the output segment to trigger out of bounds access to builtin segment error
	runner.Vm.Segments.Memory.Insert(runner.Vm.BuiltinRunners[0].Base(), NewMaybeRelocatableFelt(FeltOne()))
	err := VerifySecureRunner(runner, true, nil)
	if err == nil {
		t.Errorf("VerifySecureRunner should have failed")
	}
}

func TestVerifySecureRunnerTemporaryVal(t *testing.T) {
	runner, _ := NewCairoRunner(Program{}, "all_cairo", false)
	runner.Initialize()
	// Insert a temporary address into memory
	runner.Vm.Segments.Memory.Insert(NewRelocatable(1, 7), NewMaybeRelocatableRelocatable(NewRelocatable(-1, 0)))
	err := VerifySecureRunner(runner, true, nil)
	if err == nil {
		t.Errorf("VerifySecureRunner should have failed")
	}
}
