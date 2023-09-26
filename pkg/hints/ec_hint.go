package hints

import (
	"math/big"

	"github.com/lambdaclass/cairo-vm.go/pkg/builtins"
	"github.com/lambdaclass/cairo-vm.go/pkg/hints/hint_utils"
	. "github.com/lambdaclass/cairo-vm.go/pkg/hints/hint_utils"
	"github.com/lambdaclass/cairo-vm.go/pkg/types"
	. "github.com/lambdaclass/cairo-vm.go/pkg/types"
	"github.com/lambdaclass/cairo-vm.go/pkg/vm"
	. "github.com/lambdaclass/cairo-vm.go/pkg/vm"
)

type EcPoint struct {
	X BigInt3
	Y BigInt3
}

func EcPointFromVarName(name string, vm *VirtualMachine, idsData IdsManager) (EcPoint, error) {
	pointAddr, err := idsData.GetAddr(name, vm)
	if err != nil {
		return EcPoint{}, err
	}

	x, err := BigInt3FromBaseAddr(pointAddr, name+".x", vm)
	if err != nil {
		return EcPoint{}, err
	}

	y, err := BigInt3FromBaseAddr(pointAddr.AddUint(3), name+".y", vm)
	if err != nil {
		return EcPoint{}, err
	}

	return EcPoint{X: x, Y: y}, nil
}

/*
Implements main logic for `EC_NEGATE` and `EC_NEGATE_EMBEDDED_SECP` hints
*/

func ecNegate(vm *vm.VirtualMachine, execScopes types.ExecutionScopes, ids hint_utils.IdsManager, secpP big.Int) error {
	point, err := ids.GetRelocatable("point", vm)
	if err != nil {
		return err
	}

	pointY, err := point.AddInt(3)
	if err != nil {
		return err
	}

	yBigint3, err := BigInt3FromBaseAddr(pointY, "point.y", vm)
	if err != nil {
		return err
	}

	y := yBigint3.Pack86()
	value := new(big.Int).Neg(&y)
	value.Mod(value, &secpP)

	execScopes.AssignOrUpdateVariable("value", value)
	execScopes.AssignOrUpdateVariable("SECP_P", secpP)
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
func computeDoublingSlope(vm *VirtualMachine, execScopes ExecutionScopes, idsData IdsManager, pointAlias string, SecpP big.Int, alpha big.Int) error {
	execScopes.AssignOrUpdateVariable("SECP_P", SecpP)

	point, err := EcPointFromVarName(pointAlias, vm, idsData)
	if err != nil {
		return err
	}

	x := point.X.Pack86()
	y := point.Y.Pack86()
	doublePoint := builtins.DoublePointB{X: x, Y: y}

	value, err := builtins.EcDoubleSlope(doublePoint, alpha, SecpP)
	if err != nil {
		return err
	}

	execScopes.AssignOrUpdateVariable("value", value)
	execScopes.AssignOrUpdateVariable("slope", value)

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

func computeSlopeAndAssingSecpP(vm *VirtualMachine, execScopes ExecutionScopes, idsData IdsManager, point0Alias string, point1Alias string, secpP big.Int) error {
	execScopes.AssignOrUpdateVariable("SECP_P", secpP)
	return computeSlope(vm, execScopes, idsData, point0Alias, point1Alias)
}

func computeSlope(vm *VirtualMachine, execScopes ExecutionScopes, idsData IdsManager, point0Alias string, point1Alias string) error {
	point0, err := EcPointFromVarName(point0Alias, vm, idsData)
	if err != nil {
		return err
	}
	point1, err := EcPointFromVarName(point1Alias, vm, idsData)
	if err != nil {
		return err
	}

	secpP, err := execScopes.Get("SECP_P")
	if err != nil {
		return err
	}
	secp := secpP.(big.Int)

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

	execScopes.AssignOrUpdateVariable("value", value)
	execScopes.AssignOrUpdateVariable("slope", value)

	return nil
}

/*
Implements hint:
%{ from starkware.cairo.common.cairo_secp.secp256r1_utils import SECP256R1_ALPHA as ALPHA %}
*/

func importSecp256r1Alpha(execScopes ExecutionScopes) error {
	execScopes.AssignOrUpdateVariable("ALPHA", SECP256R1_ALPHA())
	return nil
}

/*
Implements hint:
%{ from starkware.cairo.common.cairo_secp.secp256r1_utils import SECP256R1_N as N %}
*/
func importSECP256R1N(execScopes ExecutionScopes) error {
	execScopes.AssignOrUpdateVariable("N", SECP256R1_N())
	return nil
}

/*
Implements hint:
%{
from starkware.cairo.common.cairo_secp.secp256r1_utils import SECP256R1_P as SECP_P
%}
*/

func importSECP256R1P(execScopes ExecutionScopes) error {
	execScopes.AssignOrUpdateVariable("SECP_P", SECP256R1_P())
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
func computeDoublingSlopeExternalConsts(vm VirtualMachine, execScopes ExecutionScopes, ids_data IdsManager) error {
	// ids.point
	point, err := EcPointFromVarName("point", &vm, ids_data)
	if err != nil {
		return err
	}

	secp_p_uncast, err := execScopes.Get("SECP_P")
	if err != nil {
		return err
	}
	secp_p := secp_p_uncast.(big.Int)

	alpha_uncast, err := execScopes.Get("ALPHA")
	if err != nil {
		return nil
	}

	alpha := alpha_uncast.(big.Int)
	double_point_b := builtins.DoublePointB{X: point.X.Pack86(), Y: point.Y.Pack86()}

	value, err := builtins.EcDoubleSlope(double_point_b, alpha, secp_p)
	if err != nil {
		return err
	}

	execScopes.AssignOrUpdateVariable("value", value)
	execScopes.AssignOrUpdateVariable("slope", value)
	return nil
}
