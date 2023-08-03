package vm

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

func CairoRun(data []string) error {
	return errors.New("Unimplemented")
}

func WriteEncodedTrace(relocatedTrace *[]RelocatedTraceEntry, dest io.Writer) error {
	for i, entry := range *relocatedTrace {
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
