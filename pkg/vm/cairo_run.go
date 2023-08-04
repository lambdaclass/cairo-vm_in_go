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
		// write the key
		key := make([]byte, 8)
		binary.LittleEndian.PutUint64(key, uint64(k))
		_, err := dest.Write(key)
		if err != nil {
			return encodeMemoryError(int(k), err)
		}
	}

	/*
		pub fn write_encoded_memory(
		    relocated_memory: &[Option<Felt252>],
		    dest: &mut impl Writer,
		) -> Result<(), EncodeTraceError> {
		    for (i, memory_cell) in relocated_memory.iter().enumerate() {
		        match memory_cell {
		            None => continue,
		            Some(unwrapped_memory_cell) => {
		                dest.write(&(i as u64).to_le_bytes())
		                    .map_err(|e| EncodeTraceError(i, e))?;
		                dest.write(&unwrapped_memory_cell.to_le_bytes())
		                    .map_err(|e| EncodeTraceError(i, e))?;
		            }
		        }
		    }

		    Ok(())
		}
	*/
	return errors.New("Unimplemented")
}

func encodeMemoryError(i int, err error) error {
	return errors.New(fmt.Sprintf("Failed to encode trace at position %d, serialize error: %s", i, err))
}
