package hints

import (
	"errors"
	"math/big"

	. "github.com/lambdaclass/cairo-vm.go/pkg/hints/hint_utils"
	"github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
	. "github.com/lambdaclass/cairo-vm.go/pkg/types"
	. "github.com/lambdaclass/cairo-vm.go/pkg/utils"
	. "github.com/lambdaclass/cairo-vm.go/pkg/vm"
	"github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
)

/*
Implements hint:
%{
    from starkware.cairo.common.cairo_secp.secp_utils import split

    segments.write_arg(ids.res.address_, split(value))
%}
*/

func NondetBigInt3(vm VirtualMachine, execScopes ExecutionScopes, idsData IdsManager) error {
	resRelloc, err := idsData.GetAddr("res", &vm)
	if err != nil {
		return err
	}

	valueUncast, err := execScopes.Get("value")
	if err != nil {
		return err
	}
	value, ok := valueUncast.(big.Int)
	if !ok {
		return errors.New("Could not cast value into big int")
	}

	bigint3Split, err := Bigint3Split(value)
	if err != nil {
		return err
	}

	arg := make([]memory.MaybeRelocatable, 0)

	for i := 0; i < 3; i++ {
		m := memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromBigInt(&bigint3Split[i]))
		arg = append(arg, *m)
	}

	_, loadErr := vm.Segments.LoadData(resRelloc, &arg)
	if loadErr != nil {
		return loadErr
	}

	return nil
}

/// Implements hint:
/// ```python
/// k = safe_div(res * y - x, p)
/// value = k if k > 0 else 0 - k
/// ids.flag = 1 if k > 0 else 0
/// ```

func SafeDivBigint(vm *VirtualMachine, execScopes *ExecutionScopes, idsData IdsManager) error {
	resUncast, err := execScopes.Get("res")
	if err != nil {
		return err
	}
	res, ok := resUncast.(big.Int)
	if !ok {
		return errors.New("Could not cast res value in SafeDivBigint")
	}

	yUncast, err := execScopes.Get("y")
	if err != nil {
		return err
	}
	y, ok := yUncast.(big.Int)
	if !ok {
		return errors.New("Could not cast y value in SafeDivBigint")
	}

	xUncast, err := execScopes.Get("x")
	if err != nil {
		return err
	}
	x, ok := xUncast.(big.Int)
	if !ok {
		return errors.New("Could not cast x value in SafeDivBigint")
	}

	pUncast, err := execScopes.Get("p")
	if err != nil {
		return err
	}
	p, ok := pUncast.(big.Int)
	if !ok {
		return errors.New("Could not cast p value in SafeDivBigint")
	}

	param_x := new(big.Int).Mul(&res, &y)
	param_x.Sub(param_x, &x)

	k, err := SafeDivBig(param_x, &p)
	if err != nil {
		return err
	}

	var value big.Int
	var flag lambdaworks.Felt

	// check if k is positive
	if k.Cmp(big.NewInt(0)) == 1 {
		value = *k
		flag = lambdaworks.FeltFromUint(1)
	} else {
		value = *new(big.Int).Neg(k)
		flag = lambdaworks.FeltFromUint(0)
	}

	execScopes.AssignOrUpdateVariable("k", *k)
	execScopes.AssignOrUpdateVariable("value", value)

	val := memory.NewMaybeRelocatableFelt(flag)
	idsData.Insert("flag", val, vm)

	return nil
}
