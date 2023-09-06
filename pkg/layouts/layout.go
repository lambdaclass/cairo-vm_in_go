package layouts

import "github.com/lambdaclass/cairo-vm.go/pkg/builtins"

type CairoLayout struct {
	name     string
	builtins BuiltinsInstance
	// TODO - Add when necessary:
	// cpuComponentStep uint,
	// rcUnits uint,
	// publicMemoryFraction uint,
	// memoryUnitsPerStep uint,
	// dilutedPoolInstanceDef DilutedPoolInstanceDef
	// nTraceColums uint
	// cpuInstanceDef CpuInstanceDef
}

type BuiltinsInstance struct {
	output     bool
	rangeCheck builtins.RangeCheckBuiltinRunner
	bitwise    builtins.BitwiseBuiltinRunner
	poseidon   builtins.PoseidonBuiltinRunner
}

func PlainBuiltinsInstance() BuiltinsInstance {
	instance := new(BuiltinsInstance)
	instance.output = false
	return *instance
}

func PlainLayout() CairoLayout {
	builtins := PlainBuiltinsInstance()
	return CairoLayout{"plain", builtins}
}
