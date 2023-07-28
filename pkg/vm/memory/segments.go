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
