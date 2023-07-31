package memory

import (
	"errors"
)

// Memory represents the Cairo VM's memory.
type Memory struct {
	data         map[Relocatable]MaybeRelocatable
	num_segments uint
}

func NewMemory() *Memory {
	data := make(map[Relocatable]MaybeRelocatable)
	return &Memory{data, 0}
}

func (m *Memory) Data() *map[Relocatable]MaybeRelocatable {
	return &m.data
}

func (m *Memory) NumSegments() uint {
	return m.num_segments
}

// Inserts a value in some memory address, given by a Relocatable value.
func (m *Memory) Insert(addr Relocatable, val *MaybeRelocatable) error {
	// FIXME: There should be a special handling if the key
	// segment index is negative. This is an edge
	// case, so for now let's raise an error.
	if addr.segmentIndex < 0 {
		return errors.New("Segment index of key is negative - unimplemented")
	}

	// Check that insertions are preformed within the memory bounds
	if addr.segmentIndex >= int(m.num_segments) {
		return errors.New("Error: Inserting into a non allocated segment")
	}

	// Check for possible overwrites
	prev_elem, ok := m.data[addr]
	if ok && prev_elem != *val {
		return errors.New("Memory is write-once, cannot overwrite memory value")
	}

	m.data[addr] = *val

	return nil
}

// Gets some value stored in the memory address `addr`.
func (m *Memory) Get(addr Relocatable) (*MaybeRelocatable, error) {
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
	value, ok := m.data[addr]

	if !ok {
		return nil, errors.New("Memory Get: Value not found")
	}

	return &value, nil
}
