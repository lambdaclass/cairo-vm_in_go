package hints

import (
	"math/big"

	. "github.com/lambdaclass/cairo-vm.go/pkg/hints/hint_utils"
	. "github.com/lambdaclass/cairo-vm.go/pkg/types"
	"github.com/lambdaclass/cairo-vm.go/pkg/utils"
	. "github.com/lambdaclass/cairo-vm.go/pkg/vm"
)

func divModNPacked(ids IdsManager, vm *VirtualMachine, scopes *ExecutionScopes, n *big.Int) error {
	a, err := Uint384FromVarName("a", ids, vm)
	if err != nil {
		return err
	}
	b, err := Uint384FromVarName("b", ids, vm)
	if err != nil {
		return err
	}
	packedA := a.Pack86()
	packedB := b.Pack86()

	val, err := utils.DivMod(&packedA, &packedB, n)
	if err != nil {
		return err
	}

	scopes.AssignOrUpdateVariable("a", packedA)
	scopes.AssignOrUpdateVariable("b", packedB)
	scopes.AssignOrUpdateVariable("value", *val)
	scopes.AssignOrUpdateVariable("res", *val)

	return nil
}

func divModNPackedDivMod(ids IdsManager, vm *VirtualMachine, scopes *ExecutionScopes) error {
	n, _ := new(big.Int).SetString("115792089237316195423570985008687907852837564279074904382605163141518161494337", 10)
	scopes.AssignOrUpdateVariable("N", *n)
	return divModNPacked(ids, vm, scopes, n)
}

func divModNPackedDivModExternalN(ids IdsManager, vm *VirtualMachine, scopes *ExecutionScopes) error {
	n, err := FetchScopeVar[big.Int]("N", scopes)
	if err != nil {
		return err
	}
	return divModNPacked(ids, vm, scopes, &n)
}

func divModNSafeDiv(ids IdsManager, scopes *ExecutionScopes, aAlias string, bAlias string, addOne bool) error {
	// Fetch scope variables
	a, err := FetchScopeVar[big.Int](aAlias, scopes)
	if err != nil {
		return err
	}

	b, err := FetchScopeVar[big.Int](bAlias, scopes)
	if err != nil {
		return err
	}

	res, err := FetchScopeVar[big.Int]("res", scopes)
	if err != nil {
		return err
	}

	n, err := FetchScopeVar[big.Int]("N", scopes)
	if err != nil {
		return err
	}

	// Hint logic
	value, err := utils.SafeDivBig(new(big.Int).Sub(new(big.Int).Mul(&res, &b), &a), &n)
	if err != nil {
		return err
	}
	if addOne {
		value = new(big.Int).Add(value, big.NewInt(1))
	}
	// Update scope
	scopes.AssignOrUpdateVariable("value", *value)
	return nil
}
