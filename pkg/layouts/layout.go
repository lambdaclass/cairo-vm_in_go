package layouts

import (
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
	RcUnits              uint
	PublicMemoryFraction uint
	MemoryUnitsPerStep   uint
	DilutedPoolInstance  *DilutedPoolInstanceDef
	// nTraceColums uint
	// cpuInstanceDef CpuInstanceDef
}

func NewPlainLayout() CairoLayout {
	return CairoLayout{
		Name:                 "plain",
		Builtins:             []builtins.BuiltinRunner{builtins.NewOutputBuiltinRunner()},
		RcUnits:              16,
		PublicMemoryFraction: 4,
		MemoryUnitsPerStep:   8,
		DilutedPoolInstance:  nil,
	}
}

func NewSmallLayout() CairoLayout {
	return CairoLayout{
		Name: "small",
		Builtins: []builtins.BuiltinRunner{
			builtins.NewOutputBuiltinRunner(),
			builtins.NewPedersenBuiltinRunner(),
			builtins.DefaultRangeCheckBuiltinRunner()},
		RcUnits:              16,
		PublicMemoryFraction: 4,
		MemoryUnitsPerStep:   8,
		DilutedPoolInstance:  nil,
	}
}

func NewAllCairoLayout() CairoLayout {
	return CairoLayout{
		Name: "all_cairo",
		Builtins: []builtins.BuiltinRunner{
			builtins.NewOutputBuiltinRunner(),
			builtins.DefaultPedersenBuiltinRunner(),
			builtins.DefaultRangeCheckBuiltinRunner(),
			builtins.DefaultBitwiseBuiltinRunner(),
			builtins.DefaultKeccakBuiltinRunner(),
			builtins.NewPoseidonBuiltinRunner()},
		RcUnits:              4,
		PublicMemoryFraction: 8,
		MemoryUnitsPerStep:   8,
		DilutedPoolInstance:  DefaultDilutedPoolInstance(),
	}
}
