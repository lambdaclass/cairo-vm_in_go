package memory_test

import (
	"reflect"
	"testing"

	"github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
	"github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
)

func TestMemoryInsert(t *testing.T) {
	// Instantiate memory with 3 empty segments
	data := make([][]memory.MaybeRelocatable, 3)
	mem := memory.NewMemory(data)

	// Instantiate the address where we want to insert and the value.
	// We will insert the value Int(5) in segment 1, offset 0
	key := memory.NewRelocatable(1, 0)
	val := memory.NewMaybeRelocatableInt(lambdaworks.From(5))

	// Make the insertion
	err := mem.Insert(key, val)
	if err != nil {
		t.Errorf("Insert error in test: %s", err)
	}

	// Get the value from the address back
	res_val, err := mem.Get(key)
	if err != nil {
		t.Errorf("Get error in test: %s", err)
	}

	// Check that the original and the retrieved values are the same
	if !reflect.DeepEqual(res_val, val) {
		t.Errorf("Inserted value and original value are not the same")
	}
}

func TestMemoryInsertWithHoles(t *testing.T) {
	// Instantiate memory with 3 empty segments
	data := make([][]memory.MaybeRelocatable, 3)
	mem := memory.NewMemory(data)

	// Instantiate the address where we want to insert and the value.
	// We will insert the MaybeRelocatable Int(7) in segment 1, offset 2
	key := memory.NewRelocatable(1, 2)
	val := memory.NewMaybeRelocatableInt(lambdaworks.From(7))

	// Make the insertion
	err := mem.Insert(key, val)
	if err != nil {
		t.Errorf("Insert error in test: %s", err)
	}

	// Get the value from the address back
	res_val, err := mem.Get(key)
	if err != nil {
		t.Errorf("Get error in test: %s", err)
	}

	// Check that the original and the retrieved values are the same
	if !reflect.DeepEqual(res_val, val) {
		t.Errorf("Inserted value and original value are not the same")
	}

	// Since we inserted in segment 1, offset 2 in an empty memory, now
	// the values in segment 1, offset 0 and 1 should be `nil` (memory holes)
	hole1_addr := memory.NewRelocatable(1, 0)
	hole2_addr := memory.NewRelocatable(1, 1)

	hole1, err := mem.Get(hole1_addr)
	if err != nil {
		t.Errorf("Get error in test: %s", err)
	}

	hole2, err := mem.Get(hole2_addr)
	if err != nil {
		t.Errorf("Get error in test: %s", err)
	}

	// Check that we got the holes from memory
	expected_hole := memory.NewMaybeRelocatableNil()
	if !reflect.DeepEqual(hole1, expected_hole) || !reflect.DeepEqual(hole2, expected_hole) {
		t.Errorf("Expected nil value but got another")
	}
}
