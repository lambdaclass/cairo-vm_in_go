package utils

import (
	. "github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
	. "github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"

	"github.com/pkg/errors"
)

func IsSubsequence[T comparable](subsequence []T, sequence []T) bool {
	startSeqIdx := 0
	for _, subElem := range subsequence {
		found := false
		for idx, elem := range sequence[startSeqIdx:] {
			if subElem == elem {
				startSeqIdx = idx + 1
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func CheckBuiltinsSubsequence(programBuiltins []string) error {
	orderedBuiltinNames := []string{
		"output",
		"pedersen",
		"range_check",
		"ecdsa",
		"bitwise",
		"ec_op",
		"keccak",
		"poseidon",
	}
	if !IsSubsequence(programBuiltins, orderedBuiltinNames) {
		return errors.Errorf("program builtins are not in appropiate order")
	}
	return nil
}

// Creates a new MaybeRelocatable from a uint64 value
func NewMaybeRelocatableFeltFromUint64(val uint64) *MaybeRelocatable {
	return NewMaybeRelocatableFelt(FeltFromUint64(val))
}

func NewMaybeRelocatableRelocatableParams(segment_idx int, offset uint) *MaybeRelocatable {
	return NewMaybeRelocatableRelocatable(NewRelocatable(segment_idx, offset))
}

func AddNSegments(segments MemorySegmentManager, nSegments int) MemorySegmentManager {
	for i := 0; i < nSegments; i++ {
		segments.AddSegment()
	}
	return segments
}
