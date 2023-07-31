package memory_test

import (
	"reflect"
	"testing"

	"github.com/lambdaclass/cairo-vm.go/pkg/vm"
	"github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
)

func TestComputeEffectiveSizeOneSegment(t *testing.T) {
	segments := memory.NewMemorySegmentManager()
	segments.AddSegment()
	segments.Memory.Insert(memory.NewRelocatable(0, 0), memory.NewMaybeRelocatableInt(1))
	segments.Memory.Insert(memory.NewRelocatable(0, 1), memory.NewMaybeRelocatableInt(1))
	segments.Memory.Insert(memory.NewRelocatable(0, 2), memory.NewMaybeRelocatableInt(1))

	segments.ComputeEffectiveSizes()

	expectedSizes := map[uint]uint{0: 3}
	if !reflect.DeepEqual(expectedSizes, segments.SegmentSizes) {
		t.Errorf("Segment sizes are not the same")
	}
}

func TestComputeEffectiveSizeOneSegmentWithOneGap(t *testing.T) {
	segments := memory.NewMemorySegmentManager()
	segments.AddSegment()
	segments.Memory.Insert(memory.NewRelocatable(0, 6), memory.NewMaybeRelocatableInt(1))

	segments.ComputeEffectiveSizes()

	expectedSizes := map[uint]uint{0: 7}
	if !reflect.DeepEqual(expectedSizes, segments.SegmentSizes) {
		t.Errorf("Segment sizes are not the same")
	}
}

func TestComputeEffectiveSizeOneSegmentWithMultipleGaps(t *testing.T) {
	segments := memory.NewMemorySegmentManager()
	segments.AddSegment()
	segments.Memory.Insert(memory.NewRelocatable(0, 3), memory.NewMaybeRelocatableInt(1))
	segments.Memory.Insert(memory.NewRelocatable(0, 4), memory.NewMaybeRelocatableInt(1))
	segments.Memory.Insert(memory.NewRelocatable(0, 7), memory.NewMaybeRelocatableInt(1))
	segments.Memory.Insert(memory.NewRelocatable(0, 9), memory.NewMaybeRelocatableInt(1))

	segments.ComputeEffectiveSizes()

	expectedSizes := map[uint]uint{0: 10}
	if !reflect.DeepEqual(expectedSizes, segments.SegmentSizes) {
		t.Errorf("Segment sizes are not the same")
	}
}

func TestComputeEffectiveSizeThreeSegments(t *testing.T) {
	segments := memory.NewMemorySegmentManager()
	segments.AddSegment()
	segments.AddSegment()
	segments.AddSegment()
	segments.Memory.Insert(memory.NewRelocatable(0, 0), memory.NewMaybeRelocatableInt(1))
	segments.Memory.Insert(memory.NewRelocatable(0, 1), memory.NewMaybeRelocatableInt(1))
	segments.Memory.Insert(memory.NewRelocatable(0, 2), memory.NewMaybeRelocatableInt(1))
	segments.Memory.Insert(memory.NewRelocatable(1, 0), memory.NewMaybeRelocatableInt(1))
	segments.Memory.Insert(memory.NewRelocatable(1, 1), memory.NewMaybeRelocatableInt(1))
	segments.Memory.Insert(memory.NewRelocatable(1, 2), memory.NewMaybeRelocatableInt(1))
	segments.Memory.Insert(memory.NewRelocatable(2, 0), memory.NewMaybeRelocatableInt(1))
	segments.Memory.Insert(memory.NewRelocatable(2, 1), memory.NewMaybeRelocatableInt(1))
	segments.Memory.Insert(memory.NewRelocatable(2, 2), memory.NewMaybeRelocatableInt(1))

	segments.ComputeEffectiveSizes()

	expectedSizes := map[uint]uint{0: 3, 1: 3, 2: 3}
	if !reflect.DeepEqual(expectedSizes, segments.SegmentSizes) {
		t.Errorf("Segment sizes are not the same")
	}
}

func TestComputeEffectiveSizeThreeSegmentsWithGaps(t *testing.T) {
	segments := memory.NewMemorySegmentManager()
	segments.AddSegment()
	segments.AddSegment()
	segments.AddSegment()
	segments.Memory.Insert(memory.NewRelocatable(0, 2), memory.NewMaybeRelocatableInt(1))
	segments.Memory.Insert(memory.NewRelocatable(0, 5), memory.NewMaybeRelocatableInt(1))
	segments.Memory.Insert(memory.NewRelocatable(0, 7), memory.NewMaybeRelocatableInt(1))
	segments.Memory.Insert(memory.NewRelocatable(1, 1), memory.NewMaybeRelocatableInt(1))
	segments.Memory.Insert(memory.NewRelocatable(2, 2), memory.NewMaybeRelocatableInt(1))
	segments.Memory.Insert(memory.NewRelocatable(2, 4), memory.NewMaybeRelocatableInt(1))
	segments.Memory.Insert(memory.NewRelocatable(2, 7), memory.NewMaybeRelocatableInt(1))

	segments.ComputeEffectiveSizes()

	expectedSizes := map[uint]uint{0: 8, 1: 2, 2: 8}
	if !reflect.DeepEqual(expectedSizes, segments.SegmentSizes) {
		t.Errorf("Segment sizes are not the same")
	}
}

