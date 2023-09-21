package hints

import (
	"math/big"

	"github.com/lambdaclass/cairo-vm.go/pkg/hints/hint_utils"
	. "github.com/lambdaclass/cairo-vm.go/pkg/hints/hint_utils"
	"github.com/lambdaclass/cairo-vm.go/pkg/types"
	"github.com/lambdaclass/cairo-vm.go/pkg/vm"
)

type EcPoint struct {
	X BigInt3
	Y BigInt3
}

func BigInt3FromVarName(name string, virtual_machine *vm.VirtualMachine, ids_data hint_utils.IdsManager) (EcPoint, error) {
	point_addr, err := ids_data.GetAddr(name, virtual_machine)
	if err != nil {
		return EcPoint{}, err
	}

	x, err := BigInt3FromBaseAddr(point_addr, name+".x", virtual_machine)
	if err != nil {
		return EcPoint{}, err
	}

	y, err := BigInt3FromBaseAddr(point_addr.AddUint(3), name+".y", virtual_machine)
	if err != nil {
		return EcPoint{}, err
	}

	return EcPoint{X: x, Y: y}, nil
}

/*
Implements main logic for `EC_NEGATE` and `EC_NEGATE_EMBEDDED_SECP` hints
*/
func ecNegate(virtual_machine *vm.VirtualMachine, exec_scopes types.ExecutionScopes, ids_data hint_utils.IdsManager, secp_p big.Int) error {
	point, err := ids_data.GetRelocatable("point", virtual_machine)
	if err != nil {
		return err
	}

	point_y, err := point.AddInt(3)
	if err != nil {
		return err
	}

	y_bigint3, err := BigInt3FromBaseAddr(point_y, "point.y", virtual_machine)
	if err != nil {
		return err
	}

	y := y_bigint3.Pack86()
	value := new(big.Int).Neg(&y)
	value.Mod(value, &secp_p)

	exec_scopes.AssignOrUpdateVariable("value", value)
	exec_scopes.AssignOrUpdateVariable("SECP_P", secp_p)
	return nil
}

func ecNegateImportSecpP(virtual_machine *vm.VirtualMachine, exec_scopes types.ExecutionScopes, ids_data hint_utils.IdsManager) error {
	secp_p, _ := new(big.Int).SetString("115792089237316195423570985008687907853269984665640564039457584007908834671663", 10)
	return ecNegate(virtual_machine, exec_scopes, ids_data, *secp_p)
}

/*
Implements hint:
%{
    from starkware.cairo.common.cairo_secp.secp_utils import pack
    SECP_P = 2**255-19

    y = pack(ids.point.y, PRIME) % SECP_P
    # The modulo operation in python always returns a nonnegative number.
    value = (-y) % SECP_P
%}
*/

func ecNegateEmbeddedSecpP(virtual_machine *vm.VirtualMachine, exec_scopes types.ExecutionScopes, ids_data hint_utils.IdsManager) error {
	secp_p := big.NewInt(1)
	secp_p.Lsh(secp_p, 255)
	secp_p.Sub(secp_p, big.NewInt(19))
	return ecNegate(virtual_machine, exec_scopes, ids_data, *secp_p)
}
