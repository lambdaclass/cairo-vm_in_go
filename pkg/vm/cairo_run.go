package vm

import "github.com/lambdaclass/cairo-vm.go/pkg/parser"

type RunResources struct {
	NSteps *uint
}

type BuiltinHintProcessor struct {
	//Todo: ExtraHints update to map[string][dyn func]
	ExtraHints   map[string]string
	RunResources *RunResources
}

type CairoRunConfig struct {
	Entrypoint   string
	TraceEnabled bool
	RelocateMem  bool
	Layout       string
	ProofMode    bool
	SecureRun    *bool
}

func CairoRunProgram(program parser.CompiledJson) error {
	hintExecutor := BuiltinHintProcessor{}
	cairoRunConfig := CairoRunConfig{}
	CairoRun(program, hintExecutor, cairoRunConfig)

}

func CairoRun(program parser.CompiledJson, hintExecutor BuiltinHintProcessor, cairoRunConfig CairoRunConfig) error {
	programJson := DeserializeProgramJson(program)
}
