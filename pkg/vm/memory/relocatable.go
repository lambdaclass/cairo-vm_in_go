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
	segmentIndex int
	offset       uint
}

// Creates a new Relocatable struct with the specified segment index
// and offset.
func NewRelocatable(segment_idx int, offset uint) Relocatable {
	return Relocatable{segment_idx, offset}
}

func (r *Relocatable) RelocateAddress(relocationTable *[]uint) uint {
	return (*relocationTable)[r.segmentIndex] + r.offset
}

// Int in the Cairo VM represents a value in memory that
// is not an address.
type Int struct {
	// FIXME: Here we should use Lambdaworks felt, just mocking
	// this for now.
	felt uint
}

// MaybeRelocatable is the type of the memory cells in the Cairo
// VM. For now, `inner` will hold any type but it should be
// instantiated only with `Relocatable`, `Int` or `nil` types.
// We should analyze better alternatives to this.
type MaybeRelocatable struct {
	inner any
}

// Creates a new MaybeRelocatable with an Int inner value
func NewMaybeRelocatableInt(felt uint) *MaybeRelocatable {
	return &MaybeRelocatable{inner: Int{felt}}
}

// Creates a new MaybeRelocatable with a Relocatable inner value
func NewMaybeRelocatableRelocatable(segmentIndex int, offset uint) *MaybeRelocatable {
	return &MaybeRelocatable{inner: Relocatable{segmentIndex: segmentIndex, offset: offset}}
}

// TODO: Return value should be of type (felt, error)
func (m *MaybeRelocatable) RelocateValue(relocationTable *[]uint) (int, error) {
	inner_int, ok := m.GetInt()
	if ok {
		return int(inner_int.felt), nil
	}

	inner_relocatable, ok := m.GetRelocatable()
	if ok {
		return int(inner_relocatable.RelocateAddress(relocationTable)), nil
	}

	return -1, errors.New(fmt.Sprintf("Unexpected type %T", m.inner))
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
