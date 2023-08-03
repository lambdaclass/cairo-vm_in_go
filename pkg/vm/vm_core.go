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
	Trace          []TraceEntry
	RelocatedTrace []RelocatedTraceEntry
}

func NewVirtualMachine() *VirtualMachine {
	segments := memory.NewMemorySegmentManager()
	trace := make([]TraceEntry, 0)
	relocatedTrace := make([]RelocatedTraceEntry, 0)
	return &VirtualMachine{Segments: *segments, Trace: trace, RelocatedTrace: relocatedTrace}
}

func (v *VirtualMachine) RelocateTrace(relocationTable *[]uint) error {
	if len(*relocationTable) < 2 {
		return errors.New("No relocation found for execution segment")
	}
	segment1Base := (*relocationTable)[1]

	for _, entry := range v.Trace {
		v.RelocatedTrace = append(v.RelocatedTrace, RelocatedTraceEntry{
			Pc: entry.Pc + 1,
			Ap: entry.Ap.RelocateAddress(relocationTable) + segment1Base,
			Fp: entry.Fp.RelocateAddress(relocationTable) + segment1Base,
		})
	}

	return nil
}

func (v *VirtualMachine) GetRelocatedTrace() (*[]RelocatedTraceEntry, error) {
	if len(v.RelocatedTrace) > 0 {
		return &v.RelocatedTrace, nil
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
