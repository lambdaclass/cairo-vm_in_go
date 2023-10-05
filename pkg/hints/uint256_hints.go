package hints

import (
	"math/big"

	. "github.com/lambdaclass/cairo-vm.go/pkg/hints/hint_utils"
	. "github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
	. "github.com/lambdaclass/cairo-vm.go/pkg/vm"
	. "github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
	"github.com/pkg/errors"
)

func ErrRootOOR(root *big.Int) error {
	return errors.Errorf("assert 0 <= %d < 2**128", root)
}

/*
Implements hints:
%{
    sum_low = ids.a.low + ids.b.low
    ids.carry_low = 1 if sum_low >= ids.SHIFT else 0
    sum_high = ids.a.high + ids.b.high + ids.carry_low
    ids.carry_high = 1 if sum_high >= ids.SHIFT else 0
%}
%{
    sum_low = ids.a.low + ids.b.low
    ids.carry_low = 1 if sum_low >= ids.SHIFT else 0
%}
*/

func uint256Add(ids IdsManager, vm *VirtualMachine, lowOnly bool) error {
	shift := FeltOne().Shl(128)
	a, err := ids.GetUint256("a", vm)
	if err != nil {
		return err
	}
	b, err := ids.GetUint256("b", vm)
	if err != nil {
		return err
	}

	sumLow := a.Low.Add(b.Low)
	carryLow := FeltZero()
	if sumLow.Cmp(shift) != -1 {
		carryLow = FeltOne()
	}

	if !lowOnly {
		sumHigh := a.High.Add(b.High.Add(carryLow))
		carryHigh := FeltZero()
		if sumHigh.Cmp(shift) != -1 {
			carryHigh = FeltOne()
		}
		err := ids.Insert("carry_high", NewMaybeRelocatableFelt(carryHigh), vm)
		if err != nil {
			return err
		}
	}
	return ids.Insert("carry_low", NewMaybeRelocatableFelt(carryLow), vm)

}

/*
Implements hint:
%{
    def split(num: int, num_bits_shift: int = 128, length: int = 2):
        a = []
        for _ in range(length):
            a.append( num & ((1 << num_bits_shift) - 1) )
            num = num >> num_bits_shift
        return tuple(a)

    def pack(z, num_bits_shift: int = 128) -> int:
        limbs = (z.low, z.high)
        return sum(limb << (num_bits_shift * i) for i, limb in enumerate(limbs))

    a = pack(ids.a)
    b = pack(ids.b)
    res = (a - b)%2**256
    res_split = split(res)
    ids.res.low = res_split[0]
    ids.res.high = res_split[1]
%}
*/

func uint256Sub(ids IdsManager, vm *VirtualMachine) error {
	a, err := ids.GetUint256("a", vm)
	if err != nil {
		return err
	}
	b, err := ids.GetUint256("b", vm)
	if err != nil {
		return err
	}
	var resBigInt *big.Int
	if a.ToBigInt().Cmp(b.ToBigInt()) != -1 {
		resBigInt = new(big.Int).Sub(a.ToBigInt(), b.ToBigInt())
	} else {
		mod256 := new(big.Int).Lsh(new(big.Int).SetUint64(1), 256)
		if mod256.Cmp(b.ToBigInt()) != -1 {
			resBigInt = new(big.Int).Sub(mod256, b.ToBigInt())
			resBigInt = new(big.Int).Add(resBigInt, a.ToBigInt())
		} else {
			loweredB := new(big.Int).Mod(b.ToBigInt(), mod256)
			if a.ToBigInt().Cmp(loweredB) != -1 {
				resBigInt = new(big.Int).Sub(a.ToBigInt(), loweredB)
			} else {
				resBigInt = new(big.Int).Sub(mod256, loweredB)
				resBigInt = new(big.Int).Add(resBigInt, a.ToBigInt())
			}
		}
	}

	res := ToUint256(resBigInt)
	return ids.InsertUint256("res", res, vm)
}

/*
Implements hint:

	%{
	    ids.low = ids.a & ((1<<64) - 1)
	    ids.high = ids.a >> 64

%}
*/
func split64(ids IdsManager, vm *VirtualMachine) error {
	a, err := ids.GetFelt("a", vm)
	if err != nil {
		return err
	}
	flag := (FeltOne().Shl(64)).Sub(FeltOne()) // (1 << 64) - 1
	low := a.And(flag)
	high := a.Shr(64) // a >> 64
	err = ids.Insert("low", NewMaybeRelocatableFelt(low), vm)
	if err != nil {
		return err
	}
	err = ids.Insert("high", NewMaybeRelocatableFelt(high), vm)
	if err != nil {
		return err
	}
	return nil

}

/*
Implements hint:

	%{
	    from starkware.python.math_utils import isqrt
	    n = (ids.n.high << 128) + ids.n.low
	    root = isqrt(n)
	    assert 0 <= root < 2 ** 128
	    ids.root.low = root
	    ids.root.high = 0

%}
*/
func uint256Sqrt(ids IdsManager, vm *VirtualMachine, onlyLow bool) error {
	uintN, err := ids.GetUint256("n", vm)
	if err != nil {
		return err
	}
	n := uintN.ToBigInt()
	root := new(big.Int).Sqrt(n)
	if root.BitLen() > 128 {
		return ErrRootOOR(root)
	}

	feltRoot := FeltFromBigInt(root)
	if onlyLow {
		return ids.Insert("root", NewMaybeRelocatableFelt(feltRoot), vm)
	} else {
		return ids.InsertUint256("root", Uint256{Low: feltRoot, High: FeltZero()}, vm)
	}
}

