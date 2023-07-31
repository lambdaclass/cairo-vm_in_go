package runners

import (
	"errors"

	"github.com/lambdaclass/cairo-vm.go/pkg/vm"
	"github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
)

type CairoRunner struct {
	// TODO: relocatedMemory should be of type felt
	relocatedMemory []int
}

func NewCairoRunner() *CairoRunner {
	return &CairoRunner{relocatedMemory: make([]int, 0)}
}

func (c *CairoRunner) RelocatedMemory() *[]int {
	return &c.relocatedMemory
}

func (c *CairoRunner) RelocateMemory(vm *vm.VirtualMachine, relocationTable []uint) error {
	if len(c.relocatedMemory) != 0 {
		return errors.New("Inconsistent relocation")
	}

	// Relocated addresses start at 1
	// TODO: with felts, we should use nil instead of -1
	c.relocatedMemory = append(c.relocatedMemory, -1)
	segments := vm.Segments()

	for i := uint(0); i < segments.Memory.NumSegments(); i++ {
		for j := uint(0); j < segments.SegmentSizes[i]; j++ {
			ptr := memory.NewRelocatable(int(i), j)
			cell, err := segments.Memory.Get(ptr)
			if err == nil {
				relocatedAddr := ptr.RelocateAddress(relocationTable)
				value, err := cell.RelocateValue(relocationTable)
				if err != nil {
					return err
				}
				for len(c.relocatedMemory) <= int(relocatedAddr) {
					c.relocatedMemory = append(c.relocatedMemory, -1)
				}
				c.relocatedMemory[relocatedAddr] = value
			} else {
				c.relocatedMemory = append(c.relocatedMemory, -1)
			}
		}
	}

	return nil
}
