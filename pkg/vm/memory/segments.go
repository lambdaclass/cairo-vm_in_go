package memory

import (
	"fmt"

	"github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
	"github.com/pkg/errors"
)

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
	ptr := Relocatable{int(m.Memory.numSegments), 0}
	m.Memory.numSegments += 1
	return ptr
}

// Calculates the size of each memory segment.
func (m *MemorySegmentManager) ComputeEffectiveSizes() map[uint]uint {
	if len(m.SegmentSizes) == 0 {
		for ptr := range m.Memory.Data {
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

	for i := uint(0); i < m.Memory.numSegments; i++ {
		new_addr := relocation_table[i] + m.SegmentSizes[i]
		relocation_table = append(relocation_table, new_addr)
	}
	relocation_table = relocation_table[:len(relocation_table)-1]

	return relocation_table, true
}

// Relocates the VM's memory, turning bidimensional indexes into contiguous numbers, and values
// into Felt252s. Uses the relocation_table to assign each index a number according to the value
// on its segment number.
func (s *MemorySegmentManager) RelocateMemory(relocationTable *[]uint) (map[uint]lambdaworks.Felt, error) {
	relocatedMemory := make(map[uint]lambdaworks.Felt, 0)

	for i := uint(0); i < s.Memory.numSegments; i++ {
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

func (m *MemorySegmentManager) GetSegmentUsedSize(segmentIdx uint) (uint, error) {
	size, ok := m.SegmentSizes[segmentIdx]
	if !ok {
		return 0, errors.Errorf("segment %d used size not found", segmentIdx)
	}
	return size, nil
}

// Go through each segment, calculate its size (counting holes), then count memory accesses. Substract the two and you
// get the holes for that segment. Sum each value and that's it.
// IMPORTANT: Builtin Segments DO NOT HAVE HOLES, so we don't need to count them.
// This function assumes you have already called `ComputeEffectiveSizes`, if you haven't, you'll get the wrong
// result
func (m *MemorySegmentManager) GetMemoryHoles(builtinCount uint) (uint, error) {
	var memoryHoles uint
	accessedCellsBySegment := make(map[uint]uint)

	var builtinSegmentsStart uint = 1
	var builtinSegmentsEnd uint = builtinSegmentsStart + builtinCount

	for address := range m.Memory.AccessedAddresses {
		if uint(address.SegmentIndex) > builtinSegmentsStart && uint(address.SegmentIndex) <= builtinSegmentsEnd {
			continue
		}

		accessedCellsBySegment[uint(address.SegmentIndex)]++
	}

	fmt.Println(m.SegmentSizes)
	for segmentIndex, size := range m.SegmentSizes {
		if segmentIndex > builtinSegmentsStart && segmentIndex <= builtinSegmentsEnd {
			continue
		}

		memoryHoles += size - accessedCellsBySegment[segmentIndex]
	}

	return memoryHoles, nil
}
