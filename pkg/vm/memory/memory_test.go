package memory_test

import (
	"reflect"
	"testing"

	"github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
	"github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
)

func TestMemoryInsert(t *testing.T) {
	mem_manager := memory.NewMemorySegmentManager()
	mem_manager.AddSegment()
	mem_manager.AddSegment()
	mem := &mem_manager.Memory

	// Instantiate the address where we want to insert and the value.
	// We will insert the value Felt(5) in segment 1, offset 0
	key := memory.NewRelocatable(1, 0)
	val := memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(5))

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
	mem_manager := memory.NewMemorySegmentManager()
	mem_manager.AddSegment()
	mem_manager.AddSegment()
	mem := &mem_manager.Memory

	// Instantiate the address where we want to insert and the value.
	// We will insert the MaybeRelocatable Felt(7) in segment 1, offset 2
	key := memory.NewRelocatable(1, 2)
	val := memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(5))

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

func TestMemoryInsertOverWriteSameValue(t *testing.T) {
	mem_manager := memory.NewMemorySegmentManager()
	mem := &mem_manager.Memory

	// We will insert the MaybeRelocatable Felt(7) in segment 0, offset 0
	key := mem_manager.AddSegment()
	val := memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(7))

	// Make the insertion
	err := mem.Insert(key, val)
	if err != nil {
		t.Errorf("Insert error in test: %s", err)
	}

	// Insert the same value again and check it doesn't fail
	err2 := mem.Insert(key, val)
	if err2 != nil {
		t.Errorf("Insert error in test: %s", err)
	}
}

func TestMemoryInsertOverWriteValue(t *testing.T) {
	mem_manager := memory.NewMemorySegmentManager()
	mem := &mem_manager.Memory

	// We will insert the MaybeRelocatable Felt(7) in segment 0, offset 0
	key := mem_manager.AddSegment()
	val := memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(7))

	// Make the insertion
	err := mem.Insert(key, val)
	if err != nil {
		t.Errorf("Insert error in test: %s", err)
	}

	// Insert another value into the same address and check that it fails
	val2 := memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(8))
	err2 := mem.Insert(key, val2)
	if err2 == nil {
		t.Errorf("Overwritting memory value should fail")
	}
}

func TestMemoryInsertUnallocatedSegment(t *testing.T) {
	mem_manager := memory.NewMemorySegmentManager()
	mem := &mem_manager.Memory

	// Instantiate the address where we want to insert and the value.
	// We will insert the value Felt(5) in segment 1, offset 0
	key := memory.NewRelocatable(1, 0)
	val := memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(5))

	// Make the insertion
	err := mem.Insert(key, val)
	if err == nil {
		t.Errorf("Insertion on unallocated segment should fail")
	}
}

func TestMemorySegmentsLoadDataUnallocatedSegment(t *testing.T) {
	mem_manager := memory.NewMemorySegmentManager()

	ptr := memory.NewRelocatable(1, 0)
	data := []memory.MaybeRelocatable{*memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(5))}

	// Load Data
	_, err := mem_manager.LoadData(ptr, &data)
	if err == nil {
		t.Errorf("Insertion on unallocated segment should fail")
	}
}

func TestMemorySegmentsLoadDataOneElement(t *testing.T) {
	mem_manager := memory.NewMemorySegmentManager()
	mem_manager.AddSegment()

	ptr := memory.NewRelocatable(0, 0)
	val := memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(5))
	data := []memory.MaybeRelocatable{*val}

	// Load Data
	end_ptr, err := mem_manager.LoadData(ptr, &data)
	if err != nil {
		t.Errorf("LoadData error in test: %s", err)
	}

	// Check returned ptr
	expected_end_ptr := memory.NewRelocatable(0, 1)
	if !reflect.DeepEqual(end_ptr, expected_end_ptr) {
		t.Errorf("LoadData returned wrong ptr")
	}

	// Check inserted value
	res_val, err := mem_manager.Memory.Get(ptr)
	if err != nil {
		t.Errorf("Get error in test: %s", err)
	}

	// Check that the original and the retrieved values are the same
	if !reflect.DeepEqual(res_val, val) {
		t.Errorf("Inserted value and original value are not the same")
	}
}

func TestMemorySegmentsLoadDataTwoElements(t *testing.T) {
	mem_manager := memory.NewMemorySegmentManager()
	mem_manager.AddSegment()

	ptr := memory.NewRelocatable(0, 0)
	val := memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(5))
	val2 := memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(5))
	data := []memory.MaybeRelocatable{*val, *val2}

	// Load Data
	end_ptr, err := mem_manager.LoadData(ptr, &data)
	if err != nil {
		t.Errorf("LoadData error in test: %s", err)
	}

	// Check returned ptr
	expected_end_ptr := memory.NewRelocatable(0, 2)
	if !reflect.DeepEqual(end_ptr, expected_end_ptr) {
		t.Errorf("LoadData returned wrong ptr")
	}

	// Check inserted values

	// val
	res_val, err := mem_manager.Memory.Get(ptr)
	if err != nil {
		t.Errorf("Get error in test: %s", err)
	}

	// Check that the original and the retrieved values are the same
	if !reflect.DeepEqual(res_val, val) {
		t.Errorf("Inserted value and original value are not the same")
	}

	//val2
	ptr2 := memory.NewRelocatable(0, 1)
	res_val2, err := mem_manager.Memory.Get(ptr2)
	if err != nil {
		t.Errorf("Get error in test: %s", err)
	}

	// Check that the original and the retrieved values are the same
	if !reflect.DeepEqual(res_val2, val2) {
		t.Errorf("Inserted value and original value are not the same")
	}
}
