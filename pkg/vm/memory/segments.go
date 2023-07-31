package memory

// MemorySegmentManager manages the list of memory segments.
// Also holds metadata useful for the relocation process of
// the memory at the end of the VM run.
type MemorySegmentManager struct {
	SegmentSizes map[uint]uint
	Memory       Memory
}

func NewMemorySegmentManager() *MemorySegmentManager {
	memory := NewMemory()
	return &MemorySegmentManager{make(map[uint]uint), *memory}
}

// Adds a memory segment and returns the first address of the new segment
func (m *MemorySegmentManager) AddSegment() Relocatable {
	ptr := Relocatable{int(m.Memory.num_segments), 0}
	m.Memory.num_segments += 1
	return ptr
}

func (m *MemorySegmentManager) RelocateSegments() ([]uint, bool) {
	if m.SegmentSizes == nil {
		return nil, false
	}

	first_addr := uint(1)
	relocation_table := []uint{first_addr}

	for key, value := range m.SegmentSizes {
		for uint(len(relocation_table)) <= key {
			relocation_table = append(relocation_table, relocation_table[len(relocation_table)-1])
		}
		new_addr := relocation_table[key] + value
		relocation_table = append(relocation_table, new_addr)
	}
	relocation_table = relocation_table[:len(relocation_table)-1]

	return relocation_table, true
}
