package vm

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"sort"
)

func CairoRun(data []string) error {
	return errors.New("Unimplemented")
}

// Writes the trace binary representation.
//
// Bincode encodes to little endian by default and each trace entry is composed of
// 3 usize values that are padded to always reach 64 bit size.
func WriteEncodedTrace(relocatedTrace []RelocatedTraceEntry, dest io.Writer) error {
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
	return fmt.Errorf("failed to encode trace at position %d, serialize error: %s", i, err)
}

// Writes a binary representation of the relocated memory.
//
// The memory pairs (address, value) are encoded and concatenated:
// * address -> 8-byte encoded
// * value -> 32-byte encoded
func WriteEncodedMemory(relocatedMemory map[uint]uint, dest io.Writer) error {
	// create a slice to store keys of the relocatedMemory map
	keysMap := make([]uint, 0, len(relocatedMemory))
	for k := range relocatedMemory {
		keysMap = append(keysMap, k)
	}

	// sort the keys
	sort.Slice(keysMap, func(i, j int) bool { return keysMap[i] < keysMap[j] })

	// iterate over the `relocatedMemory` map in sorted key order
	for _, k := range keysMap {

		// get relocatedMemory[k]
		value := relocatedMemory[k]
		fmt.Printf("key[%d] = %d\n", k, value)

		// write the key
		keyArray := make([]byte, 8)
		binary.LittleEndian.PutUint64(keyArray, uint64(k))
		_, err := dest.Write(keyArray)
		if err != nil {
			return encodeMemoryError(k, err)
		}

		// write the value
		valueArray := make([]byte, 8)
		binary.LittleEndian.PutUint64(valueArray, uint64(value))
		_, err = dest.Write(valueArray)
		if err != nil {
			return encodeMemoryError(k, err)
		}
	}

	return nil
}

func encodeMemoryError(i uint, err error) error {
	return fmt.Errorf("failed to encode trace at position %d, serialize error: %s", i, err)
}
