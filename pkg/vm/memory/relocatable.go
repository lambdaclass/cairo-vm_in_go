package memory

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
func NewRelocatable(segment_idx int, offset uint) *Relocatable {
	return &Relocatable{segment_idx, offset}
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

// Creates a new MaybeRelocatable with a `nil` inner value
func NewMaybeRelocatableNil() *MaybeRelocatable {
	return &MaybeRelocatable{inner: nil}
}

// Checks if inner value of MaybeRelocatable is `nil`
func (m *MaybeRelocatable) is_nil() bool {
	return m.inner == nil
}
