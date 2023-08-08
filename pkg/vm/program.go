package vm

import (
	"github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
	"github.com/lambdaclass/cairo-vm.go/pkg/parser"
	"github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
)

type Program struct {
	Data        []memory.MaybeRelocatable
	Builtins    []string
	Identifiers *map[string]parser.Identifier
}

func DeserializeProgramJson(compiledProgram parser.CompiledJson) Program {
	var program Program

	hexData := compiledProgram.Data
	for _, hexVal := range hexData {
		felt := lambdaworks.FeltFromHex(hexVal)
		program.Data = append(program.Data, *memory.NewMaybeRelocatableFelt(felt))
	}
	program.Builtins = compiledProgram.Builtins
	program.Identifiers = &compiledProgram.Identifiers

	return program
}
