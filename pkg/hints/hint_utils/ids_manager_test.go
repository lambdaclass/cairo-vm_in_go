package hint_utils_test

import (
	"testing"

	. "github.com/lambdaclass/cairo-vm.go/pkg/hints/hint_utils"
	"github.com/lambdaclass/cairo-vm.go/pkg/vm"
)

func TestIdsManagerGetAddressSimpleReference(t *testing.T) {
	ids := IdsManager{
		References: map[string]HintReference{
			"val": {
				Offset1: OffsetValue{
					Register: vm.FP,
				},
			},
		},
	}
	vm := vm.NewVirtualMachine()
	addr, err := ids.GetAddr("val", vm)
	if err != nil {
		t.Errorf("Error in test: %s", err)
	}
	if addr != vm.RunContext.Fp {
		t.Errorf("IdsManager.GetAddr returned wrong value")
	}
}
