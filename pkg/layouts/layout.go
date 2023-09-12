package layouts

import (
	"github.com/pkg/errors"

	"github.com/lambdaclass/cairo-vm.go/pkg/builtins"
)

// Representation of a cairo layout.
// Stores the layout name and the particular builtin instances and
// their configuration for it.
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

// Given a layout name, return the builtin runners configuration for that layout.
func GetLayoutBuiltinRunners(layout string) ([]builtins.BuiltinRunner, error) {
	switch layout {
	case "plain":
		return []builtins.BuiltinRunner{builtins.NewOutputBuiltinRunner()}, nil

	// FIXME: Layout "small" does not really configure all these builtins, just adding them
	// here until we have all builtins implemented.
	case "small":
		return []builtins.BuiltinRunner{
			builtins.NewOutputBuiltinRunner(),
			builtins.DefaultRangeCheckBuiltinRunner(),
			builtins.DefaultBitwiseBuiltinRunner(),
			builtins.DefaultKeccakBuiltinRunner(),
			builtins.NewPoseidonBuiltinRunner()}, nil
	default:
		return nil, errors.Errorf("layout not supported: %s", layout)
	}
}
