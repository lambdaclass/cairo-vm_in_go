package hint_utils

import (
	"fmt"
	"math/big"

	"github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
	"github.com/lambdaclass/cairo-vm.go/pkg/types"
	"github.com/lambdaclass/cairo-vm.go/pkg/vm"
	"github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
)

/*
Implements hint:
%{
    from starkware.cairo.common.cairo_secp.secp_utils import split

    segments.write_arg(ids.res.address_, split(value))
%}
*/

func NondetBigInt3(virtual_machine vm.VirtualMachine, exec_scopes types.ExecutionScopes, ids_data IdsManager) error {
	res_relloc, err := ids_data.GetAddr("res", &virtual_machine)
	if err != nil {
		return err
	}

	fmt.Println("res alloc: ", res_relloc)

	value_uncast, err := exec_scopes.Get("value")
	if err != nil {
		return err
	}
	value := value_uncast.(big.Int)
	fmt.Println("value in nond int: \n", value.Text(10))

	bigint3_split, err := Bigint3Split(value)
	if err != nil {
		return err
	}
	arg := make([]memory.MaybeRelocatable, 0)

	fmt.Println("args: ")
	for i := 0; i < 3; i++ {
		m := memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromBigInt(&bigint3_split[i]))
		fmt.Println(bigint3_split[i].Text(10))
		arg = append(arg, *m)
	}

	virtual_machine.Segments.LoadData(res_relloc, &arg)
	return nil
}
