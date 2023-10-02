package hints

import (
	"math/big"

	. "github.com/lambdaclass/cairo-vm.go/pkg/hints/hint_utils"
	"github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
	. "github.com/lambdaclass/cairo-vm.go/pkg/types"
	. "github.com/lambdaclass/cairo-vm.go/pkg/vm"
	"github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
	"github.com/pkg/errors"
)

func reduceV1(ids IdsManager, vm *VirtualMachine, scopes *ExecutionScopes) error {
	secpP := SECP_P()
	scopes.AssignOrUpdateVariable("SECP_P", secpP)
	value, err := Uint384FromVarName("x", ids, vm)
	if err != nil {
		return err
	}
	packedValue := value.Pack86()
	scopes.AssignOrUpdateVariable("value", *new(big.Int).Mod(&packedValue, &secpP))
	return nil
}

func reduceV2(ids IdsManager, vm *VirtualMachine, scopes *ExecutionScopes) error {
	secpP, err := FetchScopeVar[big.Int]("SECP_P", scopes)
	if err != nil {
		return err
	}
	value, err := Uint384FromVarName("x", ids, vm)
	if err != nil {
		return err
	}
	packedValue := value.Pack86()
	scopes.AssignOrUpdateVariable("value", *new(big.Int).Mod(&packedValue, &secpP))
	return nil
}

func reduceED25519(ids IdsManager, vm *VirtualMachine, scopes *ExecutionScopes) error {
	secpP := SECP_P_V2()
	scopes.AssignOrUpdateVariable("SECP_P", secpP)
	value, err := Uint384FromVarName("x", ids, vm)
	if err != nil {
		return err
	}
	packedValue := value.Pack86()
	scopes.AssignOrUpdateVariable("value", *new(big.Int).Mod(&packedValue, &secpP))
	return nil
}

func verifyZero(ids IdsManager, vm *VirtualMachine, scopes *ExecutionScopes, secpP big.Int) error {
	scopes.AssignOrUpdateVariable("SECP_P", secpP)
	valUnpacked, err := Uint384FromVarName("val", ids, vm)
	if err != nil {
		return err
	}
	val := valUnpacked.Pack86()
	q, r := new(big.Int).DivMod(&val, &secpP, new(big.Int))
	if r.Cmp(big.NewInt(0)) != 0 {
		return errors.Errorf("verify_zero: Invalid input %s", val.Text(10))
	}
	return ids.Insert("q", memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromBigInt(q)), vm)
}

// ids.low = (ids.x.d0 + ids.x.d1 * ids.BASE) & ((1 << 128) - 1)
func bigintToUint256(ids IdsManager, vm *VirtualMachine, constants *map[string]lambdaworks.Felt) error {
	// Fetch variables
	xD0, err := ids.GetStructFieldFelt("x", 0, vm)
	if err != nil {
		return err
	}
	xD1, err := ids.GetStructFieldFelt("x", 1, vm)
	if err != nil {
		return err
	}
	base, err := ids.GetConst("BASE", constants)
	if err != nil {
		return err
	}
	// Hint Logic
	low := xD0.Add(xD1.Mul(base)).And((lambdaworks.FeltOne().Shl(128)).Sub(lambdaworks.FeltOne()))
	return ids.Insert("low", memory.NewMaybeRelocatableFelt(low), vm)
}

func isZeroNondet(ids IdsManager, vm *VirtualMachine) error {
	x, err := ids.GetFelt("x", vm)
	if err != nil {
		return err
	}
	if x.IsZero() {
		return vm.Segments.Memory.Insert(vm.RunContext.Ap, memory.NewMaybeRelocatableFelt(lambdaworks.FeltOne()))
	}
	return vm.Segments.Memory.Insert(vm.RunContext.Ap, memory.NewMaybeRelocatableFelt(lambdaworks.FeltZero()))
}
