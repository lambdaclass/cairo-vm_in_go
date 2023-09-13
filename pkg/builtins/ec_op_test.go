package builtins_test

import (
	"math/big"
	"testing"

	"github.com/lambdaclass/cairo-vm.go/pkg/builtins"
	"github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
	"github.com/lambdaclass/cairo-vm.go/pkg/vm/cairo_run"
	"github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
)

func TestPointIsOnCurveB(t *testing.T) {
	x := lambdaworks.FeltFromDecString("3139037544796708144595053687182055617920475701120786241351436619796497072089")
	y := lambdaworks.FeltFromDecString("2119589567875935397690285099786081818522144748339117565577200220779667999801")
	alpha := lambdaworks.FeltOne()
	beta := lambdaworks.FeltFromDecString("3141592653589793238462643383279502884197169399375105820974944592307816406665")

	if !builtins.PointOnCurve(x, y, alpha, beta) {
		t.Errorf("The point is not on the curve")
	}
}

func TestPointIsNotOnCurveB(t *testing.T) {
	x := lambdaworks.FeltFromDecString("3139037544756708144595053687182055617927475701120786241351436619796497072089")
	y := lambdaworks.FeltFromDecString("2119589567875935397690885099786081818522144748339117565577200220779667999801")
	alpha := lambdaworks.FeltOne()
	beta := lambdaworks.FeltFromDecString("3141592653589793238462643383279502884197169399375105820974944592307816406665")

	if builtins.PointOnCurve(x, y, alpha, beta) {
		t.Errorf("The point should not be on the curve")
	}
}

func TestComputeEcOpImplValidA(t *testing.T) {
	partial_sum_x := lambdaworks.FeltFromDecString("3139037544796708144595053687182055617920475701120786241351436619796497072089")
	partial_sum_y := lambdaworks.FeltFromDecString("2119589567875935397690285099786081818522144748339117565577200220779667999801")
	partial_sum := builtins.PartialSum{X: partial_sum_x, Y: partial_sum_y}

	double_point_x := lambdaworks.FeltFromDecString("874739451078007766457464989774322083649278607533249481151382481072868806602")
	double_point_y := lambdaworks.FeltFromDecString("152666792071518830868575557812948353041420400780739481342941381225525861407")

	double_point := builtins.DoublePoint{X: double_point_x, Y: double_point_y}

	m := lambdaworks.FeltFromUint64(34)
	alpha := big.NewInt(1)
	heigth := 256
	const PRIME = "800000000000011000000000000000000000000000000000000000000000001"
	prime, _ := new(big.Int).SetString(PRIME, 16)

	result, err := builtins.EcOnImpl(partial_sum, double_point, m, alpha, prime, uint32(heigth))

	if err != nil {
		t.Errorf("Error computing Ec on Impl")
	}

	result_x := result.X
	result_y := result.Y

	expected_x, _ := new(big.Int).SetString("1977874238339000383330315148209250828062304908491266318460063803060754089297", 10)
	expected_y, _ := new(big.Int).SetString("2969386888251099938335087541720168257053975603483053253007176033556822156706", 10)

	if result_x.Cmp(expected_x) != 0 {
		t.Errorf("Got different X result in Ec On Impl")
	}

	if result_y.Cmp(expected_y) != 0 {
		t.Errorf("Got different Y result in Ec On Impl")
	}
}

func TestComputeEcOpImplValidB(t *testing.T) {
	partial_sum_x := lambdaworks.FeltFromDecString("2962412995502985605007699495352191122971573493113767820301112397466445942584")
	partial_sum_y := lambdaworks.FeltFromDecString("214950771763870898744428659242275426967582168179217139798831865603966154129")
	partial_sum := builtins.PartialSum{X: partial_sum_x, Y: partial_sum_y}

	double_point_x := lambdaworks.FeltFromDecString("874739451078007766457464989774322083649278607533249481151382481072868806602")
	double_point_y := lambdaworks.FeltFromDecString("152666792071518830868575557812948353041420400780739481342941381225525861407")

	double_point := builtins.DoublePoint{X: double_point_x, Y: double_point_y}

	m := lambdaworks.FeltFromUint64(34)
	alpha := big.NewInt(1)
	heigth := 256
	const PRIME = "800000000000011000000000000000000000000000000000000000000000001"
	prime, _ := new(big.Int).SetString(PRIME, 16)

	result, err := builtins.EcOnImpl(partial_sum, double_point, m, alpha, prime, uint32(heigth))

	if err != nil {
		t.Errorf("Error computing Ec on Impl")
	}

	result_x := result.X
	result_y := result.Y

	expected_x, _ := new(big.Int).SetString("2778063437308421278851140253538604815869848682781135193774472480292420096757", 10)
	expected_y, _ := new(big.Int).SetString("3598390311618116577316045819420613574162151407434885460365915347732568210029", 10)

	if result_x.Cmp(expected_x) != 0 {
		t.Errorf("Got different X result in Ec On Impl")
	}

	if result_y.Cmp(expected_y) != 0 {
		t.Errorf("Got different Y result in Ec On Impl")
	}
}

