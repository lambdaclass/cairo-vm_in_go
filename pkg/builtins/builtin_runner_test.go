package builtins_test

import (
	"testing"

	"github.com/lambdaclass/cairo-vm.go/pkg/builtins"
	"github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
	"github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
)

func TestRunSecurityChecksEmptyMemory(t *testing.T) {
	builtin := builtins.NewBitwiseBuiltinRunner(256)
	segments := memory.NewMemorySegmentManager()
	err := builtins.RunSecurityChecksForBuiltin(builtin, &segments)
	if err != nil {
		t.Errorf("RunSecurityChecks failed with error: %s", err.Error())
	}
}

func TestRunSecurityChecksMissingMemoryCells(t *testing.T) {
	builtin := builtins.NewBitwiseBuiltinRunner(256)
	segments := memory.NewMemorySegmentManager()

	builtin.InitializeSegments(&segments)
	builtinBase := builtin.Base()
	// A bitwise cell consists of 5 elements: 2 input cells & 3 output cells
	// In this test we insert cells 4-5 and leave the first input cell empty
	// This will fail the security checks, as the memory cell with offset 0 will be missing
	builtinSegment := []memory.MaybeRelocatable{
		*memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint(1)),
		*memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint(2)),
		*memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint(3)),
		*memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint(4)),
	}
	segments.LoadData(builtinBase.AddUint(1), &builtinSegment)

	err := builtins.RunSecurityChecksForBuiltin(builtin, &segments)
	if err == nil {
		t.Errorf("RunSecurityChecks should have failed")
	}
}

func TestRunSecurityChecksMissingMemoryCellsNCheck(t *testing.T) {
	builtin := builtins.NewBitwiseBuiltinRunner(256)
	segments := memory.NewMemorySegmentManager()

	builtin.InitializeSegments(&segments)
	builtinBase := builtin.Base()
	// n = max(offsets) // cellsPerInstance + 1
	// n = max[(0]) // 5 + 1 = 0 // 5 + 1 = 1
	// len(offsets) // inputCells = 1 // 2
	// This will fail the security checks, as n > len(offsets) // inputCells
	builtinSegment := []memory.MaybeRelocatable{
		*memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint(1)),
	}
	segments.LoadData(builtinBase, &builtinSegment)

	err := builtins.RunSecurityChecksForBuiltin(builtin, &segments)
	if err == nil {
		t.Errorf("RunSecurityChecks should have failed")
	}
}

func TestRunSecurityChecksValidateOutputCellsNotDeducedOk(t *testing.T) {
	builtin := builtins.NewBitwiseBuiltinRunner(256)
	segments := memory.NewMemorySegmentManager()

	builtin.InitializeSegments(&segments)
	builtinBase := builtin.Base()
	// A bitwise cell consists of 5 elements: 2 input cells & 3 output cells
	// In this test we insert the input cells (1-2), but not the output cells (3-5)
	// This will cause the security checks to run the auto-deductions for those output cells
	builtinSegment := []memory.MaybeRelocatable{
		*memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint(1)),
		*memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint(2)),
	}
	segments.LoadData(builtinBase, &builtinSegment)

	err := builtins.RunSecurityChecksForBuiltin(builtin, &segments)
	if err != nil {
		t.Errorf("RunSecurityChecks failed with error: %s", err.Error())
	}
}

func TestRunSecurityChecksValidateOutputCellsNotDeducedErr(t *testing.T) {
	builtin := builtins.NewBitwiseBuiltinRunner(256)
	segments := memory.NewMemorySegmentManager()

	builtin.InitializeSegments(&segments)
	builtinBase := builtin.Base()
	// A bitwise cell consists of 5 elements: 2 input cells & 3 output cells
	// In this test we insert the input cells (1-2), but not the output cells (3-5)
	// This will cause the security checks to run the auto-deductions for those output cells
	// As we inserted an invalid value (PRIME -1) on the first input cell, this deduction will fail
	builtinSegment := []memory.MaybeRelocatable{
		*memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromDecString("-1")),
		*memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint(2)),
	}
	segments.LoadData(builtinBase, &builtinSegment)

	err := builtins.RunSecurityChecksForBuiltin(builtin, &segments)
	if err == nil {
		t.Errorf("RunSecurityChecks should have failed")
	}
}
