package memory

import "github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"

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

// Int in the Cairo VM represents a value in memory that
// is not an address.
type Int struct {
	Felt lambdaworks.Felt
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
