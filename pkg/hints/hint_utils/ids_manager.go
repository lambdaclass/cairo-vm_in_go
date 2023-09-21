package hint_utils

import (
	"github.com/pkg/errors"

	"github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
	"github.com/lambdaclass/cairo-vm.go/pkg/parser"
	. "github.com/lambdaclass/cairo-vm.go/pkg/vm"
	. "github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
)

// Identifier Manager
// Provides methods that allow hints to interact with cairo variables given their identifier name
type IdsManager struct {
	References       map[string]HintReference
	HintApTracking   parser.ApTrackingData
	AccessibleScopes []string
}

func ErrIdsManager(err error) error {
	return errors.Wrapf(err, "IdsManager error")
}

func ErrUnknownIdentifier(name string) error {
	return ErrIdsManager(errors.Errorf("Unknown identifier %s", name))
}

func ErrIdentifierNotFelt(name string) error {
	return ErrIdsManager(errors.Errorf("Identifier %s is not a Felt", name))
}

func NewIdsManager(references map[string]HintReference, hintApTracking parser.ApTrackingData, accessibleScopes []string) IdsManager {
	return IdsManager{
		References:       references,
		HintApTracking:   hintApTracking,
		AccessibleScopes: accessibleScopes,
	}
}

// Fetches a constant used by the hint
// Searches inner modules first for name-matching constants
func (ids *IdsManager) GetConst(name string, constants *map[string]lambdaworks.Felt) (lambdaworks.Felt, error) {
	// Hints should always have accessible scopes
	if len(ids.AccessibleScopes) != 0 {
		// Accessible scopes are listed from outer to inner
		for i := len(ids.AccessibleScopes) - 1; i >= 0; i-- {
			constant, ok := (*constants)[ids.AccessibleScopes[i]+"."+name]
			if ok {
				return constant, nil
			}
		}
	}
	return lambdaworks.FeltZero(), errors.Errorf("Missing constant %s", name)
}

// Inserts value into memory given its identifier name
func (ids *IdsManager) Insert(name string, value *MaybeRelocatable, vm *VirtualMachine) error {

	addr, err := ids.GetAddr(name, vm)
	if err != nil {
		return err
	}
	return vm.Segments.Memory.Insert(addr, value)
}

// Returns the value of an identifier as a Felt
func (ids *IdsManager) GetFelt(name string, vm *VirtualMachine) (lambdaworks.Felt, error) {
	val, err := ids.Get(name, vm)
	if err != nil {
		return lambdaworks.Felt{}, err
	}
	felt, is_felt := val.GetFelt()
	if !is_felt {
		return lambdaworks.Felt{}, ErrIdentifierNotFelt(name)
	}
	return felt, nil
}

// Returns the value of an identifier as a Relocatable
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

// Returns the value of an identifier as a MaybeRelocatable
func (ids *IdsManager) Get(name string, vm *VirtualMachine) (*MaybeRelocatable, error) {
	reference, ok := ids.References[name]
	if ok {
		val, ok := getValueFromReference(&reference, ids.HintApTracking, vm)
		if ok {
			return val, nil
		}
	}
	return nil, ErrUnknownIdentifier(name)
}

// Returns the address of an identifier given its name
func (ids *IdsManager) GetAddr(name string, vm *VirtualMachine) (Relocatable, error) {
	reference, ok := ids.References[name]
	if ok {
		addr, ok := getAddressFromReference(&reference, ids.HintApTracking, vm)
		if ok {
			return addr, nil
		}
	}
	return Relocatable{}, ErrUnknownIdentifier(name)
}

/*
	 Returns the value of an ids' field (given that the identifier is a sruct)
		For example:

		struct cat {
			lives felt
			paws felt
		}

		to access each struct field, lives will be field 0 and paws will be field 1, so to access them we can use:
		ids_lives := ids.GetStructField("cat", 0, vm) or ids_lives := ids.Get("cat", vm)
		ids_paws := ids.GetStructField("cat", 1, vm)
*/
func (ids *IdsManager) GetStructField(name string, field_off uint, vm *VirtualMachine) (*MaybeRelocatable, error) {
	reference, ok := ids.References[name]
	if ok {
		val, ok := getStructFieldFromReference(&reference, field_off, ids.HintApTracking, vm)
		if ok {
			return val, nil
		}
	}
	return nil, ErrUnknownIdentifier(name)
}

