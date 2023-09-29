package hints

import (
	"errors"
	"math/big"

	"github.com/lambdaclass/cairo-vm.go/pkg/builtins"
	"github.com/lambdaclass/cairo-vm.go/pkg/hints/hint_utils"
	. "github.com/lambdaclass/cairo-vm.go/pkg/hints/hint_utils"
	. "github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
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

// Implements hint:
// from starkware.cairo.common.cairo_secp.secp_utils import SECP_P, pack
//
// slope = pack(ids.slope, PRIME)
// x = pack(ids.point.x, PRIME)
// y = pack(ids.point.y, PRIME)
//
// value = new_x = (pow(slope, 2, SECP_P) - 2 * x) % SECP_P
func ecDoubleAssignNewX(vm *VirtualMachine, execScopes ExecutionScopes, ids IdsManager, secpP big.Int) error {
	execScopes.AssignOrUpdateVariable("SECP_P", secpP)

	slope3, err := BigInt3FromVarName("slope", ids, vm)
	if err != nil {
		return err
	}
	packedSlope := slope3.Pack86()
	slope := new(big.Int).Mod(&packedSlope, Prime())
	point, err := EcPointFromVarName("point", vm, ids)
	if err != nil {
		return err
	}

	xPacked := point.X.Pack86()
	x := new(big.Int).Mod(&xPacked, Prime())
	yPacked := point.Y.Pack86()
	y := new(big.Int).Mod(&yPacked, Prime())

	value := new(big.Int).Mul(slope, slope)
	value = value.Mod(value, &secpP)

	value = value.Sub(value, x)
	value = value.Sub(value, x)
	value = value.Mod(value, &secpP)

	execScopes.AssignOrUpdateVariable("slope", slope)
	execScopes.AssignOrUpdateVariable("x", x)
	execScopes.AssignOrUpdateVariable("y", y)
	execScopes.AssignOrUpdateVariable("value", *value)
	execScopes.AssignOrUpdateVariable("new_x", *value)
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

	secpPuncast, err := execScopes.Get("SECP_P")
	if err != nil {
		return err
	}
	secpP, ok := secpPuncast.(big.Int)
	if !ok {
		return errors.New("Could not cast secp into big int")
	}

	alphaUncast, err := execScopes.Get("ALPHA")
	if err != nil {
		return nil
	}

	alpha, ok := alphaUncast.(big.Int)
	if !ok {
		return errors.New("Could not cast alpha into big int")
	}

	doublePoint_b := builtins.DoublePointB{X: point.X.Pack86(), Y: point.Y.Pack86()}

	value, err := builtins.EcDoubleSlope(doublePoint_b, alpha, secpP)
	if err != nil {
		return err
	}

	execScopes.AssignOrUpdateVariable("value", value)
	execScopes.AssignOrUpdateVariable("slope", value)
	return nil
}

/*
	%{
		from starkware.cairo.common.cairo_secp.secp_utils import SECP_P, pack

		slope = pack(ids.slope, PRIME)
		x0 = pack(ids.point0.x, PRIME)
		x1 = pack(ids.point1.x, PRIME)
		y0 = pack(ids.point0.y, PRIME)

		value = new_x = (pow(slope, 2, SECP_P) - x0 - x1) % SECP_P"

%}
*/
func fastEcAddAssignNewX(ids IdsManager, vm *VirtualMachine, execScopes *ExecutionScopes, point0Alias string, point1Alias string, secpP big.Int) error {
	execScopes.AssignOrUpdateVariable("SECP_P", secpP)

	point0, err := EcPointFromVarName(point0Alias, vm, ids)
	point1, err := EcPointFromVarName(point1Alias, vm, ids)
	if err != nil {
		return err
	}
	slopeUnpacked, err := BigInt3FromVarName("slope", ids, vm)
	if err != nil {
		return err
	}

	slope := slopeUnpacked.Pack86()
	slope = *new(big.Int).Mod(&slope, &secpP)

	x0 := point0.X.Pack86()
	x0 = *new(big.Int).Mod(&x0, &secpP)

	x1 := point1.X.Pack86()
	x1 = *new(big.Int).Mod(&x1, &secpP)

	y0 := point0.Y.Pack86()
	y0 = *new(big.Int).Mod(&y0, &secpP)

	slopeSquared := new(big.Int).Mul(&slope, &slope)
	x0PlusX1 := new(big.Int).Add(&x0, &x1)

	value := *new(big.Int).Sub(slopeSquared, x0PlusX1)
	value = *new(big.Int).Mod(&value, &secpP)

	execScopes.AssignOrUpdateVariable("slope", slope)
	execScopes.AssignOrUpdateVariable("x0", x0)
	execScopes.AssignOrUpdateVariable("y0", y0)
	execScopes.AssignOrUpdateVariable("value", value)
	execScopes.AssignOrUpdateVariable("new_x", value)

	return nil
}

/*
Implements hint:

	%{ value = new_y = (slope * (x0 - new_x) - y0) % SECP_P %}
*/
func fastEcAddAssignNewY(execScopes *ExecutionScopes) error {
	slope, err := execScopes.Get("slope")
	if err != nil {
		return err
	}
	slopeBigInt := slope.(big.Int)
	x0, err := execScopes.Get("x0")
	if err != nil {
		return err
	}
	x0BigInt := x0.(big.Int)

	newX, err := execScopes.Get("new_x")
	if err != nil {
		return err
	}
	newXBigInt := newX.(big.Int)

	y0, err := execScopes.Get("y0")
	if err != nil {
		return err
	}
	y0BigInt := y0.(big.Int)

	secpP, err := execScopes.Get("SECP_P")
	if err != nil {
		return err
	}
	secpBigInt := secpP.(big.Int)

	x0MinusNewX := *new(big.Int).Sub(&x0BigInt, &newXBigInt)
	x0MinusNewXMinusY0 := *new(big.Int).Sub(&x0MinusNewX, &y0BigInt)
	valueBeforeMod := *new(big.Int).Mul(&slopeBigInt, &x0MinusNewXMinusY0)
	value := *new(big.Int).Mod(&valueBeforeMod, &secpBigInt)

	execScopes.AssignOrUpdateVariable("value", value)
	execScopes.AssignOrUpdateVariable("new_y", value)

	return nil
}
