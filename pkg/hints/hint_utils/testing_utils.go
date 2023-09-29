package hint_utils

import (
	"reflect"
	"testing"

	"github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
	"github.com/lambdaclass/cairo-vm.go/pkg/parser"
	"github.com/lambdaclass/cairo-vm.go/pkg/types"
	. "github.com/lambdaclass/cairo-vm.go/pkg/vm"
	"github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
)

// Receives a map and builds a setup for hints tests containing ids
// For usage examples refer to testing_utils_test.go
// Builds the IdsManager & Inserts ids into memory
// Works as follows:
// Each map entry represents an ids variable
// That identifier can represent one or more elements (aka be an elment or struct)
// Some of these elements may also be missing (inserted during the hint), and are represented as a nil pointer
// Considerations:
// All references will be FP-based, so please don't update the value of FP after calling this function,
// and make sure that the memory at fp's segment is clear from its current offset onwards
func SetupIdsForTest(ids map[string][]*memory.MaybeRelocatable, vm *VirtualMachine) IdsManager {
	manager := NewIdsManager(make(map[string]HintReference), parser.ApTrackingData{}, []string{})
	base_addr := vm.RunContext.Fp
	current_offset := 0
	for name, elems := range ids {
		// Create reference
		manager.References[name] = HintReference{
			Dereference: true,
			Offset1: OffsetValue{
				ValueType: Reference,
				Value:     current_offset,
				Register:  FP,
			},
		}
		// Update current_offset
		current_offset += len(elems)

		// Insert ids variables (if present)
		for n, elem := range elems {
			if elem != nil {
				vm.Segments.Memory.Insert(base_addr.AddUint(uint(n)), elem)
			}
		}
		// Update base_addr
		base_addr.Offset += uint(len(elems))
	}
	return manager
}

// Returns a constants map accoring to the new_constants map received
// Adds a path to each constant and a matching path to the hint's accessible scopes
func SetupConstantsForTest(new_constants map[string]lambdaworks.Felt, ids *IdsManager) map[string]lambdaworks.Felt {
	constants := make(map[string]lambdaworks.Felt)
	ids.AccessibleScopes = append(ids.AccessibleScopes, "path")
	for name, constant := range new_constants {
		constants["path."+name] = constant
	}
	return constants
}

func CheckScopeVar[T any](name string, expectedVal T, scopes *types.ExecutionScopes, t *testing.T) {
	val, err := types.FetchScopeVar[T](name, scopes)
	if err != nil {
		t.Error(err.Error())
	}
	if !reflect.DeepEqual(val, expectedVal) {
		t.Errorf("Wrong scope var %s.\n Expected: %v, got: %v", name, expectedVal, val)
	}
}
