package dict_manager_test

import (
	"reflect"
	"testing"

	. "github.com/lambdaclass/cairo-vm.go/pkg/hints/dict_manager"
	. "github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
	"github.com/lambdaclass/cairo-vm.go/pkg/vm"
	. "github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
)

// DictManager

func TestDictManagerNewDictionaryGetTracker(t *testing.T) {
	dictManager := NewDictManager()
	initialDict := &map[MaybeRelocatable]MaybeRelocatable{}
	vm := vm.NewVirtualMachine()
	base := dictManager.NewDictionary(initialDict, vm)
	if base.SegmentIndex != int(vm.Segments.Memory.NumSegments())-1 {
		t.Errorf("Segment not created for DictTracker")
	}
	_, err := dictManager.GetTracker(base)
	if err != nil {
		t.Errorf("GetTracker failed: %s", err)
	}
}

func TestDictManagerNewDefaultDictionaryGetTracker(t *testing.T) {
	dictManager := NewDictManager()
	vm := vm.NewVirtualMachine()
	base := dictManager.NewDefaultDictionary(nil, vm)
	if base.SegmentIndex != int(vm.Segments.Memory.NumSegments())-1 {
		t.Errorf("Segment not created for DictTracker")
	}
	_, err := dictManager.GetTracker(base)
	if err != nil {
		t.Errorf("GetTracker failed: %s", err)
	}
}

func TestDictManagerNewDictionaryGetTrackerBadDictPtr(t *testing.T) {
	dictManager := NewDictManager()
	initialDict := &map[MaybeRelocatable]MaybeRelocatable{}
	vm := vm.NewVirtualMachine()
	base := dictManager.NewDictionary(initialDict, vm)
	_, err := dictManager.GetTracker(base.AddUint(1))
	if err == nil {
		t.Errorf("GetTracker should have failed")
	}
}

func TestDictManagerNewDictionaryGetTrackerNoTracker(t *testing.T) {
	dictManager := NewDictManager()
	_, err := dictManager.GetTracker(NewRelocatable(0, 0))
	if err == nil {
		t.Errorf("GetTracker should have failed")
	}
}

// DictTracker
func TestDictTrackerDefaultDictCopyDict(t *testing.T) {
	initialDict := &map[MaybeRelocatable]MaybeRelocatable{}
	dictTracker := NewDictTrackerForDefaultDictionary(
		NewRelocatable(0, 0),
		NewMaybeRelocatableFelt(FeltFromUint64(17)),
	)
	// Check CopyDict
	if !reflect.DeepEqual(dictTracker.CopyDictionary(), *initialDict) {
		t.Error("Wrong dict returned by CopyDictionary")
	}
}

func TestDictTrackerDefaultGetValuePresent(t *testing.T) {
	dictTracker := NewDictTrackerForDefaultDictionary(
		NewRelocatable(0, 0),
		NewMaybeRelocatableFelt(FeltFromUint64(17)),
	)
	dictTracker.InsertValue(NewMaybeRelocatableFelt(FeltFromUint64(1)), NewMaybeRelocatableFelt(FeltFromUint64(2)))
	// Check GetValue
	val, err := dictTracker.GetValue(NewMaybeRelocatableFelt(FeltFromUint64(1)))
	if err != nil {
		t.Errorf("GetValue failed with error %s", err)
	}

	if *val != *NewMaybeRelocatableFelt(FeltFromUint64(2)) {
		t.Error("Wrong value returned by GetValue")
	}
}

func TestDictTrackerDefaultGetValueNotPresent(t *testing.T) {
	dictTracker := NewDictTrackerForDefaultDictionary(
		NewRelocatable(0, 0),
		NewMaybeRelocatableFelt(FeltFromUint64(17)),
	)
	dictTracker.InsertValue(NewMaybeRelocatableFelt(FeltFromUint64(1)), NewMaybeRelocatableFelt(FeltFromUint64(2)))
	// Check that the default value is returned
	val, err := dictTracker.GetValue(NewMaybeRelocatableFelt(FeltFromUint64(2)))
	if err != nil {
		t.Errorf("GetValue failed with error %s", err)
	}

	if *val != *NewMaybeRelocatableFelt(FeltFromUint64(17)) {
		t.Error("Wrong value returned by GetValue")
	}
	// Check that the default value was written to the Dictionary
	dictCopy := dictTracker.CopyDictionary()
	if dictCopy[*NewMaybeRelocatableFelt(FeltFromUint64(2))] != *NewMaybeRelocatableFelt(FeltFromUint64(17)) {
		t.Error("Default value not written after GetValue")
	}
}

