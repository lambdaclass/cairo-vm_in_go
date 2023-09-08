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

// Asumes offsetValue is of type Reference
func getOffsetValueReference(offsetValue OffsetValue, refApTracking parser.ApTrackingData, hintApTracking parser.ApTrackingData, vm *VirtualMachine) {
	var base_addr Relocatable
	ok := true
	switch offsetValue.Register {
	case FP:
		base_addr = vm.RunContext.Fp
	case AP:
		base_addr, ok = applyApTrackingCorrection(vm.RunContext.Ap, refApTracking, hintApTracking)
	}
	if ok {

	}
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
