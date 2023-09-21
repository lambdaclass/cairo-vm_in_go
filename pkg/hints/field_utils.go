package hints

import (
	"math/big"

	. "github.com/lambdaclass/cairo-vm.go/pkg/hints/hint_utils"

	"github.com/lambdaclass/cairo-vm.go/pkg/types"
	"github.com/lambdaclass/cairo-vm.go/pkg/vm"
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

func verifyZeroWithExternalConst(virtual_machine vm.VirtualMachine, exec_scopes types.ExecutionScopes, ids_data IdsManager) error {
	secp_p_uncast, err := exec_scopes.Get("SECP_P")
	if err != nil {
		return err
	}
	secp_p := secp_p_uncast.(big.Int)
	addr, err := ids_data.GetRelocatable("val", &virtual_machine)
	if err != nil {
		return err
	}

	val, err := BigInt3FromBaseAddr(addr, virtual_machine)
	if err != nil {
		return err
	}

	return nil
}
