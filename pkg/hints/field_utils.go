package hints

import (
	"errors"
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
	secpP, ok := secpPuncast.(big.Int)
	if !ok {
		return errors.New("Could not cast secpP into big int")
	}

	addr, err := idsData.GetAddr("val", &vm)
	if err != nil {
		return err
	}

	val, err := BigInt3FromBaseAddr(addr, "val", &vm)
	if err != nil {
		return err
	}

	v := val.Pack86()
	q, r := v.DivMod(&v, &secpP, new(big.Int))

	if r.Cmp(big.NewInt(0)) != 0 {
		return errors.New("verify remainder is not zero: Invalid input")
	}

	quotient := memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromBigInt(q))
	return idsData.Insert("q", quotient, &vm)
}
