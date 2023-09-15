package hints_test

import (
	"math/big"
	"testing"

	"github.com/lambdaclass/cairo-vm.go/pkg/hints"
	"github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
)

func TestBigInt3Pack86(t *testing.T) {
	limbs1 := []lambdaworks.Felt{lambdaworks.FeltFromUint64(10), lambdaworks.FeltFromUint64(10), lambdaworks.FeltFromUint64(10)}
	bigint := hints.BigInt3{Limbs: limbs1}
	pack1 := bigint.Pack86()

	expected, _ := new(big.Int).SetString("59863107065073783529622931521771477038469668772249610", 10)

	if pack1.Cmp(expected) != 0 {
		t.Errorf("Different pack from expected")
	}

	limbs2 := []lambdaworks.Felt{lambdaworks.FeltFromDecString("773712524553362"), lambdaworks.FeltFromDecString("57408430697461422066401280"), lambdaworks.FeltFromDecString("1292469707114105")}
	bigint2 := hints.BigInt3{Limbs: limbs2}
	pack2 := bigint2.Pack86()

	expected2, _ := new(big.Int).SetString("7737125245533626718119526477371252455336267181195264773712524553362", 10)

	if pack2.Cmp(expected2) != 0 {
		t.Errorf("Different pack from expected2")
	}
}