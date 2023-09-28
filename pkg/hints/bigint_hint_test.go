package hints_test

import (
	"math/big"
	"testing"

	. "github.com/lambdaclass/cairo-vm.go/pkg/hints"
	. "github.com/lambdaclass/cairo-vm.go/pkg/hints/hint_codes"
	. "github.com/lambdaclass/cairo-vm.go/pkg/hints/hint_utils"
	"github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
	. "github.com/lambdaclass/cairo-vm.go/pkg/types"
	. "github.com/lambdaclass/cairo-vm.go/pkg/vm"
	"github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
	. "github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
)

func TestNonDetBigInt3Ok(t *testing.T) {
	vm := NewVirtualMachine()

	vm.Segments.AddSegment()
	vm.Segments.AddSegment()
	vm.Segments.AddSegment()

	vm.RunContext.Pc = NewRelocatable(0, 0)
	vm.RunContext.Ap = NewRelocatable(1, 6)
	vm.RunContext.Fp = NewRelocatable(1, 6)

	value, _ := new(big.Int).SetString("7737125245533626718119526477371252455336267181195264773712524553362", 10)
	execScopes := NewExecutionScopes()

	execScopes.AssignOrUpdateVariable("value", *value)

	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{
			"res": {nil},
		},
		vm,
	)

	hintProcessor := CairoVmHintProcessor{}
	hintData := any(HintData{
		Ids:  idsManager,
		Code: NONDET_BIGINT3_V1,
	})
	err := hintProcessor.ExecuteHint(vm, &hintData, nil, execScopes)
	if err != nil {
		t.Errorf("Non Det Big Int 3 hint test failed with error: %s", err)
	} else {
		valueInStruct0, err := idsManager.GetStructFieldFelt("res", 0, vm)
		expected0 := lambdaworks.FeltFromDecString("773712524553362")
		if err != nil {
			t.Errorf("error fetching from ids manager : %s", err)
		}
		if valueInStruct0 != expected0 {
			t.Errorf(" Incorrect field value %s, expected %s", valueInStruct0.ToBigInt().Text(10), expected0.ToBigInt().Text(10))
		}

		valueInStruct1, err := idsManager.GetStructFieldFelt("res", 1, vm)
		expected1 := lambdaworks.FeltFromDecString("57408430697461422066401280")
		if err != nil {
			t.Errorf("error fetching from ids manager : %s", err)
		}
		if valueInStruct1 != expected1 {
			t.Errorf(" Incorrect field value %s, expected %s", valueInStruct1.ToBigInt().Text(10), expected1.ToBigInt().Text(10))
		}

		valueInStruct2, err := idsManager.GetStructFieldFelt("res", 2, vm)
		expected2 := lambdaworks.FeltFromDecString("1292469707114105")
		if err != nil {
			t.Errorf("error fetching from ids manager : %s", err)
		}
		if valueInStruct2 != expected2 {
			t.Errorf(" Incorrect field value %s, expected %s", valueInStruct2.ToBigInt().Text(10), expected2.ToBigInt().Text(10))
		}
	}
}

func TestSafeDivBigInt(t *testing.T) {
	vm := NewVirtualMachine()

	vm.Segments.AddSegment()
	vm.Segments.AddSegment()
	vm.Segments.AddSegment()

	execScopes := NewExecutionScopes()

	res, _ := new(big.Int).SetString("109567829260688255124154626727441144629993228404337546799996747905569082729709", 10)
	x, _ := new(big.Int).SetString("91414600319290532004473480113251693728834511388719905794310982800988866814583", 10)
	y, _ := new(big.Int).SetString("38047400353360331012910998489219098987968251547384484838080352663220422975266", 10)
	p, _ := new(big.Int).SetString("115792089237316195423570985008687907852837564279074904382605163141518161494337", 10)

	execScopes.AssignOrUpdateVariable("res", *res)
	execScopes.AssignOrUpdateVariable("x", *x)
	execScopes.AssignOrUpdateVariable("y", *y)
	execScopes.AssignOrUpdateVariable("p", *p)

	vm.RunContext.Fp = NewRelocatable(1, 0)
	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{
			"flag": {nil},
		},
		vm,
	)

	hintProcessor := CairoVmHintProcessor{}
	hintData := any(HintData{
		Ids:  idsManager,
		Code: BIGINT_SAFE_DIV,
	})

	err := hintProcessor.ExecuteHint(vm, &hintData, nil, execScopes)

	if err != nil {
		t.Errorf("Safe Big int div hint test failed with error: %s", err)
	} else {
		expectedK, _ := new(big.Int).SetString("36002209591245282109880156842267569109802494162594623391338581162816748840003", 10)
		expectedVal, _ := new(big.Int).SetString("36002209591245282109880156842267569109802494162594623391338581162816748840003", 10)

		kUncast, err := execScopes.Get("k")
		if err != nil {
			t.Errorf("%s", err)
		}
		k, _ := kUncast.(big.Int)

		valUncast, err := execScopes.Get("value")
		if err != nil {
			t.Errorf("%s", err)
		}
		value, _ := valUncast.(big.Int)

		if expectedK.Cmp(&k) != 0 {
			t.Errorf("incorrect K value expected: %s, got: %s", expectedK.Text(10), k.Text(10))
		}

		if expectedVal.Cmp(&value) != 0 {
			t.Errorf("incorrect value expected: %s, got: %s", expectedVal.Text(10), value.Text(10))
		}

		// check memory
		addr := memory.NewRelocatable(1, 0)
		val, _ := vm.Segments.Memory.GetFelt(addr)

		if val != lambdaworks.FeltFromUint(1) {
			t.Errorf("incorrect value fetched from memory")
		}
	}
}

