package hints

import (
	"math/big"

	. "github.com/lambdaclass/cairo-vm.go/pkg/hints/hint_utils"
	. "github.com/lambdaclass/cairo-vm.go/pkg/types"
	"github.com/lambdaclass/cairo-vm.go/pkg/utils"
	. "github.com/lambdaclass/cairo-vm.go/pkg/vm"
	"github.com/pkg/errors"
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

	if packedB.Cmp(big.NewInt(0)) == 0 {
		return errors.New("Attempted to divide by zero")
	}
	val := new(big.Int).Mod(new(big.Int).Div(&packedA, &packedB), n)

	scopes.AssignOrUpdateVariable("a", packedA)
	scopes.AssignOrUpdateVariable("b", packedB)
	scopes.AssignOrUpdateVariable("val", val)
	scopes.AssignOrUpdateVariable("res", val)

	return nil
}

func divModNPackedDivMod(ids IdsManager, vm *VirtualMachine, scopes *ExecutionScopes) error {
	n, _ := new(big.Int).SetString("115792089237316195423570985008687907852837564279074904382605163141518161494337", 10)
	scopes.AssignOrUpdateVariable("N", n)
	return divModNPacked(ids, vm, scopes, n)
}

func divModNPackedDivModExternalN(ids IdsManager, vm *VirtualMachine, scopes *ExecutionScopes) error {
	nAny, err := scopes.Get("N")
	if err != nil {
		return err
	}
	n, ok := nAny.(*big.Int)
	if !ok {
		return errors.New("N not in scope")
	}
	return divModNPacked(ids, vm, scopes, n)
}

func divModNSafeDiv(ids IdsManager, scopes *ExecutionScopes, aAlias string, bAlias string, toAdd int64) error {
	// Fetch scope variables
	a, err := FetchScopeVar[*big.Int](aAlias, scopes)
	if err != nil {
		return err
	}

	b, err := FetchScopeVar[*big.Int](bAlias, scopes)
	if err != nil {
		return err
	}

	res, err := FetchScopeVar[*big.Int]("res", scopes)
	if err != nil {
		return err
	}

	n, err := FetchScopeVar[*big.Int]("N", scopes)
	if err != nil {
		return err
	}

	// Hint logic
	value, err := utils.SafeDivBig(new(big.Int).Mul(res, new(big.Int).Sub(a, b)), n)
	if err != nil {
		return err
	}
	if toAdd != 0 {
		value = new(big.Int).Add(value, big.NewInt(toAdd))
	}
	// Update scope
	scopes.AssignOrUpdateVariable("value", value)
	return nil
}
