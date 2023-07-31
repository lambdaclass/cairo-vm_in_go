package runners_test

import (
	"testing"

	"github.com/lambdaclass/cairo-vm.go/pkg/vm"
	"github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
	"github.com/lambdaclass/cairo-vm.go/pkg/vm/runners"
)

func TestRelocateMemory(t *testing.T) {
	runner := runners.NewCairoRunner()
	virtualMachine := vm.NewVirtualMachine()
	segments := &virtualMachine.Segments
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

	err := runner.RelocateMemory(virtualMachine, relocationTable)
	if err != nil {
		t.Errorf("Test failed with error: %s", err)
	}

	expectedMemory := []int{-1, 4613515612218425347, 5, 2345108766317314046, 10, 10, -1, -1, -1, 5}
	for i, v := range expectedMemory {
		actual := runner.RelocatedMemory[i]
		if actual != v {
			t.Errorf("Expected relocated memory at index %d to be %d but it's %d", i, v, actual)
		}
	}
}