func TestGetSegmentUsedSizeAfterComputingUsed(t *testing.T) {
	segments := memory.NewMemorySegmentManager()
	segments.AddSegment()
	segments.AddSegment()
	segments.AddSegment()
	segments.Memory.Insert(memory.NewRelocatable(0, 2), memory.NewMaybeRelocatableInt(1))
	segments.Memory.Insert(memory.NewRelocatable(0, 5), memory.NewMaybeRelocatableInt(1))
	segments.Memory.Insert(memory.NewRelocatable(0, 7), memory.NewMaybeRelocatableInt(1))
	segments.Memory.Insert(memory.NewRelocatable(1, 1), memory.NewMaybeRelocatableInt(1))
	segments.Memory.Insert(memory.NewRelocatable(2, 2), memory.NewMaybeRelocatableInt(1))
	segments.Memory.Insert(memory.NewRelocatable(2, 4), memory.NewMaybeRelocatableInt(1))
	segments.Memory.Insert(memory.NewRelocatable(2, 7), memory.NewMaybeRelocatableInt(1))

	segments.ComputeEffectiveSizes()

	segmentSize, ok := segments.SegmentSizes[2]
	expectedSize := 8
	if !ok || segmentSize != uint(expectedSize) {
		t.Errorf("Segment size should be %d but it's %d", expectedSize, segmentSize)
	}
}

func TestGetSegmentUsedSizeBeforeComputingUsed(t *testing.T) {
	segments := memory.NewMemorySegmentManager()

	_, ok := segments.SegmentSizes[2]
	if ok {
		t.Errorf("Expected no segment sizes loaded")
	}
}

func TestRelocateOneSegment(t *testing.T) {
	segments := memory.NewMemorySegmentManager()
	segments.AddSegment()
	segments.SegmentSizes = map[uint]uint{0: 3}
	relocationTable, ok := segments.RelocateSegments()

	if !ok {
		t.Errorf("Memory segment manager doesn't have segment sizes initialized")
	}

	expectedTable := []uint{1}
	if !reflect.DeepEqual(expectedTable, relocationTable) {
		t.Errorf("Relocation tables are not the same")
	}
}

func TestRelocateFiveSegments(t *testing.T) {
	segments := memory.NewMemorySegmentManager()
	segments.AddSegment()
	segments.AddSegment()
	segments.AddSegment()
	segments.AddSegment()
	segments.AddSegment()
	segments.SegmentSizes = map[uint]uint{0: 3, 1: 3, 2: 56, 3: 78, 4: 8}
	relocationTable, ok := segments.RelocateSegments()

	if !ok {
		t.Errorf("Memory segment manager doesn't have segment sizes initialized")
	}

	expectedTable := []uint{1, 4, 7, 63, 141}
	if !reflect.DeepEqual(expectedTable, relocationTable) {
		t.Errorf("Relocation tables are not the same")
	}
}

func TestRelocateSegmentsWithHoles(t *testing.T) {
	segments := memory.NewMemorySegmentManager()
	segments.AddSegment()
	segments.AddSegment()
	segments.AddSegment()
	segments.SegmentSizes = map[uint]uint{0: 3, 2: 3}
	relocationTable, ok := segments.RelocateSegments()

	if !ok {
		t.Errorf("Memory segment manager doesn't have segment sizes initialized")
	}

	expectedTable := []uint{1, 4, 4}
	if !reflect.DeepEqual(expectedTable, relocationTable) {
		t.Errorf("Relocation tables are not the same")
	}
}

func TestRelocateMemory(t *testing.T) {
	virtualMachine := vm.NewVirtualMachine()
	segments := virtualMachine.Segments
	for i := 0; i < 4; i++ {
		segments.AddSegment()
	}
	segments.Memory.Insert(memory.NewRelocatable(0, 0), memory.NewMaybeRelocatableInt(4613515612218425347))
	segments.Memory.Insert(memory.NewRelocatable(0, 1), memory.NewMaybeRelocatableInt(5))
	segments.Memory.Insert(memory.NewRelocatable(0, 2), memory.NewMaybeRelocatableInt(2345108766317314046))
	segments.Memory.Insert(memory.NewRelocatable(1, 0), memory.NewMaybeRelocatableRelocatable(2, 0))
	segments.Memory.Insert(memory.NewRelocatable(1, 1), memory.NewMaybeRelocatableRelocatable(3, 0))
	segments.Memory.Insert(memory.NewRelocatable(1, 5), memory.NewMaybeRelocatableInt(5))

	segments.ComputeEffectiveSizes()

	relocationTable, ok := segments.RelocateSegments()
	if !ok {
		t.Errorf("Could not create relocation table")
	}

	relocatedMemory, err := segments.RelocateMemory(&relocationTable)
	if err != nil {
		t.Errorf("Test failed with error: %s", err)
	}

	expectedMemory := map[uint]uint{1: 4613515612218425347, 2: 5, 3: 2345108766317314046, 4: 10, 5: 10, 9: 5}
	for i, v := range expectedMemory {
		actual := relocatedMemory[i]
		if actual != v {
			t.Errorf("Expected relocated memory at index %d to be %d but it's %d", i, v, actual)
		}
	}
}
