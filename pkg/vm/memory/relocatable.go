package memory

import (
	"github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
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

func (relocatable *Relocatable) SubRelocatable(other uint) (Relocatable, error) {
	if relocatable.offset < other {
		return NewRelocatable(0, 0), &SubReloctableError{Msg: "RelocatableSubUsizeNegOffset"}
	} else {
		new_offset := relocatable.offset - other
		return NewRelocatable(relocatable.segmentIndex, new_offset), nil
	}
}

func (relocatable *Relocatable) AddRelocatable(other uint) (Relocatable, error) {
	new_offset := relocatable.offset + other
	return NewRelocatable(relocatable.segmentIndex, new_offset), nil

}

// Get the the indexes of the Relocatable struct.
// Returns a tuple with both values (segment_index, offset)
func (r *Relocatable) into_indexes() (uint, uint) {
	if r.segmentIndex < 0 {
		corrected_segment_idx := uint(-(r.segmentIndex + 1))
		return corrected_segment_idx, r.offset
	}

	return uint(r.segmentIndex), r.offset
}

// Int in the Cairo VM represents a value in memory that
// is not an address.
type Int struct {
	felt lambdaworks.Felt
}

// MaybeRelocatable is the type of the memory cells in the Cairo
// VM. For now, `inner` will hold any type but it should be
// instantiated only with `Relocatable`, `Int` or `nil` types.
// We should analyze better alternatives to this.
type MaybeRelocatable struct {
	inner any
}

// Creates a new MaybeRelocatable with an Int inner value
func NewMaybeRelocatableInt(felt lambdaworks.Felt) *MaybeRelocatable {
	return &MaybeRelocatable{inner: Int{felt}}
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
