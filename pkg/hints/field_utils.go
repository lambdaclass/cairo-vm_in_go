package hints

import (
	"errors"
	"fmt"
	"math/big"

	. "github.com/lambdaclass/cairo-vm.go/pkg/hints/hint_utils"
	"github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
	"github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"

	. "github.com/lambdaclass/cairo-vm.go/pkg/types"
	. "github.com/lambdaclass/cairo-vm.go/pkg/vm"
)

/*
Implements hint:
%{
    from starkware.cairo.common.cairo_secp.secp_utils import pack

    q, r = divmod(pack(ids.val, PRIME), SECP_P)
    assert r == 0, f"verify_zero: Invalid input {ids.val.d0, ids.val.d1, ids.val.d2}."
    ids.q = q % PRIME
%}
*/

func verifyZeroWithExternalConst(vm VirtualMachine, execScopes ExecutionScopes, idsData IdsManager) error {
	secpPuncast, err := execScopes.Get("SECP_P")
	if err != nil {
		return err
	}
	secpP := secpPuncast.(big.Int)
	fmt.Println("secp: ", secpP.Text(10))
	addr, err := idsData.GetAddr("val", &vm)
	if err != nil {
		return err
	}

	fmt.Println("addr ", addr)
	val, err := BigInt3FromBaseAddr(addr, "val", &vm)
	if err != nil {
		return err
	}

	v := val.Pack86()
	fmt.Println("val in zero with external: ", v.Text(10))
	q, r := v.DivMod(&v, &secpP, new(big.Int))
	//fmt.Println(r)
	if r.Cmp(big.NewInt(0)) != 0 {
		return errors.New("verify remainder is not zero: Invalid input")
	}

	quotient := memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromBigInt(q))
	idsData.Insert("q", quotient, &vm)
	return nil
}
