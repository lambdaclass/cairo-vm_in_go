package vm

import (
	"github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
	"github.com/lambdaclass/cairo-vm.go/pkg/parser"
	"github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
)

type Program struct {
	Data []memory.MaybeRelocatable
}

func DeserializeProgramJson(compiledProgram parser.CompiledJson) Program {
	var program Program

	hexData := compiledProgram.Data
	for _, hexVal := range hexData {
		felt := lambdaworks.FeltFromHex(hexVal)
		program.Data = append(program.Data, *memory.NewMaybeRelocatableInt())
	}
}
