package memory

import (
	"errors"
	"fmt"
)

// Relocatable in the Cairo VM represents an address
// in some memory segment. When the VM finishes running,
// these values are replaced by real memory addresses,
// represented by a field element.
type Relocatable struct {
	SegmentIndex int
	Offset       uint
}

// Creates a new Relocatable struct with the specified segment index
// and offset.
func NewRelocatable(segment_idx int, offset uint) Relocatable {
	return Relocatable{segment_idx, offset}
}

// Adds a Felt value to a Relocatable
// Fails if the new offset exceeds the size of a uint
func (r *Relocatable) AddFelt(other Int) (Relocatable, error) {
	new_offset := r.Offset + other.Felt // TODO: Placeholder
	return NewRelocatable(r.SegmentIndex, new_offset), nil

}

// Performs additions if other contains a Felt value, fails otherwise
func (r *Relocatable) AddMaybeRelocatable(other MaybeRelocatable) (Relocatable, error) {
	felt, ok := other.GetInt()
	if !ok {
		return Relocatable{}, errors.New("Can't add two relocatable values")
	}
	return r.AddFelt(felt)
}

func (r *Relocatable) RelocateAddress(relocationTable *[]uint) uint {
	return (*relocationTable)[r.SegmentIndex] + r.Offset
}

// Int in the Cairo VM represents a value in memory that
// is not an address.
type Int struct {
	// FIXME: Here we should use Lambdaworks felt, just mocking
	// this for now.
	Felt uint
}

// MaybeRelocatable is the type of the memory cells in the Cairo
// VM. For now, `inner` will hold any type but it should be
// instantiated only with `Relocatable` or `Int` types.
// We should analyze better alternatives to this.
type MaybeRelocatable struct {
	inner any
}

// Creates a new MaybeRelocatable with an Int inner value
func NewMaybeRelocatableInt(felt uint) *MaybeRelocatable {
	return &MaybeRelocatable{inner: Int{felt}}
}

// Creates a new MaybeRelocatable with a Relocatable inner value
func NewMaybeRelocatableRelocatable(relocatable Relocatable) *MaybeRelocatable {
	return &MaybeRelocatable{inner: relocatable}
}

// If m is Int, returns the inner value + true, if not, returns zero + false
func (m *MaybeRelocatable) GetInt() (Int, bool) {
	int, is_type := m.inner.(Int)
	return int, is_type
}

// If m is Relocatable, returns the inner value + true, if not, returns zero + false
func (m *MaybeRelocatable) GetRelocatable() (Relocatable, bool) {
	rel, is_type := m.inner.(Relocatable)
	return rel, is_type
}

// Turns a MaybeRelocatable into a Felt252 value.
// If the inner value is an Int, it will extract the Felt252 value from it.
// If the inner value is a Relocatable, it will relocate it according to the relocation_table
// TODO: Return value should be of type (felt, error)
func (m *MaybeRelocatable) RelocateValue(relocationTable *[]uint) (uint, error) {
	inner_int, ok := m.GetInt()
	if ok {
		return inner_int.Felt, nil
	}

	inner_relocatable, ok := m.GetRelocatable()
	if ok {
		return inner_relocatable.RelocateAddress(relocationTable), nil
	}

	return 0, errors.New(fmt.Sprintf("Unexpected type %T", m.inner))
}
