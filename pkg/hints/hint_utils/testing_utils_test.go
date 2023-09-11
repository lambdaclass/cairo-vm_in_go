package hint_utils_test

import (
	"testing"

	. "github.com/lambdaclass/cairo-vm.go/pkg/hints/hint_utils"
	. "github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
	"github.com/lambdaclass/cairo-vm.go/pkg/vm"
	. "github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
)

func TestSetupIdsForTestSimpleValues(t *testing.T) {
	// ids.a = 17
	// ids.b = 7
	vm := vm.NewVirtualMachine()
	vm.Segments.AddSegment()
	ids := SetupIdsForTest(
		map[string][]*MaybeRelocatable{
			"a": {NewMaybeRelocatableFelt(FeltFromUint64(17))},
			"b": {NewMaybeRelocatableFelt(FeltFromUint64(7))},
		},
		vm,
	)
	// Check that we can fetch the values from ids
	a, err_a := ids.GetFelt("a", vm)
	b, err_b := ids.GetFelt("b", vm)
	if err_a != nil || err_b != nil {
		t.Error("Fetching ids failed")
	}
	if a != FeltFromUint64(17) || b != FeltFromUint64(7) {
		t.Error("Wrong ids values")
	}
}

func TestSetupIdsForTestSimpleValuesLeaveGap(t *testing.T) {
	// ids.a = 17
	// ids.b = 7
	vm := vm.NewVirtualMachine()
	vm.Segments.AddSegment()
	ids := SetupIdsForTest(
		map[string][]*MaybeRelocatable{
			"a": {NewMaybeRelocatableFelt(FeltFromUint64(17))},
			"b": {nil},
			"c": {NewMaybeRelocatableFelt(FeltFromUint64(7))},
		},
		vm,
	)
	// Check that we can fetch the values from ids
	a, err_a := ids.GetFelt("a", vm)
	c, err_c := ids.GetFelt("c", vm)
	if err_a != nil || err_c != nil {
		t.Error("Fetching ids failed")
	}
	if a != FeltFromUint64(17) || c != FeltFromUint64(7) {
		t.Error("Wrong ids values")
	}
	// Check that we have a gap at b (we should be able to insert into it)
	err_b := ids.Insert("b", NewMaybeRelocatableFelt(FeltFromUint64(15)), vm)
	if err_b != nil {
		t.Error("Failed to insert ids")
	}
}

func TestSetupIdsForTestStruct(t *testing.T) {
	/* struct Cat {
		lives : 7 // off 0
		paws: 4 // off 1

	}*/
	vm := vm.NewVirtualMachine()
	vm.Segments.AddSegment()
	ids := SetupIdsForTest(
		map[string][]*MaybeRelocatable{
			"cat": {NewMaybeRelocatableFelt(FeltFromUint64(7)), NewMaybeRelocatableFelt(FeltFromUint64(4))},
		},
		vm,
	)
	// Check that we can fetch the values from ids
	lives, err_lives := ids.GetStructFieldFelt("cat", 0, vm)
	paws, err_paws := ids.GetStructFieldFelt("cat", 1, vm)
	if err_lives != nil || err_paws != nil {
		t.Error("Fetching ids failed")
	}
	if lives != FeltFromUint64(7) || paws != FeltFromUint64(4) {
		t.Error("Wrong ids values")
	}
}