func TestComputeEcOpInvalidSameXCoordinate(t *testing.T) {
	partial_sum_x := lambdaworks.FeltOne()
	partial_sum_y := lambdaworks.FeltFromUint64(9)
	partial_sum := builtins.PartialSum{X: partial_sum_x, Y: partial_sum_y}

	double_point_x := lambdaworks.FeltOne()
	double_point_y := lambdaworks.FeltFromUint64(12)
	double_point := builtins.DoublePoint{X: double_point_x, Y: double_point_y}

	m := lambdaworks.FeltFromUint64(34)
	alpha := big.NewInt(1)
	heigth := 256
	const PRIME = "800000000000011000000000000000000000000000000000000000000000001"
	prime, _ := new(big.Int).SetString(PRIME, 16)

	_, err := builtins.EcOnImpl(partial_sum, double_point, m, alpha, prime, uint32(heigth))

	if err == nil {
		t.Errorf("Expected error but got None")
	}

}

func TestDeduceMemoryCellEcOpForPresetMemoryValid(t *testing.T) {
	mem := memory.NewMemorySegmentManager()
	mem.AddSegment()
	mem.AddSegment()
	mem.AddSegment()
	mem.AddSegment()
	mem.Memory.Insert(memory.NewRelocatable(3, 0), memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromDecString("2962412995502985605007699495352191122971573493113767820301112397466445942584")))
	mem.Memory.Insert(memory.NewRelocatable(3, 1), memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromDecString("214950771763870898744428659242275426967582168179217139798831865603966154129")))
	mem.Memory.Insert(memory.NewRelocatable(3, 2), memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromDecString("874739451078007766457464989774322083649278607533249481151382481072868806602")))
	mem.Memory.Insert(memory.NewRelocatable(3, 3), memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromDecString("152666792071518830868575557812948353041420400780739481342941381225525861407")))
	mem.Memory.Insert(memory.NewRelocatable(3, 4), memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(34)))
	mem.Memory.Insert(memory.NewRelocatable(3, 5), memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromDecString("2778063437308421278851140253538604815869848682781135193774472480292420096757")))

	builtin := builtins.NewEcOpBuiltinRunner(true)

	// expected value
	felt := lambdaworks.FeltFromDecString("3598390311618116577316045819420613574162151407434885460365915347732568210029")
	expected := memory.NewMaybeRelocatableFelt(felt)

	result, err := builtin.DeduceMemoryCell(memory.NewRelocatable(3, 6), &mem.Memory)

	if err != nil {
		t.Errorf("Error calculating deduced memory cell")
	}

	if *result != *expected {
		t.Errorf("Error: Results differ from expected")
	}
}

func TestDeduceMemoryCellEcOpForPresetMemoryUnfilledInputCells(t *testing.T) {
	mem := memory.NewMemorySegmentManager()
	mem.AddSegment()
	mem.AddSegment()
	mem.AddSegment()
	mem.AddSegment()
	mem.Memory.Insert(memory.NewRelocatable(3, 1), memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromDecString("214950771763870898744428659242275426967582168179217139798831865603966154129")))
	mem.Memory.Insert(memory.NewRelocatable(3, 2), memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromDecString("874739451078007766457464989774322083649278607533249481151382481072868806602")))
	mem.Memory.Insert(memory.NewRelocatable(3, 3), memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromDecString("152666792071518830868575557812948353041420400780739481342941381225525861407")))
	mem.Memory.Insert(memory.NewRelocatable(3, 4), memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(34)))
	mem.Memory.Insert(memory.NewRelocatable(3, 5), memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromDecString("2778063437308421278851140253538604815869848682781135193774472480292420096757")))

	builtin := builtins.NewEcOpBuiltinRunner(true)

	// expected value

	result, err := builtin.DeduceMemoryCell(memory.NewRelocatable(3, 6), &mem.Memory)

	if err != nil {
		t.Errorf("Error calculating deduced memory cell")
	}

	if result != nil {
		t.Errorf("Error: Results differ from nil")
	}
}

func TestDeduceMemoryCellEcOpForPresetMemoryNonIntegerInput(t *testing.T) {
	mem := memory.NewMemorySegmentManager()
	mem.AddSegment()
	mem.AddSegment()
	mem.AddSegment()
	mem.AddSegment()
	mem.Memory.Insert(memory.NewRelocatable(3, 0), memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromDecString("2962412995502985605007699495352191122971573493113767820301112397466445942584")))
	mem.Memory.Insert(memory.NewRelocatable(3, 1), memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromDecString("214950771763870898744428659242275426967582168179217139798831865603966154129")))
	mem.Memory.Insert(memory.NewRelocatable(3, 2), memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromDecString("874739451078007766457464989774322083649278607533249481151382481072868806602")))
	mem.Memory.Insert(memory.NewRelocatable(3, 3), memory.NewMaybeRelocatableRelocatable(memory.NewRelocatable(1, 2)))
	mem.Memory.Insert(memory.NewRelocatable(3, 4), memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(34)))
	mem.Memory.Insert(memory.NewRelocatable(3, 5), memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromDecString("2778063437308421278851140253538604815869848682781135193774472480292420096757")))

	builtin := builtins.NewEcOpBuiltinRunner(true)

	// expected value is an error

	_, err := builtin.DeduceMemoryCell(memory.NewRelocatable(3, 6), &mem.Memory)

	if err == nil {
		t.Errorf("Expected Error but got result")
	}
}

func TestIntegrationEcOp(t *testing.T) {
	t.Helper()
	_, err := cairo_run.CairoRun("../../cairo_programs/ec_op.json", "small", false)
	if err != nil {
		t.Errorf("TestIntegrationBitwise failed with error:\n %v", err)
	}
}
