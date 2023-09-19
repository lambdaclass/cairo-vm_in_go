package hints

import (
	"fmt"
	"math/big"

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

/*
Implements hint:

	%{
	    from starkware.cairo.common.math_utils import assert_integer
	    assert_integer(ids.div)
	    assert 0 < ids.div <= PRIME // range_check_builtin.bound, \
	        f'div={hex(ids.div)} is out of the valid range.'
	    ids.q, ids.r = divmod(ids.value, ids.div)
	%}
*/
func unsignedDivRem(ids IdsManager, vm *VirtualMachine) error {
	div, err := ids.GetFelt("div", vm)
	if err != nil {
		return err
	}
	value, err := ids.GetFelt("value", vm)
	if err != nil {
		return err
	}

	rcBound, err := vm.GetRangeCheckBound()
	if err != nil {
		return err
	}

	if rcBound.Cmp(lambdaworks.FeltZero()) == 0 {
		return errors.New("range check bound cannot be zero")
	}
	primeBoundDivision := new(big.Int).Div(lambdaworks.Prime(), rcBound.ToBigInt())

	// Check if `div` is greater than `limit`
	divGreater := div.ToBigInt().Cmp(primeBoundDivision) == 1

	if div.IsZero() || divGreater {
		return errors.Errorf("Div out of range: 0 < %d <= %d", div, rcBound)
	}

	q, r := value.DivRem(div)
	ids.Insert("q", NewMaybeRelocatableFelt(q), vm)
	ids.Insert("r", NewMaybeRelocatableFelt(r), vm)

	return nil
}

/*
Implements hint:

    %{
        from starkware.cairo.common.math_utils import as_int, assert_integer

        assert_integer(ids.div)
        assert 0 < ids.div <= PRIME // range_check_builtin.bound, \
            f'div={hex(ids.div)} is out of the valid range.'

        assert_integer(ids.bound)
        assert ids.bound <= range_check_builtin.bound // 2, \
            f'bound={hex(ids.bound)} is out of the valid range.'

        int_value = as_int(ids.value, PRIME)
        q, ids.r = divmod(int_value, ids.div)

        assert -ids.bound <= q < ids.bound, \
            f'{int_value} / {ids.div} = {q} is out of the range [{-ids.bound}, {ids.bound}).'

        ids.biased_q = q + ids.bound
    %}
*/

func signedDivRem(ids IdsManager, vm *VirtualMachine) error {
	div, err := ids.GetFelt("div", vm)
	if err != nil {
		return err
	}
	value, err := ids.GetFelt("value", vm)
	if err != nil {
		return err
	}
	bound, err := ids.GetFelt("bound", vm)
	if err != nil {
		return err
	}

	rcBound, err := vm.GetRangeCheckBound()
	if err != nil {
		return err
	}

	if rcBound.Cmp(lambdaworks.FeltZero()) == 0 {
		return errors.New("range check bound cannot be zero")
	}
	primeBoundDivision := new(big.Int).Div(lambdaworks.Prime(), rcBound.ToBigInt())

	// Check if `div` is greater than `limit` and make assertions
	divGreater := div.ToBigInt().Cmp(primeBoundDivision) == 1
	if div.IsZero() || divGreater {
		return errors.Errorf("div=%d is out of the valid range", div)
	}

	if bound.Cmp(rcBound.Shr(1)) == 1 {
		return errors.Errorf("bound=%d is out of the valid range")
	}

	sgnValue := value.ToSigned()
	sgnBound := bound.ToBigInt()
	intDiv := div.ToBigInt()

	q := new(big.Int).Div(sgnValue, intDiv)
	r := new(big.Int).Rem(sgnValue, intDiv)

	fmt.Println(sgnBound, q, r)

	return nil
}

// pub fn signed_div_rem(
//     vm: &mut VirtualMachine,
//     ids_data: &HashMap<String, HintReference>,
//     ap_tracking: &ApTracking,
// ) -> Result<(), HintError> {
//     let div = get_integer_from_var_name("div", vm, ids_data, ap_tracking)?;
//     let value = get_integer_from_var_name("value", vm, ids_data, ap_tracking)?;
//     let value = value.as_ref();
//     let bound = get_integer_from_var_name("bound", vm, ids_data, ap_tracking)?;
//     let builtin = vm.get_range_check_builtin()?;

//     match &builtin._bound {
//         Some(builtin_bound)
//             if div.is_zero() || div.as_ref() > &div_prime_by_bound(builtin_bound.clone())? =>
//         {
//             return Err(HintError::OutOfValidRange(Box::new((
//                 div.into_owned(),
//                 builtin_bound.clone(),
//             ))));
//         }
//         Some(builtin_bound) if bound.as_ref() > &(builtin_bound >> 1_u32) => {
//             return Err(HintError::OutOfValidRange(Box::new((
//                 bound.into_owned(),
//                 builtin_bound >> 1_u32,
//             ))));
//         }
//         None if div.is_zero() => {
//             return Err(HintError::OutOfValidRange(Box::new((
//                 div.into_owned(),
//                 Felt252::zero() - Felt252::one(),
//             ))));
//         }
//         _ => {}
//     }

//     let int_value = value.to_signed_felt();
//     let int_div = div.to_bigint();
//     let int_bound = bound.to_bigint();
//     let (q, r) = int_value.div_mod_floor(&int_div);

//     if int_bound.abs() < q.abs() {
//         return Err(HintError::OutOfValidRange(Box::new((
//             Felt252::new(q),
//             bound.into_owned(),
//         ))));
//     }

//     let biased_q = q + int_bound;
//     insert_value_from_var_name("r", Felt252::new(r), vm, ids_data, ap_tracking)?;
//     insert_value_from_var_name(
//         "biased_q",
//         Felt252::new(biased_q),
//         vm,
//         ids_data,
//         ap_tracking,
//     )
// }
