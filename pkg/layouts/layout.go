package layouts

import (
	"github.com/pkg/errors"

	"github.com/lambdaclass/cairo-vm.go/pkg/builtins"
)

type CairoLayout struct {
	Name     string
	Builtins []builtins.BuiltinRunner
	// TODO - Add when necessary:
	// cpuComponentStep uint,
	// rcUnits uint,
	// publicMemoryFraction uint,
	// memoryUnitsPerStep uint,
	// dilutedPoolInstanceDef DilutedPoolInstanceDef
	// nTraceColums uint
	// cpuInstanceDef CpuInstanceDef
}

func GetLayoutBuiltinRunners(layout string) ([]builtins.BuiltinRunner, error) {
	switch layout {
	case "plain":
		return []builtins.BuiltinRunner{builtins.NewOutputBuiltinRunner()}, nil
	case "small":
		return []builtins.BuiltinRunner{builtins.NewOutputBuiltinRunner(), builtins.NewBitwiseBuiltinRunner()}, nil
	default:
		return nil, errors.Errorf("layout not supported: %s", layout)
	}
}

// type BuiltinsInstance struct {
// 	Output     bool
// 	rangeCheck builtins.RangeCheckBuiltinRunner
// 	bitwise    builtins.BitwiseBuiltinRunner
// 	poseidon   builtins.PoseidonBuiltinRunner
// }

// func PlainBuiltinsInstance() BuiltinsInstance {
// 	instance := new(BuiltinsInstance)
// 	instance.output = false
// 	return *instance
// }

// func PlainLayout() CairoLayout {
// 	builtins := PlainBuiltinsInstance()
// 	return CairoLayout{"plain", builtins}
// }