/*
Implements hint:
%{ memory[ap] = 1 if 0 <= (ids.a.high % PRIME) < 2 ** 127 else 0 %}
*/
func uint256SignedNN(ids IdsManager, vm *VirtualMachine) error {
	a, err := ids.GetUint256("a", vm)
	if err != nil {
		return err
	}
	i128Max := FeltFromDecString("170141183460469231731687303715884105727")
	if a.High.Cmp(FeltZero()) != -1 && a.High.Cmp(i128Max) != 1 {
		return vm.Segments.Memory.Insert(vm.RunContext.Ap, NewMaybeRelocatableFelt(FeltOne()))
	} else {
		return vm.Segments.Memory.Insert(vm.RunContext.Ap, NewMaybeRelocatableFelt(FeltZero()))
	}
}

/*
Implements hint:

	%{
	    a = (ids.a.high << 128) + ids.a.low
	    div = (ids.div.high << 128) + ids.div.low
	    quotient, remainder = divmod(a, div)

	    ids.quotient.low = quotient & ((1 << 128) - 1)
	    ids.quotient.high = quotient >> 128
	    ids.remainder.low = remainder & ((1 << 128) - 1)
	    ids.remainder.high = remainder >> 128

%}
*/
func uint256UnsignedDivRem(ids IdsManager, vm *VirtualMachine) error {
	return uint256OfssetedUnisgnedDivRem(ids, vm, 0, 1)

}

/*
Implements hint:

	%{
	    a = (ids.a.high << 128) + ids.a.low
	    div = (ids.div.b23 << 128) + ids.div.b01
	    quotient, remainder = divmod(a, div)

	    ids.quotient.low = quotient & ((1 << 128) - 1)
	    ids.quotient.high = quotient >> 128
	    ids.remainder.low = remainder & ((1 << 128) - 1)
	    ids.remainder.high = remainder >> 128

%}
*/
func uint256ExpandedUnsignedDivRem(ids IdsManager, vm *VirtualMachine) error {
	return uint256OfssetedUnisgnedDivRem(ids, vm, 1, 3)
}

func uint256OfssetedUnisgnedDivRem(ids IdsManager, vm *VirtualMachine, divOffsetLow uint, divOffsetHigh uint) error {
	a, err := ids.GetUint256("a", vm)
	if err != nil {
		return err
	}
	divLow, err := ids.GetStructFieldFelt("div", divOffsetLow, vm)
	if err != nil {
		return err
	}
	divHigh, err := ids.GetStructFieldFelt("div", divOffsetHigh, vm)
	if err != nil {
		return err
	}
	div := Uint256{Low: divLow, High: divHigh}
	q, r := new(big.Int).DivMod(a.ToBigInt(), div.ToBigInt(), new(big.Int))

	err = ids.InsertUint256("quotient", ToUint256(q), vm)
	if err != nil {
		return err
	}
	return ids.InsertUint256("remainder", ToUint256(r), vm)

}

/*
Implements hint:

	%{
	    a = (ids.a.high << 128) + ids.a.low
	    div = (ids.div.b23 << 128) + ids.div.b01
	    quotient, remainder = divmod(a, div)

	    ids.quotient.low = quotient & ((1 << 128) - 1)
	    ids.quotient.high = quotient >> 128
	    ids.remainder.low = remainder & ((1 << 128) - 1)
	    ids.remainder.high = remainder >> 128

%}
*/
func uint256MulDivMod(ids IdsManager, vm *VirtualMachine) error {
	a, err := ids.GetUint256("a", vm)
	if err != nil {
		return err
	}
	b, err := ids.GetUint256("b", vm)
	if err != nil {
		return err
	}
	div, err := ids.GetUint256("div", vm)
	if err != nil {
		return err
	}

	if div.ToBigInt().Cmp(big.NewInt(0)) == 0 {
		return errors.Errorf("Attempted to divide by zero")
	}

	mul := new(big.Int).Mul(a.ToBigInt(), b.ToBigInt())
	quotient, rem := new(big.Int).DivMod(mul, div.ToBigInt(), new(big.Int))

	maxU128, _ := new(big.Int).SetString("340282366920938463463374607431768211455", 10)

	var quotientLow Uint256
	var quotientHigh Uint256
	var remainder Uint256
	quotientLow.Low = FeltFromBigInt(new(big.Int).And(quotient, maxU128))                         // q & maxU128
	quotientLow.High = FeltFromBigInt(new(big.Int).And(new(big.Int).Rsh(quotient, 128), maxU128)) // q >> 128 & maxU128
	quotientHigh.Low = FeltFromBigInt(new(big.Int).And(new(big.Int).Rsh(quotient, 256), maxU128)) // q >> 256 & maxU128
	quotientHigh.High = FeltFromBigInt(new(big.Int).Rsh(quotient, 384))                           // q >> 384
	remainder.Low = FeltFromBigInt(new(big.Int).And(rem, maxU128))                                // rem & maxU128
	remainder.High = FeltFromBigInt(new(big.Int).And(new(big.Int).Rsh(rem, 128), maxU128))        // rem >> 128 & maxU128

	err = ids.InsertUint256("quotient_low", quotientLow, vm)
	if err != nil {
		return err
	}
	err = ids.InsertUint256("quotient_high", quotientHigh, vm)
	if err != nil {
		return err
	}
	return ids.InsertUint256("remainder", remainder, vm)
}
