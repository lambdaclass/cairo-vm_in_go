package memory

import "sort"

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

// Calculates the size of each memory segment.
func (m *MemorySegmentManager) ComputeEffectiveSizes() map[uint]uint {
	if len(m.SegmentSizes) == 0 {
		greatestIndex := uint(0)

		for ptr := range m.Memory.Data() {
			segmentIndex := uint(ptr.segmentIndex)
			segmentMaxSize, ok := m.SegmentSizes[segmentIndex]
			segmentSize := ptr.offset + 1
			if (ok && segmentSize > segmentMaxSize) || !ok {
				m.SegmentSizes[segmentIndex] = segmentSize
			}
			if segmentIndex > greatestIndex {
				greatestIndex = segmentIndex
			}
		}

		for i := greatestIndex + 1; i < m.Memory.NumSegments(); i++ {
			m.SegmentSizes[i] = 0
		}
	}

	return m.SegmentSizes
}

func (m *MemorySegmentManager) RelocateSegments() ([]uint, bool) {
	if m.SegmentSizes == nil {
		return nil, false
	}

	first_addr := uint(1)
	relocation_table := []uint{first_addr}

	sorted_keys := make([]uint, 0, len(m.SegmentSizes))
	for key := range m.SegmentSizes {
		sorted_keys = append(sorted_keys, key)
	}
	sort.Slice(sorted_keys, func(i, j int) bool { return sorted_keys[i] < sorted_keys[j] })

	for _, key := range sorted_keys {
		for uint(len(relocation_table)) <= key {
			relocation_table = append(relocation_table, relocation_table[len(relocation_table)-1])
		}
		new_addr := relocation_table[key] + m.SegmentSizes[key]
		relocation_table = append(relocation_table, new_addr)
	}
	relocation_table = relocation_table[:len(relocation_table)-1]

	return relocation_table, true
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
