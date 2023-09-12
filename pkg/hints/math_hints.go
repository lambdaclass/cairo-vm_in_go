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

func assert_not_equal(ids IdsManager, vm *VirtualMachine) error {
	// Extract Ids Variables
	a, err := ids.Get("a", vm)
	if err != nil {
		return err
	}
	b, err := ids.Get("b", vm)
	if err != nil {
		return err
	}
	// Hint Logic
	a_rel, a_is_rel := a.GetRelocatable()
	b_rel, b_is_rel := b.GetRelocatable()
	if !((a_is_rel && b_is_rel && a_rel.SegmentIndex == b_rel.SegmentIndex) || (!a_is_rel && !b_is_rel)) {
		return errors.Errorf("assert_not_equal failed: non-comparable values: %v, %v.", a, b)
	}
	diff, err := a.Sub(*b)
	if err != nil {
		return err
	}
	if diff.IsZero() {
		return errors.Errorf("assert_not_equal failed: %v = %v.", a, b)
	}
	return nil
}
