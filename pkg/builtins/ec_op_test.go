package builtins_test

import (
	"testing"

	"github.com/lambdaclass/cairo-vm.go/pkg/builtins"
	"github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
)

func PointIsOnCurveB(t *testing.T) {

	x := lambdaworks.FeltFromDecString("3139037544796708144595053687182055617920475701120786241351436619796497072089")
	y := lambdaworks.FeltFromDecString("2119589567875935397690285099786081818522144748339117565577200220779667999801")
	alpha := lambdaworks.FeltOne()
	beta := lambdaworks.FeltFromDecString("3141592653589793238462643383279502884197169399375105820974944592307816406665")

	if !builtins.PointOnCurve(x, y, alpha, beta) {
		t.Errorf("The point is not on the curve")
	}
}

func PointIsNotOnCurveB(t *testing.T) {

	x := lambdaworks.FeltFromDecString("3139037544756708144595053687182055617927475701120786241351436619796497072089")
	y := lambdaworks.FeltFromDecString("2119589567875935397690885099786081818522144748339117565577200220779667999801")
	alpha := lambdaworks.FeltOne()
	beta := lambdaworks.FeltFromDecString("3141592653589793238462643383279502884197169399375105820974944592307816406665")

	if builtins.PointOnCurve(x, y, alpha, beta) {
		t.Errorf("The point should not be on the curve")
	}
}
