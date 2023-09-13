package hints

import (
	"github.com/lambdaclass/cairo-vm.go/pkg/builtins"
	. "github.com/lambdaclass/cairo-vm.go/pkg/hints/hint_utils"
	. "github.com/lambdaclass/cairo-vm.go/pkg/vm"
	"github.com/pkg/errors"
)

// Implements hint:
//
//	%{
//	    from starkware.cairo.common.math_utils import assert_integer
//	    assert_integer(ids.a)
//	    assert 0 <= ids.a % PRIME < range_check_builtin.bound, f'a = {ids.a} is out of range.'
//
// %}
func assert_nn(ids IdsManager, vm *VirtualMachine) error {
	a, err := ids.GetFelt("a", vm)
	if err != nil {
		return err
	}
	if a.Bits() >= builtins.RANGE_CHECK_N_PARTS*builtins.INNER_RC_BOUND_SHIFT {
		return errors.Errorf("Assertion failed, 0 <= ids.a %% PRIME < range_check_builtin.bound\n a = %s is out of range", a.ToHexString())
	}
	return nil
}

// Implements hint:from starkware.cairo.common.math.cairo
//
//	%{
//	    from starkware.cairo.common.math_utils import assert_integer
//	    assert_integer(ids.value)
//	    assert ids.value % PRIME != 0, f'assert_not_zero failed: {ids.value} = 0.'
//
// %}
func assert_not_zero(ids IdsManager, vm *VirtualMachine) error {
	value, err := ids.GetFelt("value", vm)
	if err != nil {
		return err
	}
	if value.IsZero() {
		return errors.Errorf("Assertion failed, %s %% PRIME is equal to 0", value.ToHexString())
	}
	return nil
}
