package vm

import (
	"github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
	"github.com/lambdaclass/cairo-vm.go/pkg/parser"
	"github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
)

type Identifier struct {
	FullName    string
	Members     map[string]any
	Size        int
	Decorators  []string
	PC          int
	Type        string
	CairoType   string
	Value       lambdaworks.Felt
	Destination string
}

type Program struct {
	Data             []memory.MaybeRelocatable
	Builtins         []string
	Identifiers      map[string]Identifier
	Hints            map[uint][]parser.HintParams
	ReferenceManager parser.ReferenceManager
	Start            uint
	End              uint
}

func DeserializeProgramJson(compiledProgram parser.CompiledJson) Program {
	var program Program

	hexData := compiledProgram.Data
	for _, hexVal := range hexData {
		felt := lambdaworks.FeltFromHex(hexVal)
		program.Data = append(program.Data, *memory.NewMaybeRelocatableFelt(felt))
	}
	program.Builtins = compiledProgram.Builtins
	program.Identifiers = make(map[string]Identifier)

	start := uint(compiledProgram.Identifiers["__main__.__start__"].PC)
	end := uint(compiledProgram.Identifiers["__main__.__end__"].PC)
	program.Start = start
	program.End = end

	for key, identifier := range compiledProgram.Identifiers {
		var programIdentifier Identifier
		programIdentifier.FullName = identifier.FullName
		programIdentifier.Members = identifier.Members
		programIdentifier.Size = identifier.Size
		programIdentifier.Decorators = identifier.Decorators
		programIdentifier.PC = identifier.PC
		programIdentifier.Type = identifier.Type
		programIdentifier.CairoType = identifier.CairoType
		programIdentifier.Value = lambdaworks.FeltFromDecString(identifier.Value.String())
		programIdentifier.Destination = identifier.Destination
		program.Identifiers[key] = programIdentifier
	}
	program.Hints = compiledProgram.Hints
	program.ReferenceManager = compiledProgram.ReferenceManager

	return program
}

func (p *Program) ExtractConstants() map[string]lambdaworks.Felt {
	constants := make(map[string]lambdaworks.Felt)
	for name, identifier := range p.Identifiers {
		switch identifier.Type {
		case "const":
			constants[name] = identifier.Value
		case "alias":
			val, ok := searchConstFromAlias(identifier.Destination, &p.Identifiers)
			if ok {
				constants[name] = val
			}
		}
	}

	return constants
}

func searchConstFromAlias(destination string, identifiers *map[string]Identifier) (lambdaworks.Felt, bool) {
	identifier, ok := (*identifiers)[destination]
	if ok {
		switch identifier.Type {
		case "const":
			return identifier.Value, true
		case "alias":
			return searchConstFromAlias(identifier.Destination, identifiers)
		}
	}
	return lambdaworks.Felt{}, false
}
