package hints

import (
	"errors"
	"math/big"

	"github.com/lambdaclass/cairo-vm.go/pkg/builtins"
	. "github.com/lambdaclass/cairo-vm.go/pkg/hints/hint_utils"
	"github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
	"github.com/lambdaclass/cairo-vm.go/pkg/types"
	"github.com/lambdaclass/cairo-vm.go/pkg/vm"
	"github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
)

type BigInt3 struct {
	Limbs []lambdaworks.Felt
}

type EcPoint struct {
	X BigInt3
	Y BigInt3
}

func (val *BigInt3) Pack86() big.Int {
	sum := big.NewInt(0)
	for i := 0; i < 3; i++ {
		felt := val.Limbs[i]
		signed := felt.ToSigned()
		shifed := new(big.Int).Lsh(signed, uint(i*86))
		sum.Add(sum, shifed)
	}
	return *sum
}

func BigInt3FromBaseAddr(addr memory.Relocatable, virtual_machine vm.VirtualMachine) (BigInt3, error) {
	limbs := make([]lambdaworks.Felt, 0)
	for i := 0; i < 3; i++ {
		felt, err := virtual_machine.Segments.Memory.GetFelt(addr.AddUint(uint(i)))
		if err == nil {
			limbs = append(limbs, felt)
		} else {
			return BigInt3{}, errors.New("Identifier has no member")
		}
	}
	return BigInt3{Limbs: limbs}, nil
}

func EcPointFromVarName(name string, virtual_machine vm.VirtualMachine, ids_data IdsManager) (EcPoint, error) {
	point_addr, err := ids_data.GetAddr(name, &virtual_machine)
	if err != nil {
		return EcPoint{}, err
	}

	x, err := BigInt3FromBaseAddr(point_addr, virtual_machine)
	if err != nil {
		return EcPoint{}, err
	}

	y, err := BigInt3FromBaseAddr(point_addr.AddUint(3), virtual_machine)
	if err != nil {
		return EcPoint{}, err
	}

	return EcPoint{X: x, Y: y}, nil
}

