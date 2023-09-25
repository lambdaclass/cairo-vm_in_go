package hint_utils_test

import (
	"math/big"
	"testing"

	. "github.com/lambdaclass/cairo-vm.go/pkg/types"
	. "github.com/lambdaclass/cairo-vm.go/pkg/vm"
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
}