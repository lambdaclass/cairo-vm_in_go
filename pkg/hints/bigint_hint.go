package hints

import (
	"errors"
	"math/big"

	. "github.com/lambdaclass/cairo-vm.go/pkg/hints/hint_utils"
	"github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
	. "github.com/lambdaclass/cairo-vm.go/pkg/types"
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

func NondetBigInt3(virtual_machine VirtualMachine, execScopes ExecutionScopes, idsData IdsManager) error {
	resRelloc, err := idsData.GetAddr("res", &virtual_machine)
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

	_, loadErr := virtual_machine.Segments.LoadData(resRelloc, &arg)
	if loadErr != nil {
		return loadErr
	}

	return nil
}
