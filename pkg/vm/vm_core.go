package vm

import (
	"errors"

	"github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
)

// VirtualMachine represents the Cairo VM.
// Runs Cairo assembly and produces an execution trace.
type VirtualMachine struct {
	RunContext     RunContext
	currentStep    uint
	Segments       memory.MemorySegmentManager
	traceRelocated bool
	Trace          []TraceEntry
}

func NewVirtualMachine() *VirtualMachine {
	segments := memory.NewMemorySegmentManager()
	trace := make([]TraceEntry, 0)
	return &VirtualMachine{Segments: *segments, Trace: trace}
}

func (v *VirtualMachine) RelocateTrace(relocationTable *[]uint) error {
	if len(*relocationTable) < 2 {
		return errors.New("No relocation found for execution segment")
	}
	segment1Base := (*relocationTable)[1]

	for i := range v.Trace {
		v.Trace[i].Pc++
		v.Trace[i].Ap += segment1Base
		v.Trace[i].Fp += segment1Base
	}
	v.traceRelocated = true

	return nil
}

func (v *VirtualMachine) GetRelocatedTrace() (*[]TraceEntry, error) {
	if v.traceRelocated {
		return &v.Trace, nil
	} else {
		return nil, errors.New("Trace not relocated")
	}
}

func (v *VirtualMachine) Relocate() error {
	v.Segments.ComputeEffectiveSizes()
	if len(v.Trace) == 0 {
		return nil
	}

	relocationTable, ok := v.Segments.RelocateSegments()
	// This should be unreachable
	if !ok {
		return errors.New("ComputeEffectiveSizes called but RelocateSegments still returned error")
	}

	_, err := v.Segments.RelocateMemory(&relocationTable)
	if err != nil {
		return err
	}

	v.RelocateTrace(&relocationTable)
	return nil
}
