package memory

import (
	"errors"
	"fmt"

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

// Creates a new Relocatable struct with the specified segment index
// and offset.
func NewRelocatable(segment_idx int, offset uint) Relocatable {
	return Relocatable{segment_idx, offset}

}

// Returns the index in a relocated memory (uint) from a relocatable address
func (r *Relocatable) RelocateAddress(relocationTable *[]uint) uint {
	return (*relocationTable)[r.SegmentIndex] + r.Offset
}

// Adds a Felt value to a Relocatable
// Fails if the new offset exceeds the size of a uint
func (r *Relocatable) AddFelt(other lambdaworks.Felt) (Relocatable, error) {
	new_offset_felt := lambdaworks.FeltFromUint64(uint64(r.Offset)).Add(other)
	new_offset, err := new_offset_felt.ToU64()
	if err != nil {
		return *r, err
	}
	return NewRelocatable(r.SegmentIndex, uint(new_offset)), nil
}

// Subtracts a Felt value from a Relocatable
// Performs the initial subtraction considering the offset as a Felt
// Fails if the new offset exceeds the size of a uint
func (r *Relocatable) SubFelt(other lambdaworks.Felt) (Relocatable, error) {
	new_offset_felt := lambdaworks.FeltFromUint64(uint64(r.Offset)).Sub(other)
	new_offset, err := new_offset_felt.ToU64()
	if err != nil {
		return *r, err
	}
	return NewRelocatable(r.SegmentIndex, uint(new_offset)), nil
}

// Performs additions if other contains a Felt value, fails otherwise
func (r *Relocatable) AddMaybeRelocatable(other MaybeRelocatable) (Relocatable, error) {
	felt, ok := other.GetFelt()
	if !ok {
		return Relocatable{}, errors.New("can't add two relocatable values")
	}
	return r.AddFelt(felt)
}

// Returns the distance between two relocatable values (aka the difference between their offsets)
// Fails if they have different segment indexes
func (r *Relocatable) Sub(other Relocatable) (lambdaworks.Felt, error) {
	if r.SegmentIndex != other.SegmentIndex {
		return lambdaworks.Felt{}, errors.New("can't subtract two relocatables with different segment indexes")
	}
	return lambdaworks.FeltFromUint64(uint64(r.Offset)).Sub(lambdaworks.FeltFromUint64(uint64(other.Offset))), nil
}

func (r *Relocatable) IsEqual(r1 *Relocatable) bool {
	return (r.SegmentIndex == r1.SegmentIndex && r.Offset == r1.Offset)
}

func (relocatable *Relocatable) SubUint(other uint) (Relocatable, error) {
	if relocatable.Offset < other {
		return NewRelocatable(0, 0), &SubRelocatableError{Msg: "RelocatableSubUsizeNegOffset"}
	} else {
		new_offset := relocatable.Offset - other
		return NewRelocatable(relocatable.SegmentIndex, new_offset), nil
	}
}

func (relocatable *Relocatable) AddUint(other uint) Relocatable {
	new_offset := relocatable.Offset + other
	return NewRelocatable(relocatable.SegmentIndex, new_offset)
}

func (relocatable *Relocatable) AddInt(other int) (Relocatable, error) {
	if other > 0 {
		return relocatable.AddUint(uint(other)), nil
	}
	return relocatable.SubUint(uint(-other))
}

// MaybeRelocatable is the type of the memory cells in the Cairo
// VM. For now, `inner` will hold any type but it should be
// instantiated only with `Relocatable` or `Int` types.
// We should analyze better alternatives to this.
type MaybeRelocatable struct {
	inner any
}

// Creates a new MaybeRelocatable with an Int inner value
func NewMaybeRelocatableFelt(felt lambdaworks.Felt) *MaybeRelocatable {
	return &MaybeRelocatable{inner: felt}
}

// Creates a new MaybeRelocatable with a Relocatable inner value
func NewMaybeRelocatableRelocatable(relocatable Relocatable) *MaybeRelocatable {
	return &MaybeRelocatable{inner: relocatable}
}

// If m is Felt, returns the inner value + true, if not, returns zero + false
func (m *MaybeRelocatable) GetFelt() (lambdaworks.Felt, bool) {
	felt, is_type := m.inner.(lambdaworks.Felt)
	return felt, is_type
}

// If m is Relocatable, returns the inner value + true, if not, returns zero + false
func (m *MaybeRelocatable) GetRelocatable() (Relocatable, bool) {
	rel, is_type := m.inner.(Relocatable)
	return rel, is_type
}

func (m *MaybeRelocatable) IsZero() bool {
	felt, is_int := m.GetFelt()
	return is_int && felt.IsZero()
}

// Turns a MaybeRelocatable into a Felt252 value.
// If the inner value is an Int, it will extract the Felt252 value from it.
// If the inner value is a Relocatable, it will relocate it according to the relocation_table
func (m *MaybeRelocatable) RelocateValue(relocationTable *[]uint) (lambdaworks.Felt, error) {
	inner_felt, ok := m.GetFelt()
	if ok {
		return inner_felt, nil
	}

	inner_relocatable, ok := m.GetRelocatable()
	if ok {
		return lambdaworks.FeltFromUint64(uint64(inner_relocatable.RelocateAddress(relocationTable))), nil
	}

	return lambdaworks.FeltZero(), fmt.Errorf("unexpected type %T", m.inner)
}

func (m *MaybeRelocatable) IsEqual(m1 *MaybeRelocatable) bool {
	a, a_type := m.GetFelt()
	b, b_type := m1.GetFelt()
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

func (m MaybeRelocatable) Add(other MaybeRelocatable) (MaybeRelocatable, error) {
	// check if they are felt
	m_int, m_is_int := m.GetFelt()
	other_int, other_is_int := other.GetFelt()

	if m_is_int && other_is_int {
		result := NewMaybeRelocatableFelt(m_int.Add(other_int))
		return *result, nil
	}

	// check if one is relocatable and the other int
	m_rel, is_rel_m := m.GetRelocatable()
	other_rel, is_rel_other := other.GetRelocatable()

	if is_rel_m && !is_rel_other {
		other_felt, _ := other.GetFelt()
		relocatable, err := m_rel.AddFelt(other_felt)
		if err != nil {
			return *NewMaybeRelocatableFelt(lambdaworks.FeltZero()), err
		}
		return *NewMaybeRelocatableRelocatable(relocatable), nil

	} else if !is_rel_m && is_rel_other {

		m_felt, _ := m.GetFelt()
		relocatable, err := other_rel.AddFelt(m_felt)
		if err != nil {
			return *NewMaybeRelocatableFelt(lambdaworks.FeltZero()), err
		}
		return *NewMaybeRelocatableRelocatable(relocatable), nil
	} else {
		return *NewMaybeRelocatableFelt(lambdaworks.FeltZero()), errors.New("RelocatableAdd")
	}
}

// Subtracts two MaybeRelocatable values
// Behaves as follows:
// Felt - Felt : Performs felt subtraction
// Relocatable - Felt : Subtracts the Felt value from the Relocatable's offset, fails if this subtraction results in a value bigger than uint
// Relocatable - Relocatable : Returns the difference between the two offsets, fails if the difference is negative or if the segment indexes of the two relocatables don't match
// Felt - Relocatable : Always fails, this is not supported
func (m MaybeRelocatable) Sub(other MaybeRelocatable) (MaybeRelocatable, error) {
	// check if they are felt
	m_int, m_is_int := m.GetFelt()
	other_felt, other_is_felt := other.GetFelt()

	if m_is_int && other_is_felt {
		result := NewMaybeRelocatableFelt(m_int.Sub(other_felt))
		return *result, nil
	}

	// check if one is relocatable and the other int
	m_rel, is_rel_m := m.GetRelocatable()
	other_rel, is_rel_other := other.GetRelocatable()

	if is_rel_m && !is_rel_other {
		relocatable, err := m_rel.SubFelt(other_felt)
		if err != nil {
			return *NewMaybeRelocatableFelt(lambdaworks.FeltZero()), err
		}
		return *NewMaybeRelocatableRelocatable(relocatable), nil

	} else if is_rel_m && is_rel_other {
		res, err := m_rel.Sub(other_rel)
		return *NewMaybeRelocatableFelt(res), err

	} else {
		return *NewMaybeRelocatableFelt(lambdaworks.FeltZero()), errors.New("can't sub Relocatable from Felt")
	}
}

func (m *MaybeRelocatable) ToString() string {
	felt, is_felt := m.GetFelt()
	if !is_felt {
		rel, _ := m.GetRelocatable()
		return rel.ToString()
	}
	return felt.ToSignedFeltString()
}

func (r *Relocatable) ToString() string {
	return fmt.Sprintf("{%d:%d}", r.SegmentIndex, r.Offset)
}
