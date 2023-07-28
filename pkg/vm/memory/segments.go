package memory

// MemorySegmentManager manages the list of memory segments.
// Also holds metadata useful for the relocation process of
// the memory at the end of the VM run.
type MemorySegmentManager struct {
	segmentSizes map[uint]uint
	Memory       Memory
}

func NewMemorySegmentManager() *MemorySegmentManager {
	memory := NewMemory()
	return &MemorySegmentManager{make(map[uint]uint), *memory}
}

func (m *MemorySegmentManager) Add() Relocatable {
	ptr := Relocatable{int(m.Memory.num_segments), 0}
	m.Memory.num_segments += 1
	return ptr
}

// Writes data into the memory from address ptr and returns the first address after the data.
// If any insertion fails, returns (0,0) and the memory insertion error
func (m *MemorySegmentManager) LoadData(ptr Relocatable, data *[]MaybeRelocatable) (Relocatable, error) {
	for _, val := range *data {
		err := m.Memory.Insert(ptr, &val)
		if err != nil {
			return Relocatable{0, 0}, err
		}
		ptr.offset += 1
	}
	return ptr, nil
}
