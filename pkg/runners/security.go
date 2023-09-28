package runners

import (
	"github.com/lambdaclass/cairo-vm.go/pkg/builtins"
	"github.com/pkg/errors"
)

/*
Verify that the completed run in a runner is safe to be relocated and be
used by other Cairo programs.

Checks include:
  - (Only if `verifyBuiltins` is set to true) All accesses to the builtin segments must be within the range defined by
    the builtins themselves.
  - There must not be accesses to the program segment outside the program
    data range. This check will use the `programSegmentSize` instead of the program data length if available.
  - All addresses in memory must be real (not temporary)

Note: Each builtin is responsible for checking its own segments' data.
*/
func VerifySecureRunner(runner *CairoRunner, verifyBuiltins bool, programSegmentSize *uint) error {
	programSize := uint(len(runner.Program.Data))
	if programSegmentSize != nil {
		programSize = *programSegmentSize
	}
	// Get builtin segment info
	builtinNames := make(map[int]string)
	builtinSizes := make(map[int]uint)
	if verifyBuiltins {
		for i := 0; i < len(runner.Vm.BuiltinRunners); i++ {
			base, stopPtr, err := runner.Vm.BuiltinRunners[i].GetMemorySegmentAddresses()
			if err != nil {
				return err
			}
			builtinNames[base.SegmentIndex] = runner.Vm.BuiltinRunners[i].Name()
			builtinSizes[base.SegmentIndex] = stopPtr.Offset
		}
	}
	// Run memory checks
	for addr, val := range runner.Vm.Segments.Memory.Data {
		// Check out of bound accesses to builtin segment
		size, ok := builtinSizes[addr.SegmentIndex]
		if ok && addr.Offset >= size {
			return errors.Errorf("Out of bounds access to builtin segment %s at %s", builtinNames[addr.SegmentIndex], addr.ToString())
		}
		// Check out of bound accesses to program segment
		if addr.SegmentIndex == runner.ProgramBase.SegmentIndex && addr.SegmentIndex >= int(programSize) {
			return errors.Errorf("Out of bounds access to program segment at %s", addr.ToString())
		}
		// Check non-relocated temporary addresses
		if addr.SegmentIndex < 0 {
			return errors.Errorf("Security Error: Invalid Memory Value: temporary address not relocated: %s", addr.ToString())
		}
		relVal, isRel := val.GetRelocatable()
		if isRel && relVal.SegmentIndex < 0 {
			return errors.Errorf("Security Error: Invalid Memory Value: temporary address not relocated: %s", relVal.ToString())
		}
	}
	// Run builtin-specific checks
	for i := 0; i < len(runner.Vm.BuiltinRunners); i++ {
		err := builtins.RunSecurityChecksForBuiltin(runner.Vm.BuiltinRunners[i], &runner.Vm.Segments)
		if err != nil {
			return err
		}
	}

	return nil
}
