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

func VmNew(run_context vm.RunContext, current_step uint, segments_manager memory.MemorySegmentManager) vm.VirtualMachine {
	return vm.VirtualMachine{
		RunContext:  run_context,
		CurrentStep: current_step,
		Segments:    segments_manager,
	}
}

func TestComputeOperandsAddAp(t *testing.T) {
	instruction := vm.Instruction{
		OffDst:   0,
		OffOp0:   1,
		OffOp1:   2,
		DstReg:   vm.AP,
		Op0Reg:   vm.FP,
		Op1Addr:  vm.Op1SrcAP,
		ResLogic: vm.ResAdd,
		PcUpdate: vm.PcUpdateRegular,
		ApUpdate: vm.ApUpdateRegular,
		FpUpdate: vm.FpUpdateRegular,
		Opcode:   vm.NOp,
	}

	memory_manager := memory.NewMemorySegmentManager()
	run_context := vm.RunContext{
		Ap: memory.NewRelocatable(1, 0),
		Fp: memory.NewRelocatable(1, 0),
		Pc: memory.NewRelocatable(0, 0),
	}
	vmachine := VmNew(run_context, 0, memory_manager)

	for i := 0; i < 2; i++ {
		vmachine.Segments.AddSegment()
	}

	dst_addr := memory.NewRelocatable(1, 0)
	dst_addr_value := memory.NewMaybeRelocatableInt(lambdaworks.FeltFromUint64(5))
	op0_addr := memory.NewRelocatable(1, 1)
	op0_addr_value := memory.NewMaybeRelocatableInt(lambdaworks.FeltFromUint64(2))
	op1_addr := memory.NewRelocatable(1, 2)
	op1_addr_value := memory.NewMaybeRelocatableInt(lambdaworks.FeltFromUint64(3))

	vmachine.Segments.Memory.Insert(dst_addr, dst_addr_value)
	vmachine.Segments.Memory.Insert(op0_addr, op0_addr_value)
	vmachine.Segments.Memory.Insert(op1_addr, op1_addr_value)

	expected_operands := vm.Operands{
		Dst: *dst_addr_value,
		Res: dst_addr_value,
		Op0: *op0_addr_value,
		Op1: *op1_addr_value,
	}

	operands, _, _ := vmachine.ComputeOperands(instruction)

	if operands.Dst != expected_operands.Dst {
		t.Errorf("Different Dst register")
	}
	if operands.Op0 != expected_operands.Op0 {
		t.Errorf("Different op0 register")
	}
	if operands.Op1 != expected_operands.Op1 {
		t.Errorf("Different op1 register")
	}
	if *operands.Res != *expected_operands.Res {
		t.Errorf("Different res register")
	}
}
