package cairo_run

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

type CairoRunConfig struct {
	TraceFile  *string
	MemoryFile *string
}

func CairoRun(programPath string) (runners.CairoRunner, error) {
	compiledProgram := parser.Parse(programPath)
	programJson := vm.DeserializeProgramJson(compiledProgram)

	cairoRunner := runners.NewCairoRunner(programJson)
	end, err := cairoRunner.Initialize()
	err = cairoRunner.RunUntilPC(end)
	cairoRunner.Vm.Relocate()
	return *cairoRunner, err
}

// Writes the trace binary representation.
//
// Bincode encodes to little endian by default and each trace entry is composed of
// 3 usize values that are padded to always reach 64 bit size.
func WriteEncodedTrace(relocatedTrace []vm.RelocatedTraceEntry, dest io.Writer) error {
	for i, entry := range relocatedTrace {
		ap_buffer := make([]byte, 8)
		ap, err := entry.Ap.ToU64()
		if err != nil {
			return err
		}
		binary.LittleEndian.PutUint64(ap_buffer, ap)
		_, err = dest.Write(ap_buffer)
		if err != nil {
			return encodeTraceError(i, err)
		}

		fp_buffer := make([]byte, 8)
		fp, err := entry.Fp.ToU64()
		if err != nil {
			return err
		}
		binary.LittleEndian.PutUint64(fp_buffer, fp)
		_, err = dest.Write(fp_buffer)
		if err != nil {
			return encodeTraceError(i, err)
		}

		pc_buffer := make([]byte, 8)
		pc, err := entry.Pc.ToU64()
		if err != nil {
			return err
		}
		binary.LittleEndian.PutUint64(pc_buffer, pc)
		_, err = dest.Write(pc_buffer)
		if err != nil {
			return encodeTraceError(i, err)
		}
	}

	return nil
}

func encodeTraceError(i int, err error) error {
	return errors.New(fmt.Sprintf("Failed to encode trace at position %d, serialize error: %s", i, err))
}
