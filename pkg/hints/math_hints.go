package hints

import (
	"math/big"

	"github.com/lambdaclass/cairo-vm.go/pkg/builtins"
	. "github.com/lambdaclass/cairo-vm.go/pkg/hints/hint_utils"
	"github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
	. "github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
	. "github.com/lambdaclass/cairo-vm.go/pkg/math_utils"
	. "github.com/lambdaclass/cairo-vm.go/pkg/types"
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

func assertLeFelt(ids IdsManager, vm *VirtualMachine, scopes *ExecutionScopes, constants *map[string]Felt) error {
	// Fetch constants
	primeOver3HighFelt, err := ids.GetConst("PRIME_OVER_3_HIGH", constants)
	if err != nil {
		return err
	}
	primeOver3High := primeOver3HighFelt.ToBigInt()
	primeOver2HighFelt, err := ids.GetConst("PRIME_OVER_2_HIGH", constants)
	if err != nil {
		return err
	}
	primeOver2High := primeOver2HighFelt.ToBigInt()
	// Fetch ids variables
	aFelt, err := ids.GetFelt("a", vm)
	if err != nil {
		return err
	}
	a := aFelt.ToBigInt()
	bFelt, err := ids.GetFelt("b", vm)
	if err != nil {
		return err
	}
	b := bFelt.ToBigInt()
	rangeCheckPtr, err := ids.GetRelocatable("range_check_ptr", vm)
	if err != nil {
		return err
	}
	// Hint Logic
	cairoPrime, _ := new(big.Int).SetString(CAIRO_PRIME_HEX, 0)
	halfPrime := new(big.Int).Div(cairoPrime, new(big.Int).SetUint64(2))
	thirdOfPrime := new(big.Int).Div(cairoPrime, new(big.Int).SetUint64(2))
	if a.Cmp(b) == 1 {
		return errors.Errorf("Assertion failed, %v, is not less or equal to %v", a, b)
	}
	arc1 := new(big.Int).Sub(b, a)
	arc2 := new(big.Int).Sub(new(big.Int).Sub(cairoPrime, (big.NewInt(1))), b)

	// Split lengthsAndIndices array into lenght & idxs array and mantain the same order between them
	lengths := []*big.Int{a, arc1, arc2}
	idxs := []int{0, 1, 2}
	// Sort lengths & idxs by lengths
	for i := 0; i < 3; i++ {
		for j := i; j > 0 && lengths[j-1].Cmp(lengths[j]) == 1; j-- {
			lengths[j], lengths[j-1] = lengths[j-1], lengths[j]
			idxs[j], idxs[j-1] = idxs[j-1], idxs[j]
		}
	}

	if lengths[0].Cmp(thirdOfPrime) == 1 || lengths[1].Cmp(halfPrime) == 1 {
		return errors.Errorf("Arc too big, %v must be <= %v and %v <= %v", lengths[0], thirdOfPrime, lengths[1], halfPrime)
	}
	excluded := idxs[2]
	scopes.AssignOrUpdateVariable("excluded", excluded)
	q_0, r_0 := new(big.Int).DivMod(lengths[0], primeOver3High, primeOver3High)
	q_1, r_1 := new(big.Int).DivMod(lengths[1], primeOver2High, primeOver2High)

	// Insert values into range_check_ptr
	data := []MaybeRelocatable{
		*NewMaybeRelocatableFelt(FeltFromBigInt(r_0)),
		*NewMaybeRelocatableFelt(FeltFromBigInt(q_0)),
		*NewMaybeRelocatableFelt(FeltFromBigInt(r_1)),
		*NewMaybeRelocatableFelt(FeltFromBigInt(q_1)),
	}
	_, err = vm.Segments.LoadData(rangeCheckPtr, &data)

	return err
}

// "memory[ap] = 1 if excluded != 0 else 0"
func assertLeFeltExcluded0(vm *VirtualMachine, scopes *ExecutionScopes) error {
	// Fetch scope var
	excludedAny, err := scopes.Get("excluded")
	if err != nil {
		return err
	}
	excluded, ok := excludedAny.(int)
	if !ok {
		return errors.New("exluded not in scope")
	}
	if excluded == 0 {
		return vm.Segments.Memory.Insert(vm.RunContext.Ap, NewMaybeRelocatableFelt(FeltZero()))
	}
	return vm.Segments.Memory.Insert(vm.RunContext.Ap, NewMaybeRelocatableFelt(FeltOne()))
}

// "memory[ap] = 1 if excluded != 1 else 0"
func assertLeFeltExcluded1(vm *VirtualMachine, scopes *ExecutionScopes) error {
	// Fetch scope var
	excludedAny, err := scopes.Get("excluded")
	if err != nil {
		return err
	}
	excluded, ok := excludedAny.(int)
	if !ok {
		return errors.New("exluded not in scope")
	}
	if excluded == 1 {
		return vm.Segments.Memory.Insert(vm.RunContext.Ap, NewMaybeRelocatableFelt(FeltZero()))
	}
	return vm.Segments.Memory.Insert(vm.RunContext.Ap, NewMaybeRelocatableFelt(FeltOne()))
}

// "assert excluded == 2"
func assertLeFeltExcluded2(vm *VirtualMachine, scopes *ExecutionScopes) error {
	// Fetch scope var
	excludedAny, err := scopes.Get("excluded")
	if err != nil {
		return err
	}
	excluded, ok := excludedAny.(int)
	if !ok {
		return errors.New("exluded not in scope")
	}
	if excluded != 2 {
		return errors.New("Assertion Failed: excluded == 2")
	}
	return nil
}
