package memory

import "github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"

// MemorySegmentManager manages the list of memory segments.
// Also holds metadata useful for the relocation process of
// the memory at the end of the VM run.
type MemorySegmentManager struct {
	SegmentSizes map[uint]uint
	Memory       Memory
}

func NewMemorySegmentManager() MemorySegmentManager {
	memory := NewMemory()
	return MemorySegmentManager{make(map[uint]uint), *memory}
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

		for ptr := range m.Memory.data {
			segmentIndex := uint(ptr.SegmentIndex)
			segmentMaxSize := m.SegmentSizes[segmentIndex]
			segmentSize := ptr.Offset + 1
			if segmentSize > segmentMaxSize {
				m.SegmentSizes[segmentIndex] = segmentSize
			}
		}
	}

	return m.SegmentSizes
}

// Returns a vector containing the first relocated address of each memory segment
func (m *MemorySegmentManager) RelocateSegments() ([]uint, bool) {
	if m.SegmentSizes == nil {
		return nil, false
	}

	first_addr := uint(1)
	relocation_table := []uint{first_addr}

	for i := uint(0); i < m.Memory.NumSegments(); i++ {
		new_addr := relocation_table[i] + m.SegmentSizes[i]
		relocation_table = append(relocation_table, new_addr)
	}
	relocation_table = relocation_table[:len(relocation_table)-1]

	return relocation_table, true
}

// Relocates the VM's memory, turning bidimensional indexes into contiguous numbers, and values
// into Felt252s. Uses the relocation_table to asign each index a number according to the value
// on its segment number.
func (s *MemorySegmentManager) RelocateMemory(relocationTable *[]uint) (map[uint]lambdaworks.Felt, error) {
	relocatedMemory := make(map[uint]lambdaworks.Felt, 0)

	for i := uint(0); i < s.Memory.NumSegments(); i++ {
		for j := uint(0); j < s.SegmentSizes[i]; j++ {
			ptr := NewRelocatable(int(i), j)
			cell, err := s.Memory.Get(ptr)
			if err == nil {
				relocatedAddr := ptr.RelocateAddress(relocationTable)
				value, err := cell.RelocateValue(relocationTable)
				if err != nil {
					return nil, err
				}
				relocatedMemory[relocatedAddr] = value
			}
		}
	}

	return relocatedMemory, nil
}

// Writes data into the memory from address ptr and returns the first address after the data.
// If any insertion fails, returns (0,0) and the memory insertion error
func (m *MemorySegmentManager) LoadData(ptr Relocatable, data *[]MaybeRelocatable) (Relocatable, error) {
	for _, val := range *data {
		err := m.Memory.Insert(ptr, &val)
		if err != nil {
			return Relocatable{0, 0}, err
		}
		ptr.Offset += 1
	}
	return ptr, nil
}