/*
	 Returns the value of an ids' field (given that the identifier is a sruct) as a Felt
		For example:

		struct cat {
			lives felt
			paws felt
		}

		to access each struct field, lives will be field 0 and paws will be field 1, so to access them we can use:
		ids_lives := ids.GetStructFieldFelt("cat", 0, vm) or ids_lives := ids.Get("cat", vm)
		ids_paws := ids.GetStructFieldFelt("cat", 1, vm)
*/
func (ids *IdsManager) GetStructFieldFelt(name string, field_off uint, vm *VirtualMachine) (lambdaworks.Felt, error) {
	reference, ok := ids.References[name]
	if ok {
		val, ok := getStructFieldFromReference(&reference, field_off, ids.HintApTracking, vm)
		if ok {
			felt, is_felt := val.GetFelt()
			if !is_felt {
				return lambdaworks.Felt{}, errors.Errorf("Identifier %s is not a Felt", name)
			}
			return felt, nil
		}
	}

	return lambdaworks.Felt{}, ErrUnknownIdentifier(name)
}

/*
	 Returns the value of an ids' field (given that the identifier is a sruct) as a Relocatable
		For example:

		struct shelter {
			cats cat*
			dogs dog*
		}

		to access each struct field, cats will be field 0 and dogs will be field 1, so to access them we can use:
		ids_cats := ids.GetStructFieldFelt("shelter", 0, vm) or ids_cats := ids.Get("shelter", vm)
		ids_dogs := ids.GetStructFieldFelt("shelter", 1, vm)
*/
func (ids *IdsManager) GetStructFieldRelocatable(name string, field_off uint, vm *VirtualMachine) (Relocatable, error) {
	reference, ok := ids.References[name]
	if ok {
		val, ok := getStructFieldFromReference(&reference, field_off, ids.HintApTracking, vm)
		if ok {
			rel, is_rel := val.GetRelocatable()
			if !is_rel {
				return Relocatable{}, errors.Errorf("Identifier %s is not a Relocatable", name)
			}
			return rel, nil
		}
	}

	return Relocatable{}, ErrUnknownIdentifier(name)
}

/*
	 Inserts value into an ids' field (given that the identifier is a sruct)
		For example:

		struct cat {
			lives felt
			paws felt
		}

		to access each struct field, lives will be field 0 and paws will be field 1
		, so to set the value of cat.paws we can use:
		ids.InsertStructField("cat", 1, vm)
*/
func (ids *IdsManager) InsertStructField(name string, field_off uint, value *MaybeRelocatable, vm *VirtualMachine) error {

	addr, err := ids.GetAddr(name, vm)
	if err != nil {
		return err
	}
	return vm.Segments.Memory.Insert(addr.AddUint(field_off), value)
}

// Inserts value into the address of the given identifier
func insertIdsFromReference(value *MaybeRelocatable, reference *HintReference, apTracking parser.ApTrackingData, vm *VirtualMachine) error {
	addr, ok := getAddressFromReference(reference, apTracking, vm)
	if ok {
		return vm.Segments.Memory.Insert(addr, value)
	}
	return errors.New("Failed to get ids addr")
}

func getValueFromReference(reference *HintReference, apTracking parser.ApTrackingData, vm *VirtualMachine) (*MaybeRelocatable, bool) {
	// Handle the case of  immediate
	if reference.Offset1.ValueType == Immediate {
		return NewMaybeRelocatableFelt(reference.Offset1.Immediate), true
	}
	addr, ok := getAddressFromReference(reference, apTracking, vm)
	if ok {
		if reference.Dereference {
			val, err := vm.Segments.Memory.Get(addr)
			if err == nil {
				return val, true
			}
		} else {
			return NewMaybeRelocatableRelocatable(addr), true
		}
	}
	return nil, false
}

func getStructFieldFromReference(reference *HintReference, field_off uint, apTracking parser.ApTrackingData, vm *VirtualMachine) (*MaybeRelocatable, bool) {
	addr, ok := getAddressFromReference(reference, apTracking, vm)
	if ok {
		if reference.Dereference {
			val, err := vm.Segments.Memory.Get(addr.AddUint(field_off))
			if err == nil {
				return val, true
			}
		}
	}
	return nil, false
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
		addr, err := addr.SubUint(uint(hintApTracking.Offset - refApTracking.Offset))
		if err == nil {
			return addr, true
		}
	}
	return Relocatable{}, false
}
