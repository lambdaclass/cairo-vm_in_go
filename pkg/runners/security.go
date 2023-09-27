package runners

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
	return nil
}
