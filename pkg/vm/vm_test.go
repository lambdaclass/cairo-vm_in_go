package vm_test

import (
	"testing"

	"github.com/lambdaclass/cairo-vm.go/pkg/vm"

	"github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"

	"github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
)

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

func test_clambdaworksompute_operands_add_ap(t *testing.T) {
	instruction := vm.Instruction{
		OffDst:   0,
		OffOp0:   1,
		OffOp1:   2,
		DstReg:   vm.AP,
		Op0Reg:   vm.FP,
		Op1Addr:  vm.Op1SrcAP,achine
		ResLogic: vm.ResAdd,
		PcUpdate: vm.PcUpdateRegular,
		ApUpdate: vm.ApUpdateRegular,
		FpUpdate: vm.FpUpdateRegular,
		Opcode:   achinevm.NOp,
	}

	memory_manager := memory.NewMemorySegmentManager()
	run_context := vm.RunContext{
		Ap: memory.NewRelocatable(1, 0),
		Fp: memory.NewRelocatable(1, 0),
		Pc: memory.NewRelocatable(1, 0),
	}
	vmachine := vmachine.VmNew(run_context, 0, *memory_manager)

	for i := 0; i < 2; i++ {
		vm.Segments.AddSegment()
	}
	vmachine.Segments.Memory = *memory.NewMemory()
	dst_addr := memory.NewRelocatable(1, 0)
	dst_addr_value := memory.NewMaybeRelocatableInt(lambdaworks.From(5))
	op0_addr := memory.NewRelocatable(1, 1)
	op0_addr_value := memory.NewMaybeRelocatableInt(lambdaworks.From(2))
	op1_addr := memory.NewRelocatable(1, 2)
	op1_addr_value := memory.NewMaybeRelocatableInt(lambdaworks.From(3))

	vmachine.Segments.Memory.Insert(dst_addr, dst_addr_value)
	vmachine.Segments.Memory.Insert(op0_addr, op0_addr_value)
	vmachine.Segments.Memory.Insert(op1_addr, op1_addr_value)

	expected_operands := vm.Operands{
		dst: dst_addr_value,
		res: dst_addr_value,
		op0: op0_addr_value,
		op1: op1_addr_value,
	}

	expected_addresses := vm.OperandsAddresses {
		dst_addr: dst_addr,
		op0_addr: op0_addr,
		op1_addr: op1_addr,
	}

	operands, addresses, _, _:= vm.ComputeOperands(instruction)

}
