package hints_test

import (
	"fmt"
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
	execScopes := types.NewExecutionScopes()
	err := hintProcessor.ExecuteHint(vm, &hintData, nil, execScopes)
	if err != nil {
		t.Errorf("Ec Negative hint test failed with error %s", err)
	} else {
		// Check ids.is_positive
		value, err := execScopes.Get("value")
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

func TestComputeDoublingSlopeOk(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()
	vm.Segments.AddSegment()

	vm.RunContext.Fp = NewRelocatable(1, 1)

	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{
			"point": {
				NewMaybeRelocatableFelt(FeltFromUint64(614323)),
				NewMaybeRelocatableFelt(FeltFromUint64(5456867)),
				NewMaybeRelocatableFelt(FeltFromUint64(101208)),
				NewMaybeRelocatableFelt(FeltFromUint64(773712524)),
				NewMaybeRelocatableFelt(FeltFromUint64(77371252)),
				NewMaybeRelocatableFelt(FeltFromUint64(5298795)),
			},
		},
		vm,
	)

	hintProcessor := CairoVmHintProcessor{}
	hintData := any(HintData{
		Ids:  idsManager,
		Code: EC_DOUBLE_SLOPE_V1,
	})

	exec_scopes := types.NewExecutionScopes()
	err := hintProcessor.ExecuteHint(vm, &hintData, nil, exec_scopes)
	if err != nil {
		t.Errorf("EC_DOUBLE_SLOPE_V1 hint test failed with error %s", err)
	} else {
		value, _ := exec_scopes.Get("value")
		val := value.(big.Int)

		slope_res, _ := exec_scopes.Get("slope")
		slope := slope_res.(big.Int)

		// expected values
		expectedVal, _ := new(big.Int).SetString("40442433062102151071094722250325492738932110061897694430475034100717288403728", 10)

		expectedSlope, _ := new(big.Int).SetString("40442433062102151071094722250325492738932110061897694430475034100717288403728", 10)

		if expectedVal.Cmp(&val) != 0 || expectedSlope.Cmp(&slope) != 0 {
			t.Errorf("EC_DOUBLE_SLOPE_V1 hint test incorrect value for exec_scopes.value or exec_scopes.slope")
		}
	}
}

func TestRunComputeSlopeOk(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()
	vm.Segments.AddSegment()

	vm.RunContext.Fp = NewRelocatable(1, 14)

	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{
			"point0": {
				NewMaybeRelocatableFelt(FeltFromUint64(134)),
				NewMaybeRelocatableFelt(FeltFromUint64(5123)),
				NewMaybeRelocatableFelt(FeltFromUint64(140)),
				NewMaybeRelocatableFelt(FeltFromUint64(1232)),
				NewMaybeRelocatableFelt(FeltFromUint64(4652)),
				NewMaybeRelocatableFelt(FeltFromUint64(720)),
			},
			"point1": {
				NewMaybeRelocatableFelt(FeltFromUint64(156)),
				NewMaybeRelocatableFelt(FeltFromUint64(6545)),
				NewMaybeRelocatableFelt(FeltFromUint64(100010)),
				NewMaybeRelocatableFelt(FeltFromUint64(1123)),
				NewMaybeRelocatableFelt(FeltFromUint64(1325)),
				NewMaybeRelocatableFelt(FeltFromUint64(910)),
			},
		},
		vm,
	)

	hintProcessor := CairoVmHintProcessor{}
	hintData := any(HintData{
		Ids:  idsManager,
		Code: COMPUTE_SLOPE_V1,
	})

	execScopes := types.NewExecutionScopes()
	err := hintProcessor.ExecuteHint(vm, &hintData, nil, execScopes)
	if err != nil {
		t.Errorf("EC_DOUBLE_SLOPE_V1 hint test failed with error %s", err)
	} else {
		value, _ := execScopes.Get("value")
		val := value.(big.Int)

		slope_res, _ := execScopes.Get("slope")
		slope := slope_res.(big.Int)
		fmt.Println("SLOPE: ", slope)

		// expected values
		expectedVal, _ := new(big.Int).SetString("41419765295989780131385135514529906223027172305400087935755859001910844026631", 10)

		expectedSlope, _ := new(big.Int).SetString("41419765295989780131385135514529906223027172305400087935755859001910844026631", 10)

		if expectedVal.Cmp(&val) != 0 || expectedSlope.Cmp(&slope) != 0 {
			t.Errorf("EC_DOUBLE_SLOPE_V1 hint test incorrect value for exec_scopes.value or exec_scopes.slope")
		}
	}
}