func TestDictTrackerGetValueNotPresent(t *testing.T) {
	initialDict := &map[MaybeRelocatable]MaybeRelocatable{
		*NewMaybeRelocatableFelt(FeltFromUint64(1)): *NewMaybeRelocatableFelt(FeltFromUint64(2)),
	}
	dictTracker := NewDictTrackerForDictionary(
		NewRelocatable(0, 0),
		initialDict,
	)
	// Check GetValue
	_, err := dictTracker.GetValue(NewMaybeRelocatableFelt(FeltFromUint64(1)))
	if err != nil {
		t.Error("GetValue should have failed")
	}
}

func TestDictTrackerGetValuePresent(t *testing.T) {
	initialDict := &map[MaybeRelocatable]MaybeRelocatable{
		*NewMaybeRelocatableFelt(FeltFromUint64(1)): *NewMaybeRelocatableFelt(FeltFromUint64(2)),
	}
	dictTracker := NewDictTrackerForDictionary(
		NewRelocatable(0, 0),
		initialDict,
	)
	// Check GetValue
	val, err := dictTracker.GetValue(NewMaybeRelocatableFelt(FeltFromUint64(1)))
	if err != nil {
		t.Errorf("GetValue failed with error %s", err)
	}

	if *val != *NewMaybeRelocatableFelt(FeltFromUint64(2)) {
		t.Error("Wrong value returned by GetValue")
	}
}

func TestDictTrackerInsertValue(t *testing.T) {
	initialDict := &map[MaybeRelocatable]MaybeRelocatable{}
	dictTracker := NewDictTrackerForDictionary(
		NewRelocatable(0, 0),
		initialDict,
	)
	// InsertValue
	dictTracker.InsertValue(NewMaybeRelocatableFelt(FeltFromUint64(7)), NewMaybeRelocatableFelt(FeltFromUint64(8)))
	// Check GetValue
	val, err := dictTracker.GetValue(NewMaybeRelocatableFelt(FeltFromUint64(7)))
	if err != nil {
		t.Errorf("GetValue failed with error %s", err)
	}

	if *val != *NewMaybeRelocatableFelt(FeltFromUint64(8)) {
		t.Error("Wrong value returned by GetValue")
	}
}

func TestDictTrackerCopyDict(t *testing.T) {
	initialDict := &map[MaybeRelocatable]MaybeRelocatable{
		*NewMaybeRelocatableFelt(FeltFromUint64(1)): *NewMaybeRelocatableFelt(FeltFromUint64(2)),
	}
	dictTracker := NewDictTrackerForDictionary(
		NewRelocatable(0, 0),
		initialDict,
	)
	// Check CopyDict
	if !reflect.DeepEqual(dictTracker.CopyDictionary(), *initialDict) {
		t.Error("Wrong dict returned by CopyDictionary")
	}
}

// Dictionary

func TestDictionary(t *testing.T) {
	initialDict := &map[MaybeRelocatable]MaybeRelocatable{
		*NewMaybeRelocatableFelt(FeltFromUint64(1)): *NewMaybeRelocatableFelt(FeltFromUint64(2)),
	}
	dict := NewDictionary(*&initialDict)
	// Check Get
	if *dict.Get(NewMaybeRelocatableFelt(FeltFromUint64(1))) != *NewMaybeRelocatableFelt(FeltFromUint64(2)) {
		t.Error("Wrong value returned by Get")
	}
	// InsertValue
	dict.Insert(NewMaybeRelocatableFelt(FeltFromUint64(7)), NewMaybeRelocatableFelt(FeltFromUint64(8)))
	// Check Get
	if *dict.Get(NewMaybeRelocatableFelt(FeltFromUint64(7))) != *NewMaybeRelocatableFelt(FeltFromUint64(8)) {
		t.Error("Wrong value returned by Get")
	}
}

func TestDefaultDictionary(t *testing.T) {
	dict := NewDefaultDictionary(NewMaybeRelocatableFelt(FeltFromUint64(17)))
	// InsertValue
	dict.Insert(NewMaybeRelocatableFelt(FeltFromUint64(7)), NewMaybeRelocatableFelt(FeltFromUint64(8)))
	// Check Get
	if *dict.Get(NewMaybeRelocatableFelt(FeltFromUint64(7))) != *NewMaybeRelocatableFelt(FeltFromUint64(8)) {
		t.Error("Wrong value returned by Get")
	}
	// Check Get DefaultValue
	if *dict.Get(NewMaybeRelocatableFelt(FeltFromUint64(3))) != *NewMaybeRelocatableFelt(FeltFromUint64(17)) {
		t.Error("Wrong value returned by Get")
	}
}
