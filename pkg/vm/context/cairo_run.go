package context

import (
	"github.com/lambdaclass/cairo-vm.go/pkg/parser"
	"github.com/lambdaclass/cairo-vm.go/pkg/runners"
	"github.com/lambdaclass/cairo-vm.go/pkg/vm"
)

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

func CairoRunProgram(compiledProgram parser.CompiledJson) error {
	hintExecutor := BuiltinHintProcessor{}
	cairoRunConfig := CairoRunConfig{}
	_, _, err := CairoRun(compiledProgram, hintExecutor, cairoRunConfig)
	return err
}

func CairoRun(compiledProgram parser.CompiledJson, hintExecutor BuiltinHintProcessor, cairoRunConfig CairoRunConfig) (runners.CairoRunner, vm.VirtualMachine, error) {
	programJson := vm.DeserializeProgramJson(compiledProgram)
	vm := vm.NewVirtualMachine()
	cairoRunner := runners.NewCairoRunner(programJson)
	_, err := cairoRunner.Initialize()
	//err = cairoRunner.RunUntilPC()
	//cairoRunner.EndRun()
	//vm.VerifyAutoDeductions()
	//cairoRunner.ReadReturnValues(vm)
	/*if cairoRunConfig.ProofMode {
		cairoRunner.FinalizeSegments(vm)
	}*/
	//cairoRunner.relocate(vm, cairoRunConfig.RelocateMem)
	return *cairoRunner, *vm, err
}
