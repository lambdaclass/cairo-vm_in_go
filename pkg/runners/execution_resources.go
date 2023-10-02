package runners

type ExecutionResources struct {
	NSteps                  uint
	NMemoryHoles            uint
	BuiltinsInstanceCounter map[string]uint
}
