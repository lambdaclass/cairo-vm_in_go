package vm_test

import (
	"reflect"
	"testing"

	"github.com/lambdaclass/cairo-vm.go/pkg/vm"
	"github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
)

func TestUpdatePcRegularNoImm(t *testing.T) {
	instruction := vm.Instruction{PcUpdate: vm.PcUpdateRegular, Op1Addr: vm.Op1SrcAP}
	operands := vm.Operands{}
	vm := vm.NewVirtualMachine()
	err := vm.UpdatePc(&instruction, &operands)
	if err != nil {
		t.Errorf("UpdatePc failed with error: %s", err)
	}
	if !reflect.DeepEqual(vm.RunContext.Pc, memory.Relocatable{SegmentIndex: 0, Offset: 1}) {
		t.Errorf("Wrong value after pc update")
	}

}

func TestUpdatePcRegularWithImm(t *testing.T) {
	instruction := vm.Instruction{PcUpdate: vm.PcUpdateRegular, Op1Addr: vm.Op1SrcImm}
	operands := vm.Operands{}
	vm := vm.NewVirtualMachine()
	err := vm.UpdatePc(&instruction, &operands)
	if err != nil {
		t.Errorf("UpdatePc failed with error: %s", err)
	}
	if !reflect.DeepEqual(vm.RunContext.Pc, memory.Relocatable{SegmentIndex: 0, Offset: 2}) {
		t.Errorf("Wrong value after pc update")
	}

}

// Things we are skipping for now:
// - Initializing hint_executor and passing it to `cairo_run`
// - cairo_run_config stuff
// - Asserting expected trace values
// - Asserting memory_holes
func TestFibonacci(t *testing.T) {
	// compiledProgram := parser.Parse("../../cairo_programs/fibonacci.json")

	// TODO: Uncomment test when we have the bare minimum `CairoRun`
	// err := vm.CairoRun(compiledProgram.Data)
	// if err != nil {
	// 	t.Errorf("Program execution failed with error: %s", err)
	// }
}
