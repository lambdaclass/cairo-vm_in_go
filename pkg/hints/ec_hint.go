package hints

import (
	"errors"
	"math/big"

	"github.com/lambdaclass/cairo-vm.go/pkg/hints/hint_utils"
	"github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
	"github.com/lambdaclass/cairo-vm.go/pkg/types"
	"github.com/lambdaclass/cairo-vm.go/pkg/vm"
	"github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
)

type BigInt3 struct {
	limbs []lambdaworks.Felt
}

type EcPoint struct {
	X BigInt3
	Y BigInt3
}

func (val *BigInt3) Pack86() big.Int {
	sum := big.NewInt(0)
	for i := 0; i < 3; i++ {
		felt := val.limbs[i]
		signed := felt.ToSigned()
		shifed := new(big.Int).Lsh(signed, uint(i*86))
		sum.Add(sum, shifed)
	}
	return *sum
}

func FromBaseAddr(addr memory.Relocatable, virtual_machine vm.VirtualMachine) (BigInt3, error) {
	limbs := make([]lambdaworks.Felt, 0)
	for i := 0; i < 3; i++ {
		felt, err := virtual_machine.Segments.Memory.GetFelt(addr.AddUint(uint(i)))
		if err == nil {
			limbs = append(limbs, felt)
		} else {
			return BigInt3{}, errors.New("Identifier has no member")
		}
	}
	return BigInt3{limbs: limbs}, nil
}

func FromVarName(name string, virtual_machine vm.VirtualMachine, ids_data hint_utils.IdsManager) (EcPoint, error) {
	point_addr, err := ids_data.GetAddr(name, &virtual_machine)
	if err != nil {
		return EcPoint{}, err
	}

	x, err := FromBaseAddr(point_addr, virtual_machine)
	if err != nil {
		return EcPoint{}, err
	}

	y, err := FromBaseAddr(point_addr.AddUint(3), virtual_machine)
	if err != nil {
		return EcPoint{}, err
	}

	return EcPoint{X: x, Y: y}, nil
}

/*
Implements main logic for `EC_NEGATE` and `EC_NEGATE_EMBEDDED_SECP` hints
*/
func ec_negate(virtual_machine vm.VirtualMachine, exec_scopes types.ExecutionScopes, ids_data hint_utils.IdsManager, secp_p big.Int) error {
	point, err := ids_data.GetRelocatable("point", &virtual_machine)
	if err != nil {
		return err
	}

	point_y, err := point.AddInt(3)
	if err != nil {
		return err
	}

	y_bigint3, err := FromBaseAddr(point_y, virtual_machine)
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