func TestFastEcAddAssignNewXHint(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()
	vm.Segments.AddSegment()

	vm.RunContext.Fp = NewRelocatable(1, 14)

	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{
			"point0": {
				NewMaybeRelocatableFelt(FeltFromUint64(1)),
				NewMaybeRelocatableFelt(FeltFromUint64(2)),
				NewMaybeRelocatableFelt(FeltFromUint64(3)),
				NewMaybeRelocatableFelt(FeltFromUint64(4)),
				NewMaybeRelocatableFelt(FeltFromUint64(5)),
				NewMaybeRelocatableFelt(FeltFromUint64(6)),
			},
			"point1": {
				NewMaybeRelocatableFelt(FeltFromUint64(7)),
				NewMaybeRelocatableFelt(FeltFromUint64(8)),
				NewMaybeRelocatableFelt(FeltFromUint64(9)),
				NewMaybeRelocatableFelt(FeltFromUint64(10)),
				NewMaybeRelocatableFelt(FeltFromUint64(11)),
				NewMaybeRelocatableFelt(FeltFromUint64(12)),
			},
			"slope": {
				NewMaybeRelocatableFelt(FeltFromUint64(1)),
				NewMaybeRelocatableFelt(FeltFromUint64(0)),
				NewMaybeRelocatableFelt(FeltFromUint64(0)),
			},
		},
		vm,
	)

	hintProcessor := CairoVmHintProcessor{}
	hintData := any(HintData{
		Ids:  idsManager,
		Code: FAST_EC_ADD_ASSIGN_NEW_X_V2,
	})

	execScopes := types.NewExecutionScopes()
	err := hintProcessor.ExecuteHint(vm, &hintData, nil, execScopes)
	if err != nil {
		t.Errorf("FAST_EC_ADD_ASSIGN_NEW_X_V2 hint test failed with error %s", err)
	}

	slope, _ := execScopes.Get("slope")
	slopeRes := slope.(big.Int)

	x0, _ := execScopes.Get("x0")
	x0Res := x0.(big.Int)

	y0, _ := execScopes.Get("y0")
	y0Res := y0.(big.Int)

	value, _ := execScopes.Get("value")
	valueRes := value.(big.Int)

	// expected values
	expectedSlope, _ := new(big.Int).SetString("1", 10)
	expectedX0, _ := new(big.Int).SetString("17958932119522135058886879379160190656204633450479617", 10)
	expectedY0, _ := new(big.Int).SetString("35917864239044270117773758835691633767745534082154500", 10)
	expectedVal, _ := new(big.Int).SetString("57896044618658097711785420668615475838094756785302610636461256512888400510950", 10)

	if expectedVal.Cmp(&valueRes) != 0 || expectedSlope.Cmp(&slopeRes) != 0 {
		t.Errorf("EC_DOUBLE_SLOPE_V1 hint test incorrect value for exec_scopes.value or exec_scopes.slope")
	}

	if expectedX0.Cmp(&x0Res) != 0 {
		t.Errorf("expected x0=%v, got: x0=%v", expectedX0, x0Res)
	}

	if expectedY0.Cmp(&y0Res) != 0 {
		t.Errorf("expected y0 to be %v, got: y0=%v", expectedY0, y0Res)
	}

	/*
		X0 AFTER PACK: 17958932119522135058886879379160190656204633450479617
		X1 AFTER PACK: 53876796358566405176660638292223076879286434713829383
		Y0 AFTER PACK: 35917864239044270117773758835691633767745534082154500
		VALUE: 57896044618658097711785420668615475838094756785302610636461256512888400510950
	*/
}
