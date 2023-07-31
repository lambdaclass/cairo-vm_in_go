package memory_test

import (
	"reflect"
	"testing"

	"github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
)

func TestRelocateOneSegment(t *testing.T) {
	segments := memory.NewMemorySegmentManager()
	segments.AddSegment()
	segments.SegmentSizes = map[uint]uint{0: 3}
	relocation_table, ok := segments.RelocateSegments()

	if !ok {
		t.Errorf("Memory segment manager doesn't have segment sizes initialized")
	}

	expected_table := []uint{1}
	if !reflect.DeepEqual(expected_table, relocation_table) {
		t.Errorf("Relocation tables are not the same")
	}
}

func TestRelocateFiveSegments(t *testing.T) {
	segments := memory.NewMemorySegmentManager()
	segments.AddSegment()
	segments.SegmentSizes = map[uint]uint{0: 3, 1: 3, 2: 56, 3: 78, 4: 8}
	relocation_table, ok := segments.RelocateSegments()

	if !ok {
		t.Errorf("Memory segment manager doesn't have segment sizes initialized")
	}

	expected_table := []uint{1, 4, 7, 63, 141}
	if !reflect.DeepEqual(expected_table, relocation_table) {
		t.Errorf("Relocation tables are not the same")
	}
}

func TestRelocateSegmentsWithHoles(t *testing.T) {
	segments := memory.NewMemorySegmentManager()
	segments.AddSegment()
	segments.SegmentSizes = map[uint]uint{0: 3, 2: 3}
	relocation_table, ok := segments.RelocateSegments()

	if !ok {
		t.Errorf("Memory segment manager doesn't have segment sizes initialized")
	}

	expected_table := []uint{1, 4, 4}
	if !reflect.DeepEqual(expected_table, relocation_table) {
		t.Errorf("Relocation tables are not the same")
	}
}
