package memory

import (
	"errors"
)

// A Set to store Relocatable values
type AddressSet map[Relocatable]bool

func NewAddressSet() AddressSet {
	return make(map[Relocatable]bool)
}

func (set AddressSet) Add(element Relocatable) {
	set[element] = true
}

func (set AddressSet) Contains(element Relocatable) bool {
	return set[element]
}

// A function that validates a memory address and returns a list of validated addresses
type ValidationRule func(*Memory, Relocatable) ([]Relocatable, error)

// Memory represents the Cairo VM's memory.
type Memory struct {
	data              map[Relocatable]MaybeRelocatable
	numSegments       uint
	validationRules   map[uint]ValidationRule
	validatedAdresses AddressSet
}

var ErrMissingSegmentUsize = errors.New("segment effective sizes haven't been calculated")

func NewMemory() *Memory {
	return &Memory{
		data:              make(map[Relocatable]MaybeRelocatable),
		validatedAdresses: NewAddressSet(),
		validationRules:   make(map[uint]ValidationRule),
	}
}

// Inserts a value in some memory address, given by a Relocatable value.
func (m *Memory) Insert(addr Relocatable, val *MaybeRelocatable) error {
	// FIXME: There should be a special handling if the key
	// segment index is negative. This is an edge
	// case, so for now let's raise an error.
	if addr.SegmentIndex < 0 {
		return errors.New("segment index of key is negative - unimplemented")
	}

	// Check that insertions are preformed within the memory bounds
	if addr.SegmentIndex >= int(m.NumSegments) {
		return errors.New("Error: Inserting into a non allocated segment")
	}

	// Check for possible overwrites
	prev_elem, ok := m.Data[addr]
	if ok && prev_elem != *val {
		return errors.New("Memory is write-once, cannot overwrite memory value")
	}
	m.Data[addr] = *val
	return m.ValidateAddress(addr)
}

// Gets some value stored in the memory address `addr`.
func (m *Memory) Get(addr Relocatable) (*MaybeRelocatable, error) {
	// FIXME: There should be a special handling if the key
	// segment index is negative. This is an edge
	// case, so for now let's raise an error.
	if addr.SegmentIndex < 0 {
		return nil, errors.New("segment index of key is negative - unimplemented")
	}

	// FIXME: We should create a function for this value,
	// `relocate_value()` in the future. This function should
	// check if the value is a `Relocatable` with a negative
	// segment index. Again, these are edge cases so not important
	// right now. See cairo-vm code for details.
	value, ok := m.Data[addr]

	if !ok {
		return nil, errors.New("Memory Get: Value not found")
	}

	return &value, nil
}

// Adds a validation rule for a given segment
func (m *Memory) AddValidationRule(SegmentIndex uint, rule ValidationRule) {
	m.ValidationRules[SegmentIndex] = rule
}

// Applies the validation rule for the addr's segment if any
// Skips validation if the address is temporary or if it has been previously validated
func (m *Memory) ValidateAddress(addr Relocatable) error {
	if addr.SegmentIndex < 0 || m.ValidatedAdresses.Contains(addr) {
		return nil
	}
	rule, ok := m.ValidationRules[uint(addr.SegmentIndex)]
	if !ok {
		return nil
	}
	validated_addresses, error := rule(m, addr)
	if error != nil {
		return error
	}
	for _, validated_address := range validated_addresses {
		m.ValidatedAdresses.Add(validated_address)
	}
	return nil
}

// Applies validation_rules to every memory address, if applicatble
// Skips validation if the address is temporary or if it has been previously validated
func (m *Memory) ValidateExistingMemory() error {
	for addr := range m.Data {
		err := m.ValidateAddress(addr)
		if err != nil {
			return err
		}
	}
	return nil
}