/*
Implements main logic for `EC_NEGATE` and `EC_NEGATE_EMBEDDED_SECP` hints
*/
func ecNegate(virtual_machine vm.VirtualMachine, exec_scopes types.ExecutionScopes, ids_data IdsManager, secp_p big.Int) error {
	point, err := ids_data.GetRelocatable("point", &virtual_machine)
	if err != nil {
		return err
	}

	point_y, err := point.AddInt(3)
	if err != nil {
		return err
	}

	y_bigint3, err := BigInt3FromBaseAddr(point_y, virtual_machine)
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

func ecNegateImportSecpP(virtual_machine vm.VirtualMachine, exec_scopes types.ExecutionScopes, ids_data IdsManager) error {
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

func ecNegateEmbeddedSecpP(virtual_machine vm.VirtualMachine, exec_scopes types.ExecutionScopes, ids_data IdsManager) error {
	secp_p := big.NewInt(1)
	secp_p.Lsh(secp_p, 255)
	secp_p.Sub(secp_p, big.NewInt(19))
	return ecNegate(virtual_machine, exec_scopes, ids_data, *secp_p)
}

/*
Implements hint:

	%{
	    from starkware.cairo.common.cairo_secp.secp_utils import SECP_P, pack
	    from starkware.python.math_utils import ec_double_slope

	    # Compute the slope.
	    x = pack(ids.point.x, PRIME)
	    y = pack(ids.point.y, PRIME)
	    value = slope = ec_double_slope(point=(x, y), alpha=0, p=SECP_P)

%}
*/
func computeDoublingSlope(virtual_machine vm.VirtualMachine, exec_scopes types.ExecutionScopes, ids_data IdsManager, point_alias string, secp_p big.Int, alpha big.Int) error {
	exec_scopes.AssignOrUpdateVariable("SECP_P", secp_p)

	point, err := EcPointFromVarName(point_alias, virtual_machine, ids_data)
	if err != nil {
		return err
	}

	x := point.X.Pack86()
	y := point.Y.Pack86()
	double_point := builtins.DoublePointB{X: x, Y: y}

	value, err := builtins.EcDoubleSlope(double_point, alpha, secp_p)
	if err != nil {
		return err
	}

	exec_scopes.AssignOrUpdateVariable("value", value)
	exec_scopes.AssignOrUpdateVariable("slope", value)

	return nil
}

/*
Implements hint:
%{
    from starkware.cairo.common.cairo_secp.secp_utils import SECP_P, pack
    from starkware.python.math_utils import line_slope

    # Compute the slope.
    x0 = pack(ids.point0.x, PRIME)
    y0 = pack(ids.point0.y, PRIME)
    x1 = pack(ids.point1.x, PRIME)
    y1 = pack(ids.point1.y, PRIME)
    value = slope = line_slope(point1=(x0, y0), point2=(x1, y1), p=SECP_P)
%}
*/

func computeSlopeAndAssingSecpP(virtual_machine vm.VirtualMachine, exec_scopes types.ExecutionScopes, ids_data IdsManager, point0_alias string, point1_alias string, secp_p big.Int) error {
	exec_scopes.AssignOrUpdateVariable("SECP_P", secp_p)
	return computeSlope(virtual_machine, exec_scopes, ids_data, point0_alias, point1_alias)
}

func computeSlope(virtual_machine vm.VirtualMachine, exec_scopes types.ExecutionScopes, ids_data IdsManager, point0_alias string, point1_alias string) error {
	point0, err := EcPointFromVarName(point0_alias, virtual_machine, ids_data)
	if err != nil {
		return err
	}
	point1, err := EcPointFromVarName(point1_alias, virtual_machine, ids_data)
	if err != nil {
		return err
	}

	secp_p, err := exec_scopes.Get("SECP_P")
	if err != nil {
		return err
	}
	secp := secp_p.(big.Int)

	// build partial sum
	x0 := point0.X.Pack86()
	y0 := point0.Y.Pack86()
	point_a := builtins.PartialSumB{X: x0, Y: y0}

	// build double point
	x1 := point1.X.Pack86()
	y1 := point1.Y.Pack86()
	point_b := builtins.DoublePointB{X: x1, Y: y1}

	value, err := builtins.LineSlope(point_a, point_b, secp)
	if err != nil {
		return err
	}

	exec_scopes.AssignOrUpdateVariable("value", value)
	exec_scopes.AssignOrUpdateVariable("slope", value)

	return nil
}

/*
Implements hint:
%{ from starkware.cairo.common.cairo_secp.secp256r1_utils import SECP256R1_ALPHA as ALPHA %}
*/

func importSecp256r1Alpha(exec_scopes types.ExecutionScopes) error {
	exec_scopes.AssignOrUpdateVariable("ALPHA", SECP256R1_ALPHA())
	return nil
}

/*
Implements hint:
%{ from starkware.cairo.common.cairo_secp.secp256r1_utils import SECP256R1_N as N %}
*/
func importSECP256R1N(exec_scopes types.ExecutionScopes) error {
	exec_scopes.AssignOrUpdateVariable("N", SECP256R1_N())
	return nil
}

/*
Implements hint:
%{
from starkware.cairo.common.cairo_secp.secp256r1_utils import SECP256R1_P as SECP_P
%}
*/

func importSECP256R1P(exec_scopes types.ExecutionScopes) error {
	exec_scopes.AssignOrUpdateVariable("SECP_P", SECP256R1_P())
	return nil
}

/*
Implements hint:

	%{
	    from starkware.cairo.common.cairo_secp.secp_utils import pack
	    from starkware.python.math_utils import ec_double_slope

	    # Compute the slope.
	    x = pack(ids.point.x, PRIME)
	    y = pack(ids.point.y, PRIME)
	    value = slope = ec_double_slope(point=(x, y), alpha=ALPHA, p=SECP_P)

%}
*/
func computeDoublingSlopeExternalConsts(virtual_machine vm.VirtualMachine, exec_scopes types.ExecutionScopes, ids_data IdsManager) error {
	// ids.point
	point, err := EcPointFromVarName("point", virtual_machine, ids_data)
	if err != nil {
		return err
	}

	secp_p_uncast, err := exec_scopes.Get("SECP_P")
	if err != nil {
		return err
	}
	secp_p := secp_p_uncast.(big.Int)

	alpha_uncast, err := exec_scopes.Get("ALPHA")
	if err != nil {
		return nil
	}

	alpha := alpha_uncast.(big.Int)
	double_point_b := builtins.DoublePointB{X: point.X.Pack86(), Y: point.Y.Pack86()}

	value, err := builtins.EcDoubleSlope(double_point_b, alpha, secp_p)
	if err != nil {
		return err
	}

	exec_scopes.AssignOrUpdateVariable("value", value)
	exec_scopes.AssignOrUpdateVariable("slope", value)
	return nil
}
