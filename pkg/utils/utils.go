package utils

import (
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
