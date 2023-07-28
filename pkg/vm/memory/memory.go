package memory

import (
	"errors"
)

// Memory represents the Cairo VM's memory.
type Memory struct {
	data [][]MaybeRelocatable
}

func NewMemory(data [][]MaybeRelocatable) *Memory {
	return &Memory{data}
}

// Inserts a value in some memory address, given by a Relocatable value.
func (m *Memory) Insert(addr *Relocatable, val *MaybeRelocatable) error {
	addr_idx, addr_offset := addr.into_indexes()

	// FIXME: There should be a special handling if the key
	// segment index is negative. This is an edge
	// case, so for now let's raise an error.
	if addr.segmentIndex < 0 {
		return errors.New("Segment index of key is negative - unimplemented")
	}

	segment := &m.data[addr_idx]
	segment_len := len(*segment)

	// When the offset of the insertion address is greater than the max
	// offset of the segment, memory cells are filled with `nil` in the
	// intermediate values, if any. So if segment has length 2 (last idx is 1)
	// and we want to insert something at index 4, index 2 and 3 will be filled
	// with `nil`, and index 4 will have the desired value.
	if segment_len <= int(addr_offset) {
		new_segment_len := addr_offset + 1
		for i := segment_len; i < int(new_segment_len); i++ {
			*segment = append(*segment, MaybeRelocatable{nil})
		}
	}

	// At this point, something exists at the `addr_offset` for sure.
	// Check that the value at that offset is `nil` and if it is, then
	// swap that `nil` with the desired value.
	if (*segment)[addr_offset].is_nil() {
		(*segment)[addr_offset] = *val
		// If there wasn't `nil`, then we are trying to overwrite in that
		// address. If the value we are trying to insert is not the same as
		// the one that was already in that location, raise an error.
	} else if (*segment)[addr_offset] != *val {
		return errors.New("Memory is write-once, cannot overwrite memory value")
	}

	return nil
}

// Gets some value stored in the memory address `addr`.
func (m *Memory) Get(addr *Relocatable) (*MaybeRelocatable, error) {
	addr_idx, addr_offset := addr.into_indexes()

	// FIXME: There should be a special handling if the key
	// segment index is negative. This is an edge
	// case, so for now let's raise an error.
	if addr.segmentIndex < 0 {
		return nil, errors.New("Segment index of key is negative - unimplemented")
	}

	// FIXME: We should create a function for this value,
	// `relocate_value()` in the future. This function should
	// check if the value is a `Relocatable` with a negative
	// segment index. Again, these are edge cases so not important
	// right now. See cairo-vm code for details.
	value := m.data[addr_idx][addr_offset]

	return &value, nil
}
