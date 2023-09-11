package hint_utils

import (
	"github.com/lambdaclass/cairo-vm.go/pkg/vm"
	"github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
)

// Receives a map and builds a setup for hints tests containing ids
// Builds the IdsManager & Inserts ids into memory
// Works as follows:
// Each map entry represents an ids variable
// That identifier can represent one or more elements (aka be an elment or struct)
// Some of these elements may also be missing (inserted during the hint), and are represented as a nil pointer
// Considerations:
// All references will be FP-based, so please don't update the value of FP after calling this function,
// and make sure that the memory at fp's segment is clear from its current offset onwards
func SetupIdsForTest(ids map[string][]*memory.MaybeRelocatable, vm *vm.VirtualMachine) IdsManager {
	manager := IdsManager{}
	base_addr := vm.RunContext.Fp
	current_offset := 0
	for name, elems := range ids {
		// Create reference
		manager.References[name] = HintReference{
			Dereference: true,
			Offset1: OffsetValue{
				ValueType: Reference,
				Value:     current_offset,
			},
		}
		// Update current_offset
		current_offset += len(elems)
		// Insert ids variables (if present)
		for n, elem := range elems {
			if elem != nil {
				vm.Segments.Memory.Insert(base_addr.AddUint(uint(n)), elem)
			}
		}
		// Update base_addr
		base_addr.Offset += uint(len(elems))
	}
	return manager
}
