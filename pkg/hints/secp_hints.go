package hints

import (
	"math/big"

	. "github.com/lambdaclass/cairo-vm.go/pkg/hints/hint_utils"
	. "github.com/lambdaclass/cairo-vm.go/pkg/types"
	. "github.com/lambdaclass/cairo-vm.go/pkg/vm"
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
