package hints_test

import (
	"math/big"
	"testing"

	. "github.com/lambdaclass/cairo-vm.go/pkg/hints"
	. "github.com/lambdaclass/cairo-vm.go/pkg/hints/hint_codes"
	. "github.com/lambdaclass/cairo-vm.go/pkg/hints/hint_utils"
	. "github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
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



// fn run_ec_negate_ok() {
// 	let hint_code = "from starkware.cairo.common.cairo_secp.secp_utils import SECP_P, pack\n\ny = pack(ids.point.y, PRIME) % SECP_P\n# The modulo operation in python always returns a nonnegative number.\nvalue = (-y) % SECP_P";
// 	let mut vm = vm_with_range_check!();

// 	vm.segments = segments![((1, 3), 2645i32), ((1, 4), 454i32), ((1, 5), 206i32)];
// 	//Initialize fp
// 	vm.run_context.fp = 1;
// 	//Create hint_data
// 	let ids_data = ids_data!["point"];
// 	let mut exec_scopes = ExecutionScopes::new();
// 	//Execute the hint
// 	assert_matches!(run_hint!(vm, ids_data, hint_code, &mut exec_scopes), Ok(()));
// 	//Check 'value' is defined in the vm scope
// 	assert_matches!(
// 		exec_scopes.get::<BigInt>("value"),
// 		Ok(x) if x == bigint_str!(
// 			"115792089237316195423569751828682367333329274433232027476421668138471189901786"
// 		)
// 	);
// }

func TestRunEcNegateOk(t *testing.T) {
	vm := NewVirtualMachine()
	vm.Segments.AddSegment()
	vm.Segments.AddSegment()
	vm.Segments.Memory.Insert(NewRelocatable(1,3), NewMaybeRelocatableFelt(FeltFromUint64(2645)))
	vm.Segments.Memory.Insert(NewRelocatable(1,4), NewMaybeRelocatableFelt(FeltFromUint64(454)))
	vm.Segments.Memory.Insert(NewRelocatable(1,5), NewMaybeRelocatableFelt(FeltFromUint64(206)))


	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{
			"point":       {NewMaybeRelocatableFelt(FeltFromDecString("-1"))},
			"ec_negative": {nil},
		},
		vm,
	)
	hintProcessor := CairoVmHintProcessor{}
	hintData := any(HintData{
		Ids:  idsManager,
		Code: EC_NEGATE,
	})
	err := hintProcessor.ExecuteHint(vm, &hintData, nil, nil)
	if err != nil {
		t.Errorf("Ec Negative hint test failed with error %s", err)
	}
	// Check ids.is_positive
	value, err := idsManager.GetFelt("value", vm)
	expected := FeltFromDecString("115792089237316195423569751828682367333329274433232027476421668138471189901786")

	if err != nil || value != expected {
		t.Errorf("Ec Negative hint test incorrect value for ids.value")
	}
}