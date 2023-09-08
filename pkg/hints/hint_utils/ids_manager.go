package hint_utils

import (
	"github.com/pkg/errors"

	"github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
	"github.com/lambdaclass/cairo-vm.go/pkg/parser"
	. "github.com/lambdaclass/cairo-vm.go/pkg/vm"
	. "github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
)

type IdsManager struct {
	References     map[string]HintReference
	HintApTracking parser.ApTrackingData
}

func NewIdsManager(references map[string]HintReference, hintApTracking parser.ApTrackingData) IdsManager {
	return IdsManager{
		References:     references,
		HintApTracking: hintApTracking,
	}
}

// Inserts value into an ids given its name
func (ids *IdsManager) Insert(name string, value *MaybeRelocatable, vm *VirtualMachine) error {

	addr, err := ids.GetAddr(name, vm)
	if err != nil {
		return err
	}
	return vm.Segments.Memory.Insert(addr, value)
}

// Returns the value of an ids as Felt
func (ids *IdsManager) GetFelt(name string, vm *VirtualMachine) (lambdaworks.Felt, error) {
	val, err := ids.Get(name, vm)
	if err != nil {
		return lambdaworks.Felt{}, err
	}
	felt, is_felt := val.GetFelt()
	if !is_felt {
		return lambdaworks.Felt{}, errors.Errorf("Identifier %s is not a Felt", name)
	}
	return felt, nil
}

// Returns the value of an ids as Relocatable
func (ids *IdsManager) GetRelocatable(name string, vm *VirtualMachine) (Relocatable, error) {
	val, err := ids.Get(name, vm)
	if err != nil {
		return Relocatable{}, err
	}
	relocatable, is_relocatable := val.GetRelocatable()
	if !is_relocatable {
		return Relocatable{}, errors.Errorf("Identifier %s is not a Relocatable", name)
	}
	return relocatable, nil
}

// Returns the value of an ids as MaybeRelocatable
func (ids *IdsManager) Get(name string, vm *VirtualMachine) (*MaybeRelocatable, error) {
	addr, err := ids.GetAddr(name, vm)
	if err != nil {
		return nil, err
	}
	return vm.Segments.Memory.Get(addr)
}

// Returns the address of an ids given its name
func (ids *IdsManager) GetAddr(name string, vm *VirtualMachine) (Relocatable, error) {
	reference, ok := ids.References[name]
	if ok {
		addr, ok := getAddressFromReference(&reference, ids.HintApTracking, vm)
		if ok {
			return addr, nil
		}
	}
	return Relocatable{}, errors.Errorf("Unknow identifier %s", name)
}

// Inserts value into the address of the given ids variable
func insertIdsFromReference(value *MaybeRelocatable, reference *HintReference, apTracking parser.ApTrackingData, vm *VirtualMachine) error {
	addr, ok := getAddressFromReference(reference, apTracking, vm)
	if ok {
		return vm.Segments.Memory.Insert(addr, value)
	}
	return errors.New("Failed to get ids addr")
}

// Returns the addr indicated by the reference
func getAddressFromReference(reference *HintReference, apTracking parser.ApTrackingData, vm *VirtualMachine) (Relocatable, bool) {
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
