package memory

import (
	"fmt"

	"github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
	"github.com/pkg/errors"
)

// A Set to store Relocatable values
type AddressSet map[Relocatable]bool

func MemoryError(err error) error {
	return errors.Wrapf(err, "Memory error")
}

func ErrMemoryWriteOnce(addr Relocatable, prevVal MaybeRelocatable, newVal MaybeRelocatable) error {
	addrStr := addr.ToString()
	prevValStr := prevVal.ToString()
	newValStr := newVal.ToString()

	return MemoryError(errors.Errorf("Memory is write-once, cannot overwrite memory value in %s. %s != %s", addrStr, prevValStr, newValStr))
}

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
	Data              map[Relocatable]MaybeRelocatable
	numSegments       uint
	validationRules   map[uint]ValidationRule
	validatedAdresses AddressSet
	// This is a map of addresses that were accessed during execution
	// The map is of the form `segmentIndex` -> `offset`. This is to
	// make the counting of memory holes easier
	AccessedAddresses map[Relocatable]bool
}

var ErrMissingSegmentUsize = errors.New("Segment effective sizes haven't been calculated")
var ErrInsufficientAllocatedCells = errors.New("Insufficient Allocated Memory Cells")

func InsufficientAllocatedCellsErrorWithBuiltinName(name string, used uint, size uint) error {
	return fmt.Errorf("%w, builtin: %s, used: %d, size: %d", ErrInsufficientAllocatedCells, name, used, size)
}

func InsufficientAllocatedCellsError(used uint, size uint) error {
	return fmt.Errorf("%w, used: %d, size: %d", ErrInsufficientAllocatedCells, used, size)
}

func InsufficientAllocatedCellsErrorMinStepNotReached(minStep uint, builtinName string) error {
	return fmt.Errorf("%w, Min Step not reached. minStep: %d, builtin: %s", ErrInsufficientAllocatedCells, minStep, builtinName)
}

func NewMemory() *Memory {
	return &Memory{
		Data:              make(map[Relocatable]MaybeRelocatable),
		validatedAdresses: NewAddressSet(),
		validationRules:   make(map[uint]ValidationRule),
		AccessedAddresses: make(map[Relocatable]bool),
	}
}

func (m *Memory) NumSegments() uint {
	return m.numSegments
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
	if addr.SegmentIndex >= int(m.numSegments) {
		return errors.Errorf("Error: Inserting into a non allocated segment %s", addr.ToString())
	}

	// Check for possible overwrites
	prev_elem, ok := m.Data[addr]
	if ok && prev_elem != *val {
		return ErrMemoryWriteOnce(addr, prev_elem, *val)
	}
	m.Data[addr] = *val
	return m.validateAddress(addr)
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

func (memory *Memory) GetSegment(segmentIndex int) []MaybeRelocatable {
	var ret []MaybeRelocatable

	for address, value := range memory.Data {
		if address.SegmentIndex == segmentIndex {
			ret = append(ret, value)
		}
	}

	return ret
}

// Gets the felt value stored in the memory address `addr`.
// Fails if the value doesn't exist or is not a felt
func (m *Memory) GetFelt(addr Relocatable) (lambdaworks.Felt, error) {
	elem, err := m.Get(addr)
	if err == nil {
		felt, ok := elem.GetFelt()
		if ok {
			return felt, nil
		} else {
			return lambdaworks.FeltZero(), errors.New("Memory GetFelt: Fetched value is not a Felt")
		}
	}
	return lambdaworks.FeltZero(), err
}

// Adds a validation rule for a given segment
func (m *Memory) AddValidationRule(SegmentIndex uint, rule ValidationRule) {
	m.validationRules[SegmentIndex] = rule
}

// Applies the validation rule for the addr's segment if any
// Skips validation if the address is temporary or if it has been previously validated
func (m *Memory) validateAddress(addr Relocatable) error {
	if addr.SegmentIndex < 0 || m.validatedAdresses.Contains(addr) {
		return nil
	}
	rule, ok := m.validationRules[uint(addr.SegmentIndex)]
	if !ok {
		return nil
	}
	validated_addresses, err := rule(m, addr)
	if err != nil {
		return err
	}
	for _, validated_address := range validated_addresses {
		m.validatedAdresses.Add(validated_address)
	}
	return nil
}

func (m *Memory) MarkAsAccessed(address Relocatable) {
	m.AccessedAddresses[address] = true
}

// Applies validation_rules to every memory address, if applicatble
// Skips validation if the address is temporary or if it has been previously validated
func (m *Memory) ValidateExistingMemory() error {
	for addr := range m.Data {
		err := m.validateAddress(addr)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *Memory) GetRelocatable(key Relocatable) (Relocatable, error) {
	memoryValue, err := m.Get(key)
	if err != nil {
		return Relocatable{}, err
	}

	ret, isRelocatable := memoryValue.GetRelocatable()
	if !isRelocatable {
		return Relocatable{}, errors.Errorf("Expected Relocatable value in memory at address (%d, %d)", key.SegmentIndex, key.Offset)
	}

	return ret, nil
}
