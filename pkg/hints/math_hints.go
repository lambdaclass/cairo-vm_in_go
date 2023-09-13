package hints

import (
	"github.com/lambdaclass/cairo-vm.go/pkg/builtins"
	. "github.com/lambdaclass/cairo-vm.go/pkg/hints/hint_utils"
	. "github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
	. "github.com/lambdaclass/cairo-vm.go/pkg/vm"
	. "github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
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

func is_positive(ids IdsManager, vm *VirtualMachine) error {
	value, err := ids.GetFelt("value", vm)
	if err != nil {
		return err
	}
	signedValue := value.ToSigned()
	if signedValue.BitLen() >= builtins.RANGE_CHECK_N_PARTS*builtins.INNER_RC_BOUND_SHIFT {
		return errors.Errorf("Assertion Failed: abs(val) < rc_bound, value=%s is out of the  valid range", signedValue)
	}
	is_positive := uint64(0)
	if signedValue.Sign() == 1 {
		is_positive = 1
	}
	ids.Insert("is_positive", NewMaybeRelocatableFelt(FeltFromUint64(is_positive)), vm)
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

func verify_ecdsa_signature(ids IdsManager, vm *VirtualMachine) error {
	r, err_get_r := ids.GetFelt("signature_r", vm)
	if err_get_r != nil {
		return err_get_r
	}

	s, err_get_s := ids.GetFelt("signature_s", vm)
	if err_get_s != nil {
		return err_get_s
	}

	ecdsa_ptr, err_get_ecdsa := ids.GetAddr("ecdsa_ptr", vm)
	if err_get_ecdsa != nil {
		return err_get_ecdsa
	}

	signature_builtin_interface, err_get_builtin := vm.GetBuiltinRunner("signature")
	if err_get_builtin != nil {
		return err_get_builtin
	}

	signature_builtin := (*signature_builtin_interface).(*builtins.SignatureBuiltinRunner)

	signature := builtins.Signature{
		R: r,
		S: s,
	}
	signature_builtin.AddSignature(ecdsa_ptr, signature)
	return nil
}
