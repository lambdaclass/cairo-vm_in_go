package memory

import (
	"github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
)

// Relocatable in the Cairo VM represents an address
// in some memory segment. When the VM finishes running,
// these values are replaced by real memory addresses,
// represented by a field element.
type Relocatable struct {
	SegmentIndex int
	Offset       uint
}

func (r *Relocatable) IsEqual(r1 *Relocatable) bool {
	return (r.SegmentIndex == r1.SegmentIndex && r.Offset == r1.Offset)
}

// Creates a new Relocatable struct with the specified segment index
// and offset.
func NewRelocatable(segment_idx int, offset uint) Relocatable {
	return Relocatable{segment_idx, offset}
}

func (relocatable *Relocatable) SubRelocatable(other uint) (Relocatable, error) {
	if relocatable.Offset < other {
		return NewRelocatable(0, 0), &SubReloctableError{Msg: "RelocatableSubUsizeNegOffset"}
	} else {
		new_offset := relocatable.Offset - other
		return NewRelocatable(relocatable.SegmentIndex, new_offset), nil
	}
}

func (relocatable *Relocatable) AddRelocatable(other uint) (Relocatable, error) {
	new_offset := relocatable.Offset + other
	return NewRelocatable(relocatable.SegmentIndex, new_offset), nil

}

// Get the the indexes of the Relocatable struct.
// Returns a tuple with both values (segment_index, offset)
func (r *Relocatable) into_indexes() (uint, uint) {
	if r.SegmentIndex < 0 {
		corrected_segment_idx := uint(-(r.SegmentIndex + 1))
		return corrected_segment_idx, r.Offset
	}

	return uint(r.SegmentIndex), r.Offset
}

// Int in the Cairo VM represents a value in memory that
// is not an address.
type Int struct {
	felt lambdaworks.Felt
}

// MaybeRelocatable is the type of the memory cells in the Cairo
// VM. For now, `inner` will hold any type but it should be
// instantiated only with `Relocatable` or `Int` types.
// We should analyze better alternatives to this.
type MaybeRelocatable struct {
	inner any
}

// Creates a new MaybeRelocatable with an Int inner value
func NewMaybeRelocatableInt(felt lambdaworks.Felt) *MaybeRelocatable {
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

// If m and m1 are equal, returns true, otherwise returns false
func (m *MaybeRelocatable) IsEqual(m1 *MaybeRelocatable) bool {
	a, a_type := m.GetInt()
	b, b_type := m1.GetInt()
	if a_type == b_type {
		if a_type {
			return a == b
		} else {
			a, _ := m.GetRelocatable()
			b, _ := m1.GetRelocatable()
			return a.IsEqual(&b)
		}
	} else {
		return false
	}
}
