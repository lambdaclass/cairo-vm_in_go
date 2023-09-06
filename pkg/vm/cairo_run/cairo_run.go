package cairo_run

import (
	"encoding/binary"
	"fmt"
	"io"
	"sort"

	"github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
	"github.com/lambdaclass/cairo-vm.go/pkg/parser"
	"github.com/lambdaclass/cairo-vm.go/pkg/runners"
	"github.com/lambdaclass/cairo-vm.go/pkg/vm"
	"github.com/pkg/errors"
)

type RunResources struct {
	NSteps *uint
}

type CairoRunConfig struct {
	TraceFile  *string
	MemoryFile *string
}

func CairoRunError(err error) error {
	return errors.Wrapf(err, "Cairo Run Error\n")
}

func CairoRun(programPath string) (*runners.CairoRunner, error) {
	compiledProgram, err := parser.Parse(programPath)
	if err != nil {
		return nil, CairoRunError(err)
	}
	programJson := vm.DeserializeProgramJson(compiledProgram)

	cairoRunner, err := runners.NewCairoRunner(programJson)
	if err != nil {
		return nil, err
	}
	end, err := cairoRunner.Initialize()
	if err != nil {
		return nil, err
	}
	err = cairoRunner.RunUntilPC(end)
	if err != nil {
		return nil, err
	}
	err = cairoRunner.Vm.Relocate()
	return cairoRunner, err
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

// Writes a binary representation of the relocated memory.
//
// The memory pairs (address, value) are encoded and concatenated:
// * address -> 8-byte encoded
// * value -> 32-byte encoded
func WriteEncodedMemory(relocatedMemory map[uint]lambdaworks.Felt, dest io.Writer) error {
	// create a slice to store keys of the relocatedMemory map
	keysMap := make([]uint, 0, len(relocatedMemory))
	for k := range relocatedMemory {
		keysMap = append(keysMap, k)
	}

	// sort the keys
	sort.Slice(keysMap, func(i, j int) bool { return keysMap[i] < keysMap[j] })

	// iterate over the `relocatedMemory` map in sorted key order
	for _, k := range keysMap {

		// write the key
		keyArray := make([]byte, 8)
		binary.LittleEndian.PutUint64(keyArray, uint64(k))
		_, err := dest.Write(keyArray)
		if err != nil {
			return encodeMemoryError(k, err)
		}

		// write the value
		valueArray := relocatedMemory[k].ToLeBytes()

		_, err = dest.Write(valueArray[:])
		if err != nil {
			return encodeMemoryError(k, err)
		}
	}

	return nil
}

func encodeMemoryError(i uint, err error) error {
	return fmt.Errorf("Failed to encode trace at position %d, serialize error: %s", i, err)
}
