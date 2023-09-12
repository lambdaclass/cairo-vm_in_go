package builtins_test

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/lambdaclass/cairo-vm.go/pkg/builtins"
	"github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
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

	fmt.Println("result x:   expected:  ", result_x.Text(10), expected_x)
	fmt.Println("comparison x:  ", result_x.Cmp(expected_x))
	fmt.Println("result y:   expected:  ", result_y.Text(10), expected_y)
	fmt.Println("comparison y:  ", result_y.Cmp(expected_y))

	if result_x.Cmp(expected_x) != 0 {
		t.Errorf("Got different X result in Ec On Impl")
	}

	if result_y.Cmp(expected_y) != 0 {
		t.Errorf("Got different Y result in Ec On Impl")
	}
}
