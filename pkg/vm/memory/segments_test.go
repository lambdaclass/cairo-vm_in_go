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
