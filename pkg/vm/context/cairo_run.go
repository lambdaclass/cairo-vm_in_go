package context

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"

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

// Writes the trace binary representation.
//
// Bincode encodes to little endian by default and each trace entry is composed of
// 3 usize values that are padded to always reach 64 bit size.
func WriteEncodedTrace(relocatedTrace []vm.RelocatedTraceEntry, dest io.Writer) error {
	for i, entry := range relocatedTrace {
		ap := make([]byte, 8)
		binary.LittleEndian.PutUint64(ap, uint64(entry.Ap))
		_, err := dest.Write(ap)
		if err != nil {
			return encodeTraceError(i, err)
		}

		fp := make([]byte, 8)
		binary.LittleEndian.PutUint64(fp, uint64(entry.Fp))
		_, err = dest.Write(fp)
		if err != nil {
			return encodeTraceError(i, err)
		}

		pc := make([]byte, 8)
		binary.LittleEndian.PutUint64(pc, uint64(entry.Pc))
		_, err = dest.Write(pc)
		if err != nil {
			return encodeTraceError(i, err)
		}
	}

	return nil
}

func encodeTraceError(i int, err error) error {
	return errors.New(fmt.Sprintf("Failed to encode trace at position %d, serialize error: %s", i, err))
}
