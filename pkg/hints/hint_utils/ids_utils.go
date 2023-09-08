package hint_utils

import (
	"github.com/lambdaclass/cairo-vm.go/pkg/parser"
	. "github.com/lambdaclass/cairo-vm.go/pkg/vm"
	. "github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
)

// Returns the addr indicated by the reference
func GetAddressFromReference(reference *HintReference, apTracking parser.ApTrackingData, vm *VirtualMachine) (Relocatable, bool) {
	if reference.Offset1.ValueType != Reference {
		return Relocatable{}, false
	}
	offset1 := getOffsetValueReference(reference.Offset1, reference.ApTrackingData, apTracking, vm)
	if offset1 != nil {
		offset1_rel, is_rel := offset1.GetRelocatable()
		if is_rel {
			switch reference.Offset2.ValueType {
			case Reference:
				offset2 := getOffsetValueReference(reference.Offset2, reference.ApTrackingData, apTracking, vm)
				if offset2 != nil {
					res, err := offset1_rel.AddMaybeRelocatable(*offset2)
					if err == nil {
						return res, true
					}
				}
			case Value:
				res, err := offset1_rel.AddInt(reference.Offset2.Value)
				if err == nil {
					return res, true
				}
			}
		}
	}
	return Relocatable{}, false

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