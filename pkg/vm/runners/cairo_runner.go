package runners

import (
	"errors"

	"github.com/lambdaclass/cairo-vm.go/pkg/vm"
	"github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
)

type CairoRunner struct {
	// TODO: relocatedMemory should be of type felt
	RelocatedMemory []int
}

func NewCairoRunner() *CairoRunner {
	return &CairoRunner{RelocatedMemory: make([]int, 0)}
}

func (c *CairoRunner) RelocateMemory(vm *vm.VirtualMachine, relocationTable *[]uint) error {
	if len(c.RelocatedMemory) != 0 {
		return errors.New("Inconsistent relocation")
	}

	// Relocated addresses start at 1
	// TODO: with felts, we should use nil instead of -1
	c.RelocatedMemory = append(c.RelocatedMemory, -1)

	for i := uint(0); i < vm.Segments.Memory.NumSegments(); i++ {
		for j := uint(0); j < vm.Segments.SegmentSizes[i]; j++ {
			ptr := memory.NewRelocatable(int(i), j)
			cell, err := vm.Segments.Memory.Get(ptr)
			if err == nil {
				relocatedAddr := ptr.RelocateAddress(relocationTable)
				value, err := cell.RelocateValue(relocationTable)
				if err != nil {
					return err
				}
				for len(c.RelocatedMemory) <= int(relocatedAddr) {
					c.RelocatedMemory = append(c.RelocatedMemory, -1)
				}
				c.RelocatedMemory[relocatedAddr] = value
			} else {
				c.RelocatedMemory = append(c.RelocatedMemory, -1)
			}
		}
	}

	return nil
}
