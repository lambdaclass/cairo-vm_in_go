package hint_utils

import (
	"github.com/lambdaclass/cairo-vm.go/pkg/parser"
	. "github.com/lambdaclass/cairo-vm.go/pkg/vm"
	. "github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
)

func GetAddressFromReference(reference *HintReference, apTracking parser.ApTrackingData, vm *VirtualMachine) (Relocatable, bool) {
	if reference.Offset1.ValueType != Reference {
		return nil
	}

}

// Returns the value of the offset, or nil of it can't be computed
// Asumes offsetValue is of type Reference
func getOffsetValueReference(offsetValue OffsetValue, refApTracking parser.ApTrackingData, hintApTracking parser.ApTrackingData, vm *VirtualMachine) *MaybeRelocatable {
	var base_addr Relocatable
	ok := true
	switch offsetValue.Register {
	case FP:
		base_addr = vm.RunContext.Fp
	case AP:
		base_addr, ok = applyApTrackingCorrection(vm.RunContext.Ap, refApTracking, hintApTracking)
	}
	if ok {
		base_addr, err := base_addr.AddInt(offsetValue.Value)
		if err == nil {
			if offsetValue.Dereference {
				// val will be nil if err is not nil, so we can ignore it
				val, _ := vm.Segments.Memory.Get(base_addr)
				return val
			} else {
				return NewMaybeRelocatableRelocatable(base_addr)
			}
		}
	}
	return nil
}

func applyApTrackingCorrection(addr Relocatable, refApTracking parser.ApTrackingData, hintApTracking parser.ApTrackingData) (Relocatable, bool) {
	// Reference & Hint ApTracking must belong to the same group
	if refApTracking.Group == hintApTracking.Group {
		addr, err := addr.SubUint(uint(hintApTracking.Offset - hintApTracking.Offset))
		if err == nil {
			return addr, true
		}
	}
	return Relocatable{}, false
}
