package utils_test

import (
	"testing"

	"github.com/lambdaclass/cairo-vm.go/pkg/utils"
)

func TestIsSubsequenceWorks(t *testing.T) {
	orderedBuiltinNames := []string{
		"output_builtin",
		"pedersen_builtin",
		"range_check_builtin",
		"ecdsa_builtin",
		"bitwise_builtin",
		"ec_op_builtin",
		"keccak_builtin",
		"poseidon_builtin",
	}

	subSequence := []string{
		"output_builtin",
		"poseidon_builtin",
	}

	result := utils.IsSubsequence(subSequence, orderedBuiltinNames)

	if !result {
		t.Errorf("The result of IsSubsequence should be true")
	}
}

func TestIsSubsequenceReturnsFalse(t *testing.T) {
	orderedBuiltinNames := []string{
		"output_builtin",
		"pedersen_builtin",
		"range_check_builtin",
		"ecdsa_builtin",
		"bitwise_builtin",
		"ec_op_builtin",
		"keccak_builtin",
		"poseidon_builtin",
	}

	subSequence := []string{
		"poseidon_builtin",
		"output_builtin",
	}

	result := utils.IsSubsequence(subSequence, orderedBuiltinNames)

	if result {
		t.Errorf("The result of IsSubsequence should be false")
	}
}
