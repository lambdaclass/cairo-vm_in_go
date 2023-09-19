package hints_test

import (
	"math/big"
	"testing"

	. "github.com/lambdaclass/cairo-vm.go/pkg/hints"
	. "github.com/lambdaclass/cairo-vm.go/pkg/hints/hint_codes"
	. "github.com/lambdaclass/cairo-vm.go/pkg/hints/hint_utils"
	. "github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
	"github.com/lambdaclass/cairo-vm.go/pkg/types"
	. "github.com/lambdaclass/cairo-vm.go/pkg/vm"
	. "github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
)

func TestBigInt3Pack86(t *testing.T) {
	limbs1 := []Felt{FeltFromUint64(10), FeltFromUint64(10), FeltFromUint64(10)}
	bigint := BigInt3{Limbs: limbs1}
	pack1 := bigint.Pack86()

	expected, _ := new(big.Int).SetString("59863107065073783529622931521771477038469668772249610", 10)

	if pack1.Cmp(expected) != 0 {
		t.Errorf("Different pack from expected")
	}

	limbs2 := []Felt{FeltFromDecString("773712524553362"), FeltFromDecString("57408430697461422066401280"), FeltFromDecString("1292469707114105")}
	bigint2 := BigInt3{Limbs: limbs2}
	pack2 := bigint2.Pack86()

	expected2, _ := new(big.Int).SetString("7737125245533626718119526477371252455336267181195264773712524553362", 10)

	if pack2.Cmp(expected2) != 0 {
		t.Errorf("Different pack from expected2")
	}
}

func TestRunEcNegateOk(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()
	vm.Segments.AddSegment()
	vm.Segments.Memory.Insert(NewRelocatable(1, 3), NewMaybeRelocatableFelt(FeltFromUint64(2645)))
	vm.Segments.Memory.Insert(NewRelocatable(1, 4), NewMaybeRelocatableFelt(FeltFromUint64(454)))
	vm.Segments.Memory.Insert(NewRelocatable(1, 5), NewMaybeRelocatableFelt(FeltFromUint64(206)))

	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{
			"point":       {NewMaybeRelocatableRelocatable(NewRelocatable(1, 0))},
			"ec_negative": {nil},
		},
		vm,
	)
	hintProcessor := CairoVmHintProcessor{}
	hintData := any(HintData{
		Ids:  idsManager,
		Code: EC_NEGATE,
	})
	exec_scopes := types.NewExecutionScopes()
	err := hintProcessor.ExecuteHint(vm, &hintData, nil, exec_scopes)
	if err != nil {
		t.Errorf("Ec Negative hint test failed with error %s", err)
	} else {
		// Check ids.is_positive
		value, err := exec_scopes.Get("value")
		val := value.(*big.Int)
		expected, _ := new(big.Int).SetString("115792089237316195423569751828682367333329274433232027476421668138471189901786", 10)

		if err != nil || expected.Cmp(val) != 0 {
			t.Errorf("Ec Negative hint test incorrect value for exec_scopes.value")
		}
	}
}

func TestRunEcEmbeddedSecpOk(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()
	vm.Segments.AddSegment()
	vm.Segments.Memory.Insert(NewRelocatable(1, 3), NewMaybeRelocatableFelt(FeltFromUint64(2645)))
	vm.Segments.Memory.Insert(NewRelocatable(1, 4), NewMaybeRelocatableFelt(FeltFromUint64(454)))
	vm.Segments.Memory.Insert(NewRelocatable(1, 5), NewMaybeRelocatableFelt(FeltFromUint64(206)))

	y2 := big.NewInt(206)
	y2.Lsh(y2, 86*2)

	y1 := big.NewInt(454)
	y1.Lsh(y1, 86)

	y0 := big.NewInt(2645)

	y := new(big.Int)
	y.Add(y, y2)
	y.Add(y, y1)
	y.Add(y, y0)

	vm.RunContext.Fp = NewRelocatable(1, 1)

	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{
			"point":       {NewMaybeRelocatableRelocatable(NewRelocatable(1, 0))},
			"ec_negative": {nil},
		},
		vm,
	)

	hintProcessor := CairoVmHintProcessor{}
	hintData := any(HintData{
		Ids:  idsManager,
		Code: EC_NEGATE_EMBEDDED_SECP,
	})
	exec_scopes := types.NewExecutionScopes()
	err := hintProcessor.ExecuteHint(vm, &hintData, nil, exec_scopes)
	if err != nil {
		t.Errorf("Ec Negative Embedded Sec hint test failed with error %s", err)
	} else {
		// Check ids.is_positive
		value, err := exec_scopes.Get("value")
		val := value.(*big.Int)

		// expected value
		minus_y := big.NewInt(1)
		minus_y.Lsh(minus_y, 255)
		minus_y.Sub(minus_y, big.NewInt(19))
		minus_y.Sub(minus_y, y)

		if err != nil || minus_y.Cmp(val) != 0 {
			t.Errorf("Ec Negative hint test incorrect value for exec_scopes.value")
		}
	}

}