func TestBigintPackDivModHint(t *testing.T) {
	vm := NewVirtualMachine()

	vm.Segments.AddSegment()
	vm.Segments.AddSegment()
	vm.Segments.AddSegment()
	vm.Segments.AddSegment()
	vm.Segments.AddSegment()
	vm.Segments.AddSegment()
	vm.Segments.AddSegment()
	vm.Segments.AddSegment()
	vm.Segments.AddSegment()
	vm.Segments.AddSegment()
	vm.Segments.AddSegment()

	vm.RunContext.Fp = NewRelocatable(1, 0)
	idsManager := SetupIdsForTest(
		map[string][]*MaybeRelocatable{
			"x": {
				NewMaybeRelocatableFelt(lambdaworks.FeltFromHex("0x38a23ca66202c8c2a72277")),
				NewMaybeRelocatableFelt(lambdaworks.FeltFromHex("0x6730e765376ff17ea8385")),
				NewMaybeRelocatableFelt(lambdaworks.FeltFromUint(0)),
				NewMaybeRelocatableFelt(lambdaworks.FeltFromUint(0)),
			},
			"y": {
				NewMaybeRelocatableFelt(lambdaworks.FeltFromHex("0x20a4b46d3c5e24cda81f22")),
				NewMaybeRelocatableFelt(lambdaworks.FeltFromHex("0x967bf895824330d4273d0")),
				NewMaybeRelocatableFelt(lambdaworks.FeltFromHex("0x541e10c21560da25ada4c")),
			},
			"P": {
				NewMaybeRelocatableFelt(lambdaworks.FeltFromHex("0x8a03bbfd25e8cd0364141")),
				NewMaybeRelocatableFelt(lambdaworks.FeltFromHex("0x3ffffffffffaeabb739abd")),
				NewMaybeRelocatableFelt(lambdaworks.FeltFromHex("0xfffffffffffffffffffff")),
			},
		},
		vm,
	)

	execScopes := NewExecutionScopes()

	hintProcessor := CairoVmHintProcessor{}
	hintData := any(HintData{
		Ids:  idsManager,
		Code: BIGINT_PACK_DIV_MOD,
	})

	err := hintProcessor.ExecuteHint(vm, &hintData, nil, execScopes)

	if err != nil {
		t.Errorf("BIGINT_PACK_DIV_MOD test failed with error: %s", err)
	} else {
		expected, _ := new(big.Int).SetString("109567829260688255124154626727441144629993228404337546799996747905569082729709", 10)
		yExpected, _ := new(big.Int).SetString("38047400353360331012910998489219098987968251547384484838080352663220422975266", 10)
		xExpected, _ := new(big.Int).SetString("91414600319290532004473480113251693728834511388719905794310982800988866814583", 10)
		pExpected, _ := new(big.Int).SetString("115792089237316195423570985008687907852837564279074904382605163141518161494337", 10)

		// fetch values
		resu, err := execScopes.Get("res")
		if err != nil {
			t.Errorf("BIGINT_PACK_DIV_MOD test failed with error: %s", err)
		}
		res := resu.(big.Int)

		valu, err := execScopes.Get("value")
		if err != nil {
			t.Errorf("BIGINT_PACK_DIV_MOD test failed with error: %s", err)
		}
		value := valu.(big.Int)

		yu, err := execScopes.Get("y")
		if err != nil {
			t.Errorf("BIGINT_PACK_DIV_MOD test failed with error: %s", err)
		}
		y := yu.(big.Int)

		xu, err := execScopes.Get("x")
		if err != nil {
			t.Errorf("BIGINT_PACK_DIV_MOD test failed with error: %s", err)
		}
		x := xu.(big.Int)

		pu, err := execScopes.Get("p")
		if err != nil {
			t.Errorf("BIGINT_PACK_DIV_MOD test failed with error: %s", err)
		}
		p := pu.(big.Int)

		if res.Cmp(expected) != 0 {
			t.Errorf("incorrect res expected: %s, got: %s", expected.Text(10), expected.Text(10))
		}
		if value.Cmp(expected) != 0 {
			t.Errorf("incorrect value expected: %s, got: %s", value.Text(10), expected.Text(10))
		}
		if y.Cmp(yExpected) != 0 {
			t.Errorf("incorrect y expected: %s, got: %s", yExpected.Text(10), y.Text(10))
		}
		if x.Cmp(xExpected) != 0 {
			t.Errorf("incorrect x expected: %s, got: %s", xExpected.Text(10), x.Text(10))
		}
		if p.Cmp(pExpected) != 0 {
			t.Errorf("incorrect p expected: %s, got: %s", pExpected.Text(10), p.Text(10))
		}
	}
}
