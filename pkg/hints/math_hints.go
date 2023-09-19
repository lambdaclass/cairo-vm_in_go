package hints

import (
	"github.com/lambdaclass/cairo-vm.go/pkg/builtins"
	. "github.com/lambdaclass/cairo-vm.go/pkg/hints/hint_utils"
	"github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
	. "github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
	. "github.com/lambdaclass/cairo-vm.go/pkg/math_utils"
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

	signature_builtin_interface, err_get_builtin := vm.GetBuiltinRunner(builtins.SIGNATURE_BUILTIN_NAME)
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

// Implements hint:from starkware.cairo.common.math.cairo
//
//	%{
//		from starkware.crypto.signature.signature import FIELD_PRIME
//		from starkware.python.math_utils import div_mod, is_quad_residue, sqrt
//
//		x = ids.x
//		if is_quad_residue(x, FIELD_PRIME):
//		    ids.y = sqrt(x, FIELD_PRIME)
//		else:
//		    ids.y = sqrt(div_mod(x, 3, FIELD_PRIME), FIELD_PRIME)
//
// %}
func is_quad_residue(ids IdsManager, vm *VirtualMachine) error {
	x, err := ids.GetFelt("x", vm)
	if err != nil {
		return err
	}
	if x.IsZero() || x.IsOne() {
		ids.Insert("y", NewMaybeRelocatableFelt(x), vm)

	} else if x.Pow(SignedFeltMaxValue()) == FeltOne() {
		num := x.Sqrt()
		ids.Insert("y", NewMaybeRelocatableFelt(num), vm)

	} else {
		num := (x.Div(lambdaworks.FeltFromUint64(3))).Sqrt()
		ids.Insert("y", NewMaybeRelocatableFelt(num), vm)
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

/*
Implements the hint:

	from starkware.python.math_utils import isqrt
	value = ids.value % PRIME
	assert value < 2 ** 250, f"value={value} is outside of the range [0, 2**250)."
	assert 2 ** 250 < PRIME
	ids.root = isqrt(value)
*/
func sqrt(ids IdsManager, vm *VirtualMachine) error {
	value, err := ids.GetFelt("value", vm)
	if err != nil {
		return err
	}

	if value.Bits() >= 250 {
		return errors.Errorf("Value: %v is outside of the range [0, 2**250)", value)
	}

	root_big, err := ISqrt(value.ToBigInt())
	if err != nil {
		return err
	}
	root_felt := FeltFromDecString(root_big.String())
	ids.Insert("root", NewMaybeRelocatableFelt(root_felt), vm)
	return nil
}
