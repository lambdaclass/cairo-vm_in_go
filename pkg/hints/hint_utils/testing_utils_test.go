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
			"a": []*MaybeRelocatable{NewMaybeRelocatableFelt(FeltFromUint64(17))},
			"b": []*MaybeRelocatable{NewMaybeRelocatableFelt(FeltFromUint64(7))},
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
		t.Error("Wromg ids values")
	}

}
