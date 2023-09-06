package utils

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
